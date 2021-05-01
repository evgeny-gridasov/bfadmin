package controller

const SERVER_OFFLINE_STR = `{"Status" : "OFFLINE", "Mapname" : "N/A", "Modname" : "N/A", "Gametype" : "N/A",
"Numplayers" : "-", "Maxplayers" : "-", "Tickets1" : "-", "Tickets2" : "-", "Players" : []}`
var SERVER_OFFLINE = []byte (SERVER_OFFLINE_STR)

var GAMESPYV3_REQUEST = []byte{0xFE, 0xFD, 0x00, 0x50, 0x6f, 0x4e, 0x47, 0xff, 0xff, 0xff, 0x01}
var GAMESPYV3_RESPONSE = []byte{0x00, 0x50, 0x6f, 0x4e, 0x47, 0x73, 0x70, 0x6c, 0x69, 0x74, 0x6e, 0x75, 0x6d, 0x00}

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

type StatusProtocol interface {
	GetStatus() *ServerStatus
}
