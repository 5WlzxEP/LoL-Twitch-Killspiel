package Killspiel

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var aktuellesGame *game

// StateControl sorgt für die regelmäßige (60s) Aktualisierung des GlobalConfig.State.
func StateControl(LoLId string) {
	aktuellesGame.playerId = LoLId
	for ; true; time.Sleep(1 * time.Minute) {
		res, err := http.Get(fmt.Sprintf("https://euw1.api.riotgames.com/lol/spectator/v4/active-games/by-summoner/%s?api_key=%s", LoLId, config.Lolapikey))
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
func GetLolID(lolaccount string) string {
	res, err := http.Get(fmt.Sprintf("https://euw1.api.riotgames.com/lol/summoner/v4/summoners/by-name/%s?api_key=%s", lolaccount, config.Lolapikey))
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
	return summ.Id
}

// lolidToPuuid gibt die PUUID zum game.playerId aus.
func lolidToPuuid() string {
	res, err := http.Get(fmt.Sprintf("https://euw1.api.riotgames.com/lol/summoner/v4/summoners/%s?api_key=%s", aktuellesGame.playerId, config.Lolapikey))
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
	res, err := http.Get(fmt.Sprintf("https://europe.api.riotgames.com/lol/match/v5/matches/EUW1_%d?api_key=%s", aktuellesGame.matchId, config.Lolapikey))
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
	log.Println(*arr)
	for _, v := range *arr {
		if v == ele {
			return true
		}
	}
	return false
}
