package Auswertung

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

type Global struct {
	Results      string `json:"results"`
	LogFile      string `json:"logFile"`
	BaserowTable int    `json:"baserowTable"`
	BaserowAPI   string `json:"baserowAPI"`
}

type result struct {
	MatchId  int64            `json:"matchId"`
	PlayerId string           `json:"playerId"`
	Kills    int              `json:"kills"`
	Tipps    map[int][]string `json:"Tipps"`
}

type search struct {
	Count    int            `json:"count"`
	Next     interface{}    `json:"next"`
	Previous interface{}    `json:"previous"`
	Results  []searchResult `json:"results"`
}

type searchResult struct {
	Id     int    `json:"id"`
	Order  string `json:"order"`
	Name   string `json:"Name"`
	Punkte string `json:"Punkte"`
	Spiele string `json:"Spiele"`
}

type postData struct {
	Name   string `json:"Name"`
	Punkte int    `json:"Punkte"`
	Spiele int    `json:"Spiele"`
}

var global *Global

func SetGlobal(global2 *Global) {
	global = global2
}

// Auswerten werten alle json-Dateien aus, die in dem von Results spezifiziertem Ordner liegt.
func Auswerten() (i int) {

	dir, err := os.ReadDir(global.Results)
	if err != nil {
		log.Printf("Could not read results folder: %v\n", err)
		return 0
	}
	var raw []byte
	var file *os.File
	var res result

	for _, d := range dir {
		if d.IsDir() {
			continue
		} else if !strings.HasSuffix(d.Name(), ".json") || !(len(d.Name()) == 14 || len(d.Name()) == 15) {
			continue
		}
		file, err = os.Open(path.Join(global.Results, d.Name()))
		if err != nil {
			log.Printf("error open %s, %v\n", d.Name(), err)
		}
		raw, _ = io.ReadAll(file)
		_ = file.Close()
		err = json.Unmarshal(raw, &res)
		if err != nil {
			log.Printf("cound not parse %s, %v\n", d.Name(), err)
			continue
		}
		var wg sync.WaitGroup
		for kills, players := range res.Tipps {
			for _, player := range players {
				wg.Add(1)
				//fmt.Println(player, kills)
				go func(player string, treffer bool) {
					defer wg.Done()
					baserowBase(player, treffer)
				}(player, kills == res.Kills)
			}
		}
		err = os.Remove(path.Join(global.Results, d.Name()))
		if err != nil {
			log.Printf("Cloud not remove %s, %v\n", d.Name(), err)
		}
		i++
		wg.Wait()
	}
	return
}

func baserowBase(name string, treffer bool) {

	res, err := baserowGet(fmt.Sprintf("https://api.baserow.io/api/database/rows/table/%d/?search=%s&user_field_names=true", global.BaserowTable, name))
	if err != nil {
		log.Printf("Cloud not get %s data in the db: %v\n", name, err)
	}
	raw, _ := io.ReadAll(res.Body)
	_ = res.Body.Close()
	var Search search
	_ = json.Unmarshal(raw, &Search)

	trefferInt := trefferToInt(treffer)

	if Search.Count == 0 {
		//case: Name doesn't exist
		baserowInsert(name, trefferInt)
	} else if Search.Count == 1 {
		//case: Name exist, no other name like r'Name.*'
		baserowUpdate(name, trefferInt, Search.Results[0])
	} else {
		// Name exist but multiple times.
		// e.g. name="test", but baserow have columns with "test1", "test2"
		var found bool
		for _, r := range Search.Results {
			if r.Name == name {
				found = true
				baserowUpdate(name, trefferInt, r)
				break
			}
		}
		if !found {
			log.Printf("To many result with Name=%s, no one matched. Creating new entry\n", name)
			baserowInsert(name, trefferInt)
		}

	}

	//fmt.Println(string(raw))
}

func baserowUpdate(name string, treffer int, result searchResult) {
	punkte, _ := strconv.Atoi(result.Punkte)
	spiele, _ := strconv.Atoi(result.Spiele)
	marshal, err := json.Marshal(postData{Name: name, Punkte: punkte + treffer, Spiele: spiele + 1})
	if err != nil {
		log.Printf("error jsonify data: %v\n", err)
	}

	_, err = baserowPatch(fmt.Sprintf("https://api.baserow.io/api/database/rows/table/%d/%d/?user_field_names=true", global.BaserowTable, result.Id), marshal)
	if err != nil {
		log.Printf("error while posting init for %s, %v\n", name, err)
	}
}

func baserowInsert(name string, treffer int) {
	marshal, err := json.Marshal(postData{Name: name, Punkte: treffer, Spiele: 1})
	if err != nil {
		log.Printf("error jsonify data: %v\n", err)
	}
	_, err = baserowPost(fmt.Sprintf("https://api.baserow.io/api/database/rows/table/%d/?user_field_names=true", global.BaserowTable), marshal)
	if err != nil {
		log.Printf("error while posting init for %s, %v\n", name, err)
	}
}

func baserowGet(url string) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", global.BaserowAPI))
	return client.Do(req)
}

func baserowPost(url string, data []byte) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", global.BaserowAPI))
	return client.Do(req)
}

func baserowPatch(url string, data []byte) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(data))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", global.BaserowAPI))
	return client.Do(req)
}

// trefferToInt return 1 if true else 0
func trefferToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
