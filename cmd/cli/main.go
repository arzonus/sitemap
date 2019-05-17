package main

import (
	"flag"
	"github.com/arzonus/sitemap/pkg/sitemap"
	"log"
	"os"
	"runtime"
	"time"
)

var (
	parallel       = flag.Int("parallel", runtime.NumCPU(), "count of parallel requests to sites")
	outputFilePath = flag.String("output-file", "sitemap.out", "filepath to output file")
	maxDepth       = flag.Int("max-depth", 5, "max depth of handling sites")
	timeout        = flag.Int("timeout", 30, "global timeout")
)

func main() {
	flag.Parse()
	if os.Args[1] == "" {
		log.Print("url is empty")
		os.Exit(1)
	}

	if err := sitemap.NewWalker(
		os.Args[1],
		sitemap.MaxDepthOption(*maxDepth),
		sitemap.WorkerCountOption(*parallel),
		sitemap.TimeoutOption(time.Duration(*timeout)*time.Second),
	).Run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
