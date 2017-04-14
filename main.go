package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ssimunic/jsonscraper/scraper"
)

var wg sync.WaitGroup

func init() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Missing config path(s).")
		os.Exit(1)
	}
}

func main() {
	start := time.Now()
	configPaths := os.Args[1:]
	wg.Add(len(configPaths))

	for _, configPath := range configPaths {
		go func(configPath string) {
			defer wg.Done()

			s, err := scraper.New(configPath)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}

			s.Start()
		}(configPath)
	}

	wg.Wait()
	log.Printf("Time elapsed: %vs", time.Since(start).Seconds())
}
