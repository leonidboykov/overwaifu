package overwaifu

// OverWaifu ...
type OverWaifu struct {
	Characters []Character
}

// FetchData ...
func (ow *OverWaifu) FetchData() {
	for _, c := range ow.Characters {
		c.FetchScore()
	}
}
