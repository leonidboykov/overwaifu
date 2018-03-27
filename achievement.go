package overwaifu

// Achievement holds the character achievement
type Achievement struct {
	Character string `json:"character" bson:"character"`
	Skin      string `json:"skin,omitempty" bson:"skin,omitempty"`
}
