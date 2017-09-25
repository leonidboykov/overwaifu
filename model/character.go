package model

// Sex represents character's sex, we'll use it to determine gender swaps, like
// fem-Hanzo
type Sex string

// Sex enums
const (
	Male   Sex = "male"
	Female Sex = "female"
	NA     Sex = "not_applicable"
)

// Class ...
type Class string

// Class enums
const (
	Support Class = "support"
	Tank    Class = "tank"
	Defense Class = "defense"
	Offense Class = "offense"
)

// Character contains all main data about character
type Character struct {
	Name        string `json:"name"`
	RealName    string `json:"real_name"`
	Age         int    `json:"age"`
	Nationality string `json:"nationality"`
	Class       `json:"class"`
	Sex         `json:"sex"`
	Skins       []Skin `json:"skins"`
	Score       `json:"score"`
	Tag         string `json:"tag"` // Sankaku's tag
}
