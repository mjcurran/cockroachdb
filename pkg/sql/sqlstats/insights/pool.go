// Copyright 2022 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package insights

import (
	"sync"

	"github.com/cockroachdb/cockroach/pkg/sql/clusterunique"
)

var insightPool = sync.Pool{
	New: func() interface{} {
		return new(Insight)
	},
}

func makeInsight(
	sessionID clusterunique.ID, transaction *Transaction, statement *Statement,
) *Insight {
	insight := insightPool.Get().(*Insight)
	*insight = Insight{
		Session:     Session{ID: sessionID},
		Transaction: transaction,
		Statement:   statement,
	}
	return insight
}

func releaseInsight(insight *Insight) {
	insight.Causes = insight.Causes[:0]
	*insight = Insight{Causes: insight.Causes}
	insightPool.Put(insight)
}
