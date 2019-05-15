package main

import (
	"github.com/arzonus/sitemap/pkg/sitemap"
	"log"
	"time"
)

func main() {
	if err := sitemap.NewWalker(
		"https://vk.com",
		sitemap.MaxDepthOption(5),
		sitemap.TimeoutOption(25*time.Second),
	).Run(); err != nil {
		log.Fatal(err)
	}
}
