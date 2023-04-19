package tests

import (
	"encoding/json"
	"fmt"
	db "hsk-bikeapp-solita/database"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/html"
)

func TestMain(m *testing.M) {
	// Get the current directory of the first terminal
	currDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Create a variable for the project directory one level below the current directory
	projectDir := filepath.Join(currDir, "..")

	t := new(testing.T)

	if runtime.GOOS == "darwin" {
		testerIsMacUser(t, projectDir)
	} else {
		//todo Windows handling
	}

}

func testerIsMacUser(t *testing.T, projectDir string) {
	// Launch a second terminal and execute the project in it
	cmd := exec.Command("osascript", "-e", fmt.Sprintf(`
        tell application "Terminal"
            do script "cd %s; go run main.go"
        end tell
    `, projectDir))

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	cmd = exec.Command("pgrep", "-n", "Terminal")
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error finding Terminal process: %v", err)
	}
	pid := strings.TrimSpace(string(out))

	attemptsCount := 1
	maxAttempts := 50

	for {
		if isLocalhostAvailable(t) {
			log.Println("Localhost is available")
			break
		} else {
			log.Printf("Localhost is not available yet! Retrying in 1 second (Attempt %d of %d)\n", attemptsCount, maxAttempts)
		}

		if attemptsCount == maxAttempts {
			break
		}

		attemptsCount++
		time.Sleep(time.Second)
	}

	if attemptsCount != maxAttempts {
		TestGetInvalidStationID(t)
		TestGetNonExistingStation(t)
		TestGetValidStationInfo(t)
	}

	cmd = exec.Command("kill", "-9", pid)
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error killing proces: %v", err)
	}

	if attemptsCount == maxAttempts {
		log.Fatalf("Maximum number of tries reached (%d). Aborting...", maxAttempts)
	}
}

func isLocalhostAvailable(t *testing.T) bool {
	resp, err := http.Get("http://localhost:8080/")
	if err != nil {
		log.Println("Error making GET request to localhost:", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusMovedPermanently || resp.StatusCode == http.StatusOK {
		_, err := html.Parse(resp.Body)
		if err != nil {
			log.Println("Error parsing HTML!:", err)
		}

		time.Sleep(time.Second)
		return true

	} else if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound {
		log.Println("Localhost is not available, status code:", resp.StatusCode)
		return false
	}

	log.Println("Localhost is not available.", resp.StatusCode)
	return false
}

func TestGetNonExistingStation(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/stations?id=1000")
	if err != nil {
		t.Fatalf("Error making GET request to localhost: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status code %v but got %v", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestGetInvalidStationID(t *testing.T) {

	resp, err := http.Get("http://localhost:8080/stations?id=abc")
	if err != nil {
		t.Fatalf("Error making GET request to localhost: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status code %v but got %v", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestGetValidStationInfo(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/stations?id=100")
	if err != nil {
		t.Fatalf("Error making GET request to localhost: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %v but got %v", http.StatusOK, resp.StatusCode)
	}

	var stations []db.Station
	if err := json.NewDecoder(resp.Body).Decode(&stations); err != nil {
		t.Fatalf("Error decoding JSON response: %v", err)
	}

	if stations[0].ID != 100 {
		t.Fatalf("Expected station ID to be %v but got %v", 100, stations[0].ID)
	}

}
