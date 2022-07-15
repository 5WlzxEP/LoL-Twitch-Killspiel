package Killspiel

import (
	"encoding/json"
	"fmt"
	"github.com/gempir/go-twitch-irc/v2"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type GameState uint16

const (
	Idle GameState = iota
	Wettphase
	Spielphase
	Auswertungsphase
	GameNoTrack
)

type Config struct {
	Username       string `json:"Username"`
	Oath           string `json:"Oath"`
	Wettdauer      int    `json:"Wettdauer"`
	Twitchchannel  string `json:"Twitchchannel"`
	Lolaccountname string `json:"Lolaccountname"`
	Lolapikey      string `json:"Lolapikey"`
	Joinmessage    bool   `json:"Joinmessage"`
	Logpath        string `json:"LogPath"`
	Prefix         string `json:"TwitchPrefix"`
	otp            bool
	Champs         *[]string `json:"Champions"`
	champsId       *[]int
	State          GameState
	TwitchClient   *twitch.Client
}

type result struct {
	Kills      int        `json:"Kills"`
	Teilnehmer Teilnehmer `json:"Teilnehmer"`
}
type Teilnehmer struct {
	Gewinner  []Teilnehmer2 `json:"Gewinner"`
	Verlierer []Teilnehmer2 `json:"Verlierer"`
}

type Teilnehmer2 struct {
	Name string `json:"Name"`
	Tipp int    `json:"Tipp"`
}

type fail struct {
	Teilnehmer []Teilnehmer2 `json:"teilnehmer"`
	MatchId    int64         `json:"matchId"`
	PlayerId   string        `json:"playerId"`
}

var bessereDaten map[int][]string
var config *Config
var daten map[string]int
var wettdauer time.Duration

func SetConfig(config2 *Config) {
	config = config2
	if config.otp = false; len(*config.Champs) > 0 {
		config.otp = true
		config.champsId = champNamesToId(config.Champs)
	}
	wettdauer = time.Duration(config.Wettdauer)
	bessereDaten = map[int][]string{}
	aktuellesGame = &game{}
}

func Message(messagechan chan twitch.PrivateMessage) {
	//var mess *twitch.PrivateMessage
	for true {
		message := <-messagechan
		if strings.HasPrefix(message.Message, "!vote") {
			_, value, found := strings.Cut(message.Message, " ")
			if found {
				v, err := strconv.Atoi(value)
				if err == nil {
					daten[message.User.DisplayName] = v
					//log.Println(message.User.DisplayName, v)
				}
			}
		}

	}
}

func StarteWette() {
	log.Println("Starte Wettphase")
	config.State = Wettphase
	config.TwitchClient.Say(config.Twitchchannel, config.Prefix+" Killspiel hat begonnen, nimm mit '!vote [Zahl]' teil.")
	daten = make(map[string]int)
	time.Sleep(wettdauer * time.Second)
	config.State = Spielphase

	var text string
	switch len(daten) {
	case 0:
		text = "hat keine"
	case 1:
		text = "hat eine Person"
	default:
		text = fmt.Sprintf("haben %d Personen", len(daten))
	}
	log.Printf("Wettphase angeschloßen, %d Teilnehmer", len(daten))
	config.TwitchClient.Say(config.Twitchchannel, fmt.Sprintf("%s Killspiel-Teilnahme beendet, es %s teilgenommen.", config.Prefix, text))

	// Umformatierung der Daten für eine bessere Auswertung
	for player, points := range daten {
		bessereDaten[points] = append(bessereDaten[points], player)
	}
}

func Auswertung() {
	log.Println("Starte Auswertungsphase")
	config.State = Auswertungsphase
	killd := GetKills()
	//log.Println(killd)
	//log.Println(aktuellesGame)
	var ind int = 11
	player := lolidToPuuid()
	for i := range killd.Metadata.Participants {
		if killd.Metadata.Participants[i] == player {
			ind = i
			break
		}
	}

	// check if dir exist and if not creates it
	fileInfo, err := os.Stat("results")
	if err != nil || !fileInfo.IsDir() {
		//log.Fatal(err)
		err = os.Mkdir("results", 0750)
		if err != nil {
			log.Print("Could not create dir.")
		}
	}

	if ind == 11 {
		log.Println("player not found in result")
		//log.Fatal("player not found in result")

		t := make([]Teilnehmer2, len(daten))
		failed := fail{MatchId: aktuellesGame.matchId, PlayerId: aktuellesGame.playerId, Teilnehmer: t}

		index := 0
		for k, l := range bessereDaten {
			for _, i := range l {
				failed.Teilnehmer[index] = Teilnehmer2{Name: i, Tipp: k}
				index++
			}
		}

		file, err := os.Create(fmt.Sprintf("results/error_%d.json", aktuellesGame.matchId))
		if err != nil {
			log.Fatal(err)
		}
		bites, err := json.MarshalIndent(failed, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		file.Write(bites)
		log.Printf("player not found in result. But saved tipps in results/error_%s.json\n", aktuellesGame.matchId)

	} else {

		if killd.Info.Participants[ind].TeamEarlySurrendered || killd.Info.Participants[(ind+5)%10].TeamEarlySurrendered {
			log.Println("Auswertung abgebrochen, da geremaked wurde.")
			config.TwitchClient.Say(config.Twitchchannel, config.Prefix+" Killspiel wurde abgebrochen, da Remaked wurde.")

		} else {
			kills := killd.Info.Participants[ind].Kills
			gewinner := bessereDaten[kills]
			res := result{Kills: kills, Teilnehmer: Teilnehmer{Gewinner: make([]Teilnehmer2, len(gewinner)), Verlierer: make([]Teilnehmer2, 0)}}
			for i, g := range gewinner {
				res.Teilnehmer.Gewinner[i] = Teilnehmer2{Name: g, Tipp: kills}
			}
			for k, l := range bessereDaten {
				if k == kills {
					continue
				}
				for _, i := range l {
					res.Teilnehmer.Verlierer = append(res.Teilnehmer.Verlierer, Teilnehmer2{Name: i, Tipp: k})
				}
			}

			file, err := os.Create(fmt.Sprintf("results/%d.json", aktuellesGame.matchId))
			if err != nil {
				log.Fatal(err)
			}
			bites, err := json.MarshalIndent(res, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			file.Write(bites)
			log.Printf("Ergebnis gespeichert in results/%d.json\n", aktuellesGame.matchId)

			var haben string
			switch len(gewinner) {
			case 1:
				haben = "hat"
			case 0:
				haben = "Keiner hat"
			default:
				haben = "haben"

			}

			config.TwitchClient.Say(config.Twitchchannel,
				fmt.Sprintf("%s Killspiel wurde beendet. %s hat %d Kill(s) gemacht. %s %s die richtige Killanzahl getippt.",
					config.Prefix, config.Lolaccountname, kills, strings.Join(gewinner, ", "), haben))
		}
	}
	daten = map[string]int{}
	bessereDaten = map[int][]string{}
	config.State = Idle
}
