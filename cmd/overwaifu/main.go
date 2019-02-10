package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/leonidboykov/getmoe"
	"github.com/leonidboykov/getmoe/provider"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/overwaifu/overwaifu"
	"github.com/overwaifu/overwaifu/conf"
)

var scratchFlag = false

func main() {
	flag.BoolVar(&scratchFlag, "scratch", false, "ignore date and fetch all posts")
	flag.Parse()

	config, err := conf.Load("")
	if err != nil {
		log.Fatalln(err)
	}

	connString := fmt.Sprintf("mongodb+srv://%s:%s@%s",
		config.DB.Username,
		config.DB.Password,
		config.DB.URI,
	)
	client, err := mongo.Connect(context.TODO(), connString)
	if err != nil {
		log.Fatalln(err)
	}

	postsCollection := client.Database("overwaifu").Collection("posts")
	charactersCollection := client.Database("overwaifu").Collection("posts")

	posts, err := getPosts(config)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Pushing to MongoDB")
	uploadPosts(postsCollection, posts)

	ow, err := overwaifu.New()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Calculating characters scores")
	ow.QueryScore(postsCollection, charactersCollection)
	log.Println("Calculating achievements")
	ow.QueryAchievements(charactersCollection)

	// Disconnect from MongoDB
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Fatalln(err)
	}

	data, err := json.Marshal(ow)
	if err != nil {
		log.Fatalln(err)
	}

	if err := uploadResults(config, data); err != nil {
		log.Fatalln(err)
	}

	if err := notifyNetlify(config); err != nil {
		log.Fatalln(err)
	}
}

func getPosts(config *conf.Configuration) ([]getmoe.Post, error) {
	tags := getmoe.NewTags("overwatch")
	if scratchFlag {
		log.Println("Fetching posts from scratch: -scratch flag was used")
	} else {
		date := time.Now().AddDate(0, 0, -7)
		log.Printf("Fetching posts from %s\n", date)
		// tags.AfterDate(date)
		tags.And("date:>=" + date.Format("02.01.2006"))
	}

	board := provider.AvailableBoards["chan.sankakucomplex.com"]
	board.Provider.Auth(getmoe.AuthConfiguration{
		Login:    config.SC.Username,
		Password: config.SC.Password,
	})
	board.Provider.BuildRequest(getmoe.RequestConfiguration{
		Tags: *tags,
	})
	posts, err := board.RequestAll()
	if err != nil {
		return nil, err
	}
	log.Printf("Found %d posts\n", len(posts))

	return posts, nil
}

func uploadPosts(collection *mongo.Collection, posts []getmoe.Post) {
	var models []mongo.WriteModel
	for i := range posts {
		model := mongo.NewUpdateOneModel().
			SetFilter(bson.M{"id": posts[i].ID}).
			SetUpdate(bson.M{"$set": &posts[i]}).
			SetUpsert(true)
		models = append(models, model)
	}
	_, err := collection.BulkWrite(context.TODO(), models)
	if err != nil {
		log.Fatalln(err)
	}
}

func uploadResults(config *conf.Configuration, data []byte) error {
	u := url.URL{
		Scheme: "https",
		Host:   "api.myjson.com",
		Path:   "bins/" + config.MyJSON.BucketID,
	}
	log.Printf("Uploading to %s\n", u.String())

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

	return nil
}

func notifyNetlify(config *conf.Configuration) error {
	u := url.URL{
		Scheme: "https",
		Host:   "api.netlify.com",
		Path:   "build_hooks/" + config.Netlify.BuildHook,
	}
	log.Printf("Webhook to %s\n", u.String())

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

	return nil
}
