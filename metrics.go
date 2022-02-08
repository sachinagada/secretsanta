package secretsanta

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	MParticipants = stats.Int64(
		"github.com/sachinagada/secretsanta/participant_count",
		"Number of participants per group",
		stats.UnitDimensionless,
	)

	MCommunicationLatency = stats.Int64(
		"github.com/sachinagada/secretsanta/communication_latency",
		"Time it took to send communication to all the participants in the group",
		stats.UnitMilliseconds,
	)
)

var (
	ViewParticipants = &view.View{
		Measure:     MParticipants,
		Aggregation: view.Sum(),
	}

	ViewCommunicationLatency = &view.View{
		Measure:     MCommunicationLatency,
		Aggregation: view.Sum(),
	}
)
