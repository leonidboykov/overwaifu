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
	Name       string `json:"name"`
	RealName   string `json:"real_name"`
	Age        int    `json:"age"`
	Occupation string `json:"nationality"`
	Role       `json:"role"`
	Sex        `json:"sex"`
	Skins      []Skin `json:"skins"`
	Score      `json:"score"`
	Tag        string `json:"tag"` // sankaku tag
}
