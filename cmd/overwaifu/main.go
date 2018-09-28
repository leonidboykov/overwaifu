package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/leonidboykov/getmoe"
	"github.com/leonidboykov/getmoe/board"
	gconf "github.com/leonidboykov/getmoe/conf"

	"github.com/overwaifu/overwaifu"
	"github.com/overwaifu/overwaifu/conf"
)

var scratchFlag = false

const timeFormat = "2006-01-02"

func main() {
	flag.BoolVar(&scratchFlag, "scratch", false, "ignore date and fetch all posts")
	flag.Parse()

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

	fmt.Println("Pushing to MongoDB")
	uploadPosts(postsCollection, posts)

	ow, err := overwaifu.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Calculating characters scores")
	ow.QueryScore(postsCollection, charactersCollection)
	fmt.Println("Calculating achievements")
	ow.QueryAchievements(charactersCollection)

	data, err := json.MarshalIndent(ow, "", "  ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := uploadResults(config, data); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := notifyNetlify(config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getPosts(config *conf.Configuration) ([]getmoe.Post, error) {
	var tags []string
	if scratchFlag {
		fmt.Println("Fetching posts from scratch: -scratch flag was used")
		tags = []string{"~overwatch", "~blizzard_entertainment"}
	} else {
		date := time.Now().AddDate(0, 0, -7).Format(timeFormat)
		fmt.Printf("Fetching posts from %s\n", date)
		tags = []string{
			"~overwatch",
			"~blizzard_entertainment",
			fmt.Sprintf("date:>=%s", date),
		}
	}

	board := board.AvailableBoards["chan.sankakucomplex.com"]
	board.Provider.Auth(gconf.AuthConfiguration{
		Login:    config.SC.Username,
		Password: config.SC.Password,
	})
	board.Provider.BuildRequest(gconf.RequestConfiguration{
		Tags: tags,
	})
	posts, err := board.RequestAll()
	if err != nil {
		return nil, err
	}
	fmt.Printf("Found %d posts\n", len(posts))

	return posts, nil
}

func uploadPosts(collection *mgo.Collection, posts []getmoe.Post) {
	for i := range posts {
		if _, err := collection.Upsert(bson.M{"id": posts[i].ID}, &posts[i]); err != nil {
			fmt.Println(err)
		}
	}
}

func uploadResults(config *conf.Configuration, data []byte) error {
	u := url.URL{
		Scheme: "https",
		Host:   "api.myjson.com",
		Path:   "bins/" + config.MyJSON.BucketID,
	}
	fmt.Printf("Uploading to %s\n", u.String())

	client := http.Client{}
	req, err := http.NewRequest("PUT", u.String(), bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(string(body))

	return nil
}

func notifyNetlify(config *conf.Configuration) error {
	u := url.URL{
		Scheme: "https",
		Host:   "api.netlify.com",
		Path:   "build_hooks/" + config.Netlify.BuildHook,
	}
	fmt.Printf("Webhook to %s\n", u.String())

	client := http.Client{}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	return nil
}
