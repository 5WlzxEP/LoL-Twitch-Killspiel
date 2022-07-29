package Killspiel

import "github.com/gempir/go-twitch-irc/v2"

type GameState uint16

const (
	Idle GameState = iota
	Wettphase
	Spielphase
	Auswertungsphase
	GameNoTrack
)

type champion struct {
	Type string `json:"type"`
	Data map[string]struct {
		Id   string `json:"id"`
		Key  string `json:"key"`
		Name string `json:"name"`
	} `json:"data"`
}

type game struct {
	matchId  int64
	playerId string
}

type GlobalConfig struct {
	Username           string `json:"Username"`
	Oath               string `json:"Oath"`
	Wettdauer          int    `json:"Wettdauer"`
	Twitchchannel      string `json:"Twitchchannel"`
	Lolaccountname     string `json:"Lolaccountname"`
	Lolapikey          string `json:"Lolapikey"`
	Joinmessage        bool   `json:"Joinmessage"`
	Logpath            string `json:"LogPath"`
	Prefix             string `json:"TwitchPrefix"`
	otp                bool
	Champs             *[]string `json:"Champions"`
	champsId           *[]int
	State              GameState
	TwitchClient       *twitch.Client
	CmdAfterAuswertung string `json:"CmdAfterAuswertung"`
}

type result struct {
	MatchId         int64            `json:"matchId"`
	PlayerId        string           `json:"playerId"`
	PlayerChampId   int              `json:"playerChampId"`
	PlayerChampName string           `json:"playerChampName"`
	Kills           int              `json:"kills"`
	Tipps           map[int][]string `json:"Tipps"`
}

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
			Assists              int    `json:"assists"`
			ChampionId           int    `json:"championId"`
			ChampionName         string `json:"championName"`
			Deaths               int    `json:"deaths"`
			Kills                int    `json:"kills"`
			ParticipantId        int    `json:"participantId"`
			TeamEarlySurrendered bool   `json:"teamEarlySurrendered"`
			Win                  bool   `json:"win"`
		} `json:"participants"`
	} `json:"info"`
}
