// Copyright 2021 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package scbuild_test

import (
	"bufio"
	"context"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/ccl"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/parser"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/rel"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scbuild"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scdeps/sctestdeps"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scdeps/sctestutils"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scerrors"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scpb"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/screl"
	"github.com/cockroachdb/cockroach/pkg/sql/sessiondata"
	"github.com/cockroachdb/cockroach/pkg/sql/sessiondatapb"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/serverutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/sqlutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/datadriven"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestBuildDataDriven(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	defer ccl.TestingEnableEnterprise()()

	ctx := context.Background()

	datadriven.Walk(t, testutils.TestDataPath(t), func(t *testing.T, path string) {
		for _, depsType := range []struct {
			name                string
			dependenciesWrapper func(*testing.T, serverutils.TestServerInterface, *sqlutils.SQLRunner, func(scbuild.Dependencies))
		}{
			{
				name: "sql_dependencies",
				dependenciesWrapper: func(t *testing.T, s serverutils.TestServerInterface, tdb *sqlutils.SQLRunner, fn func(scbuild.Dependencies)) {
					sctestutils.WithBuilderDependenciesFromTestServer(s, fn)
				},
			},
			{
				name: "test_dependencies",
				dependenciesWrapper: func(t *testing.T, s serverutils.TestServerInterface, tdb *sqlutils.SQLRunner, fn func(scbuild.Dependencies)) {
					// Create test dependencies and execute the schema changer.
					// The schema changer test dependencies do not hold any reference to the
					// test cluster, here the SQLRunner is only used to populate the mocked
					// catalog state.
					descriptorCatalog := sctestdeps.ReadDescriptorsFromDB(ctx, t, tdb).Catalog
					fn(
						sctestdeps.NewTestDependencies(
							sctestdeps.WithDescriptors(descriptorCatalog),
							sctestdeps.WithNamespace(sctestdeps.ReadNamespaceFromDB(t, tdb).Catalog),
							sctestdeps.WithCurrentDatabase(sctestdeps.ReadCurrentDatabaseFromDB(t, tdb)),
							sctestdeps.WithSessionData(
								sctestdeps.ReadSessionDataFromDB(
									t,
									tdb,
									func(sd *sessiondata.SessionData) {
										// For setting up a builder inside tests we will ensure that the new schema
										// changer will allow non-fully implemented operations.
										sd.NewSchemaChangerMode = sessiondatapb.UseNewSchemaChangerUnsafe
										sd.ApplicationName = ""
									},
								),
							),
							sctestdeps.WithComments(sctestdeps.ReadCommentsFromDB(t, tdb)),
							sctestdeps.WithZoneConfigs(sctestdeps.ReadZoneConfigsFromDB(t, tdb, descriptorCatalog)),
						),
					)
				},
			},
		} {
			t.Run(depsType.name, func(t *testing.T) {
				s, sqlDB, _ := serverutils.StartServer(t, base.TestServerArgs{
					DisableDefaultTestTenant: true,
				})
				defer s.Stopper().Stop(ctx)
				tdb := sqlutils.MakeSQLRunner(sqlDB)
				datadriven.RunTest(t, path, func(t *testing.T, d *datadriven.TestData) string {
					return run(ctx, t, depsType.name, d, s, tdb, depsType.dependenciesWrapper)
				})
			})
		}
	})
}

