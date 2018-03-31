package overwaifu

// Score holds all scores
type Score struct {
	All                 int     `json:"all,omitempty" bson:"all"`
	Safe                int     `json:"safe,omitempty" bson:"safe"`
	Questionable        int     `json:"questionable,omitempty" bson:"questionable"`
	Explicit            int     `json:"explicit,omitempty" bson:"explicit"`
	GenderSwaps         int     `json:"gender_swaps,omitempty" bson:"gender_swaps"`
	Lewd                float64 `json:"lewd,omitempty" bson:"lewd"`
	Pure                float64 `json:"pure,omitempty" bson:"pure"`
	VirginKillerSweater int     `json:"virgin_killer_sweater,omitempty" bson:"virgin_killer_sweater"`
	Selfie              int     `json:"selfie,omitempty" bson:"selfie"`
}
