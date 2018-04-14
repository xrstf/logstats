package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/xrstf/logstats"
)

/*
172.19.0.1 - - [02/Apr/2018:12:08:33 +0000] "GET / HTTP/1.1" 200 1043 "-" "curl/7.52.1" "-"
172.19.0.1 - - [02/Apr/2018:12:08:34 +0000] "GET / HTTP/1.1" 200 1043 "-" "curl/7.52.1" "-"
172.19.0.1 - - [02/Apr/2018:12:08:36 +0000] "GET /foo HTTP/1.1" 200 1043 "-" "curl/7.52.1" "-"
172.19.0.1 - - [02/Apr/2018:12:08:39 +0000] "GET /robots.txt HTTP/1.1" 200 84 "-" "curl/7.52.1" "-"
172.19.0.1 - - [02/Apr/2018:12:13:28 +0000] "GET /robots.txt HTTP/1.1" 200 84 "-" "curl/7.52.1" "-"
172.19.0.1 - - [02/Apr/2018:12:13:30 +0000] "GET /foo HTTP/1.1" 200 1211 "-" "curl/7.52.1" "-"
172.19.0.6 - - [02/Apr/2018:17:38:21 +0000] "GET / HTTP/1.0" 200 1189 "-" "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:59.0) Gecko/20100101 Firefox/59.0" "-"
172.19.0.6 - - [02/Apr/2018:17:38:21 +0000] "GET /favicon.ico HTTP/1.0" 200 15086 "-" "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:59.0) Gecko/20100101 Firefox/59.0" "-"
172.19.0.6 - - [02/Apr/2018:17:38:21 +0000] "GET /favicon.ico HTTP/1.0" 200 15086 "-" "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:59.0) Gecko/20100101 Firefox/59.0" "-"
172.19.0.6 - - [02/Apr/2018:17:38:24 +0000] "GET / HTTP/1.0" 200 1189 "-" "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:59.0) Gecko/20100101 Firefox/59.0" "-"
*/

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
		fmt.Println(line)
		panic("foo")
		return nil, nginxInvalidLine
	}

	parsed, err := time.Parse("02/Jan/2006:15:04:05 -0700", match[2])
	if err != nil {
		fmt.Println(line)
		panic(err)
		return nil, nginxInvalidLine
	}

	statusCode, err := strconv.ParseInt(match[6], 10, 16)
	if err != nil {
		fmt.Println(line)
		panic(err)
		return nil, nginxInvalidLine
	}

	size, err := strconv.ParseInt(match[7], 10, 64)
	if err != nil {
		fmt.Println(line)
		panic(err)
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
