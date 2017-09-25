package model

// Rarity ...
type Rarity string

// Rarity enums
const (
	Common    Rarity = "common"
	Rare      Rarity = "rare"
	Epic      Rarity = "epic"
	Legendary Rarity = "legendary"
)

// Skin contains data about skin
type Skin struct {
	Name   string `json:"name"`
	Rarity `json:"rarity"`
	Event  `json:"event"`
	Tag    string `json:"tag"`
	Score  `json:"score"`
}
