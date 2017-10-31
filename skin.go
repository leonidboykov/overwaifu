package overwaifu

import "sort"

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
	Key    string `json:"key"`
	Score  `json:"score"`
}

// SortSkinBy is the type of a "less" function that defines the ordering of its Characters arguments.
type SortSkinBy func(c1, c2 *Skin) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by SortSkinBy) Sort(skins []Skin) {
	cs := &skinSorter{
		skins: skins,
		by:    by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(cs)
}

// characterSorter joins a By function and a slice of Characters to be sorted.
type skinSorter struct {
	skins []Skin
	by    func(p1, p2 *Skin) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *skinSorter) Len() int {
	return len(s.skins)
}

// Swap is part of sort.Interface.
func (s *skinSorter) Swap(i, j int) {
	s.skins[i], s.skins[j] = s.skins[j], s.skins[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *skinSorter) Less(i, j int) bool {
	return s.by(&s.skins[i], &s.skins[j])
}
