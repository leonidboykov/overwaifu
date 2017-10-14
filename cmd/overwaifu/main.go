package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/leonidboykov/overwaifu"
)

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

func main() {
	var chars []overwaifu.Character

	for _, c := range characters {
		var char overwaifu.Character
		if _, err := toml.DecodeFile("resources/"+c+".toml", &char); err != nil {
			log.Panicln(err)
		}

		if err := char.FetchScore(); err != nil {
			log.Panicln(err)
		}

		data, err := json.MarshalIndent(char, "", "  ")
		if err != nil {
			log.Panicln(err)
		}

		if err := ioutil.WriteFile("dest/"+c+".json", data, 0644); err != nil {
			log.Panicln(err)
		}

		chars = append(chars, char)
	}

	scoreAll := func(c1, c2 *overwaifu.Character) bool {
		return c1.Score.All > c2.Score.All
	}

	overwaifu.By(scoreAll).Sort(chars)

	for _, c := range chars {
		fmt.Println(c.Name, "-", c.Score.All)
	}
	fmt.Println()
	fmt.Println(chars[0].Name, "is the hottest", chars[0].Sex, "in the Overwatch")
}
