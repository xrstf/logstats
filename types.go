package logstats

import (
	"time"
)

type LogLine struct {
	IP       string
	Date     time.Time
	Method   string
	Uri      string
	Protocol string
	Status   int
	Size     int64
}

type LogStats struct {
	TotalHits    int64
	TotalSize    int64
	AssetHits    int64
	AssetSize    int64
	DynamicHits  int64
	DynamicSize  int64
	StatusHits   map[int]int64
	MethodHits   map[string]int64
	ProtocolHits map[string]int64
	IPHits       map[string]int64
}

func NewStats() *LogStats {
	return &LogStats{
		StatusHits:   make(map[int]int64),
		MethodHits:   make(map[string]int64),
		ProtocolHits: make(map[string]int64),
		IPHits:       make(map[string]int64),
	}
}

func (s *LogStats) Count(line *LogLine) {
	s.TotalHits++
	s.TotalSize += line.Size

	if _, ok := s.StatusHits[line.Status]; !ok {
		s.StatusHits[line.Status] = 0
	}

	s.StatusHits[line.Status]++

	if _, ok := s.MethodHits[line.Method]; !ok {
		s.MethodHits[line.Method] = 0
	}

	s.MethodHits[line.Method]++

	if _, ok := s.ProtocolHits[line.Protocol]; !ok {
		s.ProtocolHits[line.Protocol] = 0
	}

	s.ProtocolHits[line.Protocol]++

	if _, ok := s.IPHits[line.IP]; !ok {
		s.IPHits[line.IP] = 0
	}

	s.IPHits[line.IP]++
}
