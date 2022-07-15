package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	res, err := http.Get("https://ddragon.leagueoflegends.com/api/versions.json")
	if err != nil {
		log.Fatal(err)
	}
	bites, _ := ioutil.ReadAll(res.Body)
	var versionen []string
	//log.Printf("%v\n", bites)
	err = json.Unmarshal(bites, &versionen)
	res.Body.Close()

	res, err = http.Get(fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/data/de_DE/champion.json", versionen[0]))
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile("champions.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	bites, _ = ioutil.ReadAll(res.Body)
	f.Write(bites)
	log.Println("Erfolgreich heruntergeladen")
}
