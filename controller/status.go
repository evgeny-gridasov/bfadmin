package controller

const SERVER_OFFLINE_STR = `{"Status" : "OFFLINE", "Mapname" : "N/A", "Modname" : "N/A", "Gametype" : "N/A",
"Numplayers" : "-", "Maxplayers" : "-", "Tickets1" : "-", "Tickets2" : "-", "Players" : []}`
var SERVER_OFFLINE = []byte (SERVER_OFFLINE_STR)

// need capital letters variables to export json
type ServerStatus struct {
	Status string
	Mapname string
	Modname string
	Gametype string
	Numplayers string
	Maxplayers string
	Tickets1 string
	Tickets2 string
	Players []Player
}

type Player struct {
	Name string
	Team string
	Score int
	Kills string
	Deaths string
	Ping string
}

