// In package intergo, some interfaces and types are defined for interleaving/multileaving algorithms implementation.
//
// Specific algorithms are implemented in its subpackages:
//
// http://godoc.org/github.com/mathetake/intergo/bm implements balanced multileaving algorithm.
//
// http://godoc.org/github.com/mathetake/intergo/om implements optimized multileaving algorithm.
//
// http://godoc.org/github.com/mathetake/intergo/tdm implements team draft multileaving algorithm.
//
// See README.md for more details.
package intergo

// type ID is used as identifier of items.
// The purpose is to remove item duplication in generated rankings.
type ID interface{}

// Ranking is the interface which all of target ranking should implement.
type Ranking interface {
	// GetIDByIndex allows interleaving/multileaving algorithms to access items' identifier
	GetIDByIndex(int) ID

	// Len is used to get the "length" of the ranking
	Len() int
}

// Result is the type of generated ranking's each entity.
type Result struct {
	// RankingIndex represents to which ranking the item belongs
	RankingIndex int

	// ItemIndex represents the item's index in the ranking declared by RankingIndex
	ItemIndex int
}

// Interleaving is the interface which every interleaving/multileaving algorithm should implement.
type Interleaving interface {
	// GetInterleavedRanking is intended to be used for ranking generation.
	//
	// First argument "num" is the expected length of a resulted ranking. Ideally,
	// if the total number of unique items in given rankings is greater than equal "num",
	// the resulted length should equal "num". However, it depends on the implementations.
	//
	// "rankings" should be your rankings from which you want to generate a multileaved ranking.
	GetInterleavedRanking(num int, rankings ...Ranking) ([]*Result, error)
}
