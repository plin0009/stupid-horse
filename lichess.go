package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const url = "https://lichess.org/api/"

type Bot struct {
	games []Game
	token string
}

type Game struct {
	moveTree MoveTree
}

type LichessEvent struct {
	Type      string
	Game      *LichessGame
	Challenge *LichessChallenge
}

type LichessGame struct {
	Id string
}

type LichessChallenge struct {
	Id          string
	Status      string
	Variant     LichessVariant
	Rated       bool
	TimeControl LichessTimeControl
	Color       string
}
type LichessVariant struct {
	Key string
}
type LichessTimeControl struct {
	Type      string
	Limit     int
	Increment int
	Show      string
}

func (b Bot) Listen() {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url+"stream/event", nil)
	if err != nil {
		panic("Could not create request")
	}

	req.Header.Add("Authorization", "Bearer "+b.token)
	resp, err := client.Do(req)
	if err != nil {
		panic("Error with request")
	}
	defer resp.Body.Close()
	for times := 0; times < 20; times++ {
		bytes := make([]byte, 8192)
		_, err := resp.Body.Read(bytes)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		start := 0
		var b []byte
		fmt.Println("Started listening")
		for i, v := range bytes {
			if v == 10 {
				b = bytes[start:i]
				if len(b) == 0 {
					break
				}
				var e LichessEvent
				err := json.Unmarshal(b, &e)
				if err != nil {
					log.Fatal("Couldn't decode JSON")
				}
				fmt.Printf("%s\n", b)
				switch e.Type {
				case "gameStart":
					// load game to bot
					fmt.Println("adding game")
				case "gameFinish":
					// remove game from bot
					fmt.Println("removing finished game")
				case "challenge":
					fmt.Printf("received %s challenge (%s)\n", e.Challenge.Variant.Key, e.Challenge.Id)
					// add challenge to bot queue
				case "challengeCanceled", "challengeDeclined":
					fmt.Printf("cancelled %s challenge (%s)\n", e.Challenge.Variant.Key, e.Challenge.Id)
					// remove existing challenge in bot
				}
				start = i + 1
			}
		}
	}
	fmt.Println("Stopped listening")
}
