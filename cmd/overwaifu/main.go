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
	// getData()
	// getLatestDate()
	// dbJob()
	testDbVersion()
}

func testDbVersion() {
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

	db := session.DB("overwaifu")

	ow, err := overwaifu.New(db)
	if err != nil {
		fmt.Println(err)
	}

	data, err := json.MarshalIndent(ow.Achievements, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))
}

func getLatestDate() {
	data, err := ioutil.ReadFile("dest/cache/cache.json")
	if err != nil {
		log.Panicln(err)
	}

	var posts []getmoe.Post
	if err = json.Unmarshal(data, &posts); err != nil {
		log.Panicln(err)
	}

	var latestDate time.Time
	for _, post := range posts {
		if post.CreatedAt.After(latestDate) {
			latestDate = post.CreatedAt
		}
	}

	fmt.Println(latestDate.Format("2006-01-02"))

	var cred overwaifu.SankakuCredentials
	if err := env.Parse(&cred); err != nil {
		log.Fatalln(err)
	}

	board := sankaku.ChanSankakuConfig
	board.BuildAuth(cred.User, cred.Password)

	board.Query = getmoe.Query{
		Tags: []string{"overwatch", fmt.Sprintf("date:>=%s", latestDate.Format("2006-01-02"))},
		Page: 1,
	}

	start := time.Now()
	println("searching for overwatch lewd images")
	newPosts, err := board.RequestAll()
	if err != nil {
		log.Fatalln(err)
	}
	println("found", len(newPosts))
	elapsed := time.Since(start)
	log.Printf("OverWaifu runtime %s", elapsed)

	for _, newPost := range newPosts {
		found := false
		for _, post := range posts {
			if newPost.Hash == post.Hash {
				found = true
			}
		}
		if !found {
			fmt.Printf("Adding post with hash %s\n", newPost.Hash)
			posts = append(posts, newPost)
		}
	}
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
		log.Panicln(err)
	}

	var posts []getmoe.Post
	if err = json.Unmarshal(data, &posts); err != nil {
		log.Panicln(err)
	}

	c := session.DB("overwaifu").C("posts")
	for i := range posts {
		fmt.Printf("Pushing %5d of %d\n", i, len(posts))
		if err := c.Insert(&posts[i]); err != nil {
			fmt.Println(err)
		}
	}

	// var hashes []string
	// if err := c.Find(nil).Distinct("hash", &hashes); err != nil {
	// 	fmt.Println(err)
	// }

	// var ids []int
	// if err := c.Find(nil).Distinct("id", &ids); err != nil {
	// 	fmt.Println(err)
	// }

	// count, err := c.Count()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println("posts count", count)
	// fmt.Println("uniq hashes", len(hashes))
	// fmt.Println("uniq ids   ", len(ids))
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
