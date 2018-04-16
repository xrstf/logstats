package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/xrstf/logstats"
)

var nginxInvalidLine = errors.New("Could not parse log line")

type nginxParser struct {
	lineRegex *regexp.Regexp
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
		lineRegex: regexp.MustCompile(regex),
	}
}

func (s *nginxParser) ParseLine(line string) (*logstats.LogLine, error) {
	match := s.lineRegex.FindStringSubmatch(line)
	if match == nil {
		return nil, nginxInvalidLine
	}

	parsed, err := time.Parse("02/Jan/2006:15:04:05 -0700", match[2])
	if err != nil {
		return nil, nginxInvalidLine
	}

	statusCode, err := strconv.ParseInt(match[6], 10, 16)
	if err != nil {
		return nil, nginxInvalidLine
	}

	size, err := strconv.ParseInt(match[7], 10, 64)
	if err != nil {
		return nil, nginxInvalidLine
	}

	return &logstats.LogLine{
		IP:       match[1],
		Date:     parsed,
		Method:   match[3],
		Uri:      match[4],
		Protocol: match[5],
		Status:   int(statusCode),
		Size:     size,
	}, nil
}
