package character

import "github.com/leonidboykov/overwaifu/model"

// Widowmaker ...
var Widowmaker = model.Character{
	Name:        "Widowmaker",
	RealName:    "Amélie Lacroix",
	Age:         33,
	Nationality: "French",
	Class:       model.Support,
	Sex:         model.Female,
	Tag:         "widowmaker_(overwatch)",
	Skins: []model.Skin{
		{
			Name:   "Classic",
			Rarity: model.Common,
		},
		{
			Name:   "Ciel",
			Rarity: model.Rare,
		},
		{
			Name:   "Nuit",
			Rarity: model.Rare,
		},
		{
			Name:   "Rose",
			Rarity: model.Rare,
		},
		{
			Name:   "Vert",
			Rarity: model.Rare,
		},
		{
			Name:   "Patina",
			Rarity: model.Epic,
		},
		{
			Name:   "Winter",
			Rarity: model.Epic,
		},
		{
			Name:   "Tricolore",
			Rarity: model.Epic,
			Event:  model.SummerGames,
		},
		{
			Name:   "Odette",
			Rarity: model.Legendary,
		},
		{
			Name:   "Odile",
			Rarity: model.Legendary,
		},
		{
			Name:   "Comtesse",
			Rarity: model.Legendary,
		},
		{
			Name:   "Huntress",
			Rarity: model.Legendary,
		},
		{
			Name:   "Talon",
			Rarity: model.Legendary,
			Event:  model.Uprising,
		},
		{
			Name:   "Côte d'Azur",
			Rarity: model.Legendary,
			Event:  model.SummerGames,
			Tag:    "cote_d'azur_widowmaker",
		},
		{
			Name:   "Noire",
			Rarity: model.Legendary,
			Tag:    "noire_widowmaker_(overwatch)",
		},
	},
}
