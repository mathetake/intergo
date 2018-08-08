# intergo

A package for interleaving / multileaving ranking generation in go

It is mainly tailored to be used for generating interleaved or multileaved ranking based on the following algorithm

- Team Draft Interleaving/Multileaving (in `github.com/mathetake/itergo/tdm` package)
- Probabilistic Team Draft Interleaving/Multileaving (TODO)
- Balanced Interleaving/Multileaving

# How to use

Note that all of your ranking satisfy the `intergo.Ranking` interface

```go
type Ranking interface {
	GetIDByIndex(int) interface{}
	Len() int
}
```

which is used for removing duplications in the list.

Anyway, the following example is self-explanatory:

```go
package main

import (
	"fmt"

	"github.com/mathetake/intergo"
	"github.com/mathetake/intergo/tdm"
)

type tRanking []int

func (rk tRanking) GetIDByIndex(i int) interface{} {
	return rk[i]
}

func (rk tRanking) Len() int {
	return len(rk)
}

var _ intergo.Ranking = tRanking{}

func main() {
	TDM := &tdm.TeamDraftMultileaving{}
	RankingA := tRanking{1, 2, 3, 4, 5,}
	RankingB := tRanking{10, 20, 30, 40, 50}

	idxToRk := map[int]tRanking{
		0: RankingA,
		1: RankingB,
	}

	res := TDM.GetInterleavedRanking(4, RankingA, RankingB)
	iRanking := tRanking{}
	for _, it := range res {
		iRanking = append(iRanking, idxToRk[it.RankingIDx][it.ItemIDx])
	}

	fmt.Println("Result: ", iRanking)
}

```

# References

1. Schuth, Anne, et al. "Multileaved comparisons for fast online evaluation." Proceedings of the 23rd ACM International Conference on Conference on Information and Knowledge Management. ACM, 2014.

2. Radlinski, Filip, Madhu Kurup, and Thorsten Joachims. "How does clickthrough data reflect retrieval quality?." Proceedings of the 17th ACM conference on Information and knowledge management. ACM, 2008.


# license

MIT

