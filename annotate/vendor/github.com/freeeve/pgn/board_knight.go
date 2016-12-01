package pgn

func (b Board) findAttackingKnight(pos Position, color Color, check bool) (Position, error) {
	count := 0
	r := pos.GetRank()
	f := pos.GetFile()
	retPos := NoPosition

	testPos := PositionFromFileRank(f+1, r+2)
	if testPos != NoPosition &&
		b.checkKnightColor(testPos, color) &&
		(!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
		count++
		retPos = testPos
	}

	testPos = PositionFromFileRank(f+1, r-2)
	if testPos != NoPosition &&
		b.checkKnightColor(testPos, color) &&
		(!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
		count++
		retPos = testPos
	}

	testPos = PositionFromFileRank(f+2, r+1)
	if testPos != NoPosition &&
		b.checkKnightColor(testPos, color) &&
		(!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
		count++
		retPos = testPos
	}

	testPos = PositionFromFileRank(f+2, r-1)
	if testPos != NoPosition &&
		b.checkKnightColor(testPos, color) &&
		(!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
		count++
		retPos = testPos
	}

	testPos = PositionFromFileRank(f-2, r-1)
	if testPos != NoPosition &&
		b.checkKnightColor(testPos, color) &&
		(!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
		count++
		retPos = testPos
	}

	testPos = PositionFromFileRank(f-2, r+1)
	if testPos != NoPosition &&
		b.checkKnightColor(testPos, color) &&
		(!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
		count++
		retPos = testPos
	}

	testPos = PositionFromFileRank(f-1, r-2)
	if testPos != NoPosition &&
		b.checkKnightColor(testPos, color) &&
		(!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
		count++
		retPos = testPos
	}

	testPos = PositionFromFileRank(f-1, r+2)
	if testPos != NoPosition &&
		b.checkKnightColor(testPos, color) &&
		(!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
		count++
		retPos = testPos
	}

	if count > 1 {
		return NoPosition, ErrAmbiguousMove
	}
	if count == 0 {
		return NoPosition, ErrAttackerNotFound
	}
	return retPos, nil
}

func (b Board) findAttackingKnightFromFile(pos Position, color Color, file File) (Position, error) {
	//fmt.Println("finding attacking knight from file:", string(file))
	count := 0
	r := pos.GetRank()
	f := pos.GetFile()
	retPos := NoPosition

	if f+1 == file {
		testPos := PositionFromFileRank(f+1, r+2)
		if testPos != NoPosition && b.checkKnightColor(testPos, color) {
			count++
			retPos = testPos
		}

		testPos = PositionFromFileRank(f+1, r-2)
		if testPos != NoPosition && b.checkKnightColor(testPos, color) {
			count++
			retPos = testPos
		}
	}

	if f+2 == file {
		testPos := PositionFromFileRank(f+2, r+1)
		if testPos != NoPosition && b.checkKnightColor(testPos, color) {
			count++
			retPos = testPos
		}

		testPos = PositionFromFileRank(f+2, r-1)
		if testPos != NoPosition && b.checkKnightColor(testPos, color) {
			count++
			retPos = testPos
		}
	}

	if f-2 == file {
		testPos := PositionFromFileRank(f-2, r-1)
		if testPos != NoPosition && b.checkKnightColor(testPos, color) {
			count++
			retPos = testPos
		}

		testPos = PositionFromFileRank(f-2, r+1)
		if testPos != NoPosition && b.checkKnightColor(testPos, color) {
			count++
			retPos = testPos
		}
	}

	if f-1 == file {
		testPos := PositionFromFileRank(f-1, r-2)
		if testPos != NoPosition && b.checkKnightColor(testPos, color) {
			count++
			retPos = testPos
		}

		testPos = PositionFromFileRank(f-1, r+2)
		if testPos != NoPosition && b.checkKnightColor(testPos, color) {
			count++
			retPos = testPos
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

func (b Board) findAttackingKnightFromRank(pos Position, color Color, rank Rank) (Position, error) {
	//fmt.Println("finding attacking knight from rank:", string(file))
	count := 0
	r := pos.GetRank()
	f := pos.GetFile()
	retPos := NoPosition

	if r+2 == rank {
		testPos := PositionFromFileRank(f+1, r+2)
		if testPos != NoPosition && b.checkKnightColor(testPos, color) {
			count++
			retPos = testPos
		}

		testPos = PositionFromFileRank(f-1, r+2)
		if testPos != NoPosition && b.checkKnightColor(testPos, color) {
			count++
			retPos = testPos
		}
	}

	if r+1 == rank {
		testPos := PositionFromFileRank(f+2, r+1)
		if testPos != NoPosition && b.checkKnightColor(testPos, color) {
			count++
			retPos = testPos
		}

		testPos = PositionFromFileRank(f-2, r+1)
		if testPos != NoPosition && b.checkKnightColor(testPos, color) {
			count++
			retPos = testPos
		}
	}

	if r-1 == rank {
		testPos := PositionFromFileRank(f-2, r-1)
		if testPos != NoPosition && b.checkKnightColor(testPos, color) {
			count++
			retPos = testPos
		}

		testPos = PositionFromFileRank(f+2, r-1)
		if testPos != NoPosition && b.checkKnightColor(testPos, color) {
			count++
			retPos = testPos
		}
	}

	if r-2 == rank {
		testPos := PositionFromFileRank(f-1, r-2)
		if testPos != NoPosition && b.checkKnightColor(testPos, color) {
			count++
			retPos = testPos
		}

		testPos = PositionFromFileRank(f+1, r-2)
		if testPos != NoPosition && b.checkKnightColor(testPos, color) {
			count++
			retPos = testPos
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

func (b Board) checkKnightColor(pos Position, color Color) bool {
	return (b.GetPiece(pos) == WhiteKnight && color == White) ||
		(b.GetPiece(pos) == BlackKnight && color == Black)
}
