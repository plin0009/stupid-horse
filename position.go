package main

import (
	"errors"
	"fmt"
)

type Position struct {
	board           Board
	turn            PieceColour
	rookSquares     [4]Square // for castling
	whiteKingMoved  bool
	blackKingMoved  bool
	enPassantSquare Square
}
type Move struct {
	from    Square
	to      Square
	capture bool
	promote PieceType
	castle  CastleDirection
}

type CastleDirection int

const (
	NoCastle CastleDirection = iota
	ASide
	HSide
)

type MoveTree struct {
	parent         *MoveTree
	move           Move
	position       Position
	candidateMoves []Move
	legalMoves     []Move
	legal          bool
	eval           int
	follow         *MoveTree
	state          State
}
type Movement [2]int
type PawnInfo struct {
	homeRank      Rank
	forward       int
	promotionRank Rank
}
type State int

const (
	Active State = iota
	Stalemate
	WhiteWon
	BlackWon
)

func (s State) String() string {
	switch s {
	case Active:
		return "active"
	case Stalemate:
		return "stalemate"
	case WhiteWon:
		return "white won"
	case BlackWon:
		return "black won"
	}
	panic("Invalid state")
}

func WinFor(pc PieceColour) State {
	if pc == White {
		return WhiteWon
	}
	return BlackWon
}

var pieceMovements = map[PieceType][]Movement{
	Bishop: {{-1, -1}, {-1, 1}, {1, -1}, {1, 1}},
	Rook:   {{-1, 0}, {1, 0}, {0, -1}, {0, 1}},
	Queen:  {{-1, -1}, {-1, 1}, {1, -1}, {1, 1}, {-1, 0}, {1, 0}, {0, -1}, {0, 1}},
	Knight: {{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2}, {1, -2}, {1, 2}, {2, -1}, {2, 1}},
	King:   {{-1, -1}, {-1, 1}, {1, -1}, {1, 1}, {-1, 0}, {1, 0}, {0, -1}, {0, 1}},
}

func (p Position) kingMoved(c PieceColour) bool {
	if c == White {
		return p.whiteKingMoved
	}
	if c == Black {
		return p.blackKingMoved
	}
	panic("Invalid colour of king")
}

func (p *Position) setKingMoved(c PieceColour, v bool) {
	if c == White {
		p.whiteKingMoved = v
		return
	}
	if c == Black {
		p.blackKingMoved = v
		return
	}
	panic("Invalid colour of king")

}

var pawnInfo = map[PieceColour]PawnInfo{
	White: PawnInfo{Rank(1), 1, Rank(7)},
	Black: PawnInfo{Rank(6), -1, Rank(0)},
}

func LoadInitialPosition(fen string) Position {
	if fen == "startpos" {
		fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	}
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
		whiteKingMoved:  false,
		blackKingMoved:  false,
		enPassantSquare: NoSquare,
	}
}

func (m Move) String() string {
	return fmt.Sprintf("%v%v%v", m.from, m.to, m.promote)
}

// converts algebraic notation to move object
func (p Position) StringToMove(moveString string) Move {
	m := Move{
		from: StringToSquare(moveString[:2]),
		to:   StringToSquare(moveString[2:4]),
	}
	toPiece := p.board[m.to]
	if toPiece != NoPiece {
		if toPiece.Colour() != p.turn {
			m.capture = true
		} else {
			// own rook is at to square -- castle
			if m.to.File() < m.from.File() {
				m.castle = ASide
			} else {
				m.castle = HSide
			}
		}
	} else {
		// check for standard variant notation of castling
		fromPiece := p.board[m.from]
		if fromPiece.Type() == King {
			switch moveString {
			case "e8c8":
				m.to = p.rookSquares[1]
				m.castle = ASide
			case "e8g8":
				m.to = p.rookSquares[2]
				m.castle = HSide
			case "e1c1":
				m.to = p.rookSquares[2]
				m.castle = ASide
			case "e1g1":
				m.to = p.rookSquares[3]
				m.castle = HSide
			}
		}
	}
	if len(moveString) == 5 {
		m.promote = StringToPieceType(string(moveString[4]))
	}
	return m
}

