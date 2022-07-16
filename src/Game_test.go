package Killspiel

import (
	"github.com/gempir/go-twitch-irc/v2"
	"testing"
)

func TestAuswertung(t *testing.T) {
	config = getConfig("config.json")

	//fix
	client := twitch.NewClient(config.Username, config.Oath)
	config.TwitchClient = client
	// doesn't write anything to twitch, hasn't joined any channel

	aktuellesGame = &game{matchId: 5967718649, playerId: "w-MRk8wehYuoOVertLqVDtHE-7EcMXkDIJ0xxycyvRg5dPU"}
	//													 "yFl1JcuA3BI5kWVh3qLjayIDvn70QNChfzMNP9RC7zfVSs0ltXytPeKIZbzQotj-6CKmP2sKGfHoSA"
	bessereDaten = map[int][]string{0: {"l00", "l01", "l02"}, 1: {"l10", "l11", "l12"}, 2: {"l20"}, 3: {"l30", "l31", "l32", "l33", "l34", "l35"}, 4: {"w0", "w1", "w2"}}
	Auswertung()
}
