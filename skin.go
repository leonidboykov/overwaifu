package overwaifu

import (
	"context"
	"encoding/json"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Skin contains data about skin
type Skin struct {
	Name   string   `json:"name" toml:"name"`
	Rarity string   `json:"rarity,omitempty" toml:"rarity"`
	Event  string   `json:"event,omitempty" toml:"event"`
	Tags   []string `json:"tags,omitempty" toml:"tags"`
	Key    string   `json:"key"`
	Score  `json:"score"`
}

// Skins provides custom marhaller for JSON
type Skins map[string]*Skin

// MarshalJSON ...
func (s Skins) MarshalJSON() ([]byte, error) {
	var skins []*Skin
	for _, v := range s {
		skins = append(skins, v)
	}
	return json.Marshal(skins)
}

// DefaultTags allows to assigs tags as `key_character`
func (s *Skin) DefaultTags(owner string) {
	if len(s.Tags) == 0 {
		key := strings.Replace(s.Key, "-", "_", -1)
		s.Tags = append(s.Tags, key+"_"+owner)
	}
}

// QueryScore ...
func (s *Skin) QueryScore(owner *Character, collection *mongo.Collection) error {
	ctx := context.TODO()

	pipeline := bson.A{
		bson.M{"$match": bson.M{"$and": bson.A{
			bson.M{"tags": bson.M{"$in": owner.Tags}},
			bson.M{"tags": bson.M{"$in": s.Tags}},
		}}},
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
		}},
		bson.M{"$project": bson.M{
			"all":          bson.M{"$arrayElemAt": bson.A{"$countAll.count", 0}},
			"safe":         bson.M{"$arrayElemAt": bson.A{"$countSafe.count", 0}},
			"questionable": bson.M{"$arrayElemAt": bson.A{"$countQuestionable.count", 0}},
			"explicit":     bson.M{"$arrayElemAt": bson.A{"$countExplicit.count", 0}},
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
	return cur.Decode(&s.Score)
}
