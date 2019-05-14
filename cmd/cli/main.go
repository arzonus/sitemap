package main

import (
	"github.com/arzonus/sitemap/pkg/sitemap"
	"log"
)

func main(){
	if err := sitemap.NewWalker("https://vk.com").Run(); err != nil {
		log.Fatal(err)
	}
}