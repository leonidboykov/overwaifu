package overwaifu

// OverWaifu ...
type OverWaifu struct {
	Characters []Character
}

// Achievements ...
type Achievements struct {
	MostPopularWaifu Character `json:"most_popular_waifu"` // best all score
	MostSexyWaifu    Character `json:"most_sexy_waifu"`    // best explicit score
	BestWaifuSkin    Character `json:"best_waifu_skin"`    // the most popular skin
}

// FetchData ...
func (ow *OverWaifu) FetchData() {
	for _, c := range ow.Characters {
		c.FetchScore()
	}
}
