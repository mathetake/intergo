package intergo

type Ranking interface {
	// GetIDByIndex ... allow algorithms to access items' identifier
	GetIDByIndex(int) interface{}

	// Len ... to get the "length" of the ranking
	Len() int
}

type Result struct {
	// Ranking ... represents to which ranking the item belongs
	RankingIndex int

	// ItemIndex ... represents the item's index in the ranking declared by RankingIndex
	ItemIndex int
}

type Interleaving interface {
	GetInterleavedRanking(int, ...Ranking) ([]Result, error)
}
