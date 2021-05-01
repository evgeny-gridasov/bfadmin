package controller

import (
	"bfadmin/util"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"
)

type GameSpy struct {
	HostPort string
	Parser func (map[string] string) *ServerStatus
}

func (gs *GameSpy) GetStatus() *ServerStatus {
	statusMap := gs.readStatus()
	if gs.Parser != nil {
		return gs.Parser(statusMap)
	}
	return nil
}

func (gs *GameSpy) readStatus() map[string] string {
	ret := make(map[string]string)
	conn, err := net.Dial("udp", gs.HostPort)
	defer conn.Close()
	util.CheckErr(err)

	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	util.CheckErr(err)

	_, err = conn.Write([]byte("\\status\\"))
	if util.CheckErr(err) {
		return ret
	}

	buff :=  make([]byte, 16384)

	// packets may be reordered, read all up to maxQueryId
	maxQueryId := 65536

	for i:=1; i < maxQueryId; i++ {
		read, err := conn.Read(buff)
		if util.CheckErr(err) {
			break
		}

		reply := string(buff[:read])
		queryId, final := convertToMap(ret, reply)
		if final { // final queryId, update maxQueryId
			maxQueryId = queryId
		}
	}
	return ret
}

// Converts string like
// \gamename\bfield1942\gamever\v1.61\active_mods\,\averageFPS\0\queryid\7.1
// and adds value to a map
func convertToMap(m map[string]string, reply string) (int, bool) {
	split := strings.Split(reply, "\\")
	queryId := 0
	final := false
	for i:=1; i < len(split); i = i + 2 {
		key := split[i]
		if i >= len(split) {
			break
		}
		value := split[i+1]
		switch {
		case key == "queryid" :
			q := strings.Split(value, ".")
			if len(q) == 2 {
				queryId = util.Atoi(q[1])
			}
		case key == "final" :
			final = true
		default:
			//fmt.Println(key, "=>", value)
			m[key] = value
		}
	}
	return queryId, final
}


func GameSpyBF1942Parser(m map[string]string) *ServerStatus {
	if len(m) == 0 {
		return nil
	}

	numplayers := util.Atoi(m["numplayers"])
	players := make([]Player, numplayers)
	for i:=0; i< numplayers; i++ {
		istr := strconv.Itoa(i)
		score := util.Atoi(m["score_"+istr])
		players[i] = Player{
			m["playername_"+istr],
			m["team_"+istr],
			score,
			m["kills_"+istr],
			m["deaths_"+istr],
			m["ping_"+istr],
		}
	}
	sort.Slice(players, func(i int, j int) bool {return players[i].Score > players[j].Score})
	response := ServerStatus{
		"ONLINE",
		strings.ToUpper(m["mapname"]),
		strings.ToUpper(m["mapId"]),
		strings.ToUpper(m["gametype"]),
		m["numplayers"],
		m["maxplayers"],
		m["tickets1"],
		m["tickets2"],
		players,
	}

	return &response
}

func GameSpyBFVParser(m map[string]string) *ServerStatus {
	if len(m) == 0 {
		return nil
	}

	numplayers := util.Atoi(m["numplayers"])
	players := make([]Player, numplayers)
	for i:=0; i< numplayers; i++ {
		istr := strconv.Itoa(i)
		score := util.Atoi(m["score_"+istr])
		players[i] = Player {
			m["player_" + istr],
			m["team_" + istr],
			score,
			m["kills_" + istr],
			m["deaths_" + istr],
			m["ping_" + istr],
		}
	}
	sort.Slice(players, func(i int, j int) bool {return players[i].Score > players[j].Score})
	response := ServerStatus{
		"ONLINE",
		strings.ToUpper(m["mapname"]),
		strings.ToUpper(m["map_id"]),
		strings.ToUpper(m["gametype"]),
		m["numplayers"],
		m["maxplayers"],
		m["tickets_t0"],
		m["tickets_t1"],
		players,
	}

	return &response
}

func GameSpyUT2004Parser(m map[string]string) *ServerStatus {
	if len(m) == 0 {
		return nil
	}

	numplayers := util.Atoi(m["numplayers"])
	var players []Player
	for i:=0; i< numplayers; i++ {
		istr := strconv.Itoa(i)
		name := m["player_"+istr]
		if name == "" {
			break
		}
		score := util.Atoi(m["frags_"+istr])
		team := m["team_"+istr]
		if team == "0" {
			team = "1"
		} else {
			team = "2"
		}
		players = append(players, Player {
			name,
			team,
			score,
			"-",
			"-",
			m["ping_" + istr],
		},
		)
	}
	sort.Slice(players, func(i int, j int) bool {return players[i].Score > players[j].Score})
	response := ServerStatus{
		"ONLINE",
		strings.ToUpper(m["mapname"]),
		"N/A",
		strings.ToUpper(m["gametype"]),
		m["numplayers"],
		m["maxplayers"],
		"-",
		"-",
		players,
	}

	return &response
}


func (gs*GameSpy) String() string {
	return gs.HostPort
}