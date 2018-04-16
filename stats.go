package logstats

import (
	"regexp"
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

	filePath string
}

func (l *LogLine) FilePath() string {
	if l.filePath == "" {
		l.filePath = queryStringSep.ReplaceAllString(l.Uri, "")
	}

	return l.filePath
}

type LogStats struct {
	Hits         map[string]int64
	Size         map[string]int64
	StatusHits   map[int]int64
	MethodHits   map[string]int64
	ProtocolHits map[string]int64
	IPHits       map[string]int64
}

func NewStats() *LogStats {
	return &LogStats{
		Hits:         make(map[string]int64),
		Size:         make(map[string]int64),
		StatusHits:   make(map[int]int64),
		MethodHits:   make(map[string]int64),
		ProtocolHits: make(map[string]int64),
		IPHits:       make(map[string]int64),
	}
}

var (
	queryStringSep = regexp.MustCompile(`[?&].*$`)
)

func (s *LogStats) Empty(config *Configuration) {
	s.Hits["total"] = 0
	s.Size["total"] = 0

	for kind := range config.Kinds {
		s.Hits[kind] = 0
		s.Size[kind] = 0
	}
}

func (s *LogStats) Count(line *LogLine, config *Configuration) {
	if config.Exclude.Matches(line) {
		return
	}

	s.countLine("total", line)

	for kind, cfg := range config.Kinds {
		if cfg.Matches(line) {
			s.countLine(kind, line)
		}
	}

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

func (s *LogStats) countLine(kind string, line *LogLine) {
	if _, ok := s.Hits[kind]; !ok {
		s.Hits[kind] = 0
	}

	s.Hits[kind]++

	if _, ok := s.Size[kind]; !ok {
		s.Size[kind] = 0
	}

	s.Size[kind] += line.Size
}
