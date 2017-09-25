package model

// Class ...
type Class string

// Class enums
const (
	Common    Class = "common"
	Rare      Class = "rare"
	Epic      Class = "epic"
	Legendary Class = "legendary"
)

// Skin contains data about skin
type Skin struct {
	Name  string `json:"name"`
	Class `json:"class"`
	Event `json:"event"`
	Tag   string `json:"tag"`
	Score `json:"score"`
}
