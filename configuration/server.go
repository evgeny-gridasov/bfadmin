package configuration

import (
	"bfadmin/util"
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const SERVERSETTINGS = "serversettings.con"
const SERVERSETTINGS_SAVED = "serversettings.con.saved"
const SERVERSETTINGS_ORIG = "serversettings.con.orig"
const MAPLIST = "maplist.con"
const MAPLIST_SAVED = "maplist.con.saved"
const MAPLIST_ORIG = "maplist.con.orig"

func ReadConfig(gameId string, dir string) ServerConfig {
	config := ServerConfig{}

	createSavedConfigs(dir)
	// serversettings
	serverSettings, err := os.Open(filepath.Join(dir, SERVERSETTINGS_SAVED))
	defer serverSettings.Close()
	if util.CheckErr(err) {
		return ServerConfig{}
	}

	reader := bufio.NewReader(serverSettings)
	for {
		readString, err := reader.ReadString('\n')
		fields := strings.Fields(readString)
		if len(fields) == 2 {
			switch fields[0] {
			case "game.serverGameTime":
				config.GameTime =util.Atoi(fields[1])
			case "game.serverMaxPlayers":
				config.MaxPlayers =util.Atoi(fields[1])
			case "game.serverNumberOfRounds":
				config.NumberOfRounds =util.Atoi(fields[1])
			case "game.serverSpawnTime":
				config.SpawnTime =util.Atoi(fields[1])
			case "game.serverGameStartDelay":
				config.GameStartDelay =util.Atoi(fields[1])
			case "game.serverTicketRatio":
				config.TicketRatio =util.Atoi(fields[1])
			case "game.serverSoldierFriendlyFire":
				config.FriendlyFire =util.Atoi(fields[1])
			case "game.serverAlliedTeamRatio":
				config.AlliedTeamRatio =util.Atoi(fields[1])
			case "game.serverAxisTeamRatio":
				config.AxisTeamRatio =util.Atoi(fields[1])
			case "game.serverCoopAiSkill":
				config.CoopAiSkill =util.Atoi(fields[1])
			}
		}
		if err != nil {
			break
		}
	}

	// maplist
	maplistCon, err := os.Open(filepath.Join(dir, MAPLIST_SAVED))
	defer maplistCon.Close()
	if util.CheckErr(err) {
		return ServerConfig{}
	}

	reader = bufio.NewReader(maplistCon)
	selectedMaps := make([]GameMap, 0)
	selectedMapsSet :=make(map[string]bool)
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
	config.SelectedMaps = selectedMaps

	// all maps
	allmapsCon, err := os.Open("levels/" + gameId)
	defer allmapsCon.Close()
	if util.CheckErr(err) {
		return ServerConfig{}
	}

	reader = bufio.NewReader(allmapsCon)
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
	config.AvailableMaps = allMaps

	return config
}

func WriteConfig(gameId string, dir string, config * ServerConfig) {
	createSavedConfigs(dir)
	// serversettings

	serverSettings, err := os.Open(filepath.Join(dir, SERVERSETTINGS_SAVED))
	defer serverSettings.Close()
	if util.CheckErr(err) {
		return
	}

	reader := bufio.NewReader(serverSettings)
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

	serverSettings.Close()
	err = ioutil.WriteFile(filepath.Join(dir, SERVERSETTINGS_SAVED), []byte(sb.String()), 0644)
	util.CheckErr(err)

	// maplist
	if len(config.SelectedMaps) == 0 {
		return
	}

	sb.Reset()
	for i:=0; i< len(config.SelectedMaps); i++ {
		split := strings.Split(config.SelectedMaps[i].Id, ":")
		if len(split) == 3 {
			sb.WriteString( "game.addLevel " + split[0] + " " + split[1] + " " + split[2] + "\n")
		}
	}
	split := strings.Split(config.SelectedMaps[0].Id, ":")
	if len(split) == 3 {
		sb.WriteString( "game.setCurrentLevel " + split[0] + " " + split[1] + " " + split[2] + "\n")
	} else {
		return
	}

	err = ioutil.WriteFile(filepath.Join(dir, MAPLIST_SAVED), []byte(sb.String()), 0644)
	util.CheckErr(err)
}

func createSavedConfigs(dir string) {
	backupOnce(filepath.Join(dir, SERVERSETTINGS), filepath.Join(dir, SERVERSETTINGS_SAVED))
	backupOnce(filepath.Join(dir, MAPLIST), filepath.Join(dir, MAPLIST_SAVED))
	backupOnce(filepath.Join(dir, SERVERSETTINGS), filepath.Join(dir, "serversettings.con.orig"))
	backupOnce(filepath.Join(dir, MAPLIST), filepath.Join(dir, "maplist.con.orig"))
}

func CopyConfigs(dir string) {
	copyFile(filepath.Join(dir, SERVERSETTINGS_SAVED), filepath.Join(dir, SERVERSETTINGS))
	copyFile(filepath.Join(dir, MAPLIST_SAVED), filepath.Join(dir, MAPLIST))
}

func RestoreConfigs(dir string) {
	copyFile(filepath.Join(dir, SERVERSETTINGS_ORIG), filepath.Join(dir, SERVERSETTINGS))
	copyFile(filepath.Join(dir, MAPLIST_ORIG), filepath.Join(dir, MAPLIST))
	copyFile(filepath.Join(dir, SERVERSETTINGS_ORIG), filepath.Join(dir, SERVERSETTINGS_SAVED))
	copyFile(filepath.Join(dir, MAPLIST_ORIG), filepath.Join(dir, MAPLIST_SAVED))
}

func backupOnce(file string, backupFile string) {
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		data, err := ioutil.ReadFile(file)
		util.CheckErr(err)
		err = ioutil.WriteFile(backupFile, data, 0644)
		util.CheckErr(err)
	}
}

func copyFile(src string, dst string) {
	if _, err := os.Stat(src); err == nil {
		data, err := ioutil.ReadFile(src)
		util.CheckErr(err)
		err = ioutil.WriteFile(dst, data, 0644)
		util.CheckErr(err)
	} else {
		util.CheckErr(err)
	}
}

