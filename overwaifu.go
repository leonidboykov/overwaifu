package overwaifu

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

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

// // OverWaifu ...
// type OverWaifu struct {
// 	posts        []getmoe.Post
// 	PostsCount   int                   `json:"posts_count"`
// 	UpdatedAt    time.Time             `json:"updated_at"`
// 	LastPostTime time.Time             `json:"last_post_time"`
// 	Waifu        map[string]*Character `json:"waifu"`
// 	Husbando     map[string]*Character `json:"husbando"`
// 	Achievements `json:"achievements"`
// }

// OverWaifu holds all overwaifu results
type OverWaifu struct {
	Characters   map[string]*Character `json:"characters"`
	Achievements struct {
		FameWaifu         Achievement     `json:"fame_waifu" bson:"fame_waifu"`
		HotWaifu          Achievement     `json:"hot_waifu" bson:"hot_waifu"`
		LewdWaifu         Achievement     `json:"lewd_waifu" bson:"lewd_waifu"`
		PureWaifu         Achievement     `json:"pure_waifu" bson:"pure_waifu"`
		FakeWaifu         Achievement     `json:"fake_waifu" bson:"fake_waifu"`
		VirginKillerWaifu Achievement     `json:"virgin_killer_waifu" bson:"virgin_killer_waifu"`
		SelfieWaifu       Achievement     `json:"selfie_waifu" bson:"selfie_waifu"`
		FameWaifuSkin     SkinAchievement `json:"fame_waifu_skin" bson:"fame_waifu_skin"`
		HotWaifuSkin      SkinAchievement `json:"hot_waifu_skin" bson:"hot_waifu_skin"`
		LewdWaifuSkin     SkinAchievement `json:"lewd_waifu_skin" bson:"lewd_waifu_skin"`
		PureWaifuSkin     SkinAchievement `json:"pure_waifu_skin" bson:"pure_waifu_skin"`
	} `json:"achievements"`
}

// Achievement holds the character achievement
type Achievement struct {
	Character string `json:"character" bson:"character"`
}

// SkinAchievement ...
type SkinAchievement struct {
	Character string `json:"character" bson:"character"`
	Skin      string `json:"skin" bson:"skin"`
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
func New(db *mgo.Database) (*OverWaifu, error) {
	postsCollection := db.C("posts")
	characters, err := getCharactersList()
	if err != nil {
		return nil, err
	}

	ow := OverWaifu{
		Characters: make(map[string]*Character),
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
		ow.Characters[characters[i]].QueryScore(postsCollection)
		ow.Characters[characters[i]].QueryScoreSkins(postsCollection)
	}

	charactersCollection := db.C("characters")
	// if err := charactersCollection.DropCollection(); err != nil {
	// 	fmt.Println(err)
	// }
	for i := range characters {
		c := ow.Characters[characters[i]]
		// if err := charactersCollection.Insert(&c); err != nil {
		// 	fmt.Println(err)
		// }
		if err := charactersCollection.Update(bson.M{"key": c.Key}, &c); err != nil {
			fmt.Println(err)
		}
	}

	if err := ow.QueryAchievements(charactersCollection); err != nil {
		return nil, err
	}

	return &ow, nil
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
