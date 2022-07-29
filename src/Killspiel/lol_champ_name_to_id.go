package Killspiel

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
)

// loadChamps erstellt eine map, die ChampionsNamen auf ChampionIds abbildet.
func loadChamps(file string) *map[string]string {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("An error occurred while opening champions.json %v\n", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("Error occured while closing champions.json: %v\n", err)
		}
	}(f)
	bites, _ := io.ReadAll(f)
	champions := &champion{}
	err = json.Unmarshal(bites, champions)
	if err != nil {
		log.Fatal(err)
	}

	champ := map[string]string{}

	for _, v := range champions.Data {
		champ[v.Name] = v.Key
	}

	return &champ
}

// champNamesToInt gibt das entsprechende int-Array zu dem gegebenen ChampionNamen-Array.
func champNamesToId(names *[]string) *[]int {
	champions := loadChamps("champions.json")
	res := make([]int, len(*names))
	var err error
	for i, name := range *names {
		res[i], err = strconv.Atoi((*champions)[name])
		if err != nil {
			log.Printf("%s konnte nicht als Held gefunden werden, wird ignoriert. Fehler: %v", name, err)
		}
	}
	return &res
}
