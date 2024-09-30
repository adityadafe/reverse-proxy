package main

import (
	"encoding/json"
	"fmt"
	"log"
	//"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var (
	mu       sync.Mutex
	networks []Network
)

type (
	dockerMap *map[string]string
)

func addInMap(container dockerMap, key, val string) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := (*container)[key]; !ok {
		(*container)[key] = val
	}
}

func deleteInMap(container dockerMap, key string) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := (*container)[key]; ok {
		delete((*container), key)
	}
}

func checkRunningContainer(container dockerMap) {
	for {
		for key, val := range *container {
			fmt.Println(key, val)
			out, _ := exec.Command("ping", val).Output()
			fmt.Println(string(out))
			if strings.Contains(string(out), "Destination Host Unreachable") {
				go deleteInMap(container, key)
			}
		}
		time.Sleep(3 * time.Second)
	}

}

func getRunningContainers(container dockerMap, networks *[]Network) {
	for {
		dockerOutput, err := exec.Command("docker", "network", "inspect", "b").Output()
		if err != nil {
			log.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}
		err = json.Unmarshal(dockerOutput, &networks)

		if err != nil {
			log.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}

		for _, vals := range *networks {
			for _, val := range vals.Containers {
				fmt.Println("Containers ", val.Name, val.IPv4Address)
				go addInMap(container, val.Name, val.IPv4Address)
			}
		}

		time.Sleep(2 * time.Second)
	}
}

func temp() {
	out, _ := exec.Command("ping", "google.com").Output()
	fmt.Println(string(out))

}

func main() {
	// mux := http.NewServeMux()
	// container := make(map[string]string)
	// go getRunningContainers(&container, &networks)
	// go checkRunningContainer(&container)
	//
	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//
	// 	http.ServeFile(w, r, "index.html")
	// })
	//
	// log.Println("Starting server at http://localhost:8080")
	// log.Fatal(http.ListenAndServe(":8080", mux))

	fmt.Println()
	temp()
}
