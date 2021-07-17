package main

import (
	"io"
	"log"
	"net/http"
)

const url = "https://lichess.org/api/"

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

type LichessGameEvent struct {
	Type string

	// gameFull
	Variant    *LichessVariant
	Rated      bool
	Clock      *LichessClock
	White      *LichessPlayer
	Black      *LichessPlayer
	InitialFen string
	State      *LichessGameEvent

	// gameState
	Moves  string
	Wtime  uint
	Btime  uint
	Winc   uint
	Binc   uint
	Status string
	Winner string

	// chatLine
	Username string
	Text     string
	Room     string
}
type LichessPlayer struct {
	Id     string
	Name   string
	Rating int
	Title  string
}
type LichessClock struct {
	Initial   uint
	Increment uint
}

func (b Bot) request(mode string, path string) (resp *http.Response, err error) {
	req, err := http.NewRequest(mode, url+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+b.token)
	return http.DefaultClient.Do(req)
}

func stream(resp *http.Response, lines chan<- []byte) {
	for {
		bytes := make([]byte, 8192)
		_, err := resp.Body.Read(bytes)
		if err == io.EOF {
			close(lines)
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		start := 0
		var line []byte
		for i, v := range bytes {
			if v == 10 {
				line = bytes[start:i]
				if len(line) == 0 {
					break
				}
				lines <- line
			}
		}
	}
}
