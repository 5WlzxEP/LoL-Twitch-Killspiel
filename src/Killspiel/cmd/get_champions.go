package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	res, err := http.Get("https://ddragon.leagueoflegends.com/api/versions.json")
	if err != nil {
		log.Fatal(err)
	}
	bites, _ := io.ReadAll(res.Body)
	var versionen []string
	//log.Printf("%v\n", bites)
	err = json.Unmarshal(bites, &versionen)
	err = res.Body.Close()
	if err != nil {
		log.Printf("Error occured closing response Body of ddragon: %v\n", err)
	}

	res, err = http.Get(fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/data/de_DE/champion.json", versionen[0]))
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile("champions.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("Error occured closing champions.json: %v\n", err)
		}
	}(f)
	bites, _ = io.ReadAll(res.Body)
	_, _ = f.Write(bites)
	log.Println("Erfolgreich heruntergeladen")
}
