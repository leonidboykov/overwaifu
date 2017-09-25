package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
)

// Character contains all main data about character
type Character struct {
	Name     string `json:"name" toml:"name"`
	RealName string `json:"real_name" toml:"realName"`
	Age      int    `json:"age" toml:"age"`
	Location string `json:"location" toml:"location"`
	Role     string `json:"role" toml:"role"`
	Sex      string `json:"sex" toml:"sex"`
	Skins    []Skin `json:"skins" toml:"skins"`
	Score    `json:"score"`
	Tag      string `json:"tag" toml:"tag"` // sankaku tag
}

const (
	sankakuURL = "https://ias.sankakucomplex.com/tag/autosuggest?tag=%s"
)

// FetchScore from Sankaku Channel
func (c *Character) FetchScore() error {
	url := fmt.Sprintf(sankakuURL, c.Tag)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var obj []interface{}
	if err = json.Unmarshal(body, &obj); err != nil {
		return err
	}

	c.Score.All, err = strconv.Atoi(obj[2].(string))
	if err != nil {
		return err
	}

	return nil
}

// By is the type of a "less" function that defines the ordering of its Planet arguments.
type By func(c1, c2 *Character) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(characters []Character) {
	cs := &characterSorter{
		characters: characters,
		by:         by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(cs)
}

// planetSorter joins a By function and a slice of Characters to be sorted.
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