func (s Square) Move(m Movement) (Square, error) {
	newFile, newRank := int(s.File())+m[0], int(s.Rank())+m[1]
	if newFile >= 8 || newFile < 0 {
		return NoSquare, errors.New("Out of bounds")
	}
	if newRank >= 8 || newRank < 0 {
		return NoSquare, errors.New("Out of bounds")
	}
	return ToSquare(File(newFile), Rank(newRank)), nil
}

func BetweenSquares(s1, s2 Square) []Square {
	s1Rank := s1.Rank()
	s1File := s1.File()
	squares := []Square{}
	if s2.Rank() != s1Rank {
		if s2.File() != s1File {
			panic(fmt.Sprintf("Squares not the rank nor file: %v %v", s1, s2))
		}
		s2Rank := s2.Rank()
		if s1Rank > s2Rank {
			for r := s1Rank - 1; r > s2Rank; r-- {
				squares = append(squares, ToSquare(s1File, r))
			}
		} else {
			for r := s1Rank + 1; r < s2Rank; r++ {
				squares = append(squares, ToSquare(s1File, r))
			}
		}
	} else {
		s2File := s2.File()
		if s1File > s2File {
			for f := s1File - 1; f > s2File; f-- {
				squares = append(squares, ToSquare(f, s1Rank))
			}
		} else {
			for f := s1File + 1; f < s2File; f++ {
				squares = append(squares, ToSquare(f, s1Rank))
			}
		}
	}
	return squares
}

func Diff(s1, s2 Square) (files int, ranks int) {
	files = int(s2.File()) - int(s1.File())
	ranks = int(s2.Rank()) - int(s1.Rank())
	if files < 0 {
		files = -files
	}
	if ranks < 0 {
		ranks = -ranks
	}
	return files, ranks
}

func GetCastleSquares(kingSquare, rookSquare Square) (Square, Square) {
	rank := kingSquare.Rank()
	if rookSquare.File() < kingSquare.File() {
		// a-side castling
		return ToSquare(File(2), rank), ToSquare(File(3), rank)
	} else {
		// h-side castling
		return ToSquare(File(6), rank), ToSquare(File(5), rank)
	}
}

func (p Position) ProcessMove(m Move) Position {
	fromPiece := p.board[m.from]
	toPiece := p.board[m.to]
	p.board[m.from] = NoPiece
	// if captured a rook, remove it from unmoved rooks
	if toPiece != NoPiece && toPiece.Type() == Rook {
		for i, square := range p.rookSquares {
			if square == m.to {
				p.rookSquares[i] = NoSquare
				break
			}
		}
	}
	p.board[m.to] = fromPiece
	p.enPassantSquare = NoSquare
	switch fromPiece.Type() {
	case Rook:
		// if moved an unmoved rook, remove it from unmoved rooks
		for i, square := range p.rookSquares {
			if square == m.from {
				p.rookSquares[i] = NoSquare
				break
			}
		}
	case King:
		p.setKingMoved(p.turn, true)
		if m.castle != NoCastle {
			p.board[m.to] = NoPiece
			kingSquare, rookSquare := GetCastleSquares(m.from, m.to)
			p.board[kingSquare] = fromPiece
			p.board[rookSquare] = toPiece
		}
	case Pawn:
		_, ranks := Diff(m.from, m.to)
		if ranks == 2 {
			p.enPassantSquare = BetweenSquares(m.from, m.to)[0]
		}
		// en passant
		if toPiece == NoPiece {
			victim := ToSquare(m.to.File(), m.from.Rank())
			p.board[victim] = NoPiece
		}
		if m.promote != NoPieceType {
			p.board[m.to] = CreatePiece(p.turn, m.promote)
		}
	}
	p.turn = p.turn.Flip()
	return p
}

