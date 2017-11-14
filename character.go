package overwaifu

import (
	"log"

	"github.com/leonidboykov/getmoe"
)

const (
	genderswapTag          = "genderswap"
	virginKillerSweaterTag = "virgin_killer_sweater"
	selfieTag              = "selfie"
)

const (
	minScoreForChar = 500
	minScoreForTag  = 75
)

// Character contains all main data about character
type Character struct {
	Name     string           `json:"name" toml:"name"`
	RealName string           `json:"real_name" toml:"realName"`
	Age      int              `json:"age" toml:"age"`
	Role     string           `json:"role" toml:"role"`
	Sex      string           `json:"sex" toml:"sex"`
	Skins    map[string]*Skin `json:"skins" toml:"skins"`
	Tag      string           `json:"tag" toml:"tag"` // sankaku tag
	Key      string           `json:"key"`
	Score    `json:"score"`
}

// UpdateSkinKey ...
func (c *Character) UpdateSkinKey() {
	for i := range c.Skins {
		c.Skins[i].Key = i
	}
}

// CalcScore ...
func (c *Character) CalcScore(posts []getmoe.Post) {
	// Score.All is how much pictures with this character on Sankaku Channel
	c.Score.All = len(posts)

	for i := range posts {
		if posts[i].HasTag(genderswapTag) {
			c.Score.GenderSwaps++
		}
		if posts[i].HasTag(virginKillerSweaterTag) {
			c.Score.VirginKillerSweater++
		}
		if posts[i].HasTag(selfieTag) {
			c.Score.Selfie++
		}
		switch posts[i].Rating {
		case "s":
			c.Score.Safe++
		case "q":
			c.Score.Questionable++
		case "e":
			c.Score.Explicit++
		default:
			log.Printf("overwaifu: got unknown rating %s", posts[i].Rating)
		}

		for j := range c.Skins {
			if posts[i].HasTag(c.Skins[j].Tag) {
				c.Skins[j].Score.All++
				switch posts[i].Rating {
				case "s":
					c.Skins[j].Score.Safe++
				case "q":
					c.Skins[j].Score.Questionable++
				case "e":
					c.Skins[j].Score.Explicit++
				}
			}
		}
	}

	if c.Score.All > minScoreForChar {
		c.Score.Lewd = float64(c.Score.Explicit) / float64(c.Score.All)
		c.Score.Pure = float64(c.Score.Safe) / float64(c.Score.All)
	}

	for i := range c.Skins {
		if c.Skins[i].Score.All > minScoreForTag {
			c.Skins[i].Score.Lewd = float64(c.Skins[i].Score.Explicit) / float64(c.Skins[i].Score.All)
			c.Skins[i].Score.Pure = float64(c.Skins[i].Score.Safe) / float64(c.Skins[i].Score.All)
		}
	}
}

// FameSkin ...
func (c *Character) FameSkin() (fame int, skin string) {
	for key := range c.Skins {
		if c.Skins[key].Score.All > fame {
			fame = c.Skins[key].Score.All
			skin = key
		}
	}
	return
}

// HotSkin ...
func (c *Character) HotSkin() (hot int, skin string) {
	for key := range c.Skins {
		if c.Skins[key].Score.Explicit > hot {
			hot = c.Skins[key].Score.Explicit
			skin = key
		}
	}
	return
}

// PureSkin ...
func (c *Character) PureSkin() (pure float64, skin string) {
	for key := range c.Skins {
		if c.Skins[key].Score.Pure > pure {
			pure = c.Skins[key].Score.Pure
			skin = key
		}
	}
	return
}

// LewdSkin ...
func (c *Character) LewdSkin() (lewd float64, skin string) {
	for key := range c.Skins {
		if c.Skins[key].Score.Lewd > lewd {
			lewd = c.Skins[key].Score.Lewd
			skin = key
		}
	}
	return
}
