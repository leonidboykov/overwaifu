package overwaifu

import (
	"github.com/BurntSushi/toml"
	"github.com/leonidboykov/getmoe"
)

var repository []getmoe.Post

// Supported characters, comment to disable
var characters = []string{
	"ana",
	"bastion",
	"doomfist",
	"dva",
	"genji",
	"hanzo",
	"junkrat",
	"lucio",
	"mccree",
	"mei",
	"mercy",
	"orisa",
	"pharah",
	"reaper",
	"reinhardt",
	"roadhog",
	"soldier76",
	"sombra",
	"symmetra",
	"torbjorn",
	"tracer",
	"widowmaker",
	"winston",
	"zarya",
	"zenyatta",
}

// OverWaifu ...
type OverWaifu struct {
	posts        []getmoe.Post
	Waifu        []Character `json:"waifu"`
	Husbando     []Character `json:"husbando"`
	Achievements `json:"achievements"`
}

// Achievements ...
type Achievements struct {
	FameWaifu         string          `json:"fame_waifu"`
	HotWaifu          string          `json:"hot_waifu"`
	LewdWaifu         string          `json:"lewd_waifu"`
	PureWaifu         string          `json:"pure_waifu"`
	FakeWaifu         string          `json:"fake_waifu"`
	VirginKillerWaifu string          `json:"virgin_killer_waifu"`
	FameWaifuSkin     SkinAchievement `json:"fame_waifu_skin"`
	HotWaifuSkin      SkinAchievement `json:"hot_waifu_skin"`
	LewdWaifuSkin     SkinAchievement `json:"lewd_waifu_skin"`
	PureWaifuSkin     SkinAchievement `json:"pure_waifu_skin"`
}

// SkinAchievement ...
type SkinAchievement struct {
	Owner string `json:"owner"`
	Skin  string `json:"skin"`
}

// FetchData ...
func (ow *OverWaifu) FetchData() {
	for i := range ow.Waifu {
		var p []getmoe.Post
		for j := range ow.posts {
			if ow.posts[j].HasTag(ow.Waifu[i].Tag) {
				p = append(p, ow.posts[j])
			}
		}
		ow.Waifu[i].CalcScore(p)
	}

	for i := range ow.Husbando {
		var p []getmoe.Post
		for j := range ow.posts {
			if ow.posts[j].HasTag(ow.Husbando[i].Tag) {
				p = append(p, ow.posts[j])
			}
		}
		ow.Husbando[i].CalcScore(p)
	}
}

// Analyse ...
func (ow *OverWaifu) Analyse() {
	SortCharacterBy(sortByAll).Sort(ow.Waifu)
	ow.Achievements.FameWaifu = ow.Waifu[0].Slug

	SortCharacterBy(sortByExplicit).Sort(ow.Waifu)
	ow.Achievements.HotWaifu = ow.Waifu[0].Slug

	SortCharacterBy(sortByLewd).Sort(ow.Waifu)
	ow.Achievements.LewdWaifu = ow.Waifu[0].Slug

	SortCharacterBy(sortByPure).Sort(ow.Waifu)
	ow.Achievements.PureWaifu = ow.Waifu[0].Slug

	SortCharacterBy(sortByGenderSwaps).Sort(ow.Husbando)
	ow.Achievements.FakeWaifu = ow.Husbando[0].Slug

	SortCharacterBy(sortByVirginKiller).Sort(ow.Waifu)
	ow.Achievements.VirginKillerWaifu = ow.Waifu[0].Slug

	SortCharacterBy(sortSkinsByAll).Sort(ow.Waifu)
	ow.Achievements.FameWaifuSkin = SkinAchievement{
		Owner: ow.Waifu[0].Slug,
		Skin:  ow.Waifu[0].Skins[0].Slug,
	}

	SortCharacterBy(sortSkinsByExplicit).Sort(ow.Waifu)
	ow.Achievements.HotWaifuSkin = SkinAchievement{
		Owner: ow.Waifu[0].Slug,
		Skin:  ow.Waifu[0].Skins[0].Slug,
	}

	SortCharacterBy(sortSkinsByLewd).Sort(ow.Waifu)
	ow.Achievements.LewdWaifuSkin = SkinAchievement{
		Owner: ow.Waifu[0].Slug,
		Skin:  ow.Waifu[0].Skins[0].Slug,
	}

	SortCharacterBy(sortSkinsByPure).Sort(ow.Waifu)
	ow.Achievements.PureWaifuSkin = SkinAchievement{
		Owner: ow.Waifu[0].Slug,
		Skin:  ow.Waifu[0].Skins[0].Slug,
	}
}

// New ...
func New(posts []getmoe.Post) (*OverWaifu, error) {
	ow := OverWaifu{
		posts: posts,
	}
	for i := range characters {
		var c Character
		_, err := toml.DecodeFile("resources/"+characters[i]+".toml", &c)
		if err != nil {
			return nil, err
		}
		if c.Sex == "female" {
			ow.Waifu = append(ow.Waifu, c)
		} else {
			ow.Husbando = append(ow.Husbando, c)
		}
	}

	return &ow, nil
}

func sortByAll(c1, c2 *Character) bool {
	return c1.Score.All > c2.Score.All
}

func sortByExplicit(c1, c2 *Character) bool {
	return c1.Score.Explicit > c2.Score.Explicit
}

func sortByLewd(c1, c2 *Character) bool {
	return c1.Score.Lewd > c2.Score.Lewd
}

func sortByPure(c1, c2 *Character) bool {
	return c1.Score.Pure > c2.Score.Pure
}

func sortByGenderSwaps(c1, c2 *Character) bool {
	return c1.Score.GenderSwaps > c2.Score.GenderSwaps
}

func sortByVirginKiller(c1, c2 *Character) bool {
	return c1.Score.VirginKillerSweater > c2.Score.VirginKillerSweater
}

func sortSkinsByAll(c1, c2 *Character) bool {
	sortSkin := func(s1, s2 *Skin) bool {
		return s1.Score.All > s2.Score.All
	}
	SortSkinBy(sortSkin).Sort(c1.Skins)
	SortSkinBy(sortSkin).Sort(c2.Skins)

	return c1.Skins[0].Score.All > c2.Skins[0].Score.All
}

func sortSkinsByExplicit(c1, c2 *Character) bool {
	sortSkin := func(s1, s2 *Skin) bool {
		return s1.Score.Explicit > s2.Score.Explicit
	}
	SortSkinBy(sortSkin).Sort(c1.Skins)
	SortSkinBy(sortSkin).Sort(c2.Skins)

	return c1.Skins[0].Score.Explicit > c2.Skins[0].Score.Explicit
}

func sortSkinsByLewd(c1, c2 *Character) bool {
	sortSkin := func(s1, s2 *Skin) bool {
		return s1.Score.Lewd > s2.Score.Lewd
	}
	SortSkinBy(sortSkin).Sort(c1.Skins)
	SortSkinBy(sortSkin).Sort(c2.Skins)

	return c1.Skins[0].Score.Lewd > c2.Skins[0].Score.Lewd
}

func sortSkinsByPure(c1, c2 *Character) bool {
	sortSkin := func(s1, s2 *Skin) bool {
		return s1.Score.Pure > s2.Score.Pure
	}
	SortSkinBy(sortSkin).Sort(c1.Skins)
	SortSkinBy(sortSkin).Sort(c2.Skins)

	return c1.Skins[0].Score.Pure > c2.Skins[0].Score.Pure
}
