package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/leonidboykov/getmoe"
	"github.com/leonidboykov/getmoe/board/sankaku"
	"github.com/leonidboykov/overwaifu"
)

func main() {
	getCache()
	getData()
}

func getData() {
	data, err := ioutil.ReadFile("dest/cache/cache.json")
	if err != nil {
		log.Panicln(err)
	}

	var posts []getmoe.Post
	if err = json.Unmarshal(data, &posts); err != nil {
		log.Panicln(err)
	}

	ow, err := overwaifu.New(posts)
	if err != nil {
		log.Panicln(err)
	}

	ow.FetchData()
	ow.Analyse()

	data, err = json.MarshalIndent(ow, "", "  ")
	if err != nil {
		log.Panicln(err)
	}

	if err := ioutil.WriteFile("dest/overwaifu.json", data, 0644); err != nil {
		log.Panicln(err)
	}
}

func getCache() {
	board := sankaku.ChanSankakuConfig
	board.BuildAuth("xxx", "xxx")

	board.Query = getmoe.Query{
		Tags: []string{"overwatch"},
		Page: 1,
	}

	start := time.Now()
	println("searching for overwatch lewd images")
	posts, err := board.RequestAll()
	if err != nil {
		log.Panicln(err)
	}
	println("found", len(posts))
	elapsed := time.Since(start)
	log.Printf("OverWaifu runtime %s", elapsed)

	data, err := json.Marshal(posts)
	if err != nil {
		log.Panicln(err)
	}

	if err := ioutil.WriteFile("dest/cache/cache.json", data, 0644); err != nil {
		log.Panicln(err)
	}
}
