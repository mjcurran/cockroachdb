// Copyright 2021 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package ssmemstorage

import (
	"context"
	"time"
	"unsafe"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/sql/execstats"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlstats"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlstats/insights"
	"github.com/cockroachdb/cockroach/pkg/util"
	"github.com/cockroachdb/errors"
)

var (
	// ErrMemoryPressure is returned from the Container when we have reached
	// the memory limit allowed.
	ErrMemoryPressure = errors.New("insufficient sql stats memory")

	// ErrFingerprintLimitReached is returned from the Container when we have
	// more fingerprints than the limit specified in the cluster setting.
	ErrFingerprintLimitReached = errors.New("sql stats fingerprint limit reached")

	// ErrExecStatsFingerprintFlushed is returned from the Container when the
	// stats object for the fingerprint has been flushed to system table before
	// the roachpb.ExecStats can be recorded.
	ErrExecStatsFingerprintFlushed = errors.New("stmtStats flushed before execution stats can be recorded")
)

var timestampSize = int64(unsafe.Sizeof(time.Time{}))

var _ sqlstats.Writer = &Container{}

func getStatus(statementError error) insights.Statement_Status {
	if statementError == nil {
		return insights.Statement_Completed
	}

	return insights.Statement_Failed
}

// RecordStatement implements sqlstats.Writer interface.
// RecordStatement saves per-statement statistics.
//
// samplePlanDescription can be nil, as these are only sampled periodically
// per unique fingerprint.
// RecordStatement always returns a valid stmtFingerprintID corresponding to the given
// stmt regardless of whether the statement is actually recorded or not.
//
// If the statement is not actually recorded due to either:
// 1. the memory budget has been exceeded
// 2. the unique statement fingerprint limit has been exceeded
// and error is being returned.
// Note: This error is only related to the operation of recording the statement
// statistics into in-memory structs. It is unrelated to the stmtErr in the
// arguments.
func (s *Container) RecordStatement(
	ctx context.Context, key roachpb.StatementStatisticsKey, value sqlstats.RecordedStmtStats,
) (roachpb.StmtFingerprintID, error) {
	createIfNonExistent := true
	// If the statement is below the latency threshold, or stats aren't being
	// recorded we don't need to create an entry in the stmts map for it. We do
	// still need stmtFingerprintID for transaction level metrics tracking.
	t := sqlstats.StatsCollectionLatencyThreshold.Get(&s.st.SV)
	if !sqlstats.StmtStatsEnable.Get(&s.st.SV) || (t > 0 && t.Seconds() >= value.ServiceLatency) {
		createIfNonExistent = false
	}

	// Get the statistics object.
	stats, statementKey, stmtFingerprintID, created, throttled := s.getStatsForStmt(
		key.Query,
		key.ImplicitTxn,
		key.Database,
		key.Failed,
		key.PlanHash,
		key.TransactionFingerprintID,
		createIfNonExistent,
	)

	// This means we have reached the limit of unique fingerprintstats. We don't
	// record anything and abort the operation.
	if throttled {
		return stmtFingerprintID, ErrFingerprintLimitReached
	}

	// This statement was below the latency threshold or sql stats aren't being
	// recorded. Either way, we don't need to record anything in the stats object
	// for this statement, though we do need to return the statement fingerprint ID for
	// transaction level metrics collection.
	if !createIfNonExistent {
		return stmtFingerprintID, nil
	}

	// Collect the per-statement statisticstats.
	stats.mu.Lock()
	defer stats.mu.Unlock()

	stats.mu.data.Count++
	if key.Failed {
		stats.mu.data.SensitiveInfo.LastErr = value.StatementError.Error()
	}
	// Only update MostRecentPlanDescription if we sampled a new PlanDescription.
	if value.Plan != nil {
		stats.mu.data.SensitiveInfo.MostRecentPlanDescription = *value.Plan
		stats.mu.data.SensitiveInfo.MostRecentPlanTimestamp = s.getTimeNow()
		s.setLogicalPlanLastSampled(statementKey.sampledPlanKey, stats.mu.data.SensitiveInfo.MostRecentPlanTimestamp)
	}
	if value.AutoRetryCount == 0 {
		stats.mu.data.FirstAttemptCount++
	} else if int64(value.AutoRetryCount) > stats.mu.data.MaxRetries {
		stats.mu.data.MaxRetries = int64(value.AutoRetryCount)
	}

	stats.mu.data.SQLType = value.StatementType.String()
	stats.mu.data.NumRows.Record(stats.mu.data.Count, float64(value.RowsAffected))
	stats.mu.data.IdleLat.Record(stats.mu.data.Count, value.IdleLatency)
	stats.mu.data.ParseLat.Record(stats.mu.data.Count, value.ParseLatency)
	stats.mu.data.PlanLat.Record(stats.mu.data.Count, value.PlanLatency)
	stats.mu.data.RunLat.Record(stats.mu.data.Count, value.RunLatency)
	stats.mu.data.ServiceLat.Record(stats.mu.data.Count, value.ServiceLatency)
	stats.mu.data.OverheadLat.Record(stats.mu.data.Count, value.OverheadLatency)
	stats.mu.data.BytesRead.Record(stats.mu.data.Count, float64(value.BytesRead))
	stats.mu.data.RowsRead.Record(stats.mu.data.Count, float64(value.RowsRead))
	stats.mu.data.RowsWritten.Record(stats.mu.data.Count, float64(value.RowsWritten))
	stats.mu.data.LastExecTimestamp = s.getTimeNow()
	stats.mu.data.Nodes = util.CombineUniqueInt64(stats.mu.data.Nodes, value.Nodes)
	stats.mu.data.PlanGists = util.CombineUniqueString(stats.mu.data.PlanGists, []string{value.PlanGist})
	stats.mu.data.IndexRecommendations = value.IndexRecommendations
	stats.mu.data.Indexes = util.CombineUniqueString(stats.mu.data.Indexes, value.Indexes)

	// Note that some fields derived from tracing statements (such as
	// BytesSentOverNetwork) are not updated here because they are collected
	// on-demand.
	// TODO(asubiotto): Record the aforementioned fields here when always-on
	//  tracing is a thing.
	stats.mu.vectorized = key.Vec
	stats.mu.distSQLUsed = key.DistSQL
	stats.mu.fullScan = key.FullScan
	stats.mu.database = key.Database
	stats.mu.querySummary = key.QuerySummary

	if created {
		// stats size + stmtKey size + hash of the statementKey
		estimatedMemoryAllocBytes := stats.sizeUnsafe() + statementKey.size() + 8

		// We also accounts for the memory used for s.sampledPlanMetadataCache.
		// timestamp size + key size + hash.
		estimatedMemoryAllocBytes += timestampSize + statementKey.sampledPlanKey.size() + 8
		s.mu.Lock()
		defer s.mu.Unlock()

		// If the monitor is nil, we do not track memory usage.
		if s.mu.acc.Monitor() == nil {
			return stats.ID, nil
		}

		// We attempt to account for all the memory we used. If we have exceeded our
		// memory budget, delete the entry that we just created and report the error.
		if err := s.mu.acc.Grow(ctx, estimatedMemoryAllocBytes); err != nil {
			delete(s.mu.stmts, statementKey)
			return stats.ID, ErrMemoryPressure
		}
	}

	var autoRetryReason string
	if value.AutoRetryReason != nil {
		autoRetryReason = value.AutoRetryReason.Error()
	}

	var contention *time.Duration
	var contentionEvents []roachpb.ContentionEvent
	if value.ExecStats != nil {
		contention = &value.ExecStats.ContentionTime
		contentionEvents = value.ExecStats.ContentionEvents
	}

	s.insights.ObserveStatement(value.SessionID, &insights.Statement{
		ID:                   value.StatementID,
		FingerprintID:        stmtFingerprintID,
		LatencyInSeconds:     value.ServiceLatency,
		Query:                value.Query,
		Status:               getStatus(value.StatementError),
		StartTime:            value.StartTime,
		EndTime:              value.EndTime,
		FullScan:             value.FullScan,
		PlanGist:             value.PlanGist,
		Retries:              int64(value.AutoRetryCount),
		AutoRetryReason:      autoRetryReason,
		RowsRead:             value.RowsRead,
		RowsWritten:          value.RowsWritten,
		Nodes:                value.Nodes,
		Contention:           contention,
		ContentionEvents:     contentionEvents,
		IndexRecommendations: value.IndexRecommendations,
		Database:             value.Database,
	})

	return stats.ID, nil
}

