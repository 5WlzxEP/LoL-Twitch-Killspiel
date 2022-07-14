package Killspiel

import (
	"testing"
)

func TestAuswertung(t *testing.T) {
	config = getConfig("config.json")
	aktuellesGame = &game{matchId: 5967718649, playerId: "yFl1JcuA3BI5kWVh3qLjayIDvn70QNChfzMNP9RC7zfVSs0ltXytPeKIZbzQotj-6CKmP2sKGfHoSA"}
	bessereDaten = map[int][]string{0: {"l00", "l01", "l02"}, 1: {"l10", "l11", "l12"}, 2: {"l20"}, 3: {"l30", "l31", "l32", "l33", "l34", "l35"}, 4: {"w0", "w1", "w2"}}
	Auswertung()
}
