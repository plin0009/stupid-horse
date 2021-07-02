package main

type Piece byte
type PieceType byte
type PieceColour bool

type Square byte
type File byte
type Rank byte

type Board [64]Piece

const (
	Pawn PieceType = iota
	Knight
	Bishop
	Rook
	Queen
	King
)

const (
	White PieceColour = true
	Black             = false
)

const NoSquare = Square(64)

var fenPieces = map[rune]Piece{
	'P': 10,
	'N': 11,
	'B': 12,
	'R': 13,
	'Q': 14,
	'K': 15,
	'p': 20,
	'n': 21,
	'b': 22,
	'r': 23,
	'q': 24,
	'k': 25,
}

var pieceSymbols = map[Piece]string{
	10: "♙",
	11: "♘",
	12: "♗",
	13: "♖",
	14: "♕",
	15: "♔",
	20: "♟︎",
	21: "♞",
	22: "♝",
	23: "♜",
	24: "♛",
	25: "♚",
}

func (p Piece) Type() PieceType {
	return PieceType(p % 10)
}
func (p Piece) Colour() PieceColour {
	return PieceColour((p/10)%2 == 0)
}
func (p Piece) Value() int {
	switch p.Type() {
	case Pawn:
		return 10
	case Knight:
		return 70
	case Bishop:
		return 30
	case Rook:
		return 50
	case Queen:
		return 90
	case King:
		return 999
	}
	return 0
}

func (pt PieceType) String() string {
	switch pt {
	case Pawn:
		return "pawn"
	case Knight:
		return "knight"
	case Bishop:
		return "bishop"
	case Rook:
		return "rook"
	case Queen:
		return "queen"
	case King:
		return "king"
	}
	return ""
}
func (pc PieceColour) String() string {
	if pc == White {
		return "white"
	} else {
		return "black"
	}
}

func (b Board) String() string {
	str := ""
	for i, tile := range b {
		if i%8 == 0 {
			str += "\n"
		}
		if tile == 0 {
			str += " "
			continue
		} else {
			str += pieceSymbols[tile]
		}
	}
	return str
}

func (s Square) File() File {
	return File(s % 8)
}
func (s Square) Rank() Rank {
	return Rank(7 - s/8)
}
func (f File) String() string {
	return string(f + 'a')
}
func (r Rank) String() string {
	return string(r + '1')
}
func (s Square) String() string {
	return s.File().String() + s.Rank().String()
}
