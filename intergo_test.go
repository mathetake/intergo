// put bench marks on implemented algorithms
package intergo_test

import (
	"testing"

	"fmt"

	"github.com/mathetake/intergo"
	"github.com/mathetake/intergo/bm"
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
	inputRankingItemNum       int
	interleavedRankingItemNum int
}

var fixtures = []fixture{
	{inputRankingItemNum: 10, interleavedRankingItemNum: 5},
	{inputRankingItemNum: 200, interleavedRankingItemNum: 50},
	{inputRankingItemNum: 200, interleavedRankingItemNum: 200},
	{inputRankingItemNum: 1000, interleavedRankingItemNum: 200},
}

var inputRankingNumToBenchFunction = map[int]func(fx fixture, il intergo.Interleaving, b *testing.B){
	2:  benchmarkWith2Input,
	5:  benchmarkWith5Input,
	10: benchmarkWith10Input,
}

func BenchmarkMultileaving(b *testing.B) {
	for inputNum, benchmark := range inputRankingNumToBenchFunction {
		for n, fx := range fixtures {
			fxx := fx
			b.ReportAllocs()
			fmt.Println("")
			fmt.Printf(
				"inputRankingNum: %d, inputRankingItemNum: %d, interleavedRankingItemNum: %d\n",
				inputNum, fxx.inputRankingItemNum, fxx.interleavedRankingItemNum,
			)

			b.Run(fmt.Sprintf("[[%d-th bench on Team Draft Multileaving]", n), func(b *testing.B) {
				benchmark(fxx, &tdm.TeamDraftMultileaving{}, b)
			})

			b.Run(fmt.Sprintf("[[%d-th bench on Balanced Multileaving]]", n), func(b *testing.B) {
				benchmark(fxx, &bm.BalancedMultileaving{}, b)
			})

			for _, samplingSize := range []int{2, 10, 50, 100} {
				b.Run(fmt.Sprintf("[[%d-th bench on Optimized Multileaving with sampling size: %d]]", n, samplingSize), func(b *testing.B) {
					benchmark(fxx, &om.OptimizedMultiLeaving{samplingSize}, b)
				})
			}
			fmt.Println("")
		}
	}
}

func benchmarkWith2Input(fx fixture, il intergo.Interleaving, b *testing.B) {
	rks := getRankings(fx, 2)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		il.GetInterleavedRanking(fx.interleavedRankingItemNum, rks[0], rks[1])
	}
}

func benchmarkWith5Input(fx fixture, il intergo.Interleaving, b *testing.B) {
	rks := getRankings(fx, 5)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		il.GetInterleavedRanking(fx.interleavedRankingItemNum, rks[0], rks[1], rks[2], rks[3], rks[4], rks[4])
	}
}

func benchmarkWith10Input(fx fixture, il intergo.Interleaving, b *testing.B) {
	rks := getRankings(fx, 10)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		il.GetInterleavedRanking(
			fx.interleavedRankingItemNum, rks[0], rks[1], rks[2], rks[3],
			rks[4], rks[5], rks[6], rks[7], rks[8], rks[9],
		)
	}
}

func getRankings(fx fixture, inputRankingNum int) []tRanking {
	rks := make([]tRanking, inputRankingNum)
	for i := 0; i < inputRankingNum; i++ {
		rk := tRanking{}
		for j := 0; j < fx.inputRankingItemNum; j++ {
			rk = append(rk, i*fx.inputRankingItemNum+j)
		}
		rks[i] = rk
	}
	return rks
}
