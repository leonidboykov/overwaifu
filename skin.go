package overwaifu

import (
	"encoding/json"
	"strings"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
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
func (s *Skin) QueryScore(owner *Character, collection *mgo.Collection) error {
	query := []bson.M{
		{"$match": bson.M{"$and": []bson.M{
			{"tags": bson.M{"$in": owner.Tags}},
			{"tags": bson.M{"$in": s.Tags}},
		}}},
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
		}},
		{"$project": bson.M{
			"all":          bson.M{"$arrayElemAt": []interface{}{"$countAll.count", 0}},
			"safe":         bson.M{"$arrayElemAt": []interface{}{"$countSafe.count", 0}},
			"questionable": bson.M{"$arrayElemAt": []interface{}{"$countQuestionable.count", 0}},
			"explicit":     bson.M{"$arrayElemAt": []interface{}{"$countExplicit.count", 0}},
		}},
		{"$addFields": bson.M{
			"pure": bson.M{"$divide": []string{"$safe", "$all"}},
			"lewd": bson.M{"$divide": []string{"$explicit", "$all"}},
		}},
	}

	return collection.Pipe(query).One(&s.Score)
}
