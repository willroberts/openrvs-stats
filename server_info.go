package main

import (
	beacon "github.com/willroberts/openrvs-beacon"
	registry "github.com/willroberts/openrvs-registry"
)

type ServerInfo struct {
	ServerName     string       `json:"server_name"`
	CurrentPlayers int          `json:"current_players"`
	MaxPlayers     int          `json:"max_players"`
	IP             string       `json:"ip_address"`
	Port           int          `json:"port"`
	MapName        string       `json:"current_map"`
	GameMode       string       `json:"game_mode"`
	MOTD           string       `json:"motd"`
	Players        []Player     `json:"players"`
	Maps           []Map        `json:"maps"`
	PVPSettings    PVPSettings  `json:"pvp_settings"`
	CoopSettings   CoopSettings `json:"coop_settings"`
}

type Player struct {
	Name  string `json:"name"`
	Kills int    `json:"kills"`
	Time  string `json:"time"`
}

type Map struct {
	Name string `json:"name"`
	Mode string `json:"mode"`
}

type PVPSettings struct {
	AutoTeamBalance   bool `json:"auto_team_balance"`
	BombTimer         int  `json:"bomb_timer"`
	FriendlyFire      bool `json:"friendly_fire"`
	RoundsPerMatch    int  `json:"rounds_per_match"`
	TimePerRound      int  `json:"time_per_round"`
	TimeBetweenRounds int  `json:"time_between_rounds"`
}

type CoopSettings struct {
	AIBackup           bool `json:"ai_backup"`
	FriendlyFire       bool `json:"friendly_fire"`
	TerroristCount     int  `json:"terrorist_count"`
	RotateMapOnSuccess bool `json:"rotate_map_on_success"`
	RoundsPerMatch     int  `json:"rounds_per_match"`
	TimePerRound       int  `json:"time_per_round"`
	TimeBetweenRounds  int  `json:"time_between_rounds"`
}

func populateBeaconData(hp HostPort) (ServerInfo, error) {
	b, err := beacon.GetServerReport(hp.IP, hp.Port+1000, beaconTimeout)
	if err != nil {
		return ServerInfo{}, err
	}
	report, err := beacon.ParseServerReport(hp.IP, b)
	if err != nil {
		return ServerInfo{}, err
	}
	info := ServerInfo{
		ServerName:     report.ServerName,
		CurrentPlayers: report.NumPlayers,
		MaxPlayers:     report.MaxPlayers,
		IP:             report.IPAddress,
		Port:           report.Port,
		MapName:        report.CurrentMap,
		GameMode:       FriendlyGameModes[report.CurrentMode],
		MOTD:           report.MOTD,
	}
	var players = make([]Player, 0)
	for i := 0; i < len(report.ConnectedPlayerNames); i++ {
		var p Player
		p.Name = report.ConnectedPlayerNames[i]
		p.Kills = report.ConnectedPlayerKills[i]
		p.Time = report.ConnectedPlayerTimes[i]
		players = append(players, p)
	}
	info.Players = players
	var maps = make([]Map, 0)
	for i := 0; i < len(report.MapRotation); i++ {
		var m Map
		m.Name = report.MapRotation[i]
		m.Mode = "?"
		if i < len(report.ModeRotation) {
			mode, ok := FriendlyGameModes[report.ModeRotation[i]]
			if ok {
				m.Mode = mode
			}
		}
		maps = append(maps, m)
	}
	info.Maps = maps
	if registry.GameTypes[report.CurrentMode] == "adv" {
		var pvp PVPSettings
		pvp.AutoTeamBalance = report.AutoTeamBalance
		if report.CurrentMode == "Bomb" {
			pvp.BombTimer = report.BombTimer
		}
		pvp.FriendlyFire = report.FriendlyFire
		pvp.RoundsPerMatch = report.RoundsPerMatch
		pvp.TimePerRound = report.TimePerRound
		pvp.TimeBetweenRounds = report.TimeBetweenRounds
		info.PVPSettings = pvp
	} else {
		var coop CoopSettings
		coop.AIBackup = report.AIBackup
		coop.FriendlyFire = report.FriendlyFire
		coop.TerroristCount = report.NumTerrorists
		coop.RotateMapOnSuccess = report.RotateMapOnSuccess
		coop.RoundsPerMatch = report.RoundsPerMatch
		coop.TimePerRound = report.TimePerRound
		coop.TimeBetweenRounds = report.TimeBetweenRounds
		info.CoopSettings = coop
	}
	return info, nil
}
