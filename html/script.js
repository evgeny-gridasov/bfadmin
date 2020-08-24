var timeout = 5000;

function id(id) {
    return document.getElementById(id);
}
function add() {
    let selected = id("selected");
    let options = id("maps").options;
    for (let i = 0; i < options.length; i++) {
        let item = options.item(i);
        if (item.selected) {
            i--;
            item.selected = false;
            selected.appendChild(item);
        }
    }
    return false;
}

function remove() {
    let maps = id("maps");
    let options = id("selected").options;
    for (let i = 0; i < options.length; i++) {
        let item = options.item(i);
        if (item.selected) {
            i--;
            item.selected = false;
            maps.appendChild(item);
        }
    }
    return false;
}

function save() {
    let options = id("selected").options;
    if (options.length == 0) {
        showAlert("Please add at least one map.");
        return false;
    }
    let selectedMaps = [];
    for (let i = 0; i< options.length; i++) {
        selectedMaps.push({
            Id: options[i].value,
        });
    }

    send(JSON.stringify({
        MaxPlayers : parseInt(id("maxPlayers").value),
        GameTime : parseInt(id("roundTimeLimit").value),
        NumberOfRounds : parseInt(id("nrOfRounds").value),
        SpawnTime : parseInt(id("spawnTime").value),
        GameStartDelay : parseInt(id("gameStartDelay").value),
        TicketRatio : parseInt(id("tickets").value),
        FriendlyFire : parseInt(id("friendlyFire").value),
        AlliedTeamRatio : parseInt(id("teamRatioBlue").value),
        AxisTeamRatio : parseInt(id("teamRatioRed").value),
        CoopAiSkill : parseInt(id("aiSkills").value),
        SelectedMaps : selectedMaps

    }), "config?gameId=" + gameId);
    return false;
}

function startServer() {
    send(JSON.stringify({action : "start"}), "control?gameId=" + gameId, loadConfig);
    timeout = 1000;
    window.setTimeout(loadStatus, timeout);
    return false;
}

function stopServer() {
    send(JSON.stringify({action : "stop"}), "control?gameId=" + gameId);
    timeout = 1000;
    window.setTimeout(loadStatus, timeout);
    return false;
}

function restore() {
    send(JSON.stringify({action : "restore"}), "control?gameId=" + gameId, loadConfig);
    return false;
}

function send(request, url, func) {
    let ajax = new XMLHttpRequest();
    ajax.open("POST", url);
    ajax.responseType = "text";
    ajax.setRequestHeader('Content-type', 'application/json; charset=utf-8');
    ajax.onload = function () {
        showAlert(ajax.response.substr(0, 50));
        if (func != null) {
            func();
        }
    }
    ajax.send(request);
}

function showAlert(message) {
    id("serverMessage").innerText = message;
    id("coolAlert").style.display = "block";
}

function closeAlert() {
    id("coolAlert").style.display = "none";
}

function loadConfig() {
    let ajax = new XMLHttpRequest();
    ajax.open("GET", "config?gameId=" + gameId);
    ajax.responseType = "json";
    ajax.send();
    ajax.onload = function () {
        id("maxPlayers").value = ajax.response["MaxPlayers"];
        id("roundTimeLimit").value = ajax.response["GameTime"];
        id("nrOfRounds").value = ajax.response["NumberOfRounds"];
        id("spawnTime").value = ajax.response["SpawnTime"];
        id("gameStartDelay").value = ajax.response["GameStartDelay"];
        id("tickets").value = ajax.response["TicketRatio"];
        id("friendlyFire").value = ajax.response["FriendlyFire"];
        id("teamRatioBlue").value = ajax.response["AlliedTeamRatio"];
        id("teamRatioRed").value = ajax.response["AxisTeamRatio"];
        id("aiSkills").value = ajax.response["CoopAiSkill"];

        let availableList = id("maps");
        availableList.options.length = 0;
        let availableMaps = ajax.response["AvailableMaps"];
        for (let i=0; i<availableMaps.length; i++) {
            let m = createOption(availableMaps[i]["Id"], availableMaps[i]["Name"]);
            availableList.add(m);
        }

        let selectedList = id("selected");
        selectedList.options.length = 0;
        let selectedMaps = ajax.response["SelectedMaps"];
        for (let i=0; i<selectedMaps.length; i++) {
            let m = createOption(selectedMaps[i]["Id"], selectedMaps[i]["Name"]);
            selectedList.add(m);
        }
    }
}

function createOption(value, text) {
    let opt = document.createElement("option");
    opt.value = value;
    opt.text = text;
    return opt;
}

function loadStatus() {
    let ajax = new XMLHttpRequest();
    ajax.open("GET", "status?gameId=" + gameId);
    ajax.responseType = "json";
    ajax.onload = function () {
        let status = ajax.response["Status"];
        id("status").innerText = status;
        if (status == "ONLINE") {
            id("status").style.color = "gold";
            id("startButton").disabled = true;
            id("stopButton").disabled = false;
        } else if (status == "STARTING" || status == "STOPPING" ) {
            id("status").style.color = "green";
            id("startButton").disabled = true;
            id("stopButton").disabled = false;
        } else {
            id("status").style.color = "black";
            id("startButton").disabled = false;
            id("stopButton").disabled = true;
        }
        id("mod").innerText = ajax.response["Modname"];
        id("map").innerText = ajax.response["Mapname"];
        id("gameType").innerText = ajax.response["Gametype"];
        id("players").innerText = ajax.response["Numplayers"] + "/" + ajax.response["Maxplayers"];
        id("tickets1").innerText = ajax.response["Tickets1"];
        id("tickets2").innerText = ajax.response["Tickets2"];

        let playerList = id("playerList");
        while(playerList.rows.length > 1) {
            playerList.deleteRow(1);
        }

        let players = ajax.response["Players"];
        for (let i=0; i<players.length; i++) {
            let player = players[i];
            let tr = playerList.insertRow(i+1);
            let playerName = tr.insertCell(0);
            playerName.innerText = player["Name"];
            if (player["Team"] == "1") {
                playerName.className = "red";
            }
            if (player["Team"] == "2") {
                playerName.className = "blue";
            }
            tr.insertCell(1).innerText = player["Score"];
            tr.insertCell(2).innerText = player["Kills"];
            tr.insertCell(3).innerText = player["Deaths"];
            tr.insertCell(4).innerText = player["Ping"];
        }

        window.setTimeout(loadStatus, timeout);
        if (timeout < 5000) {
            timeout = timeout + 1000;
        }
    }
    ajax.onerror = function() {
        window.setTimeout(loadStatus, 10000);
    }
    ajax.send();
}

window.onload = function () {
    loadConfig();
    loadStatus();
}