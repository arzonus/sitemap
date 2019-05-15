package main

import (
	"github.com/arzonus/sitemap/pkg/sitemap"
	"log"
	"time"
)

func main() {
	if err := sitemap.NewWalker(
		"https://example.com",
		sitemap.MaxDepthOption(5),
		sitemap.TimeoutOption(10*time.Second),
	).Run(); err != nil {
		log.Fatal(err)
	}
}