func (mt *MoveTree) FindMoves(depth int, tt *TranspositionTable, f func(*MoveTree, int, *TranspositionTable) int) int {
	// if searched before
	cachedMoveTree := tt.Get(mt.position)
	if cachedMoveTree != nil {
		mt.legal = true
		if mt.parent != nil {
			mt.parent.legalMoves = append(mt.parent.legalMoves, mt.move)
		}
		mt.candidateMoves = cachedMoveTree.candidateMoves
	} else if mt.legal {
		mt.legalMoves = nil
		mt.follow = nil
		mt.eval = 0
	} else {
		p := mt.position
		moves := make([]Move, 0, 40)
		for _, square := range Squares {
			piece := p.board[square]
			if piece == NoPiece || piece.Colour() != p.turn {
				continue
			}
			pieceType := piece.Type()
			switch pieceType {
			case Bishop, Rook, Queen, Knight, King: // regular movements
				for _, m := range pieceMovements[pieceType] {
					curSquare := square
					for ok := true; ok; ok = (pieceType != Knight && pieceType != King) {
						newSquare, err := curSquare.Move(m)
						if err != nil {
							break
						}
						// check if parent castling is valid
						increment := 0
						if mt.move.castle != NoCastle {
							rank := mt.move.from.Rank()
							if newSquare.Rank() == rank {
								// TODO: optimize checking squares by avoiding loop
								if mt.move.castle == ASide {
									increment = -1
								} else {
									increment = 1
								}
								for f := mt.move.from.File(); ; f = File(int(f) + increment) {
									passSquare := ToSquare(f, rank)
									if passSquare == newSquare {
										mt.legal = false
										return 0
									}
									passPiece := p.board[passSquare]
									if passPiece != NoPiece && passPiece.Type() == King {
										break
									}
								}
							}
						}

						curSquare = newSquare
						pieceAt := p.board[curSquare]
						if pieceAt == NoPiece {
							moves = append(moves, Move{
								from: square, to: curSquare,
							})
							continue
						} else if pieceAt.Colour() != p.turn {
							if pieceAt.Type() == King {
								mt.legal = false
								return 0
							}
							moves = append(moves, Move{
								from: square, to: curSquare, capture: true,
							})
						}
						break
					}
				}
			case Pawn:
				pi := pawnInfo[p.turn]
				// one square forward
				oneSquare, err := square.Move(Movement{0, pi.forward})
				if err != nil {
					panic("Pawn can somehow step forward off the board")
				}
				if p.board[oneSquare] == NoPiece {
					// promotion
					if oneSquare.Rank() == pi.promotionRank {
						for _, promotionPieceType := range []PieceType{Queen, Rook, Bishop, Knight} {
							moves = append(moves, Move{
								from: square, to: oneSquare, promote: promotionPieceType,
							})
						}
					} else {
						moves = append(moves, Move{
							from: square, to: oneSquare,
						})
						// two squares forward
						if square.Rank() == pi.homeRank {
							twoSquare, err := square.Move(Movement{0, pi.forward * 2})
							if err != nil {
								panic("Pawn can somehow step forward off the board")
							}
							if p.board[twoSquare] == NoPiece {
								moves = append(moves, Move{
									from: square, to: twoSquare,
								})
							}
						}
					}
				}
				// capture
				for _, m := range []Movement{Movement{-1, pi.forward}, Movement{1, pi.forward}} {
					captureSquare, err := square.Move(m)
					if err != nil {
						continue
					}
					// check if parent castling is valid
					increment := 0
					if mt.move.castle != NoCastle {
						rank := mt.move.from.Rank()
						if captureSquare.Rank() == rank {
							// TODO: optimize checking squares by avoiding loop
							if mt.move.castle == ASide {
								increment = -1
							} else {
								increment = 1
							}
							for f := mt.move.from.File(); ; f = File(int(f) + increment) {
								passSquare := ToSquare(f, rank)
								if passSquare == captureSquare {
									mt.legal = false
									return 0
								}
								passPiece := p.board[passSquare]
								if passPiece != NoPiece && passPiece.Type() == King {
									break
								}
							}
						}
					}
					pieceAt := p.board[captureSquare]
					if pieceAt == NoPiece {
						// en passant
						if captureSquare == p.enPassantSquare {
							moves = append(moves, Move{
								from: square, to: captureSquare, capture: true,
							})
						}
						continue
					}
					if pieceAt.Colour() != p.turn {
						if pieceAt.Type() == King {
							mt.legal = false
							return 0
						}
						// promotion
						if captureSquare.Rank() == pi.promotionRank {
							for _, promotionPieceType := range []PieceType{Queen, Rook, Bishop, Knight} {
								moves = append(moves, Move{
									from: square, to: captureSquare, capture: true, promote: promotionPieceType,
								})
							}
						} else {
							moves = append(moves, Move{
								from: square, to: captureSquare, capture: true,
							})
						}
					}
				}
			}
			if pieceType == King {
				if !p.kingMoved(p.turn) {
					for _, rookSquare := range p.rookSquares {
						if rookSquare == NoSquare || p.board[rookSquare].Colour() != p.turn {
							continue
						}
						blocked := false
						for _, betweenSquare := range BetweenSquares(square, rookSquare) {
							if p.board[betweenSquare] != NoPiece {
								blocked = true
								break
							}
						}
						if blocked {
							continue
						}
						kingToSquare, rookToSquare := GetCastleSquares(square, rookSquare)
						if p.board[kingToSquare] != NoPiece && kingToSquare != square && kingToSquare != rookSquare {
							continue
						}
						if p.board[rookToSquare] != NoPiece && rookToSquare != square && rookToSquare != rookSquare {
							continue
						}
						kingBetweenSquares := BetweenSquares(square, kingToSquare)
						for _, betweenSquare := range kingBetweenSquares {
							if p.board[betweenSquare] != NoPiece && betweenSquare != rookSquare {
								blocked = true
								break
							}
						}
						if blocked {
							continue
						}

						castleDirection := ASide
						if rookSquare.File() > square.File() {
							castleDirection = HSide
						}

						moves = append(moves, Move{
							from: square, to: rookSquare, castle: castleDirection,
						})
					}
				}
			}
		}
		mt.legal = true
		if mt.parent != nil {
			mt.parent.legalMoves = append(mt.parent.legalMoves, mt.move)
		}
		mt.candidateMoves = moves
		tt.Add(mt)
	}
	rv := f(mt, depth, tt)
	if len(mt.legalMoves) == 0 && depth > 0 {
		// check if stalemate or checkmate
		lastMove := mt.position.turn.Flip()
		testMt := new(MoveTree)
		testMt.position = mt.position
		// act if we skipped a turn--can opponent take the king?
		testMt.position.turn = lastMove
		testMt.FindMoves(0, tt,
			func(_ *MoveTree, _ int, _ *TranspositionTable) int {
				return 0
			})
		if testMt.legal {
			mt.state = Stalemate
		} else {
			mt.state = WinFor(lastMove)
		}
	}

	return rv
}

func (root *MoveTree) FindAllMoves(startDepth int) int {
	var c func(*MoveTree, int, *TranspositionTable) int
	c = func(mt *MoveTree, depth int, tt *TranspositionTable) int {
		if depth == 0 {
			return 1
		}
		nodes := 0
		for _, move := range mt.candidateMoves {
			child := new(MoveTree)
			child.parent = mt
			child.move = move
			child.position = mt.position.ProcessMove(move)
			childNodes := child.FindMoves(depth-1, tt, c)
			if child.legal {
				nodes += childNodes
			}
		}
		return nodes
	}
	tt := NewTranspositionTable(TTMaxSize)
	return root.FindMoves(startDepth, tt, c)
}

func (root *MoveTree) Peek() {
	var f func(*MoveTree, int, *TranspositionTable) int
	f = func(mt *MoveTree, depth int, tt *TranspositionTable) int {
		if depth == 0 {
			return 0
		}
		for _, move := range mt.candidateMoves {
			child := new(MoveTree)
			child.parent = mt
			child.move = move
			child.position = mt.position.ProcessMove(move)
			child.FindMoves(depth-1, tt, f)
		}
		return 0
	}
	root.FindMoves(1, nil, f)
}
