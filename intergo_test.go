// put bench marks on implemented algorithms
package intergo_test

import (
	"testing"

	"fmt"

	"github.com/mathetake/intergo"
	"github.com/mathetake/intergo/om"
	"github.com/mathetake/intergo/tdm"
)

type tRanking []int

func (rk tRanking) GetIDByIndex(i int) interface{} {
	return rk[i]
}

func (rk tRanking) Len() int {
	return len(rk)
}

type fixture struct {
	inputRankingNum           int
	inputRankingItemNum       int
	interleavedRankingItemNUm int
}

func getCase(fx fixture) []tRanking {
	rks := make([]tRanking, fx.inputRankingNum)
	for i := 0; i < fx.inputRankingNum; i++ {
		rk := tRanking{}
		for j := 0; j < fx.inputRankingItemNum; j++ {
			rk = append(rk, i*fx.inputRankingNum+j)
		}
		rks = append(rks, rk)
	}
	return rks
}

func benchmark(fx fixture, il intergo.Interleaving, b *testing.B) {
	rks := getCase(fx)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		il.GetInterleavedRanking(fx.interleavedRankingItemNUm, rks[0], rks[1])
	}
}

var fixtures = []fixture{
	{inputRankingItemNum: 10, inputRankingNum: 2, interleavedRankingItemNUm: 5},
	{inputRankingItemNum: 100, inputRankingNum: 2, interleavedRankingItemNUm: 50},
	{inputRankingItemNum: 1000, inputRankingNum: 2, interleavedRankingItemNUm: 200},
	{inputRankingItemNum: 10, inputRankingNum: 5, interleavedRankingItemNUm: 5},
	{inputRankingItemNum: 100, inputRankingNum: 5, interleavedRankingItemNUm: 50},
	{inputRankingItemNum: 1000, inputRankingNum: 5, interleavedRankingItemNUm: 200},
	{inputRankingItemNum: 10, inputRankingNum: 10, interleavedRankingItemNUm: 5},
	{inputRankingItemNum: 100, inputRankingNum: 10, interleavedRankingItemNUm: 50},
	{inputRankingItemNum: 1000, inputRankingNum: 10, interleavedRankingItemNUm: 200},
}

func BenchmarkMultileaving(b *testing.B) {
	for n, fx := range fixtures {
		b.ReportAllocs()

		fxx := fx
		fmt.Println("")
		fmt.Printf(
			"inputRankingNum: %d, inputRankingItemNum: %d, interleavedRankingItemNUm: %d\n",
			fxx.inputRankingNum, fxx.inputRankingItemNum, fxx.interleavedRankingItemNUm,
		)

		b.Run(fmt.Sprintf("%d-th bench on Team Draft Multileaving", n), func(b *testing.B) {
			benchmark(fxx, &tdm.TeamDraftMultileaving{}, b)
		})

		for _, samplingSize := range []int{2, 10, 50, 100} {
			b.Run(fmt.Sprintf("%d-th bench on Optimized Multileaving with sampling size: %d", n, samplingSize), func(b *testing.B) {
				benchmark(fxx, &om.OptimizedMultiLeaving{samplingSize}, b)
			})
		}
		fmt.Println("")
	}
}
