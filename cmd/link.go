package cmd

import (
	"domain/lib"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
)

// temporary cos i cant for the life of me understand why i cant access the lib folder funcs
func getData() (string, string, string, error) {
	// Open the configuration file
	file, err := os.Open("/config/config.cfg")
	if err != nil {
		return "", "", "", fmt.Errorf("error opening config file: %v", err)
	}
	defer file.Close()

	// Read the file content
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", "", "", fmt.Errorf("error reading config file: %v", err)
	}

	// Unmarshal the JSON data into the Config struct
	var config lib.Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return "", "", "", fmt.Errorf("error parsing config file: %v", err)
	}

	// Return the token and email from the configuration
	return config.Token, config.Email, config.ZoneID, nil
}

// linkCmd represents the link command
var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "adds a dns record to the cloudflare account.",
	Args:  cobra.ExactArgs(3),
	Long:  `you give in the domain name and the ip address and the command will add the dns record to the cloudflare account.`,
	Run: func(cmd *cobra.Command, args []string) {
		token, email, zoneid, err := getData()
		if err != nil {
			fmt.Println(err)
		}
		typeFlag, _ := cmd.Flags().GetString("type")
		fmt.Println("creating a " + typeFlag + "record on " + args[0] + "." + args[1])
		client := &http.Client{}
		req, err := http.NewRequest("POST", "https://api.cloudflare.com/client/v4/zones/"+zoneid+"/dns_records", nil)
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Auth-Email", email)
		req.Header.Set("X-Auth-Key", token)
		req.Body = ioutil.NopCloser(bytes.NewReader([]byte(`{"type":"` + typeFlag + `","name":"` + args[0] + "." + args[1] + `","content":"` + args[2] + `","ttl":120,"proxied":false}`)))
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		if resp.StatusCode != 200 {
			fmt.Println("error")
		} else {
			fmt.Println("created dns record")
		}
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	linkCmd.PersistentFlags().String("type", "", "your dns record type")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// linkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
