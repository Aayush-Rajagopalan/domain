package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Config represents the structure of the configuration data
type Config struct {
	Token  string `json:"token"`
	Email  string `json:"email"`
	ZoneID string `json:"zoneid"`
}

// getData reads the configuration from /config/config.cfg and returns the token and email
func getData() (string, string, error) {
	// Open the configuration file
	file, err := os.Open("/config/config.cfg")
	if err != nil {
		return "", "", fmt.Errorf("error opening config file: %v", err)
	}
	defer file.Close()

	// Read the file content
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", "", fmt.Errorf("error reading config file: %v", err)
	}

	// Unmarshal the JSON data into the Config struct
	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return "", "", fmt.Errorf("error parsing config file: %v", err)
	}

	// Return the token and email from the configuration
	return config.Token, config.Email, nil
}

func verify(token string, email string) bool {
	// Check if the file exists
	if _, err := os.Stat("/config/confing.cfg"); os.IsNotExist(err) {
		return false
	} else {
		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://api.cloudflare.com/client/v4/accounts", nil)
		if err != nil {
			return false
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Auth-Email", email)
		req.Header.Set("X-Auth-Key", token)
		resp, err := client.Do(req)
		if err != nil {
			return false
		}
		if resp.StatusCode != 200 {
			return false
		} else {
			return true
		}
	}
}
