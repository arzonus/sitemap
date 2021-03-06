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
	timeout        = flag.Int("timeout", 30, "timeout")
)

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 || len(flag.Args()) != 0 && flag.Args()[0] == "" {
		log.Print("url is empty")
		os.Exit(1)
	}

	w := sitemap.NewWalker(
		sitemap.MaxDepthOption(*maxDepth),
		sitemap.WorkerCountOption(*parallel),
		sitemap.TimeoutOption(time.Duration(*timeout)*time.Second),
	)

	log.Println("**********************")
	log.Printf("* url      %s", flag.Args()[0])
	log.Printf("* file     %s", *outputFilePath)
	log.Printf("* depth    %d", *maxDepth)
	log.Printf("* timeout  %ds", *timeout)
	log.Printf("* parallel %d", *parallel)
	log.Println("**********************")

	start := time.Now()

	n, err := w.Walk(flag.Args()[0])
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	log.Print("\n", n.Tree())
	log.Printf("time: %s", time.Since(start))

	if err := ioutil.WriteFile(*outputFilePath, n.TreeBytes(), os.ModePerm); err != nil {
		log.Print(err)
		os.Exit(1)
	}
	log.Printf("result saved to file %s", *outputFilePath)
}
