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
	readFlag  = flag.Int("read", 0, "MiB to read from the end of the file (speeds up reading huge logs)")
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

	skipFirstLine := false

	if *readFlag > 0 {
		info, err := file.Stat()
		if err != nil {
			log.Fatalf("Failed to stat file: %v", err)
		}

		totalSize := info.Size()
		jumpTo := totalSize - int64(*readFlag*(1024*1024))

		if jumpTo > 0 {
			log.Printf("File is %d KiB in total, seeking to offset %d KiB.", totalSize/1024, jumpTo/1024)

			_, err = file.Seek(jumpTo, 0)
			if err != nil {
				log.Fatalf("Failed to seek to offset %d in file: %v", jumpTo, err)
			}

			skipFirstLine = true
		} else {
			log.Printf("File is %d KiB in total (smaller than -read size). Not seeking anywhere.", totalSize/1024)
		}
	}

	parser := parser.NewNginxParser()
	stats := logstats.NewStats()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++

		if skipFirstLine && lineNumber == 1 {
			continue
		}

		l := strings.TrimSpace(scanner.Text())

		line, err := parser.ParseLine(l)
		if err != nil {
			log.Printf("[%06d] %v", lineNumber, err)
			continue
		}

		// not in range (do not exit when reaching the end of the range, because
		// logs might not be in perfect chronological order)
		if !line.Date.After(start) || line.Date.After(end) {
			continue
		}

		stats.Count(line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	formatter := output.NewJSONFormatter()
	fmt.Println(formatter.Format(stats))
}
