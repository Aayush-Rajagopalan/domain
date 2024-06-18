package lib

import (
	"net/http"
	"os"
)

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
