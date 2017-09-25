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
	Name   string `json:"name" toml:"name"`
	Rarity `json:"rarity" toml:"rarity"`
	Event  string `json:"event" toml:"event"`
	Tag    string `json:"tag" toml:"tag"`
	Score  `json:"score"`
}