func run(
	ctx context.Context,
	t *testing.T,
	depsTypeName string,
	d *datadriven.TestData,
	s serverutils.TestServerInterface,
	tdb *sqlutils.SQLRunner,
	withDependencies func(*testing.T, serverutils.TestServerInterface, *sqlutils.SQLRunner, func(scbuild.Dependencies)),
) string {
	sqlutils.VerifyStatementPrettyRoundtrip(t, d.Input)
	switch d.Cmd {
	case "setup":
		stmts, err := parser.Parse(d.Input)
		require.NoError(t, err)
		for _, stmt := range stmts {
			tableName := sctestutils.TableNameFromStmt(stmt)
			tdb.Exec(t, stmt.SQL)
			if len(tableName) > 0 {
				var tableID descpb.ID
				tdb.QueryRow(t, fmt.Sprintf(`SELECT '%s'::REGCLASS::INT`, tableName)).Scan(&tableID)
				if tableID == 0 {
					d.Fatalf(t, "failed to read ID of new table %s", tableName)
				}
				t.Logf("created relation with id %d", tableID)
			}
		}
		return ""
	case "build":
		if a := d.CmdArgs; len(a) > 0 && a[0].Key == "skip" {
			for _, v := range a[0].Vals {
				if v == depsTypeName {
					return d.Expected
				}
			}
		}
		var output scpb.CurrentState
		withDependencies(t, s, tdb, func(deps scbuild.Dependencies) {
			stmts, err := parser.Parse(d.Input)
			require.NoError(t, err)
			for i := range stmts {
				output, err = scbuild.Build(ctx, deps, output, stmts[i].AST)
				require.NoErrorf(t, err, "%s: %s", d.Pos, stmts[i].SQL)
			}
		})
		return marshalState(t, output)

	case "unimplemented":
		withDependencies(t, s, tdb, func(deps scbuild.Dependencies) {
			stmts, err := parser.Parse(d.Input)
			require.NoError(t, err)
			require.NotEmpty(t, stmts)

			for _, stmt := range stmts {
				_, err = scbuild.Build(ctx, deps, scpb.CurrentState{}, stmt.AST)
				expected := scerrors.NotImplementedError(nil)
				require.Errorf(t, err, "%s: expected %T instead of success for", stmt.SQL, expected)
				require.Truef(t, scerrors.HasNotImplemented(err), "%s: expected %T instead of %v", stmt.SQL, expected, err)
			}
		})
		return ""

	default:
		return fmt.Sprintf("unknown command: %s", d.Cmd)
	}
}

// indentText indents text for formatting out marshaled data.
func indentText(input string, tab string) string {
	result := strings.Builder{}
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		result.WriteString(tab)
		result.WriteString(line)
		result.WriteString("\n")
	}
	return result.String()
}

// marshalState marshals a scpb.CurrentState to YAML.
func marshalState(t *testing.T, state scpb.CurrentState) string {
	var sortedEntries nodeEntries
	for i, status := range state.Current {
		node := screl.Node{
			Target:        &state.Targets[i],
			CurrentStatus: status,
		}

		entry := strings.Builder{}
		entry.WriteString("- ")
		entry.WriteString(screl.NodeString(&node))
		entry.WriteString("\n")
		entry.WriteString(indentText(string(formatElementForDisplay(t, node.Element())), "  "))
		sortedEntries = append(sortedEntries, nodeEntry{
			node:  node,
			entry: entry.String(),
		})
	}
	// Sort the output buffer of state for determinism.
	result := strings.Builder{}
	sort.Sort(sortedEntries)
	for _, entry := range sortedEntries {
		result.WriteString(entry.entry)
	}
	return result.String()
}

type nodeEntry struct {
	node  screl.Node
	entry string
}

type nodeEntries []nodeEntry

func (n nodeEntries) Len() int { return len(n) }

func (n nodeEntries) Less(i, j int) bool {
	less, _ := screl.Schema.CompareOn([]rel.Attr{
		screl.DescID, screl.ColumnID, screl.IndexID, screl.ConstraintID,
		screl.ColumnFamilyID, screl.ReferencedDescID,
	}, &n[i].node, &n[j].node)
	return less
}

func (n nodeEntries) Swap(i, j int) { n[i], n[j] = n[j], n[i] }

var _ sort.Interface = nodeEntries{}

func formatElementForDisplay(t *testing.T, e scpb.Element) []byte {
	marshaled, err := sctestutils.ProtoToYAML(
		e, false /* emitDefaults */, nil, /* rewrites */
	)
	require.NoError(t, err)
	dec := yaml.NewDecoder(strings.NewReader(marshaled))
	dec.KnownFields(true)
	var n yaml.Node
	require.NoError(t, dec.Decode(&n))
	walkYaml(&n, func(node *yaml.Node) { node.Style = yaml.FlowStyle })
	data, err := yaml.Marshal(&n)
	require.NoError(t, err)
	return data
}

func walkYaml(root *yaml.Node, f func(node *yaml.Node)) {
	var walk func(node *yaml.Node)
	walk = func(node *yaml.Node) {
		f(node)
		for _, child := range node.Content {
			walk(child)
		}
	}
	walk(root)
}
