// Copyright 2020 Daniel Erat. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flag]...\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "Reads a Zoom chat log from stdin and writes SubRip subtitles to stdout.")
		fmt.Fprintln(flag.CommandLine.Output(), "Flags:")
		flag.PrintDefaults()
	}
	firstNames := flag.Bool("first-names", false, "Only include first names in subtitles")
	showSec := flag.Int("show-sec", 8, "Seconds to show subtitles")
	startStr := flag.String("start-time", "", "Video start time as HH:MM:SS")
	flag.Parse()

	if *startStr == "" {
		fmt.Fprintln(os.Stderr, `Must supply video start time with "-start-time HH:MM:SS"`)
		os.Exit(2)
	}
	start, err := time.Parse("15:04:05", *startStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Bad start time %q: %v\n", *startStr, err)
		os.Exit(2)
	}

	showDur := time.Duration(*showSec) * time.Second

	cnt := 0
	re := regexp.MustCompile(`^(\d\d:\d\d:\d\d)\s+From\s+([^:]+):(.+)`)
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		ms := re.FindStringSubmatch(sc.Text())
		if ms == nil {
			log.Fatalf("Failed parsing line %q", sc.Text())
		}

		tm, err := time.Parse("15:04:05", ms[1])
		if err != nil {
			log.Fatalf("Bad line time %q", ms[1])
		}
		if tm.Before(start) {
			log.Fatalf("Line time %v precedes start time %v", tm, start)
		}
		offset := tm.Sub(start)

		who := strings.TrimSpace(ms[2])
		if *firstNames {
			who = strings.Fields(who)[0]
		}

		text := strings.TrimSpace(ms[3])

		fmt.Println(cnt + 1)
		fmt.Printf("%v --> %v\n", formatDuration(offset), formatDuration(offset+showDur))
		fmt.Printf("%v: %v\n", who, text)
		fmt.Println()

		cnt++
	}
	if sc.Err() != nil {
		log.Fatalf("Failed reading input: %v", sc.Err())
	}
}

// formatDuration stringifies d into the "HH:MM:SS,mmm" format expected by SubRip.
func formatDuration(d time.Duration) string {
	return fmt.Sprintf("%02d:%02d:%02d,%03d",
		int(d.Hours()), int(d.Minutes())%60, int(d.Seconds())%60, int(d.Milliseconds()%1000))
}
