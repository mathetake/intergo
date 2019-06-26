# intergo 
[![CircleCI](https://circleci.com/gh/mathetake/intergo.svg?style=shield&circle-token=89a8a65229dd121bd61be11222cdc2a0416cef22)](https://circleci.com/gh/mathetake/intergo)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)


A package for interleaving / multileaving ranking generation in go

It is mainly tailored to be used for generating interleaved or multileaved ranking based on the following algorithm

- Balanced Interleaving/Multileaving (in `github.com/mathetake/itergo/bm` package)
- Optimized Multileaving (in `github.com/mathetake/intergo/om` package)
- Team Draft Interleaving/Multileaving (in `github.com/mathetake/itergo/tdm` package)

__NOTE:__ this package aims only at generating a single combined ranking and does not implement the evaluation functions of the given rankings.

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
	rankingA := tRanking{1, 2, 3, 4, 5,}
	rankingB := tRanking{10, 20, 30, 40, 50}

	idxToRk := map[int]tRanking{
		0: rankingA,
		1: rankingB,
	}

	res, _ := TDM.GetInterleavedRanking(4, rankingA, rankingB)
	iRanking := tRanking{}
	for _, it := range res {
		iRanking = append(iRanking, idxToRk[it.RankingIndex][it.ItemIndex])
	}

	fmt.Println("Result: ", iRanking)
}
```

# References

1. Radlinski, Filip, Madhu Kurup, and Thorsten Joachims. "How does clickthrough data reflect retrieval quality?." Proceedings of the 17th ACM conference on Information and knowledge management. ACM, 2008.

2. Schuth, Anne, et al. "Multileaved comparisons for fast online evaluation." Proceedings of the 23rd ACM International Conference on Conference on Information and Knowledge Management. ACM, 2014.

3. Manabe, Tomohiro, et al. "A comparative live evaluation of multileaving methods on a commercial cqa search." Proceedings of the 40th International ACM SIGIR Conference on Research and Development in Information Retrieval. ACM, 2017.


# license

MIT

