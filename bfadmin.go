package main

import (
	"bfadmin/configuration"
	"bfadmin/controller"
	"bfadmin/util"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Game struct {
	gameId string
	poller *controller.Poller
	runner *controller.Runner
	settingsDir string
}

var games = make(map[string] * Game)

func control(w http.ResponseWriter, req * http.Request) {
	all, err := ioutil.ReadAll(req.Body)
	if util.CheckErr(err) {
		w.Write([]byte("Network Error."))
	}
	var action map[string]interface{}
	json.Unmarshal(all, &action)

	gameId := getGameId(req)
	game := games[gameId]

	if game != nil {
		if action["action"] == "start" {
			game.runner.Start()
			w.Write([]byte("Startup Initiated"))
		} else if action["action"] == "stop"  {
			game.runner.Stop()
			w.Write([]byte("Shutdown Initiated"))
		} else if action["action"] == "restore" {
			configuration.RestoreConfigs(game.settingsDir)
			w.Write([]byte("Configuration Restored"))
		}
	} else {
		w.Write([]byte("Not Implemented (yet)"))
	}
}

func getStatus(w http.ResponseWriter, req * http.Request) {
	if req.Method != "GET" {
		makeHttpError(w)
		return
	}
	setJsonHeader(w)

	gameId := getGameId(req)
	game := games[gameId]

	if game != nil {
		w.Write(game.poller.GetStatusJson())
	} else {
		w.Write(controller.SERVER_OFFLINE)
	}
}

func processConfig(w http.ResponseWriter, req * http.Request) {
	gameId := getGameId(req)
	game := games[gameId]
	if game != nil {
		if req.Method == "POST" {
			setConfig(game, w, req)
		} else {
			getConfig(game, w, req)
		}
	}
}

func setConfig(game *Game, w http.ResponseWriter, req * http.Request) {
	all, err := ioutil.ReadAll(req.Body)
	if util.CheckErr(err) {
		w.Write([]byte("Network Error"))
	}
	config := configuration.ServerConfig{}
	err = json.Unmarshal(all, &config)
	if util.CheckErr(err) {
		w.Write([]byte("JSON Error"))
		return
	}
    configuration.WriteConfig(game.gameId, game.settingsDir, & config)

	w.Write([]byte("Configuration Saved"))
	return
}

func getConfig(game *Game, w http.ResponseWriter, req * http.Request) {
	setJsonHeader(w)
    cfg := configuration.ReadConfig(game.gameId, game.settingsDir)
	json, _ := json.Marshal(cfg)
	w.Write(json)
}

func setJsonHeader(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
}

func makeHttpError(w http.ResponseWriter) {
	w.WriteHeader(500)
	w.Write([]byte("Invalid Request"))
}

func getGameId(req *http.Request) string {
	gameId := req.URL.Query().Get("gameId")
	game := games[gameId]
	if game != nil {
		return game.gameId
	}
	return ""
}

func main() {
	conf := util.ReadPropertiesFile("bfadmin.conf")

	maxRunTimeMin := util.Atoi(conf["maxRunTimeMin"])

	// BF1942
	if conf["bf1942Exec"] != "" {
		bf1942Runner := controller.NewRunner("bf1942", conf["bf1942Exec"], conf["bf1942SettingsDir"],
			nil, maxRunTimeMin)
		bf1942GameSpy := controller.GameSpy{
			HostPort: conf["bf1942GameSpy"],
			Parser:   controller.GameSpyBF1942Parser,
		}
		bf1942Poller := controller.NewPoller(&bf1942GameSpy, bf1942Runner)
		bf1942Poller.StartPolling()

		games["bf1942"] = & Game{
			"bf1942",
			bf1942Poller,
			bf1942Runner,
			conf["bf1942SettingsDir"],
		}
	}

	// BFVietnam
	if conf["bfvExec"] != "" {
		bfvRunner := controller.NewRunner("bfv", conf["bfvExec"], conf["bfvSettingsDir"],
			nil, maxRunTimeMin)
		bfvGameSpy := controller.GameSpy{
			HostPort: conf["bfvGameSpy"],
			Parser:   controller.GameSpyBFVParser,
		}
		bfvPoller := controller.NewPoller(&bfvGameSpy, bfvRunner)
		bfvPoller.StartPolling()

		games["bfv"] = &Game{
			"bfv",
			bfvPoller,
			bfvRunner,
			conf["bfvSettingsDir"],
		}
	}

	// BF2
	if conf["bf2Exec"] != "" {
		bf2Runner := controller.NewRunner("bf2", conf["bf2Exec"], conf["bf2SettingsDir"],
			util.ParseCommandLine(conf["bf2ExecArgs"]), maxRunTimeMin)
		bf2GameSpy := controller.GameSpyV3{
			HostPort: conf["bf2GameSpy"],
			Parser:   controller.GameSpyBF2Parser,
		}
		bf2Poller := controller.NewPoller(&bf2GameSpy, bf2Runner)
		bf2Poller.StartPolling()

		games["bf2"] = &Game{
			"bf2",
			bf2Poller,
			bf2Runner,
			conf["bf2SettingsDir"],
		}
	}

	// BF2142
	if conf["bf2142Exec"] != "" {
		bf2142Runner := controller.NewRunner("bf2142", conf["bf2142Exec"], conf["bf2142SettingsDir"],
			util.ParseCommandLine(conf["bf2142ExecArgs"]), maxRunTimeMin)
		bf2142GameSpy := controller.GameSpyV3{
			HostPort: conf["bf2142GameSpy"],
			Parser:   controller.GameSpyBF2Parser,
		}
		bf2142Poller := controller.NewPoller(&bf2142GameSpy, bf2142Runner)
		bf2142Poller.StartPolling()

		games["bf2142"] = &Game{
			"bf2142",
			bf2142Poller,
			bf2142Runner,
			conf["bf2142SettingsDir"],
		}
	}


	// BF2 - Project Reality
	if conf["prBf2Exec"] != "" {
		prBf2Runner := controller.NewRunner("prbf2", conf["prBf2Exec"], conf["prBf2SettingsDir"],
			util.ParseCommandLine(conf["prBf2ExecArgs"]), maxRunTimeMin)
		prBf2GameSpy := controller.GameSpyV3{
			HostPort: conf["prBf2GameSpy"],
			Parser:   controller.GameSpyBF2Parser,
		}
		prBf2Poller := controller.NewPoller(&prBf2GameSpy, prBf2Runner)
		prBf2Poller.StartPolling()

		games["prbf2"] = &Game{
			"prbf2",
			prBf2Poller,
			prBf2Runner,
			conf["prBf2SettingsDir"],
		}
	}

	// UT2004
	if conf["ut2004Exec"] != "" {
		ut2004Runner := controller.NewRunner("ut2004", conf["ut2004Exec"], conf["ut2004SettingsDir"],
			util.ParseCommandLine(conf["ut2004ExecArgs"]), maxRunTimeMin)
		ut2004GameSpy := controller.GameSpy{
			HostPort: conf["ut2004GameSpy"],
			Parser:   controller.GameSpyUT2004Parser,
		}
		ut2004Poller := controller.NewPoller(&ut2004GameSpy, ut2004Runner)
		ut2004Poller.StartPolling()

		games["ut2004"] = &Game{
			"ut2004",
			ut2004Poller,
			ut2004Runner,
			conf["ut2004SettingsDir"],
		}
	}

	http.Handle("/", http.FileServer(http.Dir("html")))
	http.HandleFunc("/status", getStatus)
	http.HandleFunc("/config", processConfig)
	http.HandleFunc("/control", control)

	log.Printf("Listening on %s", conf["listenAddress"])
	log.Fatal(http.ListenAndServe(conf["listenAddress"], nil))
}
