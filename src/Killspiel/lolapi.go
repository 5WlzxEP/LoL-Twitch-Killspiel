package Killspiel

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var aktuellesGame *game

// StateControl sorgt für die regelmäßige (60s) Aktualisierung des GlobalConfig.State.
func StateControl(LoLId string) {
	aktuellesGame.playerId = LoLId
	for ; true; time.Sleep(1 * time.Minute) {
		res, err := http.Get(fmt.Sprintf("https://%s.api.riotgames.com/lol/spectator/v4/active-games/by-summoner/%s?api_key=%s", config.LoLRegion, LoLId, config.Lolapikey))
		if err != nil {
			log.Fatal(err)
		}
		bites, _ := io.ReadAll(res.Body)
		err = res.Body.Close()
		sp := &spectatorStruct{}
		err = json.Unmarshal(bites, sp)
		if sp.Status.StatusCode != 0 && sp.Status.StatusCode != 404 {
			log.Printf("Riot API responded an error: %d, %s\n", sp.Status.StatusCode, sp.Status.Message)
		}
		switch config.State {
		case Idle:
			if sp.GameId != 0 && sp.GameLength < 120 && !config.otp {
				aktuellesGame.matchId = sp.GameId
				go StarteWette()
			} else if sp.GameId != 0 && sp.GameLength < 120 && config.otp {
				config.State = GameNoTrack
				for _, participant := range sp.Participants {
					if participant.SummonerId == aktuellesGame.playerId {
						if isElementOfArray[int](config.champsId, participant.ChampionId) {
							aktuellesGame.matchId = sp.GameId
							go StarteWette()
							break
						}
					}
				}
			}
		case Spielphase:
			if sp.GameId == 0 {
				//config.State = Auswertungsphase
				go Auswertung()
			}
		case GameNoTrack:
			if sp.GameId == 0 {
				config.State = Idle
			}
		}
	}
}

// GetLolID gibt die ID zu einem LoL-Account aus
func GetLolID(lolaccount string) (string, error) {
	res, err := http.Get(fmt.Sprintf("https://%s.api.riotgames.com/lol/summoner/v4/summoners/by-name/%s?api_key=%s", config.LoLRegion, lolaccount, config.Lolapikey))
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error occured closing summoners api response: %v\n", err)
		}
	}(res.Body)
	bites, _ := io.ReadAll(res.Body)
	summ := &summoner{}
	err = json.Unmarshal(bites, summ)

	if summ.Status.StatusCode != 0 {
		return "", errors.New(summ.Status.Message)
	}

	return summ.Id, nil
}

// lolidToPuuid gibt die PUUID zum game.playerId aus.
func lolidToPuuid() string {
	res, err := http.Get(fmt.Sprintf("https://%s.api.riotgames.com/lol/summoner/v4/summoners/%s?api_key=%s", config.LoLRegion, aktuellesGame.playerId, config.Lolapikey))
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error occured closing summoners api response: %v\n", err)
		}
	}(res.Body)
	bites, _ := io.ReadAll(res.Body)
	summ := &summoner{}
	err = json.Unmarshal(bites, summ)
	return summ.Puuid
}

// GetKills gibt killData zu game.matchId aus
func GetKills() *killData {
	res, err := http.Get(fmt.Sprintf("https://%s.api.riotgames.com/lol/match/v5/matches/%s_%d?api_key=%s", config.LolServer, strings.ToUpper(config.LoLRegion), aktuellesGame.matchId, config.Lolapikey))
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(res.Body)
	bites, _ := io.ReadAll(res.Body)
	killData := &killData{}
	err = json.Unmarshal(bites, killData)
	return killData
}

// isElementOfArray checks if an element is part of an unsorted array.
func isElementOfArray[T comparable](arr *[]T, ele T) bool {
	//log.Println(*arr)
	for _, v := range *arr {
		if v == ele {
			return true
		}
	}
	return false
}

func LoLRegionToServer(region string) (string, LolServer, bool) {
	region = strings.ToLower(region)
	switch region {
	case "las":
		return "la2", America, false
	case "lan":
		return "la1", America, false
	}

	if isElementOfArray[string](&[]string{"br", "eun", "euw", "jp", "na", "oc", "tr"}, region) {
		region = fmt.Sprintf("%s1", region)
	}

	switch {
	case isElementOfArray[string](&[]string{"na1", "br1"}, region):
		return region, America, false
	case isElementOfArray[string](&[]string{"jp1", "kr"}, region):
		return region, Asia, false
	case isElementOfArray[string](&[]string{"eun1", "euw1", "ru", "tr1"}, region):
		return region, Europe, false
	case region == "oc1":
		return region, Sea, false
	}
	return "", "", true
}

// America
//BR1
//NA1

// Asia
//JP1
//KR

// Europe
//EUN1
//EUW1
//RU
//TR1

// Sea
//OC1
