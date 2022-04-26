package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	registryURL        string
	defaultRegistryURL = "http://127.0.0.1:8080/servers"

	beaconTimeout  = 5 * time.Second
	beaconInterval = 10 * time.Second

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

	// Global server cache.
	Servers = make([]ServerInfo, 0)
)

func init() {
	flag.StringVar(&registryURL, "registry-url", defaultRegistryURL, "Full URL for openrvs-registry")
	flag.Parse()
}

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
		// Sleep early to avoid fast iterations on registry failures.
		time.Sleep(beaconInterval)

		// Retrieve healthy servers from openrvs-registry over HTTP.
		healthyServers, err := getServersFromRegistry()
		if err != nil {
			log.Println(err)
			continue
		}

		// Retrieve game info for healthy servers over UDP.
		newServers := make([]ServerInfo, 0)
		var wg sync.WaitGroup
		for _, hp := range healthyServers {
			wg.Add(1)
			go func(hp HostPort) {
				// Get fresh server data.
				info, err := populateBeaconData(hp)
				if err != nil {
					log.Println("beacon error:", err)
					wg.Done()
					return
				}
				// Filter empty servers.
				if info.CurrentPlayers > 0 {
					newServers = append(newServers, info)
				}
				wg.Done()
			}(hp)
		}
		wg.Wait()

		// Rebuild server cache.
		Servers = newServers
	}
}
