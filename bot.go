package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Bot struct {
	id    string
	token string
	game  *Game
}

// converts Lichess game player data to struct used by the bot
func (b Bot) ToPlayer(lp LichessPlayer) Player {
	return Player{
		id:     lp.Id,
		name:   lp.Name,
		rating: lp.Rating,
		title:  lp.Title,
		me:     lp.Id == b.id,
	}
}

func StartBot() {
	err := godotenv.Load()
	if err != nil {
		panic("could not load .env file")
	}
	b := Bot{id: os.Getenv("LICHESS_BOT_ID"), token: os.Getenv("LICHESS_KEY")}
	b.Listen()
}

func (b *Bot) Listen() {
	resp, err := b.request("GET", "stream/event")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println("Started listening")
	lines := make(chan []byte)
	go stream(resp, lines)
	for line := range lines {
		var e LichessEvent
		err := json.Unmarshal(line, &e)
		if err != nil {
			log.Fatal("Couldn't decode JSON")
		}
		fmt.Printf("%s\n", line)
		switch e.Type {
		case "gameStart":
			// load game to bot
			go b.loadGame(e.Game)
		case "gameFinish":
			// remove game from bot
			fmt.Println("removing finished game")
			b.game = nil
		case "challenge":
			b.considerChallenge(e.Challenge)
		case "challengeCanceled", "challengeDeclined":
			fmt.Printf("cancelled %s challenge (%s)\n", e.Challenge.Variant.Key, e.Challenge.Id)
			// TODO: remove existing challenge in bot queue
		}
	}
	fmt.Println("Stopped listening")
}

func (b *Bot) considerChallenge(c *LichessChallenge) {
	fmt.Printf("considering %s challenge (%s)\n", c.Variant.Key, c.Id)
	// TODO: instead of ignoring challenges, add them to a queue or decline
	switch {
	case b.game != nil: // game in progress
		return
	case c.Variant.Key != "standard" && c.Variant.Key != "chess960":
		return
	default:
		resp, err := b.request("POST", "challenge/"+c.Id+"/accept")
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
	}
}

func (b *Bot) loadGame(g *LichessGame) {
	fmt.Printf("adding game %v\n", g.Id)
	fmt.Println(*g)
	resp, err := b.request("GET", "bot/game/stream/"+g.Id)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println("Started listening")
	lines := make(chan []byte)
	go stream(resp, lines)
	for line := range lines {
		var e LichessGameEvent
		err := json.Unmarshal(line, &e)
		if err != nil {
			log.Fatal("Couldn't decode JSON")
		}
		fmt.Printf("%s\n", line)
		switch e.Type {
		case "gameFull":
			clock := map[PieceColour]uint{}
			if e.Clock != nil {
				clock = map[PieceColour]uint{
					White: e.Clock.Initial,
					Black: e.Clock.Initial,
				}
			}
			b.game = NewGame(e.InitialFen,
				map[PieceColour]Player{
					White: b.ToPlayer(*e.White),
					Black: b.ToPlayer(*e.Black),
				}, clock)
			b.game.id = g.Id
			// read state field
			b.ProcessGameState(*e.State)
		case "gameState":
			b.ProcessGameState(e)
		case "chatLine":
			fmt.Printf("%v said \"%v\" in %v", e.Username, e.Text, e.Room)
		}
	}
	fmt.Println("Stopped listening")
}

func (b *Bot) ProcessGameState(s LichessGameEvent) {
	if s.Type != "gameState" {
		log.Fatal("Not a gameState")
	}
	fmt.Println(s.Status)
	if s.Status != "started" {
		return
	}
	if b.game == nil {
		fmt.Println("Game does not exist anymore")
		return
	}
	// update timers
	currentMs := time.Now().UnixNano() / 1000
	b.game.timers[White] = uint(currentMs) + s.Wtime
	b.game.timers[Black] = uint(currentMs) + s.Btime
	// update moves
	oldMoves := b.game.moves
	curMoves := s.Moves
	fmt.Println("Old:", oldMoves)
	fmt.Println("New:", curMoves)
	if curMoves != oldMoves {
		if len(curMoves) < len(oldMoves) {
			log.Fatal("New moves do not build on old moves")
		}
		if curMoves[:len(oldMoves)] != oldMoves {
			log.Fatal("New moves do not build on old moves")
		}
		newMoves := curMoves[len(oldMoves):]
		b.game.AddMoves(newMoves)
	}
	fmt.Println("Now", b.game.moveTree.position.turn, "to move")
	if b.game.players[b.game.moveTree.position.turn].me {
		b.Think()
	}
}

func (b *Bot) Think() {
	Think(b.game.moveTree, 6)
	fmt.Println(b.game.moveTree.follow)
	b.MakeMove(b.game.moveTree.follow.move)
}

func (b *Bot) MakeMove(m Move) {
	resp, err := b.request("POST", "bot/game/"+b.game.id+"/move/"+m.String())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println("Made move", m)
}
