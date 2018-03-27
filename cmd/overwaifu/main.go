package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/leonidboykov/getmoe"
	"github.com/leonidboykov/getmoe/board/sankaku"

	"github.com/overwaifu/overwaifu"
	"github.com/overwaifu/overwaifu/conf"
)

func main() {
	config, err := conf.Load("")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dialInfo := &mgo.DialInfo{
		Addrs:          config.DB.URI,
		ReplicaSetName: config.DB.ReplicaSetName,
		Username:       config.DB.Username,
		Password:       config.DB.Password,
		Source:         config.DB.Source,
		DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), &tls.Config{})
			return conn, err
		},
		Timeout: time.Second * 10,
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer session.Close()

	postsCollection := session.DB("overwaifu").C("posts")
	charactersCollection := session.DB("overwaifu").C("characters")

	posts, err := getPosts(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	uploadPosts(postsCollection, posts)

	ow, err := overwaifu.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ow.QueryScore(postsCollection, charactersCollection)
	ow.QueryAchievements(charactersCollection)

	data, err := json.MarshalIndent(ow, "", "  ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

func getPosts(config *conf.Configuration) ([]getmoe.Post, error) {
	board := sankaku.ChanSankakuConfig
	board.BuildAuth(config.SC.Username, config.SC.Password)

	board.Query = getmoe.Query{
		Tags: []string{"overwatch"},
		Page: 1,
	}

	start := time.Now()
	fmt.Println("searching for overwatch lewd images")
	posts, err := board.RequestAll()
	if err != nil {
		return nil, err
	}
	fmt.Println("found", len(posts))
	elapsed := time.Since(start)
	fmt.Printf("OverWaifu runtime %s", elapsed)

	return posts, nil
}

func uploadPosts(collection *mgo.Collection, posts []getmoe.Post) {
	for i := range posts {
		fmt.Printf("Pushing %5d of %d\n", i+1, len(posts))
		if _, err := collection.Upsert(bson.M{"id": posts[i].ID}, &posts[i]); err != nil {
			fmt.Println(err)
		}
	}
}
