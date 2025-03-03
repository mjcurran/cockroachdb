// Copyright 2022 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

// Code generated by sctestgen, DO NOT EDIT.

package schemachanger

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/sctest"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

func TestEndToEndSideEffects_add_column(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.EndToEndSideEffects(t, "pkg/sql/schemachanger/testdata/end_to_end/add_column", sctest.SingleNodeCluster)
}
func TestExecuteWithDMLInjection_add_column(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.ExecuteWithDMLInjection(t, "pkg/sql/schemachanger/testdata/end_to_end/add_column", sctest.SingleNodeCluster)
}
func TestGenerateSchemaChangeCorpus_add_column(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.GenerateSchemaChangeCorpus(t, "pkg/sql/schemachanger/testdata/end_to_end/add_column", sctest.SingleNodeCluster)
}
func TestPause_add_column(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Pause(t, "pkg/sql/schemachanger/testdata/end_to_end/add_column", sctest.SingleNodeCluster)
}
func TestRollback_add_column(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Rollback(t, "pkg/sql/schemachanger/testdata/end_to_end/add_column", sctest.SingleNodeCluster)
}
func TestEndToEndSideEffects_add_column_default_seq(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.EndToEndSideEffects(t, "pkg/sql/schemachanger/testdata/end_to_end/add_column_default_seq", sctest.SingleNodeCluster)
}
func TestExecuteWithDMLInjection_add_column_default_seq(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.ExecuteWithDMLInjection(t, "pkg/sql/schemachanger/testdata/end_to_end/add_column_default_seq", sctest.SingleNodeCluster)
}
func TestGenerateSchemaChangeCorpus_add_column_default_seq(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.GenerateSchemaChangeCorpus(t, "pkg/sql/schemachanger/testdata/end_to_end/add_column_default_seq", sctest.SingleNodeCluster)
}
func TestPause_add_column_default_seq(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Pause(t, "pkg/sql/schemachanger/testdata/end_to_end/add_column_default_seq", sctest.SingleNodeCluster)
}
func TestRollback_add_column_default_seq(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Rollback(t, "pkg/sql/schemachanger/testdata/end_to_end/add_column_default_seq", sctest.SingleNodeCluster)
}
func TestEndToEndSideEffects_add_column_default_unique(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.EndToEndSideEffects(t, "pkg/sql/schemachanger/testdata/end_to_end/add_column_default_unique", sctest.SingleNodeCluster)
}
func TestExecuteWithDMLInjection_add_column_default_unique(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.ExecuteWithDMLInjection(t, "pkg/sql/schemachanger/testdata/end_to_end/add_column_default_unique", sctest.SingleNodeCluster)
}
func TestGenerateSchemaChangeCorpus_add_column_default_unique(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.GenerateSchemaChangeCorpus(t, "pkg/sql/schemachanger/testdata/end_to_end/add_column_default_unique", sctest.SingleNodeCluster)
}
func TestPause_add_column_default_unique(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Pause(t, "pkg/sql/schemachanger/testdata/end_to_end/add_column_default_unique", sctest.SingleNodeCluster)
}
func TestRollback_add_column_default_unique(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Rollback(t, "pkg/sql/schemachanger/testdata/end_to_end/add_column_default_unique", sctest.SingleNodeCluster)
}
func TestEndToEndSideEffects_add_column_no_default(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.EndToEndSideEffects(t, "pkg/sql/schemachanger/testdata/end_to_end/add_column_no_default", sctest.SingleNodeCluster)
}
func TestExecuteWithDMLInjection_add_column_no_default(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.ExecuteWithDMLInjection(t, "pkg/sql/schemachanger/testdata/end_to_end/add_column_no_default", sctest.SingleNodeCluster)
}
func TestGenerateSchemaChangeCorpus_add_column_no_default(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.GenerateSchemaChangeCorpus(t, "pkg/sql/schemachanger/testdata/end_to_end/add_column_no_default", sctest.SingleNodeCluster)
}
func TestPause_add_column_no_default(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Pause(t, "pkg/sql/schemachanger/testdata/end_to_end/add_column_no_default", sctest.SingleNodeCluster)
}
func TestRollback_add_column_no_default(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Rollback(t, "pkg/sql/schemachanger/testdata/end_to_end/add_column_no_default", sctest.SingleNodeCluster)
}
func TestEndToEndSideEffects_alter_table_add_check_vanilla(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.EndToEndSideEffects(t, "pkg/sql/schemachanger/testdata/end_to_end/alter_table_add_check_vanilla", sctest.SingleNodeCluster)
}
func TestExecuteWithDMLInjection_alter_table_add_check_vanilla(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.ExecuteWithDMLInjection(t, "pkg/sql/schemachanger/testdata/end_to_end/alter_table_add_check_vanilla", sctest.SingleNodeCluster)
}
func TestGenerateSchemaChangeCorpus_alter_table_add_check_vanilla(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.GenerateSchemaChangeCorpus(t, "pkg/sql/schemachanger/testdata/end_to_end/alter_table_add_check_vanilla", sctest.SingleNodeCluster)
}
func TestPause_alter_table_add_check_vanilla(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Pause(t, "pkg/sql/schemachanger/testdata/end_to_end/alter_table_add_check_vanilla", sctest.SingleNodeCluster)
}
func TestRollback_alter_table_add_check_vanilla(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Rollback(t, "pkg/sql/schemachanger/testdata/end_to_end/alter_table_add_check_vanilla", sctest.SingleNodeCluster)
}
func TestEndToEndSideEffects_alter_table_add_check_with_seq_and_udt(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.EndToEndSideEffects(t, "pkg/sql/schemachanger/testdata/end_to_end/alter_table_add_check_with_seq_and_udt", sctest.SingleNodeCluster)
}
func TestExecuteWithDMLInjection_alter_table_add_check_with_seq_and_udt(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.ExecuteWithDMLInjection(t, "pkg/sql/schemachanger/testdata/end_to_end/alter_table_add_check_with_seq_and_udt", sctest.SingleNodeCluster)
}
func TestGenerateSchemaChangeCorpus_alter_table_add_check_with_seq_and_udt(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.GenerateSchemaChangeCorpus(t, "pkg/sql/schemachanger/testdata/end_to_end/alter_table_add_check_with_seq_and_udt", sctest.SingleNodeCluster)
}
func TestPause_alter_table_add_check_with_seq_and_udt(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Pause(t, "pkg/sql/schemachanger/testdata/end_to_end/alter_table_add_check_with_seq_and_udt", sctest.SingleNodeCluster)
}
func TestRollback_alter_table_add_check_with_seq_and_udt(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Rollback(t, "pkg/sql/schemachanger/testdata/end_to_end/alter_table_add_check_with_seq_and_udt", sctest.SingleNodeCluster)
}
func TestEndToEndSideEffects_alter_table_add_primary_key_drop_rowid(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.EndToEndSideEffects(t, "pkg/sql/schemachanger/testdata/end_to_end/alter_table_add_primary_key_drop_rowid", sctest.SingleNodeCluster)
}
func TestExecuteWithDMLInjection_alter_table_add_primary_key_drop_rowid(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.ExecuteWithDMLInjection(t, "pkg/sql/schemachanger/testdata/end_to_end/alter_table_add_primary_key_drop_rowid", sctest.SingleNodeCluster)
}
func TestGenerateSchemaChangeCorpus_alter_table_add_primary_key_drop_rowid(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.GenerateSchemaChangeCorpus(t, "pkg/sql/schemachanger/testdata/end_to_end/alter_table_add_primary_key_drop_rowid", sctest.SingleNodeCluster)
}
func TestPause_alter_table_add_primary_key_drop_rowid(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Pause(t, "pkg/sql/schemachanger/testdata/end_to_end/alter_table_add_primary_key_drop_rowid", sctest.SingleNodeCluster)
}
func TestRollback_alter_table_add_primary_key_drop_rowid(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Rollback(t, "pkg/sql/schemachanger/testdata/end_to_end/alter_table_add_primary_key_drop_rowid", sctest.SingleNodeCluster)
}
func TestEndToEndSideEffects_alter_table_alter_primary_key_drop_rowid(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.EndToEndSideEffects(t, "pkg/sql/schemachanger/testdata/end_to_end/alter_table_alter_primary_key_drop_rowid", sctest.SingleNodeCluster)
}
func TestExecuteWithDMLInjection_alter_table_alter_primary_key_drop_rowid(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.ExecuteWithDMLInjection(t, "pkg/sql/schemachanger/testdata/end_to_end/alter_table_alter_primary_key_drop_rowid", sctest.SingleNodeCluster)
}
func TestGenerateSchemaChangeCorpus_alter_table_alter_primary_key_drop_rowid(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.GenerateSchemaChangeCorpus(t, "pkg/sql/schemachanger/testdata/end_to_end/alter_table_alter_primary_key_drop_rowid", sctest.SingleNodeCluster)
}
func TestPause_alter_table_alter_primary_key_drop_rowid(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Pause(t, "pkg/sql/schemachanger/testdata/end_to_end/alter_table_alter_primary_key_drop_rowid", sctest.SingleNodeCluster)
}
func TestRollback_alter_table_alter_primary_key_drop_rowid(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Rollback(t, "pkg/sql/schemachanger/testdata/end_to_end/alter_table_alter_primary_key_drop_rowid", sctest.SingleNodeCluster)
}
func TestEndToEndSideEffects_alter_table_alter_primary_key_vanilla(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.EndToEndSideEffects(t, "pkg/sql/schemachanger/testdata/end_to_end/alter_table_alter_primary_key_vanilla", sctest.SingleNodeCluster)
}
func TestExecuteWithDMLInjection_alter_table_alter_primary_key_vanilla(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.ExecuteWithDMLInjection(t, "pkg/sql/schemachanger/testdata/end_to_end/alter_table_alter_primary_key_vanilla", sctest.SingleNodeCluster)
}
func TestGenerateSchemaChangeCorpus_alter_table_alter_primary_key_vanilla(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.GenerateSchemaChangeCorpus(t, "pkg/sql/schemachanger/testdata/end_to_end/alter_table_alter_primary_key_vanilla", sctest.SingleNodeCluster)
}
func TestPause_alter_table_alter_primary_key_vanilla(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Pause(t, "pkg/sql/schemachanger/testdata/end_to_end/alter_table_alter_primary_key_vanilla", sctest.SingleNodeCluster)
}
func TestRollback_alter_table_alter_primary_key_vanilla(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Rollback(t, "pkg/sql/schemachanger/testdata/end_to_end/alter_table_alter_primary_key_vanilla", sctest.SingleNodeCluster)
}
func TestEndToEndSideEffects_create_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.EndToEndSideEffects(t, "pkg/sql/schemachanger/testdata/end_to_end/create_index", sctest.SingleNodeCluster)
}
func TestExecuteWithDMLInjection_create_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.ExecuteWithDMLInjection(t, "pkg/sql/schemachanger/testdata/end_to_end/create_index", sctest.SingleNodeCluster)
}
func TestGenerateSchemaChangeCorpus_create_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.GenerateSchemaChangeCorpus(t, "pkg/sql/schemachanger/testdata/end_to_end/create_index", sctest.SingleNodeCluster)
}
func TestPause_create_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Pause(t, "pkg/sql/schemachanger/testdata/end_to_end/create_index", sctest.SingleNodeCluster)
}
func TestRollback_create_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Rollback(t, "pkg/sql/schemachanger/testdata/end_to_end/create_index", sctest.SingleNodeCluster)
}
func TestEndToEndSideEffects_drop_column_basic(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.EndToEndSideEffects(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_column_basic", sctest.SingleNodeCluster)
}
func TestExecuteWithDMLInjection_drop_column_basic(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.ExecuteWithDMLInjection(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_column_basic", sctest.SingleNodeCluster)
}
func TestGenerateSchemaChangeCorpus_drop_column_basic(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.GenerateSchemaChangeCorpus(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_column_basic", sctest.SingleNodeCluster)
}
func TestPause_drop_column_basic(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Pause(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_column_basic", sctest.SingleNodeCluster)
}
func TestRollback_drop_column_basic(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Rollback(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_column_basic", sctest.SingleNodeCluster)
}
func TestEndToEndSideEffects_drop_column_computed_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.EndToEndSideEffects(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_column_computed_index", sctest.SingleNodeCluster)
}
func TestExecuteWithDMLInjection_drop_column_computed_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.ExecuteWithDMLInjection(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_column_computed_index", sctest.SingleNodeCluster)
}
func TestGenerateSchemaChangeCorpus_drop_column_computed_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.GenerateSchemaChangeCorpus(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_column_computed_index", sctest.SingleNodeCluster)
}
func TestPause_drop_column_computed_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Pause(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_column_computed_index", sctest.SingleNodeCluster)
}
func TestRollback_drop_column_computed_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Rollback(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_column_computed_index", sctest.SingleNodeCluster)
}
func TestEndToEndSideEffects_drop_column_create_index_separate_statements(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.EndToEndSideEffects(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_column_create_index_separate_statements", sctest.SingleNodeCluster)
}
func TestExecuteWithDMLInjection_drop_column_create_index_separate_statements(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.ExecuteWithDMLInjection(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_column_create_index_separate_statements", sctest.SingleNodeCluster)
}
func TestGenerateSchemaChangeCorpus_drop_column_create_index_separate_statements(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.GenerateSchemaChangeCorpus(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_column_create_index_separate_statements", sctest.SingleNodeCluster)
}
func TestPause_drop_column_create_index_separate_statements(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Pause(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_column_create_index_separate_statements", sctest.SingleNodeCluster)
}
func TestRollback_drop_column_create_index_separate_statements(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Rollback(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_column_create_index_separate_statements", sctest.SingleNodeCluster)
}
func TestEndToEndSideEffects_drop_column_unique_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.EndToEndSideEffects(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_column_unique_index", sctest.SingleNodeCluster)
}
func TestExecuteWithDMLInjection_drop_column_unique_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.ExecuteWithDMLInjection(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_column_unique_index", sctest.SingleNodeCluster)
}
func TestGenerateSchemaChangeCorpus_drop_column_unique_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.GenerateSchemaChangeCorpus(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_column_unique_index", sctest.SingleNodeCluster)
}
func TestPause_drop_column_unique_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Pause(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_column_unique_index", sctest.SingleNodeCluster)
}
func TestRollback_drop_column_unique_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Rollback(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_column_unique_index", sctest.SingleNodeCluster)
}
func TestEndToEndSideEffects_drop_column_with_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.EndToEndSideEffects(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_column_with_index", sctest.SingleNodeCluster)
}
func TestExecuteWithDMLInjection_drop_column_with_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.ExecuteWithDMLInjection(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_column_with_index", sctest.SingleNodeCluster)
}
func TestGenerateSchemaChangeCorpus_drop_column_with_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.GenerateSchemaChangeCorpus(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_column_with_index", sctest.SingleNodeCluster)
}
func TestPause_drop_column_with_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Pause(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_column_with_index", sctest.SingleNodeCluster)
}
func TestRollback_drop_column_with_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Rollback(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_column_with_index", sctest.SingleNodeCluster)
}
func TestEndToEndSideEffects_drop_index_hash_sharded_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.EndToEndSideEffects(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_index_hash_sharded_index", sctest.SingleNodeCluster)
}
func TestExecuteWithDMLInjection_drop_index_hash_sharded_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.ExecuteWithDMLInjection(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_index_hash_sharded_index", sctest.SingleNodeCluster)
}
func TestGenerateSchemaChangeCorpus_drop_index_hash_sharded_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.GenerateSchemaChangeCorpus(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_index_hash_sharded_index", sctest.SingleNodeCluster)
}
func TestPause_drop_index_hash_sharded_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Pause(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_index_hash_sharded_index", sctest.SingleNodeCluster)
}
func TestRollback_drop_index_hash_sharded_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Rollback(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_index_hash_sharded_index", sctest.SingleNodeCluster)
}
func TestEndToEndSideEffects_drop_index_partial_expression_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.EndToEndSideEffects(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_index_partial_expression_index", sctest.SingleNodeCluster)
}
func TestExecuteWithDMLInjection_drop_index_partial_expression_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.ExecuteWithDMLInjection(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_index_partial_expression_index", sctest.SingleNodeCluster)
}
func TestGenerateSchemaChangeCorpus_drop_index_partial_expression_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.GenerateSchemaChangeCorpus(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_index_partial_expression_index", sctest.SingleNodeCluster)
}
func TestPause_drop_index_partial_expression_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Pause(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_index_partial_expression_index", sctest.SingleNodeCluster)
}
func TestRollback_drop_index_partial_expression_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Rollback(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_index_partial_expression_index", sctest.SingleNodeCluster)
}
func TestEndToEndSideEffects_drop_index_vanilla_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.EndToEndSideEffects(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_index_vanilla_index", sctest.SingleNodeCluster)
}
func TestExecuteWithDMLInjection_drop_index_vanilla_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.ExecuteWithDMLInjection(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_index_vanilla_index", sctest.SingleNodeCluster)
}
func TestGenerateSchemaChangeCorpus_drop_index_vanilla_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.GenerateSchemaChangeCorpus(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_index_vanilla_index", sctest.SingleNodeCluster)
}
func TestPause_drop_index_vanilla_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Pause(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_index_vanilla_index", sctest.SingleNodeCluster)
}
func TestRollback_drop_index_vanilla_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Rollback(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_index_vanilla_index", sctest.SingleNodeCluster)
}
func TestEndToEndSideEffects_drop_multiple_columns_separate_statements(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.EndToEndSideEffects(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_multiple_columns_separate_statements", sctest.SingleNodeCluster)
}
func TestExecuteWithDMLInjection_drop_multiple_columns_separate_statements(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.ExecuteWithDMLInjection(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_multiple_columns_separate_statements", sctest.SingleNodeCluster)
}
func TestGenerateSchemaChangeCorpus_drop_multiple_columns_separate_statements(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.GenerateSchemaChangeCorpus(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_multiple_columns_separate_statements", sctest.SingleNodeCluster)
}
func TestPause_drop_multiple_columns_separate_statements(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Pause(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_multiple_columns_separate_statements", sctest.SingleNodeCluster)
}
func TestRollback_drop_multiple_columns_separate_statements(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Rollback(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_multiple_columns_separate_statements", sctest.SingleNodeCluster)
}
func TestEndToEndSideEffects_drop_schema(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.EndToEndSideEffects(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_schema", sctest.SingleNodeCluster)
}
func TestExecuteWithDMLInjection_drop_schema(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.ExecuteWithDMLInjection(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_schema", sctest.SingleNodeCluster)
}
func TestGenerateSchemaChangeCorpus_drop_schema(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.GenerateSchemaChangeCorpus(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_schema", sctest.SingleNodeCluster)
}
func TestPause_drop_schema(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Pause(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_schema", sctest.SingleNodeCluster)
}
func TestRollback_drop_schema(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Rollback(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_schema", sctest.SingleNodeCluster)
}
func TestEndToEndSideEffects_drop_table(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.EndToEndSideEffects(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_table", sctest.SingleNodeCluster)
}
func TestExecuteWithDMLInjection_drop_table(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.ExecuteWithDMLInjection(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_table", sctest.SingleNodeCluster)
}
func TestGenerateSchemaChangeCorpus_drop_table(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.GenerateSchemaChangeCorpus(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_table", sctest.SingleNodeCluster)
}
func TestPause_drop_table(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Pause(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_table", sctest.SingleNodeCluster)
}
func TestRollback_drop_table(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	sctest.Rollback(t, "pkg/sql/schemachanger/testdata/end_to_end/drop_table", sctest.SingleNodeCluster)
}
