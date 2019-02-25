package overwaifu

import (
	"context"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	genderswapTag          = "genderswap"
	virginKillerSweaterTag = "virgin_killer_sweater"
	selfieTag              = "selfie"
)

const tagTemplate = "%s_(overwatch)"

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

// UpdateSkinDefaults ...
func (c *Character) UpdateSkinDefaults() {
	for k := range c.Skins {
		c.Skins[k].Key = k
		c.Skins[k].DefaultTags(c.Key)
	}
}

// QueryScore ...
func (c *Character) QueryScore(collection *mongo.Collection) error {
	ctx := context.TODO()
	pipeline := bson.A{
		bson.M{"$match": bson.M{"tags": bson.M{"$in": c.Tags}}},
		bson.M{"$facet": bson.M{
			"countAll": bson.A{bson.M{"$count": "count"}},
			"countSafe": bson.A{
				bson.M{"$match": bson.M{"rating": "s"}},
				bson.M{"$count": "count"},
			},
			"countQuestionable": bson.A{
				bson.M{"$match": bson.M{"rating": "q"}},
				bson.M{"$count": "count"},
			},
			"countExplicit": bson.A{
				bson.M{"$match": bson.M{"rating": "e"}},
				bson.M{"$count": "count"},
			},
			"countGenderSwaps": bson.A{
				bson.M{"$match": bson.M{"tags": bson.M{"$in": []string{"genderswap"}}}},
				bson.M{"$count": "count"},
			},
			"countVirginKillerSweater": bson.A{
				bson.M{"$match": bson.M{"tags": bson.M{"$in": []string{"virgin_killer_sweater"}}}},
				bson.M{"$count": "count"},
			},
			"countSelfie": bson.A{
				bson.M{"$match": bson.M{"tags": bson.M{"$in": []string{"selfie", "snapchat"}}}},
				bson.M{"$count": "count"},
			},
		}},
		bson.M{"$project": bson.M{
			"all":                   bson.M{"$arrayElemAt": bson.A{"$countAll.count", 0}},
			"safe":                  bson.M{"$arrayElemAt": bson.A{"$countSafe.count", 0}},
			"questionable":          bson.M{"$arrayElemAt": bson.A{"$countQuestionable.count", 0}},
			"explicit":              bson.M{"$arrayElemAt": bson.A{"$countExplicit.count", 0}},
			"gender_swaps":          bson.M{"$arrayElemAt": bson.A{"$countGenderSwaps.count", 0}},
			"virgin_killer_sweater": bson.M{"$arrayElemAt": bson.A{"$countVirginKillerSweater.count", 0}},
			"selfie":                bson.M{"$arrayElemAt": bson.A{"$countSelfie.count", 0}},
		}},
		bson.M{"$addFields": bson.M{
			"pure": bson.M{"$divide": []string{"$safe", "$all"}},
			"lewd": bson.M{"$divide": []string{"$explicit", "$all"}},
		}},
	}

	cur, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	cur.Next(ctx) // Go to the first element
	return cur.Decode(&c.Score)
}

// QueryScoreSkins ...
func (c *Character) QueryScoreSkins(collection *mongo.Collection) {
	for i := range c.Skins {
		// Skip error handing for skins
		c.Skins[i].QueryScore(c, collection)
	}
}
