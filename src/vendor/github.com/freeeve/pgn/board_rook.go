package pgn

func (b Board) findAttackingRook(pos Position, color Color, check bool) (Position, error) {
	count := 0
	retPos := NoPosition

	r := pos.GetRank()
	f := pos.GetFile()
	for {
		f--
		testPos := PositionFromFileRank(f, r)
		if b.checkRookColor(testPos, color) && (!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
			retPos = testPos
			count++
			break
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	r = pos.GetRank()
	f = pos.GetFile()
	for {
		f++
		testPos := PositionFromFileRank(f, r)
		if b.checkRookColor(testPos, color) && (!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
			retPos = testPos
			count++
			break
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	r = pos.GetRank()
	f = pos.GetFile()
	for {
		r++
		testPos := PositionFromFileRank(f, r)
		if b.checkRookColor(testPos, color) && (!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
			retPos = testPos
			count++
			break
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	r = pos.GetRank()
	f = pos.GetFile()
	for {
		r--
		testPos := PositionFromFileRank(f, r)
		if b.checkRookColor(testPos, color) && (!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
			retPos = testPos
			count++
			break
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	if count > 1 {
		return NoPosition, ErrAmbiguousMove
	}
	if count == 0 {
		return NoPosition, ErrAttackerNotFound
	}
	return retPos, nil
}

func (b Board) findAttackingRookFromFile(pos Position, color Color, file File) (Position, error) {
	count := 0
	retPos := NoPosition

	r := pos.GetRank()
	f := pos.GetFile()
	for {
		f--
		testPos := PositionFromFileRank(f, r)
		if file == f && b.checkRookColor(testPos, color) {
			retPos = testPos
			count++
			break
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	r = pos.GetRank()
	f = pos.GetFile()
	for {
		f++
		testPos := PositionFromFileRank(f, r)
		if file == f && b.checkRookColor(testPos, color) {
			retPos = testPos
			count++
			break
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	if file == pos.GetFile() {
		r = pos.GetRank()
		f = pos.GetFile()
		for {
			r++
			testPos := PositionFromFileRank(f, r)
			if b.checkRookColor(testPos, color) {
				retPos = testPos
				count++
				break
			} else if testPos == NoPosition || b.containsPieceAt(testPos) {
				break
			}
		}

		r = pos.GetRank()
		f = pos.GetFile()
		for {
			r--
			testPos := PositionFromFileRank(f, r)
			if b.checkRookColor(testPos, color) {
				retPos = testPos
				count++
				break
			} else if testPos == NoPosition || b.containsPieceAt(testPos) {
				break
			}
		}
	}
	if count > 1 {
		return NoPosition, ErrAmbiguousMove
	}
	if count == 0 {
		return NoPosition, ErrAttackerNotFound
	}
	return retPos, nil
}

func (b Board) findAttackingRookFromRank(pos Position, color Color, rank Rank) (Position, error) {
	count := 0
	retPos := NoPosition

	r := pos.GetRank()
	f := pos.GetFile()
	for {
		r--
		testPos := PositionFromFileRank(f, r)
		if rank == r && b.checkRookColor(testPos, color) {
			retPos = testPos
			count++
			break
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	r = pos.GetRank()
	f = pos.GetFile()
	for {
		r++
		testPos := PositionFromFileRank(f, r)
		if rank == r && b.checkRookColor(testPos, color) {
			retPos = testPos
			count++
			break
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	if rank == pos.GetRank() {
		r := pos.GetRank()
		f := pos.GetFile()
		for {
			f--
			testPos := PositionFromFileRank(f, r)
			if b.checkRookColor(testPos, color) {
				retPos = testPos
				count++
				break
			} else if testPos == NoPosition || b.containsPieceAt(testPos) {
				break
			}
		}

		r = pos.GetRank()
		f = pos.GetFile()
		for {
			f++
			testPos := PositionFromFileRank(f, r)
			if b.checkRookColor(testPos, color) {
				retPos = testPos
				count++
				break
			} else if testPos == NoPosition || b.containsPieceAt(testPos) {
				break
			}
		}
	}

	if count > 1 {
		return NoPosition, ErrAmbiguousMove
	}
	if count == 0 {
		return NoPosition, ErrAttackerNotFound
	}
	return retPos, nil
}

func (b Board) checkRookColor(pos Position, color Color) bool {
	return (b.GetPiece(pos) == WhiteRook && color == White) ||
		(b.GetPiece(pos) == BlackRook && color == Black)
}
