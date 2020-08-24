package configuration

type ServerConfig struct {
	MaxPlayers      int
	GameTime        int
	NumberOfRounds  int
	SpawnTime       int
	GameStartDelay  int
	TicketRatio     int
	FriendlyFire    int
	AlliedTeamRatio int
	AxisTeamRatio   int
	CoopAiSkill     int
	AvailableMaps   []GameMap
	SelectedMaps    []GameMap
}

type GameMap struct {
	Id string
	Name string
}
