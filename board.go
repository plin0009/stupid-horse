package main

type Piece byte
type PieceType byte
type PieceColour bool

type Square byte
type File byte
type Rank byte

type Board [64]Piece

const (
	NoPieceType PieceType = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

const (
	White PieceColour = false
	Black             = true
)

const NoPiece = Piece(0)
const NoSquare = Square(64)

var Squares = [64]Square{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63}

var fenPieces = map[rune]Piece{
	'P': CreatePiece(White, Pawn),
	'N': CreatePiece(White, Knight),
	'B': CreatePiece(White, Bishop),
	'R': CreatePiece(White, Rook),
	'Q': CreatePiece(White, Queen),
	'K': CreatePiece(White, King),
	'p': CreatePiece(Black, Pawn),
	'n': CreatePiece(Black, Knight),
	'b': CreatePiece(Black, Bishop),
	'r': CreatePiece(Black, Rook),
	'q': CreatePiece(Black, Queen),
	'k': CreatePiece(Black, King),
}

var pieceSymbols = map[Piece]string{
	CreatePiece(White, Pawn):   "♙",
	CreatePiece(White, Knight): "♘",
	CreatePiece(White, Bishop): "♗",
	CreatePiece(White, Rook):   "♖",
	CreatePiece(White, Queen):  "♕",
	CreatePiece(White, King):   "♔",
	CreatePiece(Black, Pawn):   "♟︎",
	CreatePiece(Black, Knight): "♞",
	CreatePiece(Black, Bishop): "♝",
	CreatePiece(Black, Rook):   "♜",
	CreatePiece(Black, Queen):  "♛",
	CreatePiece(Black, King):   "♚",
}

func CreatePiece(pc PieceColour, pt PieceType) Piece {
	if pc == White {
		return Piece(10 + int(pt))
	}
	return Piece(20 + int(pt))
}

func (p Piece) Type() PieceType {
	if p == NoPiece {
		return NoPieceType
	}
	return PieceType(p % 10)
}
func (p Piece) Colour() PieceColour {
	if p == NoPiece {
		panic("No piece")
	}
	return PieceColour((p/10)%2 == 0)
}
func (pt PieceType) String() string {
	switch pt {
	case Pawn:
		return "p"
	case Knight:
		return "n"
	case Bishop:
		return "b"
	case Rook:
		return "r"
	case Queen:
		return "q"
	case King:
		return "k"
	}
	return ""
}
func StringToPieceType(str string) PieceType {
	switch str {
	case "p":
		return Pawn
	case "n":
		return Knight
	case "b":
		return Bishop
	case "r":
		return Rook
	case "q":
		return Queen
	case "k":
		return King
	}
	return NoPieceType
}
func (pc PieceColour) String() string {
	if pc == White {
		return "white"
	} else {
		return "black"
	}
}
func (p Piece) String() string {
	return pieceSymbols[p]
}

func (b Board) String() string {
	str := ""
	for i, piece := range b {
		if i > 0 && i%8 == 0 {
			str += "\n"
		}
		if piece == NoPiece {
			str += " "
			continue
		} else {
			str += pieceSymbols[piece]
		}
	}
	return str
}

func ToSquare(f File, r Rank) Square {
	return Square((7-byte(r))*8 + byte(f))
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
func StringToSquare(str string) Square {
	if len(str) != 2 {
		return NoSquare
	}
	fileNum := rune(str[0]) - 'a'
	rankNum := rune(str[1]) - '1'
	if fileNum < 0 || fileNum > 7 {
		return NoSquare
	}
	if rankNum < 0 || rankNum > 7 {
		return NoSquare
	}
	return ToSquare(File(fileNum), Rank(rankNum))
}

func (pc PieceColour) Flip() PieceColour {
	if pc == White {
		return Black
	}
	return White
}
