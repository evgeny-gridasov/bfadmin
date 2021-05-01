package configuration

import (
	"bfadmin/util"
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
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
	switch gameId {
	case "bfv", "bf1942":
		readRefactor1Config(reader, &config)
	case "bf2":
		readRefactor2Config(reader, &config)
	}

	// maplist
	maplistCon, err := os.Open(filepath.Join(dir, MAPLIST_SAVED))
	defer maplistCon.Close()
	if util.CheckErr(err) {
		return ServerConfig{}
	}

	reader = bufio.NewReader(maplistCon)
	selectedMapsSet := make(map[string]bool)

	switch gameId {
	case "bfv", "bh1942":
		config.SelectedMaps = getSelectedRefactor1Maps(reader, selectedMapsSet)
	case "bf2":
		config.SelectedMaps = getSelectedRefactor2Maps(reader, selectedMapsSet)
	}

	// all maps
	allmapsCon, err := os.Open("levels/" + gameId)
	defer allmapsCon.Close()
	if util.CheckErr(err) {
		return ServerConfig{}
	}

	reader = bufio.NewReader(allmapsCon)
	switch gameId {
	case "bfv", "bf1942":
		config.AvailableMaps = getAllRefactor1Maps(reader, selectedMapsSet)
	case "bf2":
		config.AvailableMaps = getAllRefactor2Maps(reader, selectedMapsSet)
	}

	return config
}


func WriteConfig(gameId string, dir string, config *ServerConfig) {
	createSavedConfigs(dir)
	// serversettings

	serverSettings, err := os.Open(filepath.Join(dir, SERVERSETTINGS_SAVED))
	defer serverSettings.Close()
	if util.CheckErr(err) {
		return
	}

	reader := bufio.NewReader(serverSettings)
	var newConfig string
	switch gameId {
	case "bfv", "bf1942":
		newConfig = renderRefactor1Config(reader, config)
	case "bf2":
		newConfig = renderRefactor2Config(reader, config)
	default:
		return
	}

	serverSettings.Close()
	err = ioutil.WriteFile(filepath.Join(dir, SERVERSETTINGS_SAVED), []byte(newConfig), 0644)
	util.CheckErr(err)

	// maplist
	if len(config.SelectedMaps) == 0 {
		return
	}

	var newMapConfig string
	switch gameId {
	case "bfv", "bf1942":
		newMapConfig = renderRefactor1MapsList(config)
	case "bf2":
		newMapConfig = renderRefactor2MapsList(config)
	}

	err = ioutil.WriteFile(filepath.Join(dir, MAPLIST_SAVED), []byte(newMapConfig), 0644)
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
