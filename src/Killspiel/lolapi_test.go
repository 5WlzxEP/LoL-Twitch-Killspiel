package Killspiel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestGetKills(t *testing.T) {
	config = getConfig("config.json")
	aktuellesGame = &game{matchId: 5967718649, playerId: "w-MRk8wehYuoOVertLqVDtHE-7EcMXkDIJ0xxycyvRg5dPU"}
	expected := &killData{
		Metadata: struct {
			MatchId      string   `json:"matchId"`
			Participants []string `json:"participants"`
		}{MatchId: "EUW1_5967718649", Participants: []string{
			"yFl1JcuA3BI5kWVh3qLjayIDvn70QNChfzMNP9RC7zfVSs0ltXytPeKIZbzQotj-6CKmP2sKGfHoSA",
			"iNz7a5wXvEeZs6zCypR1JBrjP5DoxzOrQydR0xGLklIiweDdr0MsQLRj-RCrQs1vrT3f0-M348iD5w",
			"PP5QVYTxyXwETh1JmvNAsH1ilWqdGt7BEqP-6TwWg_ybmikizwpDUlMC6CGxYGNHqMMAuD2yu7WIcg",
			"Gpy5gbG-Pa-ZL0j5Hco2ux2BSBH8LuOYxV_tZatuAlml0j0NtubzJxQPiCqSij6jN88V3jgk3r6tkA",
			"tjimD1D-JIFwwqg-bmpTPi-81hiUzM8LqtPeCwJgzD5GbKYeKL1DrwD948Ke-ollTxwNDyiJWFAJjg",
			"HtkOWPumXvPqCBpg18dhJfcx-xyEEzQNN5GAeLilnyRhaZZmt5g1X8YSxTiwZn0PIYjL_4nBSwMkUA",
			"xq-5os9GdyIYSBrjD8ZN4W7TKOz3zZ8YM0-SFZ881MDN6DkQO6d5UOaifCpk-0ujv1Rh5ISSPhi7vA",
			"YYjiuLRzcPnJEgxYiA8n0utBwVFm9pQ8j-lE-miWWZVF9rfD_t7kg0duTWwcq20JMcJml9EhAOI4Zg",
			"8PsjIE-1eQk_6D-xR-JuQolFGXu_Dg5-qzmMT84YUshZHMDZjCAU3kO-xLLdd7fXvO1IluypbNW_NQ",
			"pYwqz_dojNsSgaOmNqqBerklP0rlueYNR_WRSTnz_ty7O_YailAFTvOmlmENvFkgNqxx7dt0ttQ2lQ"}},
		Info: struct {
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
		}{GameId: 5967718649, GameMode: "CLASSIC", GameType: "MATCHED_GAME", Participants: []struct {
			Assists              int  `json:"assists"`
			Deaths               int  `json:"deaths"`
			Kills                int  `json:"kills"`
			ParticipantId        int  `json:"participantId"`
			TeamEarlySurrendered bool `json:"teamEarlySurrendered"`
			Win                  bool `json:"win"`
		}{{Assists: 8, Deaths: 1, Kills: 4, ParticipantId: 1, TeamEarlySurrendered: false, Win: true},
			{Assists: 5, Deaths: 0, Kills: 14, ParticipantId: 2, TeamEarlySurrendered: false, Win: true},
			{Assists: 8, Deaths: 2, Kills: 5, ParticipantId: 3, TeamEarlySurrendered: false, Win: true},
			{Assists: 13, Deaths: 7, Kills: 5, ParticipantId: 4, TeamEarlySurrendered: false, Win: true},
			{Assists: 19, Deaths: 5, Kills: 2, ParticipantId: 5, TeamEarlySurrendered: false, Win: true},
			{Assists: 2, Deaths: 5, Kills: 2, ParticipantId: 6, TeamEarlySurrendered: false, Win: false},
			{Assists: 6, Deaths: 6, Kills: 2, ParticipantId: 7, TeamEarlySurrendered: false, Win: false},
			{Assists: 3, Deaths: 5, Kills: 4, ParticipantId: 8, TeamEarlySurrendered: false, Win: false},
			{Assists: 4, Deaths: 10, Kills: 4, ParticipantId: 9, TeamEarlySurrendered: false, Win: false},
			{Assists: 7, Deaths: 4, Kills: 3, ParticipantId: 10, TeamEarlySurrendered: false, Win: false}}},
	}
	kills := GetKills()

	if !reflect.DeepEqual(*kills, *expected) {
		t.Log("error, should be but got\n", expected, "\n", kills)
		t.Fail()
	}
	if fmt.Sprintf("%v", kills) != fmt.Sprintf("%v", expected) {
		t.Fail()
	}
}

func getConfig(file string) *Config {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("Error in Test, probaly no config file. err: %v\n", err)
		}
	}(f)
	bites, _ := ioutil.ReadAll(f)
	conf := &Config{}
	err = json.Unmarshal(bites, conf)
	if err != nil {
		log.Fatal(err)
	}
	return conf
}
