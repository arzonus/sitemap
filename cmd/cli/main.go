package main

import (
	"flag"
	"github.com/arzonus/sitemap/pkg/sitemap"
	"io/ioutil"
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

	w := sitemap.NewWalker(
		sitemap.MaxDepthOption(*maxDepth),
		sitemap.WorkerCountOption(*parallel),
		sitemap.TimeoutOption(time.Duration(*timeout)*time.Second),
	)

	node, err := w.Walk(os.Args[1])
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	log.Print(node.String())
	ioutil.WriteFile(*outputFilePath, []byte(node.String()), os.ModePerm)
}
