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

var tests = []func(*testing.T){
	TestGetInvalidStationID,
	TestGetNonExistingStation,
	TestGetValidStationInfo,
}

func TestMain(m *testing.M) {
	currDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	t := new(testing.T)

	projectDir := filepath.Join(currDir, "..")

	runTestsByOS(t, projectDir)

}

func runTestsByOS(t *testing.T, projectDir string) {
	tester, err := getTesterByOS(runtime.GOOS)
	if err != nil {
		log.Fatalf("failed to get tester: %v", err)
	}

	err = tester(t, projectDir)
	if err != nil {
		t.Fatalf("failed to run test: %v", err)
	}
}

func getTesterByOS(os string) (func(*testing.T, string) error, error) {
	switch os {
	case "darwin":
		return testerIsMacUser, nil
	case "windows":
		return testerIsWindowUser, nil
	default:
		return nil, fmt.Errorf("unsupported operating system: %s", os)
	}
}

func testerIsMacUser(t *testing.T, projectDir string) error {
	cmd := exec.Command("osascript", "-e", fmt.Sprintf(`
        tell application "Terminal"
            do script "cd %s; go run main.go"
        end tell
    `, projectDir))

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to start process: %v", err)
	}

	cmd = exec.Command("pgrep", "-n", "Terminal")
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error finding Terminal process: %v", err)
	}

	success := retryTestsUntilSuccess(t, 50, tests)
	if !success {
		return fmt.Errorf("tests failed")
	}

	pid := strings.TrimSpace(string(out))

	cmd = exec.Command("kill", "-9", pid)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error killing process: %v", err)
	}

	return nil
}

func testerIsWindowUser(t *testing.T, projectDir string) error {
	cmd := exec.Command("cmd", "/c", "start", "cmd", "/k", "go", "run", "main.go")
	cmd.Dir = projectDir
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start process: %v", err)
	}

	success := retryTestsUntilSuccess(t, 50, tests)
	if !success {
		return fmt.Errorf("tests failed")
	}

	time.Sleep(time.Second)

	cmd = exec.Command("taskkill", "/IM", "cmd.exe", "/T", "/F")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to kill process: %v", err)
	}

	return nil
}

func retryTestsUntilSuccess(t *testing.T, maxAttempts int, tests []func(*testing.T)) bool {
	attemptsCount := 1

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
		for _, test := range tests {
			test(t)
		}
		return true
	}

	log.Fatalf("Maximum number of tries reached (%d). Aborting...", maxAttempts)
	return false
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

	var station db.Station
	if err := json.NewDecoder(resp.Body).Decode(&station); err != nil {
		t.Fatalf("Error decoding JSON response: %v", err)
	}

	if station.ID != 100 {
		t.Fatalf("Expected station ID to be %v but got %v", 100, station.ID)
	}

}
