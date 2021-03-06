# LoL-Twitch-Killspiel

Automatisierung, um im Twitchchat das *Killspiel* spielen zu können.  
Das *Killspiel* besteht daraus, dass zu Beginn eines League of Legends Spieles der Twitchchat rät, wie viele Kills der Streamer erzielten wird.  
  
Dazu wird zum Beginn eines Spieles (bis zu 2 min später) die Wettphase automatisch gestartet.  
![Beispiel Beginn vom Spiel](img/Beginn.png)  
In der Zeit, die in der [config](#config) geändert werden, können die Chatteilnehmer mit ` !vote [Zahl] `, als z.B. `!vote 5` abstimmen, wenn sie denken, der Streamer erzielt 5 Kills.  
![Beispiel !vote](img/vote%20example.png)  
Spieler können innerhalb der Zeit beliebig oft anstimmen, nur der letzte Vote zählt. Nach der Zeit wird das Spiel automatisch beendet. Dies wird im Chat bekannt gegeben und die Anzahl an Teilnehmern bekannt gegeben.  
![](img/Ende%20Wettphase.png)  
Nach dem LoL Spiel wird die Anzahl an erzielten Kills automatisch aus der League-Api besorgt und die Spieler, die richtig getippt haben werden im Chat ausgegeben. Zudem wird eine Json-Datei mit allen Teilnehmern und ihren Tipps in einen *results*-Ordner gespeichert.  
![Beispiel Ende](img/Ende.png)  
![Ende2](img/Ende2.png)
![Ende meherere Gewinner](img/Ende%20mehrere%20Gewinner.png)

## config

```json
{
    "Username": "5w_lzxep", 
    "Oath": "oauth:bcgf6ogc6svu319nmeqprjgdtdizgw", 
    "Wettdauer": 120, 
    "Twitchchannel": "5w_lzxep", 
    "Lolaccountname": "5w_lzxep", 
    "Lolapikey": "RGAPI-cm5wf8rs-akq5-xrqh-is5p-4skbgcv1ekjg",
    "LoLRegion": "EUW1",
    "Joinmessage": true, 
    "LogPath": "./", 
    "TwitchPrefix": "/announce",
    "Champions": ["Nilah", "Bel'Veth", "Miss Fortune", "Shaco"],
    "CmdAfterAuswertung": "./Auswertung"
}
```

| Schlüsselwort       | Bedeutung                                                                                                                                                                                                                                                                                                                                | Required | Standardwert                 |
|---------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------|------------------------------|
| Username            | Nutzername für Twitch. Muss Moderator sein, wenn /announce o.ä. genutzt werden soll.                                                                                                                                                                                                                                                     | ✓        | /                            |
| OAuth               | Twitch Auth Token                                                                                                                                                                                                                                                                                                                        | ✓        | /                            |
| Wettdauer           | Sekunden, die der Chat Zeit hat abzustimmen                                                                                                                                                                                                                                                                                              | ×        | 120 Sekunden                 | 
| Twitchchannel       | Twitchchannel, auf dem die Nachrichten kommen und die Votes ausgelesen werden                                                                                                                                                                                                                                                            | ✓        | /                            |
| Lolaccountname      | Accountname der getrackt wird                                                                                                                                                                                                                                                                                                            | ✓        | /                            | 
| Lolapikey           | LoL api zugang                                                                                                                                                                                                                                                                                                                           | ✓        | /                            |
| LoLRegion           | Region des LoL-Accountes                                                                                                                                                                                                                                                                                                                 | ×        | `euw1`                       |
| Joinmessage         | Ob eine Nachricht geschickt werden soll, wenn der Bot verbunden ist.                                                                                                                                                                                                                                                                     | ×        | Nein/false                   | 
| LogPath             | Pfad der Log-Datei, ./ ist das aktuelle Verzeichnis                                                                                                                                                                                                                                                                                      | ×        | `./` / Aktuelles Verzeichnis |
| TwitchPrefix        | Prefix der Beginn, Ende und Auswertungsnachricht. Achtung kann Twitchkommands ausführen, wie z.B. "/announce", "/me", aber auch "/ban 5W_lzxEP".                                                                                                                                                                                         | ×        | ` ` / nichts                 |
| Champions           | Liste an Champions, bei denen der Bot triggrt. Falls leer wird jeder Champ akzeptiert.<br/> Die genaue Schreibweise der Champs kannst du [hier](Champs.md) nachgucken <br/> Es muss eine champions.json Datei im selben Verzeichnis liegen, sofern die LListe nicht leer ist. Dies kann über die `Champion_aktualisieren` gamacht werden | ×        | alle                         |
| CmdAfterAuswertung  | Shellcommand der nach der Erstellung von einem Result ausgeführt wird. Z.B. falls eine Aktualisierung einer Datenbank vorgenommen werden soll, kann man diesen Prozess hiermit auslösen                                                                                                                                                  | ×        | leer / macht nichts          |

### LoL-Rgionen

Folgende LoL-Regionen gibt es:

- BR1
- EUN1
- EUW1
- JP1
- KR
- LA1 / LAN
- LA2 / LAS
- NA1
- OC1
- RU
- TR1

## Beschaffung der API-Token

### OAuth

Entweder über [Twitch](https://dev.twitch.tv/docs/authentication/getting-tokens-oauth/) selber oder, ich nutze immer [diese Drittseite](https://twitchapps.com/tmi/).

### LoL Apikey

Über [Riot](https://developer.riotgames.com/)
- Developer Key. Dieser ist jedoch nur 24h nutzbar
- Register Product &rarr; *Personal API Key* sollte dauerhaften Key liefern. 

## Selber kompilieren

[go(lang)](https://go.dev/dl/) installieren

```bash
go build cmd/main.go
```

## Auswertung

Wie man die Daten auswerten kann, die dabei einstehen, kann man z.B. [hier](src/Auswertung) sehen.