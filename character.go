package overwaifu

import (
	"log"
	"sort"

	"github.com/leonidboykov/getmoe"
)

const (
	genderswapTag          = "genderswap"
	virginKillerSweaterTag = "virgin_killer_sweater"
)

const (
	minScoreForChar = 500
	minScoreForTag  = 75
)

// Character contains all main data about character
type Character struct {
	Name     string `json:"name" toml:"name"`
	Slug     string `json:"slug" toml:"slug"`
	RealName string `json:"real_name" toml:"realName"`
	Age      int    `json:"age" toml:"age"`
	Role     string `json:"role" toml:"role"`
	Sex      string `json:"sex" toml:"sex"`
	Skins    []Skin `json:"skins" toml:"skins"`
	Score    `json:"score"`
	Tag      string `json:"tag" toml:"tag"` // sankaku tag
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

// SortCharacterBy is the type of a "less" function that defines the ordering of its Characters arguments.
type SortCharacterBy func(c1, c2 *Character) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by SortCharacterBy) Sort(characters []Character) {
	cs := &characterSorter{
		characters: characters,
		by:         by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(cs)
}

// characterSorter joins a By function and a slice of Characters to be sorted.
type characterSorter struct {
	characters []Character
	by         func(p1, p2 *Character) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *characterSorter) Len() int {
	return len(s.characters)
}

// Swap is part of sort.Interface.
func (s *characterSorter) Swap(i, j int) {
	s.characters[i], s.characters[j] = s.characters[j], s.characters[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *characterSorter) Less(i, j int) bool {
	return s.by(&s.characters[i], &s.characters[j])
}
