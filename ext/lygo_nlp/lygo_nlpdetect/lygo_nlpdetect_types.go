package lygo_nlpdetect

type tuple struct {
	Code  string
	Count float64
}

type tuples []tuple

func (t tuples) Len() int {
	return len(t)
}

func (t tuples) Less(i, j int) bool {
	return t[i].Count < t[j].Count
}

func (t tuples) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

type languages map[string]trigrams
type trigrams map[string]int
