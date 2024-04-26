package types

func DeepEqual(a, b MalValue) Bool {
	seqA, seqErrA := Sequence(a)
	seqB, seqErrB := Sequence(b)

	if seqErrA != nil && seqErrB != nil {
		return a == b
	}

	if seqErrA != nil || seqErrB != nil || len(seqA) != len(seqB) {
		return false
	}

	for i := range seqA {
		if !DeepEqual(seqA[i], seqB[i]) {
			return false
		}
	}

	return true
}
