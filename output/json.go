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
	type typeOutput struct {
		Total   int64 `json:"total"`
		Assets  int64 `json:"assets"`
		Dynamic int64 `json:"dynamic"`
	}

	type output struct {
		Hits      typeOutput       `json:"hits"`
		Traffic   typeOutput       `json:"traffic"`
		Status    map[int]int64    `json:"status"`
		Methods   map[string]int64 `json:"methods"`
		Protocols map[string]int64 `json:"protocols"`
		UniqueIPs int64            `json:"uniqueIPs"`
	}

	out := output{
		Hits: typeOutput{
			Total:   stats.TotalHits,
			Assets:  stats.AssetHits,
			Dynamic: stats.DynamicHits,
		},
		Traffic: typeOutput{
			Total:   stats.TotalSize,
			Assets:  stats.AssetSize,
			Dynamic: stats.DynamicSize,
		},
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
