package main

import (
	"github.com/arzonus/sitemap/pkg/sitemap"
	"log"
	"time"
)

func main() {
	if err := sitemap.NewWalker(
		"https://vk.com",
		sitemap.MaxDepthOption(10),
		sitemap.TimeoutOption(30*time.Second),
	).Run(); err != nil {
		log.Fatal(err)
	}
}
