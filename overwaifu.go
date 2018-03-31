package overwaifu

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/leonidboykov/getmoe"
)

var repository []getmoe.Post

const (
	resourceFolder = "resources/overwatch/"
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
	Characters   map[string]*Character   `json:"characters"`
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

// New createn new OverWaifu instance
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
		ow.Characters[characters[i]].UpdateSkinKey()
	}

	return &ow, nil
}

// QueryScore calculates score
func (ow *OverWaifu) QueryScore(postsCollection, charactersCollection *mgo.Collection) {
	// query meta information
	count, err := postsCollection.Count()
	if err != nil {
		fmt.Println(err)
	}
	ow.PostsCount = count

	// query characters score
	for k := range ow.Characters {
		ow.Characters[k].QueryScore(postsCollection)
		ow.Characters[k].QueryScoreSkins(postsCollection)
	}
	ow.saveScores(charactersCollection)
}

func (ow *OverWaifu) saveScores(collection *mgo.Collection) {
	for k := range ow.Characters {
		c := ow.Characters[k]
		if _, err := collection.Upsert(bson.M{"key": c.Key}, &c); err != nil {
			fmt.Println(err)
		}
	}
}

// QueryAchievements calculates achievements
func (ow *OverWaifu) QueryAchievements(collection *mgo.Collection) error {
	// Build query
	query := []bson.M{
		{"$addFields": bson.M{"skins": bson.M{"$objectToArray": "$skins"}}},
		{"$facet": bson.M{
			"fameWaifu": []bson.M{
				{"$match": bson.M{"sex": "female"}},
				{"$sort": bson.M{"score.all": -1}},
				{"$limit": 1},
				{"$project": bson.M{
					"_id":       0,
					"character": "$key",
				}},
			},
			"hotWaifu": []bson.M{
				{"$match": bson.M{"sex": "female"}},
				{"$sort": bson.M{"score.explicit": -1}},
				{"$limit": 1},
				{"$project": bson.M{
					"_id":       0,
					"character": "$key",
				}},
			},
			"lewdWaifu": []bson.M{
				{"$match": bson.M{"sex": "female"}},
				{"$match": bson.M{"score.all": bson.M{"$gt": minScoreForCharacter}}},
				{"$sort": bson.M{"score.lewd": -1}},
				{"$limit": 1},
				{"$project": bson.M{
					"_id":       0,
					"character": "$key",
				}},
			},
			"pureWaifu": []bson.M{
				{"$match": bson.M{"sex": "female"}},
				{"$match": bson.M{"score.all": bson.M{"$gt": minScoreForCharacter}}},
				{"$sort": bson.M{"score.pure": -1}},
				{"$limit": 1},
				{"$project": bson.M{
					"_id":       0,
					"character": "$key",
				}},
			},
			"fameWaifuSkin": []bson.M{
				{"$match": bson.M{"sex": "female"}},
				{"$unwind": "$skins"},
				{"$sort": bson.M{"skins.v.score.all": -1}},
				{"$limit": 1},
				{"$project": bson.M{
					"_id":       0,
					"character": "$key",
					"skin":      "$skins.k",
				}},
			},
			"hotWaifuSkin": []bson.M{
				{"$match": bson.M{"sex": "female"}},
				{"$unwind": "$skins"},
				{"$sort": bson.M{"skins.v.score.explicit": -1}},
				{"$limit": 1},
				{"$project": bson.M{
					"_id":       0,
					"character": "$key",
					"skin":      "$skins.k",
				}},
			},
			"lewdWaifuSkin": []bson.M{
				{"$match": bson.M{"sex": "female"}},
				{"$unwind": "$skins"},
				{"$match": bson.M{"skins.v.score.all": bson.M{"$gt": minScoreForSkin}}},
				{"$sort": bson.M{"skins.v.score.lewd": -1}},
				{"$limit": 1},
				{"$project": bson.M{
					"_id":       0,
					"character": "$key",
					"skin":      "$skins.k",
				}},
			},
			"pureWaifuSkin": []bson.M{
				{"$match": bson.M{"sex": "female"}},
				{"$unwind": "$skins"},
				{"$match": bson.M{"skins.v.score.all": bson.M{"$gt": minScoreForSkin}}},
				{"$sort": bson.M{"skins.v.score.pure": -1}},
				{"$limit": 1},
				{"$project": bson.M{
					"_id":       0,
					"character": "$key",
					"skin":      "$skins.k",
				}},
			},
			"fakeWaifu": []bson.M{
				{"$match": bson.M{"sex": "male"}},
				{"$sort": bson.M{"score.gender_swaps": -1}},
				{"$limit": 1},
				{"$project": bson.M{
					"_id":       0,
					"character": "$key",
				}},
			},
			"virginKillerWaifu": []bson.M{
				{"$match": bson.M{"sex": "female"}},
				{"$sort": bson.M{"score.virgin_killer_sweater": -1}},
				{"$limit": 1},
				{"$project": bson.M{
					"_id":       0,
					"character": "$key",
				}},
			},
			"selfieWaifu": []bson.M{
				{"$match": bson.M{"sex": "female"}},
				{"$sort": bson.M{"score.selfie": -1}},
				{"$limit": 1},
				{"$project": bson.M{
					"_id":       0,
					"character": "$key",
				}},
			},
		}},
		{"$project": bson.M{
			"fame_waifu":          bson.M{"$arrayElemAt": []interface{}{"$fameWaifu", 0}},
			"hot_waifu":           bson.M{"$arrayElemAt": []interface{}{"$hotWaifu", 0}},
			"lewd_waifu":          bson.M{"$arrayElemAt": []interface{}{"$lewdWaifu", 0}},
			"pure_waifu":          bson.M{"$arrayElemAt": []interface{}{"$pureWaifu", 0}},
			"fame_waifu_skin":     bson.M{"$arrayElemAt": []interface{}{"$fameWaifuSkin", 0}},
			"hot_waifu_skin":      bson.M{"$arrayElemAt": []interface{}{"$hotWaifuSkin", 0}},
			"lewd_waifu_skin":     bson.M{"$arrayElemAt": []interface{}{"$lewdWaifuSkin", 0}},
			"pure_waifu_skin":     bson.M{"$arrayElemAt": []interface{}{"$pureWaifuSkin", 0}},
			"fake_waifu":          bson.M{"$arrayElemAt": []interface{}{"$fakeWaifu", 0}},
			"virgin_killer_waifu": bson.M{"$arrayElemAt": []interface{}{"$virginKillerWaifu", 0}},
			"selfie_waifu":        bson.M{"$arrayElemAt": []interface{}{"$selfieWaifu", 0}},
		}},
	}

	return collection.Pipe(query).One(&ow.Achievements)
}
