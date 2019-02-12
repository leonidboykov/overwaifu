package overwaifu

import (
	"context"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo/options"

	"github.com/BurntSushi/toml"
	"github.com/leonidboykov/getmoe"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

var repository []getmoe.Post

const (
	resourceFolder = "assets/overwatch/"
	resourceExt    = ".toml"
)

const (
	minScoreForCharacter = 1500
	minScoreForSkin      = 100
)

// OverWaifu holds all overwaifu results
type OverWaifu struct {
	UpdatedAt    time.Time               `json:"updated_at"`
	PostsCount   int                     `json:"posts_count"`
	Characters   Characters              `json:"characters"`
	Achievements map[string]*Achievement `json:"achievements"`
}

func getCharactersList() ([]string, error) {
	files, err := ioutil.ReadDir(resourceFolder)
	if err != nil {
		return nil, err
	}

	var result []string

	for _, f := range files {
		basename := f.Name()
		if path.Ext(basename) == resourceExt {
			basename = strings.TrimSuffix(basename, filepath.Ext(basename))
			result = append(result, basename)
		}
	}
	return result, nil
}

// New creates a new OverWaifu instance
func New() (*OverWaifu, error) {
	characters, err := getCharactersList()
	if err != nil {
		return nil, err
	}

	ow := OverWaifu{
		UpdatedAt:    time.Now(),
		Characters:   make(map[string]*Character),
		Achievements: make(map[string]*Achievement),
	}
	for i := range characters {
		var c Character
		_, err := toml.DecodeFile(resourceFolder+characters[i]+resourceExt, &c)
		if err != nil {
			return nil, err
		}
		ow.Characters[characters[i]] = &c
		ow.Characters[characters[i]].Key = characters[i]
		ow.Characters[characters[i]].UpdateSkinDefaults()
	}

	return &ow, nil
}

// QueryScore calculates score
func (ow *OverWaifu) QueryScore(postsCollection, charactersCollection *mongo.Collection) {
	// Query meta information
	count, err := postsCollection.Count(context.TODO(), bson.M{})
	if err != nil {
		log.Println(err)
	}
	ow.PostsCount = int(count)

	// Query characters score
	for k := range ow.Characters {
		if err := ow.Characters[k].QueryScore(postsCollection); err != nil {
			log.Println(err)
		}
		ow.Characters[k].QueryScoreSkins(postsCollection)
	}
	ow.saveScores(charactersCollection)
}

func (ow *OverWaifu) saveScores(collection *mongo.Collection) {
	// var models []mongo.WriteModel
	for _, c := range ow.Characters {
		if _, err := collection.UpdateOne(context.TODO(), bson.M{"key": c.Key}, bson.M{"$set": &c}, options.Update().SetUpsert(true)); err != nil {
			log.Fatalln(err)
		}
		// model := mongo.NewUpdateOneModel().
		// 	SetFilter(bson.M{"key": c.Key}).
		// 	SetUpdate(bson.M{"$set": &c}).
		// 	SetUpsert(true)
		// models = append(models, model)
	}
	// _, err := collection.BulkWrite(context.TODO(), models)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
}

// QueryAchievements calculates achievements
func (ow *OverWaifu) QueryAchievements(collection *mongo.Collection) error {
	ctx := context.TODO()
	pipeline := bson.A{
		bson.M{"$addFields": bson.M{"skins": bson.M{"$objectToArray": "$skins"}}},
		bson.M{"$facet": bson.M{
			"fameWaifu": bson.A{
				bson.M{"$match": bson.M{"sex": "female"}},
				bson.M{"$sort": bson.M{"score.all": -1}},
				bson.M{"$limit": 1},
				bson.M{"$project": bson.M{
					"_id":       0,
					"character": "$key",
				}},
			},
			"hotWaifu": bson.A{
				bson.M{"$match": bson.M{"sex": "female"}},
				bson.M{"$sort": bson.M{"score.explicit": -1}},
				bson.M{"$limit": 1},
				bson.M{"$project": bson.M{
					"_id":       0,
					"character": "$key",
				}},
			},
			"lewdWaifu": bson.A{
				bson.M{"$match": bson.M{"sex": "female"}},
				bson.M{"$match": bson.M{"score.all": bson.M{"$gt": minScoreForCharacter}}},
				bson.M{"$sort": bson.M{"score.lewd": -1}},
				bson.M{"$limit": 1},
				bson.M{"$project": bson.M{
					"_id":       0,
					"character": "$key",
				}},
			},
			"pureWaifu": bson.A{
				bson.M{"$match": bson.M{"sex": "female"}},
				bson.M{"$match": bson.M{"score.all": bson.M{"$gt": minScoreForCharacter}}},
				bson.M{"$sort": bson.M{"score.pure": -1}},
				bson.M{"$limit": 1},
				bson.M{"$project": bson.M{
					"_id":       0,
					"character": "$key",
				}},
			},
			"fameWaifuSkin": bson.A{
				bson.M{"$match": bson.M{"sex": "female"}},
				bson.M{"$unwind": "$skins"},
				bson.M{"$sort": bson.M{"skins.v.score.all": -1}},
				bson.M{"$limit": 1},
				bson.M{"$project": bson.M{
					"_id":       0,
					"character": "$key",
					"skin":      "$skins.k",
				}},
			},
			"hotWaifuSkin": bson.A{
				bson.M{"$match": bson.M{"sex": "female"}},
				bson.M{"$unwind": "$skins"},
				bson.M{"$sort": bson.M{"skins.v.score.explicit": -1}},
				bson.M{"$limit": 1},
				bson.M{"$project": bson.M{
					"_id":       0,
					"character": "$key",
					"skin":      "$skins.k",
				}},
			},
			"lewdWaifuSkin": bson.A{
				bson.M{"$match": bson.M{"sex": "female"}},
				bson.M{"$unwind": "$skins"},
				bson.M{"$match": bson.M{"skins.v.score.all": bson.M{"$gt": minScoreForSkin}}},
				bson.M{"$sort": bson.M{"skins.v.score.lewd": -1}},
				bson.M{"$limit": 1},
				bson.M{"$project": bson.M{
					"_id":       0,
					"character": "$key",
					"skin":      "$skins.k",
				}},
			},
			"pureWaifuSkin": bson.A{
				bson.M{"$match": bson.M{"sex": "female"}},
				bson.M{"$unwind": "$skins"},
				bson.M{"$match": bson.M{"skins.v.score.all": bson.M{"$gt": minScoreForSkin}}},
				bson.M{"$sort": bson.M{"skins.v.score.pure": -1}},
				bson.M{"$limit": 1},
				bson.M{"$project": bson.M{
					"_id":       0,
					"character": "$key",
					"skin":      "$skins.k",
				}},
			},
			"fakeWaifu": bson.A{
				bson.M{"$match": bson.M{"sex": "male"}},
				bson.M{"$sort": bson.M{"score.gender_swaps": -1}},
				bson.M{"$limit": 1},
				bson.M{"$project": bson.M{
					"_id":       0,
					"character": "$key",
				}},
			},
			"virginKillerWaifu": bson.A{
				bson.M{"$match": bson.M{"sex": "female"}},
				bson.M{"$sort": bson.M{"score.virgin_killer_sweater": -1}},
				bson.M{"$limit": 1},
				bson.M{"$project": bson.M{
					"_id":       0,
					"character": "$key",
				}},
			},
			"selfieWaifu": bson.A{
				bson.M{"$match": bson.M{"sex": "female"}},
				bson.M{"$sort": bson.M{"score.selfie": -1}},
				bson.M{"$limit": 1},
				bson.M{"$project": bson.M{
					"_id":       0,
					"character": "$key",
				}},
			},
		}},
		bson.M{"$project": bson.M{
			"fame_waifu":          bson.M{"$arrayElemAt": bson.A{"$fameWaifu", 0}},
			"hot_waifu":           bson.M{"$arrayElemAt": bson.A{"$hotWaifu", 0}},
			"lewd_waifu":          bson.M{"$arrayElemAt": bson.A{"$lewdWaifu", 0}},
			"pure_waifu":          bson.M{"$arrayElemAt": bson.A{"$pureWaifu", 0}},
			"fame_waifu_skin":     bson.M{"$arrayElemAt": bson.A{"$fameWaifuSkin", 0}},
			"hot_waifu_skin":      bson.M{"$arrayElemAt": bson.A{"$hotWaifuSkin", 0}},
			"lewd_waifu_skin":     bson.M{"$arrayElemAt": bson.A{"$lewdWaifuSkin", 0}},
			"pure_waifu_skin":     bson.M{"$arrayElemAt": bson.A{"$pureWaifuSkin", 0}},
			"fake_waifu":          bson.M{"$arrayElemAt": bson.A{"$fakeWaifu", 0}},
			"virgin_killer_waifu": bson.M{"$arrayElemAt": bson.A{"$virginKillerWaifu", 0}},
			"selfie_waifu":        bson.M{"$arrayElemAt": bson.A{"$selfieWaifu", 0}},
		}},
	}

	cur, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	cur.Next(ctx) // Go to the first element
	return cur.Decode(&ow.Achievements)
}
