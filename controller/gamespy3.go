package controller

import (
	"bfadmin/util"
	"bytes"
	"net"
	"sort"
	"strings"
	"time"
)

type GameSpyV3 struct {
	HostPort string
	Parser func (map[string] string) *ServerStatus
}

func (gs *GameSpyV3) GetStatus() *ServerStatus {
	statusMap := gs.readStatus()
	if gs.Parser != nil {
		return gs.Parser(statusMap)
	}
	return nil
}

func (gs *GameSpyV3) readStatus() map[string] string {
	packets := make(map[int][]byte)
	ret := make(map[string]string)

	conn, err := net.Dial("udp", gs.HostPort)
	defer conn.Close()
	util.CheckErr(err)

	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	util.CheckErr(err)

	_, err = conn.Write(GAMESPYV3_REQUEST)
	if util.CheckErr(err) {
		return ret
	}

	buff :=  make([]byte, 16384)

	// packets may be reordered, read all up to maxSplitNum
	maxSplitNum := 65536

	// collect packets
	for i:=0; i < maxSplitNum; i++ {
		read, err := conn.Read(buff)
		if util.CheckErr(err) || read < 16 {
			break
		}
		if bytes.Compare(buff[:14], GAMESPYV3_RESPONSE) != 0 {
			return ret
		}
		splitNum := int(buff[14])
		final := (splitNum & 0x80) > 0
		if final { // final splitNum, update maxSplitNum
			maxSplitNum = splitNum & 0x0f
		}
		buffCopy := make([]byte, read-16)
		copy(buffCopy, buff[16:read])
		packets[splitNum] = buffCopy
	}

	// reorder packets
	keys := make([]int, 0, len(packets))
	for k := range packets {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _,k := range keys {
		convertToMapV3(ret, packets[k])
	}

	return ret
}


func convertToMapV3(m map[string]string, reply []byte) {
	for pos:=0; pos < len(reply); {
		var key string
		var value string
		var subpos int
		key, subpos = getNextString(reply[pos:])
		pos += subpos
		if key != "" {
			// array format is: fieldname 0x00 0x00 data 0x00 data 0x00 data ... 0x00 data 0x00 0x00
			if strings.HasSuffix(key, "_") || strings.HasSuffix(key, "_t"){
				pos++
				strArray := m[key]
				for pos < len(reply) {
					value, subpos = getNextString(reply[pos:])
					pos += subpos
					if value == "" {
						break
					}
					if pos < len(reply) {
						if strArray != "" && value != "" {
							strArray = strArray + ","
						}
						strArray += value
					}
				}
				value = strArray
			} else {
				value, subpos = getNextString(reply[pos:])
				pos += subpos
			}
			//fmt.Println(key, "=>", value)
			m[key] = value
		}
	}
}

func getNextString(data []byte) (string, int) {
	wordEnd := -1
	for i:=0; i< len(data); i++ {
		if data[i] == 0x00 || data[i] == 0x01 || data[i] == 0x02 {
			wordEnd = i
			break
		}
	}
	if wordEnd >=0 {
		return string(data[:wordEnd]), wordEnd + 1
	}
	return string(data), len(data)
}

func GameSpyBF2Parser(m map[string]string) *ServerStatus {
	if len(m) == 0 {
		return nil
	}
	aiBot_ := strings.Split(m["AIBot_"], ",")
	player_ := strings.Split(m["player_"], ",")
	team_ := strings.Split(m["team_"], ",")
	score_ := strings.Split(m["score_"], ",")
	skill_ := strings.Split(m["skill_"], ",")
	deaths_ := strings.Split(m["deaths_"], ",")
	ping_ := strings.Split(m["ping_"], ",")
	players := make([]Player, 0)
	for i:=0; i< len(player_); i++ {
		if aiBot_[i] == "0"  {
			players = append(players, Player{
				player_[i],
				team_[i],
				util.Atoi(score_[i]),
				skill_[i],
				deaths_[i],
				ping_[i],
			})
		}
	}
	sort.Slice(players, func(i int, j int) bool {return players[i].Score > players[j].Score})
	response := ServerStatus{
		"ONLINE",
		strings.ToUpper(m["mapname"]),
		strings.ToUpper(m["gamevariant"]),
		strings.ToUpper(m["gametype"]),
		m["numplayers"],
		m["maxplayers"],
		"-",
		"-",
		players,
	}

	return &response
}

func (gs*GameSpyV3) String() string {
	return gs.HostPort
}