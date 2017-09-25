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

// Role ...
type Role string

// Class enums
const (
	Support Role = "support"
	Tank    Role = "tank"
	Defense Role = "defense"
	Offense Role = "offense"
)

// Character contains all main data about character
type Character struct {
	Name     string `json:"name" toml:"name"`
	RealName string `json:"real_name" toml:"realName"`
	Age      int    `json:"age" toml:"age"`
	Location string `json:"location" toml:"location"`
	Role     `json:"role" toml:"role"`
	Sex      `json:"sex" toml:"sex"`
	Skins    []Skin `json:"skins" toml:"skins"`
	Score    `json:"score"`
	Tag      string `json:"tag" toml:"tag"` // sankaku tag
}
