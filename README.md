# stupid-horse
A bot that plays standard/960 chess on [Lichess](https://lichess.org), all time controls. Still a work in progress!

### How it works
- Set up a Lichess bot account
- Make an `.env` file with bot token and ID
- Run the bot (`go run .`), and it listens for incoming challenges and ongoing games.

### To do
- Make the bot care about time controls. Right now, the bot thinks at a certain depth no matter the time left.
- (Possibly) create a web interface to look at bot evaluations in live-time

### Known issues
- Bot does not put importance to sooner checkmate patterns. In [this game](https://lichess.org/7uwbkBlXsVR7), it played into a draw by repetition instead of a mate in 1.
- Bot cannot convert clearly winning positions (examples: [one](https://lichess.org/4Minxwys65gE))
- Bot sometimes doesn't hear opponent move notifications causing it to wait for a move forever (needs to be restarted to fix)

## My journey

### Why was this bot created?
- I enjoy chess
- I have always wanted to create a chess bot
- I saw it as a daunting challenge as a child and now I am equipped with the skill and experience to avenge myself
- I saw this as a nice exercise with a lots of room to learn

### My process
- Browsed through Go tutorials
- Created functions and structs to parse FEN, and set the logic involving board, pieces
- Set up move searching to get all valid moves for a given position (state of chess board in a game)
- Wrote tests and benchmark functions along the way to look for optimizations
- Implemented the minimax algorithm (with α-β) to choose best move after considering all possible outcomes at a certain depth (e.g. 3 moves ahead)
- Hooked up the Lichess API and made the bot accept challenges and play moves!
- To be continnued...

### Things I learned
- Coding in Go! I found that I really enjoy Go due to its readability, strong typing, intuitive keywords, built-in test writing, etc...
- Writing modular components and tests & benchmarks, making it easy to add new features and optimize parts of the project
- Consuming an API in Go, especially handling streams in [ndjson](http://ndjson.org/) format
- The fastest stalemate from the standard starting position, researched while creating tests for game state
- Honestly a lot more that I'd love to chat about
