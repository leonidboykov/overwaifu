package character

import "github.com/leonidboykov/overwaifu/model"

// Mercy containt predefined data for OW's Mercy
var Mercy = model.Character{
	Name:        "Mercy",
	RealName:    "Angela Ziegler",
	Age:         37,
	Nationality: "Swiss",
	Class:       model.Support,
	Sex:         model.Female,
	Tag:         "mercy_(overwatch)",
	Skins: []model.Skin{
		{
			Name:   "Classic",
			Rarity: model.Common,
		},
		{
			Name:   "Celestial",
			Rarity: model.Rare,
		},
		{
			Name:   "Mist",
			Rarity: model.Rare,
		},
		{
			Name:   "Orchid",
			Rarity: model.Rare,
		},
		{
			Name:   "Verdant",
			Rarity: model.Rare,
		},
		{
			Name:   "Amber",
			Rarity: model.Epic,
		},
		{
			Name:   "Cobalt",
			Rarity: model.Epic,
			Tag:    "cobalt_mercy_(overwatch)",
		},
		{
			Name:   "Fortune",
			Rarity: model.Epic,
			Tag:    "fortune_mercy_(overwatch)",
		},
		{
			Name:   "Eidgenossin",
			Rarity: model.Epic,
			Event:  model.SummerGames,
			Tag:    "eidgenossin_mercy_(overwatch)",
		},
		{
			Name:   "Sigrún",
			Rarity: model.Legendary,
			Tag:    "sigrun_mercy_(overwatch)",
		},
		{
			Name:   "Valkyrie",
			Rarity: model.Legendary,
			Tag:    "valkyrie_mercy_(overwatch)",
		},
		{
			Name:   "Devil",
			Rarity: model.Legendary,
			Tag:    "devil_mercy_(overwatch)",
		},
		{
			Name:   "Imp",
			Rarity: model.Legendary,
			Tag:    "imp_mercy",
		},
		{
			Name:   "Witch",
			Rarity: model.Legendary,
			Event:  model.HelloweenTerror,
			Tag:    "witch_mercy_(overwatch)",
		},
		{
			Name:   "Combat Medic Ziegler",
			Rarity: model.Legendary,
			Event:  model.Uprising,
			Tag:    "combat_medic_ziegler_(overwatch)",
		},
		{
			Name:   "Winged Victory",
			Rarity: model.Legendary,
			Event:  model.SummerGames,
			Tag:    "winged_victory_mercy_(overwatch)",
		},
	},
}
