package intergo

type Ranking interface {
	GetIDByIndex(int) interface{}
	Len() int
}

type Res struct {
	// Ranking ... represents to which ranking the item belongs
	RankingIDx int

	// ItemIDx ... represents the item's index in the ranking declared by RankingIDx
	ItemIDx int
}

type Interleaving interface {
	GetInterleavedRanking(int, ...Ranking) []Res
}
