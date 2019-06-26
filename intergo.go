package intergo

// ID ... identifier of items
type ID interface{}

type Ranking interface {
	// GetIDByIndex ... allow algorithms to access items' identifier
	GetIDByIndex(int) ID

	// Len ... to get the "length" of the ranking
	Len() int
}

type Result struct {
	// RankingIndex ... represents to which ranking the item belongs
	RankingIndex int

	// ItemIndex ... represents the item's index in the ranking declared by RankingIDx
	ItemIndex int
}

type Interleaving interface {
	GetInterleavedRanking(int, ...Ranking) ([]Result, error)
}
