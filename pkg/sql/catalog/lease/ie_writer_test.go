// Copyright 2021 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package lease

import (
	"context"
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/kv"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/systemschema"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlutil"
	"github.com/cockroachdb/errors"
)

type ieWriter struct {
	insertQuery string
	deleteQuery string
	ie          sqlutil.InternalExecutor
}

func newInternalExecutorWriter(ie sqlutil.InternalExecutor, tableName string) *ieWriter {
	if systemschema.TestSupportMultiRegion() {
		const (
			deleteLease = `
DELETE FROM %s
      WHERE (crdb_region, "descID", version, "nodeID", expiration)
            = ($1, $2, $3, $4, $5);`
			insertLease = `
INSERT
  INTO %s (crdb_region, "descID", version, "nodeID", expiration)
VALUES ($1, $2, $3, $4, $5)`
		)
		return &ieWriter{
			ie:          ie,
			insertQuery: fmt.Sprintf(insertLease, tableName),
			deleteQuery: fmt.Sprintf(deleteLease, tableName),
		}
	}
	const (
		deleteLease = `
DELETE FROM %s
      WHERE ("descID", version, "nodeID", expiration)
            = ($1, $2, $3, $4);`
		insertLease = `
INSERT
  INTO %s ("descID", version, "nodeID", expiration)
VALUES ($1, $2, $3, $4)`
	)
	return &ieWriter{
		ie:          ie,
		insertQuery: fmt.Sprintf(insertLease, tableName),
		deleteQuery: fmt.Sprintf(deleteLease, tableName),
	}
}

func (w *ieWriter) deleteLease(ctx context.Context, txn *kv.Txn, l leaseFields) error {
	if systemschema.TestSupportMultiRegion() {
		_, err := w.ie.Exec(
			ctx,
			"lease-release",
			nil, /* txn */
			w.deleteQuery,
			l.regionPrefix, l.descID, l.version, l.instanceID, &l.expiration,
		)
		return err
	}
	_, err := w.ie.Exec(
		ctx,
		"lease-release",
		nil, /* txn */
		w.deleteQuery,
		l.descID, l.version, l.instanceID, &l.expiration,
	)
	return err
}

func (w *ieWriter) insertLease(ctx context.Context, txn *kv.Txn, l leaseFields) error {
	if systemschema.TestSupportMultiRegion() {
		count, err := w.ie.Exec(ctx, "lease-insert", txn, w.insertQuery,
			l.regionPrefix, l.descID, l.version, l.instanceID, &l.expiration,
		)
		if err != nil {
			return err
		}
		if count != 1 {
			return errors.Errorf("%s: expected 1 result, found %d", w.insertQuery, count)
		}
		return nil
	}
	count, err := w.ie.Exec(ctx, "lease-insert", txn, w.insertQuery,
		l.descID, l.version, l.instanceID, &l.expiration,
	)
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.Errorf("%s: expected 1 result, found %d", w.insertQuery, count)
	}
	return nil
}

var _ writer = (*ieWriter)(nil)