// RecordStatementExecStats implements sqlstats.Writer interface.
func (s *Container) RecordStatementExecStats(
	key roachpb.StatementStatisticsKey, stats execstats.QueryLevelStats,
) error {
	stmtStats, _, _, _, _ :=
		s.getStatsForStmt(
			key.Query,
			key.ImplicitTxn,
			key.Database,
			key.Failed,
			key.PlanHash,
			key.TransactionFingerprintID,
			false, /* createIfNotExists */
		)
	if stmtStats == nil {
		return ErrExecStatsFingerprintFlushed
	}
	stmtStats.recordExecStats(stats)
	return nil
}

// ShouldSample implements sqlstats.Writer interface.
func (s *Container) ShouldSample(
	fingerprint string, implicitTxn bool, database string,
) (previouslySampled, savePlanForStats bool) {
	lastSampled, previouslySampled := s.getLogicalPlanLastSampled(sampledPlanKey{
		stmtNoConstants: fingerprint,
		implicitTxn:     implicitTxn,
		database:        database,
	})
	savePlanForStats = s.shouldSaveLogicalPlanDescription(lastSampled)
	return previouslySampled, savePlanForStats
}

// RecordTransaction implements sqlstats.Writer interface and saves
// per-transaction statistics.
func (s *Container) RecordTransaction(
	ctx context.Context, key roachpb.TransactionFingerprintID, value sqlstats.RecordedTxnStats,
) error {
	s.recordTransactionHighLevelStats(value.TransactionTimeSec, value.Committed, value.ImplicitTxn)

	if !sqlstats.TxnStatsEnable.Get(&s.st.SV) {
		return nil
	}
	// Do not collect transaction statistics if the stats collection latency
	// threshold is set, since our transaction UI relies on having stats for every
	// statement in the transaction.
	t := sqlstats.StatsCollectionLatencyThreshold.Get(&s.st.SV)
	if t > 0 {
		return nil
	}

	// Get the statistics object.
	stats, created, throttled := s.getStatsForTxnWithKey(key, value.StatementFingerprintIDs, true /* createIfNonexistent */)

	if throttled {
		return ErrFingerprintLimitReached
	}

	// Collect the per-transaction statistics.
	stats.mu.Lock()
	defer stats.mu.Unlock()

	// If we have created a new entry successfully, we check if we have reached
	// the memory limit. If we have, then we delete the newly created entry and
	// return the memory allocation error.
	// If the entry is not created, this means we have reached the limit of unique
	// fingerprints for this app. We also abort the operation and return an error.
	if created {
		estimatedMemAllocBytes :=
			stats.sizeUnsafe() + key.Size() + 8 /* hash of transaction key */
		s.mu.Lock()

		// If the monitor is nil, we do not track memory usage.
		if s.mu.acc.Monitor() != nil {
			if err := s.mu.acc.Grow(ctx, estimatedMemAllocBytes); err != nil {
				delete(s.mu.txns, key)
				s.mu.Unlock()
				return ErrMemoryPressure
			}
		}
		s.mu.Unlock()
	}

	stats.mu.data.Count++

	stats.mu.data.NumRows.Record(stats.mu.data.Count, float64(value.RowsAffected))
	stats.mu.data.ServiceLat.Record(stats.mu.data.Count, value.ServiceLatency.Seconds())
	stats.mu.data.RetryLat.Record(stats.mu.data.Count, value.RetryLatency.Seconds())
	stats.mu.data.CommitLat.Record(stats.mu.data.Count, value.CommitLatency.Seconds())
	stats.mu.data.IdleLat.Record(stats.mu.data.Count, value.IdleLatency.Seconds())
	if value.RetryCount > stats.mu.data.MaxRetries {
		stats.mu.data.MaxRetries = value.RetryCount
	}
	stats.mu.data.RowsRead.Record(stats.mu.data.Count, float64(value.RowsRead))
	stats.mu.data.RowsWritten.Record(stats.mu.data.Count, float64(value.RowsWritten))
	stats.mu.data.BytesRead.Record(stats.mu.data.Count, float64(value.BytesRead))

	if value.CollectedExecStats {
		stats.mu.data.ExecStats.Count++
		stats.mu.data.ExecStats.NetworkBytes.Record(stats.mu.data.ExecStats.Count, float64(value.ExecStats.NetworkBytesSent))
		stats.mu.data.ExecStats.MaxMemUsage.Record(stats.mu.data.ExecStats.Count, float64(value.ExecStats.MaxMemUsage))
		stats.mu.data.ExecStats.ContentionTime.Record(stats.mu.data.ExecStats.Count, value.ExecStats.ContentionTime.Seconds())
		stats.mu.data.ExecStats.NetworkMessages.Record(stats.mu.data.ExecStats.Count, float64(value.ExecStats.NetworkMessages))
		stats.mu.data.ExecStats.MaxDiskUsage.Record(stats.mu.data.ExecStats.Count, float64(value.ExecStats.MaxDiskUsage))
	}

	var retryReason string
	if value.AutoRetryReason != nil {
		retryReason = value.AutoRetryReason.Error()
	}

	s.insights.ObserveTransaction(value.SessionID, &insights.Transaction{
		ID:              value.TransactionID,
		FingerprintID:   key,
		UserPriority:    value.Priority.String(),
		ImplicitTxn:     value.ImplicitTxn,
		Contention:      &value.ExecStats.ContentionTime,
		StartTime:       value.StartTime,
		EndTime:         value.EndTime,
		User:            value.SessionData.User().Normalized(),
		ApplicationName: value.SessionData.ApplicationName,
		RowsRead:        value.RowsRead,
		RowsWritten:     value.RowsWritten,
		RetryCount:      value.RetryCount,
		AutoRetryReason: retryReason,
	})
	return nil
}

func (s *Container) recordTransactionHighLevelStats(
	transactionTimeSec float64, committed bool, implicit bool,
) {
	if !sqlstats.TxnStatsEnable.Get(&s.st.SV) {
		return
	}
	s.txnCounts.recordTransactionCounts(transactionTimeSec, committed, implicit)
}
