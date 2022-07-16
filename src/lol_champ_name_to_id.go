package Killspiel

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type champion struct {
	Type string `json:"type"`
	Data map[string]struct {
		Id   string `json:"id"`
		Key  string `json:"key"`
		Name string `json:"name"`
	} `json:"data"`
}

func loadChamps(file string) *map[string]string {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("An error occurred while opening champions.json %v\n", err)
	}
	defer f.Close()
	bites, _ := ioutil.ReadAll(f)
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
