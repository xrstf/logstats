package parser

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	pcre "github.com/gijsbers/go-pcre"
	"github.com/xrstf/logstats"
)

var nginxInvalidLine = errors.New("Could not parse log line")

type nginxParser struct {
	lineRegex pcre.Regexp
}

func NewNginxParser() *nginxParser {
	ipRegex := `[a-f0-9.:]+`
	userRegex := `[^ ]+`
	timeRegex := `.+?`
	methodRegex := `[^ ]+`
	uriRegex := `[^ ]+`
	protocolRegex := `HTTP/[0-9.]+`
	statusRegex := `[0-9]+`
	sizeRegex := `[0-9]+`

	regex := fmt.Sprintf(`^(%s) - %s \[(%s)\] "(%s) (%s) (%s)" (%s) (%s)`, ipRegex, userRegex, timeRegex, methodRegex, uriRegex, protocolRegex, statusRegex, sizeRegex)

	return &nginxParser{
		lineRegex: pcre.MustCompile(regex, 0),
	}
}

func (s *nginxParser) ParseLine(line string) (*logstats.LogLine, error) {
	matcher := s.lineRegex.MatcherString(line, 0)
	if !matcher.Matches() {
		return nil, nginxInvalidLine
	}

	parsed, err := time.Parse("02/Jan/2006:15:04:05 -0700", matcher.GroupString(2))
	if err != nil {
		return nil, nginxInvalidLine
	}

	statusCode, err := strconv.ParseInt(matcher.GroupString(6), 10, 16)
	if err != nil {
		return nil, nginxInvalidLine
	}

	size, err := strconv.ParseInt(matcher.GroupString(7), 10, 64)
	if err != nil {
		return nil, nginxInvalidLine
	}

	return &logstats.LogLine{
		IP:       matcher.GroupString(1),
		Date:     parsed,
		Method:   matcher.GroupString(3),
		Uri:      matcher.GroupString(4),
		Protocol: matcher.GroupString(5),
		Status:   int(statusCode),
		Size:     size,
	}, nil
}
