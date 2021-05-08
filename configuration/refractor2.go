package configuration

import (
	"bfadmin/util"
	"bufio"
	"strconv"
	"strings"
)

// Refractor 2 Engine config: BF2, BF2142
//
func readRefractor2Config(reader *bufio.Reader, config *ServerConfig) {
	for {
		readString, err := reader.ReadString('\n')
		fields := strings.Fields(readString)
		if len(fields) == 2 {
			switch fields[0] {
			case "sv.timeLimit":
				config.GameTime = util.Atoi(fields[1])
			case "sv.maxPlayers":
				config.MaxPlayers = util.Atoi(fields[1])
			case "sv.roundsPerMap":
				config.NumberOfRounds = util.Atoi(fields[1])
			case "sv.spawnTime":
				config.SpawnTime = util.Atoi(fields[1])
			case "sv.startDelay":
				config.GameStartDelay = util.Atoi(fields[1])
			case "sv.ticketRatio":
				config.TicketRatio = util.Atoi(fields[1])
			case "sv.soldierFriendlyFire":
				config.FriendlyFire = util.Atoi(fields[1])
			case "sv.coopBotRatio":
				config.AlliedTeamRatio = util.Atoi(fields[1])
			case "sv.coopBotCount":
				config.AxisTeamRatio = util.Atoi(fields[1])
			case "sv.coopBotDifficulty":
				config.CoopAiSkill = util.Atoi(fields[1])
			case "sv.botSkill":
				config.CoopAiSkill = int(100 * util.Atof(fields[1]))
			}
		}
		if err != nil {
			break
		}
	}
}

func getSelectedRefractor2Maps(reader *bufio.Reader, selectedMapsSet map[string]bool) []GameMap {
	selectedMaps := make([]GameMap, 0)
	for {
		readString, err := reader.ReadString('\n')
		fields := strings.Fields(readString)

		if len(fields) == 4 {
			if fields[0] == "mapList.append" {
				selectedMaps = append(selectedMaps, GameMap{
					util.MakeId(fields[1], fields[2], fields[3]),
					util.MakeNameRefractor2(fields[1], fields[2], fields[3]),
				})
				selectedMapsSet[fields[1]] = true
			}
		}
		if err != nil {
			break
		}
	}
	return selectedMaps
}

func getAllRefractor2Maps(reader *bufio.Reader, selectedMapsSet map[string]bool) []GameMap {
	allMaps := make([]GameMap, 0)
	for {
		readString, err := reader.ReadString('\n')
		fields := strings.Fields(readString)

		if len(fields) == 4 {
			if fields[0] == "mapList.append" && !selectedMapsSet[fields[1]] {
				allMaps = append(allMaps, GameMap{
					util.MakeId(fields[1], fields[2], fields[3]),
					util.MakeNameRefractor2(fields[1], fields[2], fields[3]),
				})
			}
		}
		if err != nil {
			break
		}
	}
	return allMaps
}

func renderRefractor2Config(reader *bufio.Reader, config *ServerConfig) string {
	sb := strings.Builder{}
	for {
		readString, err := reader.ReadString('\n')
		fields := strings.Fields(readString)
		if len(fields) == 2 {
			sb.WriteString(fields[0] + " ")
			switch fields[0] {
			case "sv.timeLimit":
				sb.WriteString(strconv.Itoa(config.GameTime))
			case "sv.maxPlayers":
				sb.WriteString(strconv.Itoa(config.MaxPlayers))
			case "sv.roundsPerMap":
				sb.WriteString(strconv.Itoa(config.NumberOfRounds))
			case "sv.spawnTime":
				sb.WriteString(strconv.Itoa(config.SpawnTime))
			case "sv.startDelay":
				sb.WriteString(strconv.Itoa(config.GameStartDelay))
			case "sv.ticketRatio":
				sb.WriteString(strconv.Itoa(config.TicketRatio))
			case "sv.soldierFriendlyFire":
				sb.WriteString(strconv.Itoa(config.FriendlyFire))
			case "sv.coopBotRatio":
				sb.WriteString(strconv.Itoa(config.AlliedTeamRatio))
			case "sv.coopBotCount":
				sb.WriteString(strconv.Itoa(config.AxisTeamRatio))
			case "sv.coopBotDifficulty": // bf2
				sb.WriteString(strconv.Itoa(config.CoopAiSkill))
			case "sv.botSkill": // bf2142
				sb.WriteString(strconv.FormatFloat(float64(config.CoopAiSkill) / 100.0, 'f', 2, 64))
			default:
				sb.WriteString(fields[1])
			}
			sb.WriteString("\n")
		} else {
			sb.WriteString(readString)
		}
		if err != nil {
			break
		}
	}
	return sb.String()
}

func renderRefractor2MapsList(config *ServerConfig) string {
	sb := strings.Builder{}
	for i := 0; i < len(config.SelectedMaps); i++ {
		split := strings.Split(config.SelectedMaps[i].Id, ":")
		if len(split) == 3 {
			sb.WriteString("mapList.append " + split[0] + " " + split[1] + " " + split[2] + "\n")
		}
	}
	return sb.String()
}