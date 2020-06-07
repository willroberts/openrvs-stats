package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	beacon "github.com/ijemafe/openrvs-beacon"
	registry "github.com/ijemafe/openrvs-registry"
)

var Servers = make([]ServerInfo, 0)
var beaconTimeout = 5 * time.Second
var beaconInterval = 15 * time.Second

func main() {
	go pollServers()

	http.HandleFunc("/stats.json", func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(Servers)
		if err != nil {
			log.Println("marshal error:", err)
			w.Write([]byte("error"))
		}
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		w.Write(b)
	})
	log.Println("listening on http://127.0.0.1:8081")
	log.Fatal(http.ListenAndServe("127.0.0.1:8081", nil))
}

func pollServers() {
	for {
		inputs, err := getHostPorts()
		if err != nil {
			log.Println(err)
			continue
		}

		var wg sync.WaitGroup
		var lock = sync.RWMutex{}
		for _, input := range inputs {
			wg.Add(1)

			go func(input Input) {
				info, err := populateBeaconData(input)
				if err != nil {
					log.Println("beacon error:", err)
					wg.Done()
					return
				}
				lock.Lock()
				for i, s := range Servers {
					// Server is already in the list, update or remove it.
					if info.IP == s.IP && info.Port == s.Port {
						if info.CurrentPlayers > 0 {
							Servers[i] = info
						} else {
							Servers = append(Servers[:i], Servers[i+1:]...)
						}
						lock.Unlock()
						wg.Done()
						return
					}
				}
				// Server is not in the list, add it.
				if info.CurrentPlayers > 0 {
					Servers = append(Servers, info)
				}
				lock.Unlock()
				wg.Done()
			}(input)
		}
		wg.Wait()
		log.Println("server info updated")
		time.Sleep(beaconInterval)
	}
}

func getHostPorts() ([]Input, error) {
	var inputs = make([]Input, 0)
	resp, err := http.Get("http://127.0.0.1:8080/servers")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	lines := bytes.Split(bytes.TrimSuffix(b, []byte{'\n'}), []byte{'\n'})
	for i := 1; i < len(lines); i++ {
		fields := bytes.Split(lines[i], []byte{','})
		host := string(fields[1])
		portBytes := fields[2]
		port, err := strconv.Atoi(string(portBytes))
		if err != nil {
			log.Println("atoi error:", err)
			continue
		}
		inputs = append(inputs, Input{IP: host, Port: port})
	}
	return inputs, nil
}

func populateBeaconData(input Input) (ServerInfo, error) {
	b, err := beacon.GetServerReport(input.IP, input.Port+1000, beaconTimeout)
	if err != nil {
		return ServerInfo{}, err
	}
	report, err := beacon.ParseServerReport(input.IP, b)
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
		GameMode:       report.CurrentMode,
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
		m.Mode = "unknown"
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

type Input struct {
	IP   string
	Port int
}

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

var FriendlyGameModes = map[string]string{
	"RGM_BombAdvMode":           "Bomb",
	"RGM_DeathmatchMode":        "Survival",
	"RGM_EscortAdvMode":         "Escort the Pilot",
	"RGM_HostageRescueAdvMode":  "Hostage",
	"RGM_HostageRescueCoopMode": "Hostage Rescue",
	"RGM_MissionMode":           "Mission",
	"RGM_TeamDeathmatchMode":    "Team Survival",
	"RGM_TerroristHuntCoopMode": "Terrorist Hunt",
}
