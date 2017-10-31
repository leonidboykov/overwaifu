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
	Waifu        map[string]*Character `json:"waifu"`
	Husbando     map[string]*Character `json:"husbando"`
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

// New ...
func New(posts []getmoe.Post) (*OverWaifu, error) {
	ow := OverWaifu{
		posts:    posts,
		Waifu:    make(map[string]*Character),
		Husbando: make(map[string]*Character),
	}
	for i := range characters {
		var c Character
		_, err := toml.DecodeFile("resources/"+characters[i]+".toml", &c)
		if err != nil {
			return nil, err
		}
		if c.Sex == "female" {
			ow.Waifu[characters[i]] = &c
		} else {
			ow.Husbando[characters[i]] = &c
		}
	}

	return &ow, nil
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
	var fame, hot, virginKiller int
	var pure, lewd float64
	var fameSkin, hotSkin int
	var pureSkin, lewdSkin float64
	for i := range ow.Waifu {
		if ow.Waifu[i].Score.All > fame {
			ow.Achievements.FameWaifu = i
			fame = ow.Waifu[i].Score.All
		}

		if ow.Waifu[i].Score.Explicit > hot {
			ow.Achievements.HotWaifu = i
			hot = ow.Waifu[i].Score.Explicit
		}

		if ow.Waifu[i].Score.Pure > pure {
			ow.Achievements.PureWaifu = i
			pure = ow.Waifu[i].Score.Pure
		}

		if ow.Waifu[i].Score.Lewd > lewd {
			ow.Achievements.LewdWaifu = i
			lewd = ow.Waifu[i].Score.Lewd
		}

		if ow.Waifu[i].Score.VirginKillerSweater > virginKiller {
			ow.Achievements.VirginKillerWaifu = i
			virginKiller = ow.Waifu[i].Score.VirginKillerSweater
		}

		// Skins
		f, fSkin := ow.Waifu[i].FameSkin()
		if f > fameSkin {
			ow.Achievements.FameWaifuSkin = SkinAchievement{
				Owner: i,
				Skin:  fSkin,
			}
			fameSkin = f
		}

		h, hSkin := ow.Waifu[i].HotSkin()
		if h > hotSkin {
			ow.Achievements.HotWaifuSkin = SkinAchievement{
				Owner: i,
				Skin:  hSkin,
			}
			hotSkin = h
		}

		p, pSkin := ow.Waifu[i].PureSkin()
		if p > pureSkin {
			ow.Achievements.PureWaifuSkin = SkinAchievement{
				Owner: i,
				Skin:  pSkin,
			}
			pureSkin = p
		}

		l, lSkin := ow.Waifu[i].LewdSkin()
		if l > lewdSkin {
			ow.Achievements.LewdWaifuSkin = SkinAchievement{
				Owner: i,
				Skin:  lSkin,
			}
			lewdSkin = l
		}
	}

	var genderSwaps int
	for i := range ow.Husbando {
		if ow.Husbando[i].Score.GenderSwaps > genderSwaps {
			ow.Achievements.FakeWaifu = i
			genderSwaps = ow.Husbando[i].Score.GenderSwaps
		}
	}
}
