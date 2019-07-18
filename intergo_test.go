// put bench marks on implemented algorithms
package intergo_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/mathetake/intergo"
	"github.com/mathetake/intergo/bm"
	"github.com/mathetake/intergo/gom"
	"github.com/mathetake/intergo/tdm"
)

type tRanking []int

func (rk tRanking) GetIDByIndex(i int) intergo.ID {
	return intergo.ID(strconv.Itoa(rk[i]))
}

func (rk tRanking) Len() int {
	return len(rk)
}

type fixture struct {
	inputRankingItemNum       int
	interleavedRankingItemNum int
}

var fixtures = []fixture{
	{inputRankingItemNum: 10, interleavedRankingItemNum: 5},
	{inputRankingItemNum: 200, interleavedRankingItemNum: 50},
	{inputRankingItemNum: 200, interleavedRankingItemNum: 200},
	{inputRankingItemNum: 1000, interleavedRankingItemNum: 200},
}

func BenchmarkMultileaving(b *testing.B) {
	for _, inputRankingNum := range []int{2, 10, 50, 100} {
		for n, fx := range fixtures {
			fxx := fx
			b.ReportAllocs()
			fmt.Println("")
			fmt.Printf(
				"inputRankingNum: %d, inputRankingItemNum: %d, interleavedRankingItemNum: %d\n",
				inputRankingNum, fxx.inputRankingItemNum, fxx.interleavedRankingItemNum,
			)

			b.Run(fmt.Sprintf("[[%d-th bench on Team Draft Multileaving]]", n), func(b *testing.B) {
				benchmarkInputNum(fxx, inputRankingNum, &tdm.TeamDraftMultileaving{}, b)
			})

			b.Run(fmt.Sprintf("[[%d-th bench on Balanced Multileaving]]", n), func(b *testing.B) {
				benchmarkInputNum(fxx, inputRankingNum, &bm.BalancedMultileaving{}, b)
			})

			for _, samplingSize := range []int{2, 10, 50, 100} {
				b.Run(fmt.Sprintf("[[%d-th bench on Greedy Optimized Multileaving with sampling size: %d]]", n, samplingSize), func(b *testing.B) {
					benchmarkInputNum(fxx, inputRankingNum, &gom.GreedyOptimizedMultiLeaving{
						NumSampling: samplingSize, CreditLabel: 0, Alpha: 0,
					}, b)
				})
			}
			fmt.Println("")
		}
	}
}

func benchmarkInputNum(fx fixture, inputRankingNum int, il intergo.Interleaving, b *testing.B) {
	rks := getRankings(fx, inputRankingNum)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		il.GetInterleavedRanking(fx.interleavedRankingItemNum, rks...)
	}
}

func getRankings(fx fixture, inputRankingNum int) []intergo.Ranking {
	rks := make([]intergo.Ranking, inputRankingNum)
	for i := 0; i < inputRankingNum; i++ {
		rk := tRanking{}
		for j := 0; j < fx.inputRankingItemNum; j++ {
			rk = append(rk, i*fx.inputRankingItemNum+j)
		}
		rks[i] = rk
	}
	return rks
}
