package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"time"
)

var RFC3339Regex = regexp.MustCompile("((?:(\\d{4}-\\d{2}-\\d{2})T(\\d{2}:\\d{2}:\\d{2}(?:\\.\\d+)?))(Z|[\\+-]\\d{2}:\\d{2})?)")

type truncateGranularity string

var (
	secondGranularity = truncateGranularity("second")
	minuteGranularity = truncateGranularity("minute")
	hourGranularity   = truncateGranularity("hour")
	dayGranularity    = truncateGranularity("day")
	monthGranularity  = truncateGranularity("month")
	YearGranularity   = truncateGranularity("year")
)
var granularityChoices = map[truncateGranularity]interface{}{
	secondGranularity: nil,
	minuteGranularity: nil,
	hourGranularity:   nil,
	dayGranularity:    nil,
	monthGranularity:  nil,
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Granularity must be provided as the first argument.")
		os.Exit(1)
	}
	granularity := truncateGranularity(os.Args[1])
	if _, ok := granularityChoices[granularity]; !ok {
		log.Fatalln("Invalid granularity choice")
		os.Exit(1)
	}

	var lines *bufio.Reader
	if len(os.Args) > 2 {
		fp, err := os.Open(os.Args[2])
		if err != nil {
			log.Fatalf("Error opening file: %v\n", err)
		}
		lines = bufio.NewReader(fp)
	} else {
		lines = bufio.NewReader(os.Stdin)
	}

	buffer := bufio.NewWriterSize(os.Stdout, 16*1024)
	for {
		line, err := lines.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Failed to read line: %v\n", err)
		}
		matches := RFC3339Regex.FindAllIndex(line, -1)
		offset := 0
		for _, indexes := range matches {
			start, end := indexes[0], indexes[1]
			dt, err := time.Parse(time.RFC3339, string(line[start:end]))
			if err != nil {
				log.Fatalf("Failed to parse date: %v\n", err)
			}
			buffer.Write(line[offset:start])
			buffer.WriteString(truncateDate(dt, granularity).Format(time.RFC3339))
			offset = end
		}
		buffer.WriteByte('\n')
	}
	buffer.Flush()
}

func truncateDate(dt time.Time, granularity truncateGranularity) time.Time {
	dt = dt.Add(time.Nanosecond * time.Duration(dt.Nanosecond()) * -1)
	if granularity == secondGranularity {
		return dt
	}
	dt = dt.Add(time.Second * time.Duration(dt.Second()) * -1)
	if granularity == minuteGranularity {
		return dt
	}
	dt = dt.Add(time.Minute * time.Duration(dt.Minute()) * -1)
	if granularity == hourGranularity {
		return dt
	}
	dt = dt.Add(time.Hour * time.Duration(dt.Hour()) * -1)
	if granularity == dayGranularity {
		return dt
	}
	dt = dt.Add(time.Hour * 24 * time.Duration(dt.Day()-1) * -1)
	if granularity == monthGranularity {
		return dt
	}

	return dt
}

func getOffset(dt time.Time) time.Duration {
	return (time.Hour * time.Duration(dt.Hour()) * time.Duration(-1)) +
		(time.Minute * time.Duration(dt.Minute()) * time.Duration(-1)) +
		(time.Second * time.Duration(dt.Second()) * time.Duration(-1))
}
