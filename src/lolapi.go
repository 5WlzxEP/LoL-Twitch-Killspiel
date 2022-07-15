package Killspiel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type summoner struct {
	Id            string `json:"id"`
	AccountId     string `json:"accountId"`
	Puuid         string `json:"puuid"`
	Name          string `json:"name"`
	ProfileIconId int    `json:"profileIconId"`
	RevisionDate  int    `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}

type spectatorStruct struct {
	GameId       int64  `json:"gameId"`
	GameMode     string `json:"gameMode"`
	GameType     string `json:"gameType"`
	Participants [10]struct {
		ChampionId   int    `json:"championId"`
		SummonerName string `json:"summonerName"`
		SummonerId   string `json:"summonerId"`
	} `json:"participants"`
	GameLength int `json:"gameLength"`
	Status     struct {
		Message    string `json:"message"`
		StatusCode int    `json:"status_code"`
	} `json:"status"`
}

type killData struct {
	Metadata struct {
		MatchId      string   `json:"matchId"`
		Participants []string `json:"participants"`
	} `json:"metadata"`
	Info struct {
		GameId       int64  `json:"gameId"`
		GameMode     string `json:"gameMode"`
		GameType     string `json:"gameType"`
		Participants []struct {
			Assists              int  `json:"assists"`
			Deaths               int  `json:"deaths"`
			Kills                int  `json:"kills"`
			ParticipantId        int  `json:"participantId"`
			TeamEarlySurrendered bool `json:"teamEarlySurrendered"`
			Win                  bool `json:"win"`
		} `json:"participants"`
	} `json:"info"`
}

type game struct {
	matchId  int64
	playerId string
}

var aktuellesGame *game

func Statecontroll(lolid string) {
	//log.Println("Updating State")
	aktuellesGame.playerId = lolid
	for ; true; time.Sleep(1 * time.Minute) {
		res, err := http.Get(fmt.Sprintf("https://euw1.api.riotgames.com/lol/spectator/v4/active-games/by-summoner/%s?api_key=%s", lolid, config.Lolapikey))
		if err != nil {
			log.Fatal(err)
		}
		bites, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		sp := &spectatorStruct{}
		json.Unmarshal(bites, sp)
		switch config.State {
		case Idle:
			if sp.GameId != 0 && sp.GameLength < 120 && !config.otp {
				aktuellesGame.matchId = sp.GameId
				go StarteWette()
			} else if sp.GameId != 0 && sp.GameLength < 120 && config.otp {
				config.State = GameNoTrack
				for _, participant := range sp.Participants {
					if participant.SummonerId == aktuellesGame.playerId {
						if isElementOfArray(*config.champsId, participant.ChampionId) {
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

func GetLolID(lolaccount string) string {
	//log.Println("Getting account id")
	res, err := http.Get(fmt.Sprintf("https://euw1.api.riotgames.com/lol/summoner/v4/summoners/by-name/%s?api_key=%s", lolaccount, config.Lolapikey))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	bites, _ := ioutil.ReadAll(res.Body)
	summ := &summoner{}
	json.Unmarshal(bites, summ)
	return summ.Id
}

func lolidToPuuid() string {
	res, err := http.Get(fmt.Sprintf("https://euw1.api.riotgames.com/lol/summoner/v4/summoners/%s?api_key=%s", aktuellesGame.playerId, config.Lolapikey))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	bites, _ := ioutil.ReadAll(res.Body)
	summ := &summoner{}
	json.Unmarshal(bites, summ)
	return summ.Puuid
}

func GetKills() *killData {
	res, err := http.Get(fmt.Sprintf("https://europe.api.riotgames.com/lol/match/v5/matches/EUW1_%d?api_key=%s", aktuellesGame.matchId, config.Lolapikey))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	bites, _ := ioutil.ReadAll(res.Body)
	killd := &killData{}
	json.Unmarshal(bites, killd)
	//log.Printf("%v\n", killd)
	return killd
}

func isElementOfArray[T comparable](arr []T, ele T) bool {
	for _, v := range arr {
		if v == ele {
			return true
		}
	}
	return false
}
