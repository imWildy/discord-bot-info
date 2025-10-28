package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func SendRequest(method string, url string, token string) (*http.Response, error) {
	var resp *http.Response
	var req *http.Request
	var err error

	client := &http.Client{}
	req, err = http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bot "+token)
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp, err
}

func main() {
	var token, id, username, discriminator string
	fmt.Print("Bot Token: ")
	fmt.Scanln(&token)
	url := "https://discord.com/api/v9/"

	var resp *http.Response
	var r []byte
	var err error

	// Get bot info (Valid Token?, ID, Username, etc)
	resp, err = SendRequest("GET", url+"users/@me", token)
	if err != nil {
		log.Fatal(err)
	} else if resp.StatusCode == 401 {
		fmt.Println("Invalid Token!")
		main()
		os.Exit(0)
	} else if resp.StatusCode != 200 {
		log.Fatal(resp.Status)
	}
	defer resp.Body.Close()

	r, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var respData map[string]any
	if err := json.Unmarshal(r, &respData); err != nil {
		log.Fatal(err)
	}
	id = respData["id"].(string)
	username = respData["username"].(string)
	discriminator = respData["discriminator"].(string)

	fmt.Println("-----Bot Info-----")
	fmt.Printf("ID: %s\n", id)
	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Discriminator: %s\n", discriminator)

	// Get guilds bot is in and some info (id, name, etc)
	resp, err = SendRequest("GET", url+"users/@me/guilds", token)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	r, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var data []any
	if err := json.Unmarshal(r, &data); err != nil {
		log.Fatal(err)
	}

	fmt.Println("-----Guild Info-----")
	fmt.Printf("Bot is in %d guilds.\n", len(data))

	for i, guild := range data {
		g := guild.(map[string]any)

		// Get member count of guild
		resp, err = SendRequest("GET", fmt.Sprintf("%sguilds/%s?with_counts=true", url, g["id"]), token)
		if err != nil {
			log.Fatal(err)
		}

		r, err = io.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		var data map[string]any
		if err := json.Unmarshal(r, &data); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Guild %d\n", i+1)
		fmt.Printf("  ID: %s\n", g["id"])
		fmt.Printf("  Name: %s\n", g["name"])
		fmt.Printf("  Members: %.0f\n", data["approximate_member_count"])
	}
}
