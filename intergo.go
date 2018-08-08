package intergo

type Items []Item

type Item interface {
	GetID() interface{}
}

type Res struct {
	// Ranking ... represents to which ranking the item belongs
	RankingIDx int

	// ItemIDx ... represents the item's index in the ranking declared by RankingIDx
	ItemIDx int
}

type Interleaving interface {
	GetInterleavedRanking([]Items, int) []Res
}
