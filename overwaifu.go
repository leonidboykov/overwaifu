package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/leonidboykov/overwaifu/model"
)

// Supported characters, comment to disable
var characters = []string{
	"ana",
	"dva",
	"mei",
	"mercy",
	// "orisa", // no nsfw cows!
	"pharah",
	"sombra",
	"symmetra",
	"tracer",
	"widowmaker",
	"zarya",
}

func main() {
	var chars []model.Character

	for _, c := range characters {
		var char model.Character
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

		if err := ioutil.WriteFile("results/"+c+".json", data, 0644); err != nil {
			log.Panicln(err)
		}

		chars = append(chars, char)
	}

	scoreAll := func(c1, c2 *model.Character) bool {
		return c1.Score.All > c2.Score.All
	}

	model.By(scoreAll).Sort(chars)

	for _, c := range chars {
		fmt.Println(c.Name, "-", c.Score.All)
	}
	fmt.Println()
	fmt.Println(chars[0].Name, "is the hottest", chars[0].Sex, "in the Overwatch")
}
