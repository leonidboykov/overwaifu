package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/caarlos0/env"
	"github.com/globalsign/mgo"
	"github.com/joho/godotenv"
	"github.com/leonidboykov/getmoe"
	"github.com/leonidboykov/getmoe/board/sankaku"

	"github.com/overwaifu/overwaifu"
)

func main() {
	// .env file is used for the local development
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}

	// getCache()
	getData()
}

func dbJob() {
	var mgoConfig overwaifu.MongoDBConfig
	err := env.Parse(&mgoConfig)
	if err != nil {
		log.Fatalln(err)
	}

	dialInfo := &mgo.DialInfo{
		Addrs:          mgoConfig.URI,
		ReplicaSetName: mgoConfig.ReplicaSetName,
		Username:       mgoConfig.User,
		Password:       mgoConfig.Password,
		Source:         mgoConfig.Source,
		DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), &tls.Config{})
			return conn, err
		},
		Timeout: time.Second * 10,
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		fmt.Println(err)
	}
	defer session.Close()

	data, err := ioutil.ReadFile("dest/cache/cache.json")
	if err != nil {
		log.Println(err)
	}

	var posts []getmoe.Post
	if err = json.Unmarshal(data, &posts); err != nil {
		log.Println(err)
	}

	c := session.DB("overwaifu").C("posts")
	// for i := range posts {
	// 	if err := c.Insert(&posts[i]); err != nil {
	// 		fmt.Println(err)
	// 	}
	// }

	var hashes []string
	if err := c.Find(nil).Distinct("author", &hashes); err != nil {
		fmt.Println(err)
	}

	fmt.Println(len(posts))
	fmt.Println(len(hashes))
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
	var cred overwaifu.SankakuCredentials
	err := env.Parse(&cred)
	if err != nil {
		log.Fatalln(err)
	}

	board := sankaku.ChanSankakuConfig
	board.BuildAuth(cred.User, cred.Password)

	board.Query = getmoe.Query{
		Tags: []string{"overwatch"},
		Page: 1,
	}

	start := time.Now()
	println("searching for overwatch lewd images")
	posts, err := board.RequestAll()
	if err != nil {
		log.Fatalln(err)
	}
	println("found", len(posts))
	elapsed := time.Since(start)
	log.Printf("OverWaifu runtime %s", elapsed)

	data, err := json.Marshal(posts)
	if err != nil {
		log.Fatalln(err)
	}

	if err := ioutil.WriteFile("dest/cache/cache.json", data, 0644); err != nil {
		log.Fatalln(err)
	}
}
