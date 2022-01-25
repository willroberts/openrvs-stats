package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

const RegistryURL = "http://127.0.0.1:8080/servers"

var (
	beaconTimeout     = 5 * time.Second
	beaconInterval    = 15 * time.Second
	FriendlyGameModes = map[string]string{
		"RGM_BombAdvMode":           "Bomb",
		"RGM_DeathmatchMode":        "Survival",
		"RGM_EscortAdvMode":         "Escort the Pilot",
		"RGM_HostageRescueAdvMode":  "Hostage",
		"RGM_HostageRescueCoopMode": "Hostage Rescue",
		"RGM_MissionMode":           "Mission",
		"RGM_TeamDeathmatchMode":    "Team Survival",
		"RGM_TerroristHuntCoopMode": "Terrorist Hunt",
	}
)

// Global server cache.
var Servers = make([]ServerInfo, 0)

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

// Continuously refreshes beacon data every `beaconInterval` seconds.
func pollServers() {
	for {
		hostports, err := getHostPorts()
		if err != nil {
			log.Println(err)
			continue
		}

		var wg sync.WaitGroup
		var lock = sync.RWMutex{}
		for _, hp := range hostports {
			wg.Add(1)

			go func(hp HostPort) {
				info, err := populateBeaconData(hp)
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
			}(hp)
		}
		wg.Wait()
		log.Println("server info updated")
		time.Sleep(beaconInterval)
	}
}
