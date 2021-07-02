package main

type Position struct {
	board           Board
	turn            PieceColour
	rookSquares     [4]Square // for castling
	kingMoved       [2]bool
	enPassantSquare Square
}
type Move struct {
	piece    Piece
	from     Square
	to       Square
	captured Square
	castle   bool
}

func LoadInitialPosition(fen string) Position {
	var b Board
	var rs [4]Square
	nr := 0
	s := 0
	for _, char := range fen {
		if char-'1' >= 0 && '8'-char >= 0 {
			s += int(char-'1') + 1
			continue
		}
		if char == '/' {
			continue
		}
		if char == ' ' {
			break
		}
		b[s] = fenPieces[char]
		if b[s].Type() == Rook {
			rs[nr] = Square(s)
			nr++
		}
		s++
	}
	return Position{
		board:           b,
		turn:            White,
		rookSquares:     rs,
		kingMoved:       [2]bool{false, false},
		enPassantSquare: NoSquare,
	}
}
