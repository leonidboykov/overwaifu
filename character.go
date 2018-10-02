package overwaifu

import (
	"encoding/json"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const (
	genderswapTag          = "genderswap"
	virginKillerSweaterTag = "virgin_killer_sweater"
	selfieTag              = "selfie"
)

// Character contains all main data about character
type Character struct {
	Name     string   `json:"name" toml:"name"`
	RealName string   `json:"real_name" toml:"realName"`
	Age      int      `json:"age" toml:"age"`
	Role     string   `json:"role" toml:"role"`
	Sex      string   `json:"sex" toml:"sex"`
	Skins    Skins    `json:"skins" toml:"skins"`
	Tags     []string `json:"tags" toml:"tags"` // sankaku tag
	Key      string   `json:"key"`
	Score    `json:"score"`
}

// Characters provides custom marhaller for JSON
type Characters map[string]*Character

// MarshalJSON ...
func (c Characters) MarshalJSON() ([]byte, error) {
	var characters []*Character
	for _, v := range c {
		characters = append(characters, v)
	}
	return json.Marshal(characters)
}

// UpdateSkinKey ...
func (c *Character) UpdateSkinKey() {
	for k := range c.Skins {
		c.Skins[k].Key = k
	}
}

// QueryScore ...
func (c *Character) QueryScore(collection *mgo.Collection) error {
	// Build query
	query := []bson.M{
		{"$match": bson.M{"tags": bson.M{"$in": c.Tags}}},
		{"$facet": bson.M{
			"countAll": []bson.M{{"$count": "count"}},
			"countSafe": []bson.M{
				{"$match": bson.M{"rating": "s"}},
				{"$count": "count"},
			},
			"countQuestionable": []bson.M{
				{"$match": bson.M{"rating": "q"}},
				{"$count": "count"},
			},
			"countExplicit": []bson.M{
				{"$match": bson.M{"rating": "e"}},
				{"$count": "count"},
			},
			"countGenderSwaps": []bson.M{
				{"$match": bson.M{"tags": bson.M{"$in": []string{"genderswap"}}}},
				{"$count": "count"},
			},
			"countVirginKillerSweater": []bson.M{
				{"$match": bson.M{"tags": bson.M{"$in": []string{"virgin_killer_sweater"}}}},
				{"$count": "count"},
			},
			"countSelfie": []bson.M{
				{"$match": bson.M{"tags": bson.M{"$in": []string{"selfie", "snapchat"}}}},
				{"$count": "count"},
			},
		}},
		{"$project": bson.M{
			"all":                   bson.M{"$arrayElemAt": []interface{}{"$countAll.count", 0}},
			"safe":                  bson.M{"$arrayElemAt": []interface{}{"$countSafe.count", 0}},
			"questionable":          bson.M{"$arrayElemAt": []interface{}{"$countQuestionable.count", 0}},
			"explicit":              bson.M{"$arrayElemAt": []interface{}{"$countExplicit.count", 0}},
			"gender_swaps":          bson.M{"$arrayElemAt": []interface{}{"$countGenderSwaps.count", 0}},
			"virgin_killer_sweater": bson.M{"$arrayElemAt": []interface{}{"$countVirginKillerSweater.count", 0}},
			"selfie":                bson.M{"$arrayElemAt": []interface{}{"$countSelfie.count", 0}},
		}},
		{"$addFields": bson.M{
			"pure": bson.M{"$divide": []string{"$safe", "$all"}},
			"lewd": bson.M{"$divide": []string{"$explicit", "$all"}},
		}},
	}

	return collection.Pipe(query).One(&c.Score)
}

// QueryScoreSkins ...
func (c *Character) QueryScoreSkins(collection *mgo.Collection) {
	for i := range c.Skins {
		c.Skins[i].QueryScore(c, collection)
	}
}
