package character

import "github.com/leonidboykov/overwaifu/model"

// Mercy containt predefined data for OW's Mercy
var Mercy = model.Character{
	Name:     "Mercy",
	RealName: "Angela Ziegler",
	Sex:      model.Female,
	Skins: []model.Skin{
		{
			Name:  "Classic",
			Class: model.Common,
		},
		{
			Name:  "Celestial",
			Class: model.Rare,
		},
		{
			Name:  "Mist",
			Class: model.Rare,
		},
		{
			Name:  "Orchid",
			Class: model.Rare,
		},
		{
			Name:  "Verdant",
			Class: model.Rare,
		},
		{
			Name:  "Amber",
			Class: model.Epic,
		},
		{
			Name:  "Cobalt",
			Class: model.Epic,
			Tag:   "cobalt_mercy_(overwatch)",
		},
		{
			Name:  "Fortune",
			Class: model.Epic,
			Tag:   "fortune_mercy_(overwatch)",
		},
		{
			Name:  "Eidgenossin",
			Class: model.Epic,
			Event: model.SummerGames,
			Tag:   "eidgenossin_mercy_(overwatch)",
		},
		{
			Name:  "Sigrún",
			Class: model.Legendary,
		},
		{
			Name:  "Valkyrie",
			Class: model.Legendary,
		},
		{
			Name:  "Devil",
			Class: model.Legendary,
		},
		{
			Name:  "Imp",
			Class: model.Legendary,
			Tag:   "imp_mercy",
		},
		{
			Name:  "Witch",
			Class: model.Legendary,
			Event: model.HelloweenTerror,
			Tag:   "witch_mercy_(overwatch)",
		},
		{
			Name:  "Combat Medic Ziegler",
			Class: model.Legendary,
			Event: model.Uprising,
		},
		{
			Name:  "Winged Victory",
			Class: model.Legendary,
			Event: model.SummerGames,
			Tag:   "winged_victory_mercy_(overwatch)",
		},
	},
	Tag: "mercy_(overwatch)",
}
