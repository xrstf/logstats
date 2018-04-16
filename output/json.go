package output

import (
	"encoding/json"

	"github.com/xrstf/logstats"
)

type jsonFormatter struct{}

func NewJSONFormatter() *jsonFormatter {
	return &jsonFormatter{}
}

func (f *jsonFormatter) Format(stats *logstats.LogStats) string {
	type output struct {
		Hits      map[string]int64 `json:"hits"`
		Traffic   map[string]int64 `json:"traffic"`
		Status    map[int]int64    `json:"status"`
		Methods   map[string]int64 `json:"methods"`
		Protocols map[string]int64 `json:"protocols"`
		UniqueIPs int64            `json:"uniqueIPs"`
	}

	out := output{
		Hits:      stats.Hits,
		Traffic:   stats.Size,
		Status:    stats.StatusHits,
		Methods:   stats.MethodHits,
		Protocols: stats.ProtocolHits,
		UniqueIPs: int64(len(stats.IPHits)),
	}

	encoded, err := json.Marshal(out)
	if err != nil {
		panic(err)
	}

	return string(encoded)
}
