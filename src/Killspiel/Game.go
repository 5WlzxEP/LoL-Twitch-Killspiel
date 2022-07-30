package Killspiel

import (
	"encoding/json"
	"fmt"
	"github.com/gempir/go-twitch-irc/v2"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// bessereDaten hat den Vorteil, dass die Auswertung, wer Gewonnen und wer Verloren hat, leichter ist.
var bessereDaten map[int][]string
var config *GlobalConfig

// daten hat den Vorteil, dass während der Einsendephase die Zuschauer ihren Tipp beliebig ändern können.
var daten map[string]int
var wettdauer time.Duration

func SetConfig(config2 *GlobalConfig) {
	config = config2
	if config.otp = false; len(*config.Champs) > 0 {
		config.otp = true
		config.champsId = champNamesToId(config.Champs)
	}
	wettdauer = time.Duration(config.Wettdauer)
	bessereDaten = map[int][]string{}
	aktuellesGame = &game{}
}

// Message verarbeitet die eingehenden Nachrichten in der Zeit, in der die Wettphase läuft
func Message(messages chan twitch.PrivateMessage) {
	//var mess *twitch.PrivateMessage
	for true {
		message := <-messages
		if strings.HasPrefix(message.Message, "!vote ") {
			_, value, found := strings.Cut(message.Message, " ")
			if found {
				v, err := strconv.Atoi(value)
				if err == nil {
					daten[message.User.DisplayName] = v
				}
			}
		}

	}
}

// StarteWette startet die Wette. D.h. es wird für die Zeit von GlobalConfig.Wettdauer können die Zuschauer im Twitchchat
// ihre Tipps angeben. Danach wird die Wette automatisch geschlossen und es werden keine weiteren Tippe angenommen.
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
		text = "hat keiner"
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
	daten = map[string]int{}

	// save data as tmp_{game.matchId}.json for the case, that the programm crashes before the data get Ausgewertet
	f, err := os.Create(fmt.Sprintf("results/tmp_%d.json", aktuellesGame.matchId))
	if err != nil {
		log.Printf("Error occured while saving data to tmp-File. %v\n", err)
		return
	}
	defer func() {
		err := f.Close()
		if err != nil {
			return
		}
	}()
	dataBytes, err := json.Marshal(bessereDaten)
	if err != nil {
		log.Printf("Error occured while Marshalling data: %v\n", err)
		return
	}
	_, err = f.Write(dataBytes)
	if err != nil {
		log.Printf("Error occured while writing data: %v\n", err)
	}

}

// Auswertung wertet die bessereDaten aus. Dazu speichert es die Ergebnisse in {game.matchId}.json
func Auswertung() {
	log.Println("Starte Auswertungsphase")
	config.State = Auswertungsphase
	killd := GetKills()

	var ind = 11
	for i, v := range killd.Metadata.Participants {
		if v == config.lolPUUID {
			ind = i
			break
		}
	}

	// check if dir exist and if not creates it
	fileInfo, err := os.Stat("results")
	if err != nil || !fileInfo.IsDir() {
		err = os.Mkdir("results", 0750)
		if err != nil {
			log.Print("Could not create dir.")
		}
	}

	// if player is not found in the result
	if ind == 11 {
		res := result{MatchId: aktuellesGame.matchId, PlayerId: aktuellesGame.playerId, Kills: -1, Tipps: bessereDaten}

		file, err := os.Create(fmt.Sprintf("results/error_%d.json", aktuellesGame.matchId))
		if err != nil {
			log.Fatal(err)
		}
		bites, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		_, err = file.Write(bites)
		if err != nil {
			log.Printf("Error occurded writing error res: %v\n", err)
			return
		}
		log.Printf("player not found in result. But saved tipps in results/error_%d.json\n", aktuellesGame.matchId)

	} else {

		// check if remake
		if killd.Info.Participants[ind].TeamEarlySurrendered || killd.Info.Participants[(ind+5)%10].TeamEarlySurrendered {
			log.Println("Auswertung abgebrochen, da geremaked wurde.")
			config.TwitchClient.Say(config.Twitchchannel, config.Prefix+" Killspiel wurde abgebrochen, da Remaked wurde.")

		} else {
			kills := killd.Info.Participants[ind].Kills

			// important for the twitch message
			gewinner := bessereDaten[kills]

			res := result{
				MatchId:         aktuellesGame.matchId,
				PlayerId:        aktuellesGame.playerId,
				PlayerChampId:   killd.Info.Participants[ind].ChampionId,
				PlayerChampName: killd.Info.Participants[ind].ChampionName,
				Kills:           kills,
				Tipps:           bessereDaten,
			}

			// write result into file
			file, err := os.Create(fmt.Sprintf("results/%d.json", aktuellesGame.matchId))
			if err != nil {
				log.Fatal(err)
			}
			bites, err := json.MarshalIndent(res, "", "  ")
			if err != nil {
				log.Printf("Error occurded writing res: %v\n", err)
			}
			_, err = file.Write(bites)
			if err != nil {
				return
			}
			log.Printf("Ergebnis gespeichert in results/%d.json\n", aktuellesGame.matchId)

			// write twitch message
			var haben string
			switch len(gewinner) {
			case 0:
				haben = "Keiner hat"
			case 1:
				haben = "hat"
			default:
				haben = "haben"

			}
			config.TwitchClient.Say(config.Twitchchannel,
				fmt.Sprintf("%s Killspiel wurde beendet. %s hat %d Kill(s) gemacht. %s %s die richtige Killanzahl getippt.",
					config.Prefix, config.Lolaccountname, kills, strings.Join(gewinner, ", "), haben))
		}
	}

	// remove tmp-file
	err = os.Remove(fmt.Sprintf("results/tmp_%d.json", aktuellesGame.matchId))
	if err != nil {
		log.Printf("An error occured while deleting results/tmp_%d.json: %v\n", aktuellesGame.matchId, err)
	}

	// reset
	bessereDaten = map[int][]string{}
	config.State = Idle
	aktuellesGame.matchId = 0

	// run shell-command, e.g. to analyse the result and update a database
	if config.CmdAfterAuswertung != "" {
		exec.Command(config.CmdAfterAuswertung)
	}
}
