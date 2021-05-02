package configuration

import (
	"bfadmin/util"
	"bufio"
	"strconv"
	"strings"
)

//
// Refractor 1 Engine config: BFV, BF1942
//
func readRefractor1Config(reader *bufio.Reader, config *ServerConfig) {
	for {
		readString, err := reader.ReadString('\n')
		fields := strings.Fields(readString)
		if len(fields) == 2 {
			switch fields[0] {
			case "game.serverGameTime":
				config.GameTime = util.Atoi(fields[1])
			case "game.serverMaxPlayers":
				config.MaxPlayers = util.Atoi(fields[1])
			case "game.serverNumberOfRounds":
				config.NumberOfRounds = util.Atoi(fields[1])
			case "game.serverSpawnTime":
				config.SpawnTime = util.Atoi(fields[1])
			case "game.serverGameStartDelay":
				config.GameStartDelay = util.Atoi(fields[1])
			case "game.serverTicketRatio":
				config.TicketRatio = util.Atoi(fields[1])
			case "game.serverSoldierFriendlyFire":
				config.FriendlyFire = util.Atoi(fields[1])
			case "game.serverAlliedTeamRatio":
				config.AlliedTeamRatio = util.Atoi(fields[1])
			case "game.serverAxisTeamRatio":
				config.AxisTeamRatio = util.Atoi(fields[1])
			case "game.serverCoopAiSkill":
				config.CoopAiSkill = util.Atoi(fields[1])
			}
		}
		if err != nil {
			break
		}
	}
}

func getSelectedRefractor1Maps(reader *bufio.Reader, selectedMapsSet map[string]bool) []GameMap {
	selectedMaps := make([]GameMap, 0)
	for {
		readString, err := reader.ReadString('\n')
		fields := strings.Fields(readString)

		if len(fields) == 4 {
			if fields[0] == "game.addLevel" {
				selectedMaps = append(selectedMaps, GameMap{
					util.MakeId(fields[1], fields[2], fields[3]),
					util.MakeName(fields[1], fields[3]),
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

func getAllRefractor1Maps(reader *bufio.Reader, selectedMapsSet map[string]bool) []GameMap {
	allMaps := make([]GameMap, 0)
	for {
		readString, err := reader.ReadString('\n')
		fields := strings.Fields(readString)

		if len(fields) == 4 {
			if fields[0] == "game.addLevel" && !selectedMapsSet[fields[1]] {
				allMaps = append(allMaps, GameMap{
					util.MakeId(fields[1], fields[2], fields[3]),
					util.MakeName(fields[1], fields[3]),
				})
			}
		}
		if err != nil {
			break
		}
	}
	return allMaps
}

func renderRefractor1Config(reader *bufio.Reader, config *ServerConfig) string {
	sb := strings.Builder{}
	for {
		readString, err := reader.ReadString('\n')
		fields := strings.Fields(readString)
		if len(fields) == 2 {
			sb.WriteString(fields[0] + " ")
			switch fields[0] {
			case "game.serverGameTime":
				sb.WriteString(strconv.Itoa(config.GameTime))
			case "game.serverMaxPlayers":
				sb.WriteString(strconv.Itoa(config.MaxPlayers))
			case "game.serverNumberOfRounds":
				sb.WriteString(strconv.Itoa(config.NumberOfRounds))
			case "game.serverSpawnTime":
				sb.WriteString(strconv.Itoa(config.SpawnTime))
			case "game.serverGameStartDelay":
				sb.WriteString(strconv.Itoa(config.GameStartDelay))
			case "game.serverTicketRatio":
				sb.WriteString(strconv.Itoa(config.TicketRatio))
			case "game.serverSoldierFriendlyFire":
				sb.WriteString(strconv.Itoa(config.FriendlyFire))
			case "game.serverAlliedTeamRatio":
				sb.WriteString(strconv.Itoa(config.AlliedTeamRatio))
			case "game.serverAxisTeamRatio":
				sb.WriteString(strconv.Itoa(config.AxisTeamRatio))
			case "game.serverCoopAiSkill":
				sb.WriteString(strconv.Itoa(config.CoopAiSkill))
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

func renderRefractor1MapsList(config *ServerConfig) string {
	sb := strings.Builder{}
	for i := 0; i < len(config.SelectedMaps); i++ {
		split := strings.Split(config.SelectedMaps[i].Id, ":")
		if len(split) == 3 {
			sb.WriteString("game.addLevel " + split[0] + " " + split[1] + " " + split[2] + "\n")
		}
	}
	split := strings.Split(config.SelectedMaps[0].Id, ":")
	if len(split) == 3 {
		sb.WriteString("game.setCurrentLevel " + split[0] + " " + split[1] + " " + split[2] + "\n")
	}
	return sb.String()
}