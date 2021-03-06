package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/xrstf/logstats"
	"github.com/xrstf/logstats/output"
	"github.com/xrstf/logstats/parser"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	start := time.Now()

	if len(os.Args) < 3 {
		log.Fatalln("No config file and/or log file given.")
	}

	/////////////////////////////////////////////////////////////////////////////
	// setup files and config

	configFile := os.Args[1]
	logFile := os.Args[2]

	config, err := readConfigFile(configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	/////////////////////////////////////////////////////////////////////////////
	// open file

	file, err := os.Open(logFile)
	if err != nil {
		log.Fatalf("Failed to open %s: %v", logFile, err)
	}
	defer file.Close()

	/////////////////////////////////////////////////////////////////////////////
	// process file

	stats, err := processFile(file, config, time.Now())
	if err != nil {
		log.Fatalln(err)
	}

	/////////////////////////////////////////////////////////////////////////////
	// print result

	formatter := output.NewJSONFormatter()
	fmt.Println(formatter.Format(stats))

	if len(os.Getenv("LOGSTATS_DEBUG")) > 0 {
		log.Printf("Total time: %s", time.Since(start))
	}
}

func processFile(file *os.File, config *logstats.Configuration, now time.Time) (*logstats.LogStats, error) {
	/////////////////////////////////////////////////////////////////////////////
	// prepare range

	end := now
	start := end.Add(-config.Range)

	/////////////////////////////////////////////////////////////////////////////
	// seek to file offset

	skipFirstLine, err := seekToFileOffset(file, config)
	if err != nil {
		return nil, err
	}

	/////////////////////////////////////////////////////////////////////////////
	// setup parsing logic

	parser := parser.NewNginxParser()
	stats := logstats.NewStats()

	config.Compile()
	stats.Empty(config)

	/////////////////////////////////////////////////////////////////////////////
	// here we go

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

		stats.Count(line, config)
	}

	if err := scanner.Err(); err != nil {
		return stats, fmt.Errorf("Failed to read file: %v", err)
	}

	return stats, nil
}

func readConfigFile(filename string) (*logstats.Configuration, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &logstats.Configuration{}
	err = yaml.Unmarshal(content, config)

	return config, nil
}

func seekToFileOffset(file *os.File, config *logstats.Configuration) (bool, error) {
	if config.Read > 0 {
		info, err := file.Stat()
		if err != nil {
			return false, fmt.Errorf("Failed to stat file: %v", err)
		}

		totalSize := info.Size()
		jumpTo := totalSize - int64(config.Read*(1024*1024))

		if jumpTo > 0 {
			_, err = file.Seek(jumpTo, 0)
			if err != nil {
				err = fmt.Errorf("Failed to seek to offset %d in file: %v", jumpTo, err)
			}

			return true, err
		}
	}

	return false, nil
}
