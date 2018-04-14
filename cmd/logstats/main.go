package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/xrstf/logstats"
	"github.com/xrstf/logstats/output"
	"github.com/xrstf/logstats/parser"
)

var (
	rangeFlag = flag.String("range", "1m", "time range to scan")
)

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatalln("No filename given.")
	}

	timeRange, err := time.ParseDuration(*rangeFlag)
	if err != nil {
		log.Fatalf("Invalid -range flag: %v", err)
	}

	end := time.Now()
	start := end.Add(-timeRange)

	log.Printf("Range start: %s", start.Format("2006-01-02 15:04:05"))
	log.Printf("Range end:   %s", end.Format("2006-01-02 15:04:05"))

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatalf("Failed to open %s: %v", flag.Arg(0), err)
	}
	defer file.Close()

	parser := parser.NewNginxParser()
	stats := logstats.NewStats()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++

		l := strings.TrimSpace(scanner.Text())

		line, err := parser.ParseLine(l)
		if err != nil {
			log.Printf("[%06d] %v", lineNumber, err)
			continue
		}

		// too early
		if !line.Date.After(start) {
			continue
		}

		stats.Count(line)

		// reached end of range
		if line.Date.After(end) {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	formatter := output.NewJSONFormatter()

	fmt.Printf("%#v\n", stats)
	fmt.Printf("%s\n", formatter.Format(stats))
}
