package bm_test

import (
	"testing"

	"github.com/mathetake/intergo"

	"gotest.tools/assert"

	"fmt"

	"github.com/mathetake/intergo/bm"
)

type tRanking []int

func (rk tRanking) GetIDByIndex(i int) interface{} {
	return rk[i]
}

func (rk tRanking) Len() int {
	return len(rk)
}

var _ intergo.Ranking = tRanking{}

func TestBalancedMultileaving(t *testing.T) {
	bMultileaving := &bm.BalancedMultileaving{}

	cases := []struct {
		inputRks         []intergo.Ranking
		num              int
		expectedPatterns [][]intergo.Res
	}{
		{
			inputRks: []intergo.Ranking{
				tRanking{1, 2, 3, 4, 5},
				tRanking{10, 20, 30, 40, 50},
			},
			num: 2,
			expectedPatterns: [][]intergo.Res{
				{
					intergo.Res{RankingIDx: 0, ItemIDx: 0},
					intergo.Res{RankingIDx: 1, ItemIDx: 0},
				},
				{
					intergo.Res{RankingIDx: 0, ItemIDx: 0},
					intergo.Res{RankingIDx: 0, ItemIDx: 1},
				},
				{
					intergo.Res{RankingIDx: 1, ItemIDx: 0},
					intergo.Res{RankingIDx: 0, ItemIDx: 0},
				},
				{
					intergo.Res{RankingIDx: 1, ItemIDx: 0},
					intergo.Res{RankingIDx: 1, ItemIDx: 1},
				},
			},
		},
		{
			inputRks: []intergo.Ranking{
				tRanking{1, 2, 3, 4, 5},
				tRanking{1, 20, 30, 40, 50},
			},
			num: 2,
			expectedPatterns: [][]intergo.Res{
				{
					intergo.Res{RankingIDx: 0, ItemIDx: 0},
					intergo.Res{RankingIDx: 1, ItemIDx: 1},
				},
				{
					intergo.Res{RankingIDx: 0, ItemIDx: 0},
					intergo.Res{RankingIDx: 0, ItemIDx: 1},
				},
				{
					intergo.Res{RankingIDx: 1, ItemIDx: 0},
					intergo.Res{RankingIDx: 0, ItemIDx: 1},
				},
				{
					intergo.Res{RankingIDx: 1, ItemIDx: 0},
					intergo.Res{RankingIDx: 1, ItemIDx: 1},
				},
			},
		},
		{
			inputRks: []intergo.Ranking{
				tRanking{1, 2, 3, 4, 5},
				tRanking{1, 1, 30, 40, 50},
			},
			num: 2,
			expectedPatterns: [][]intergo.Res{
				{
					intergo.Res{RankingIDx: 0, ItemIDx: 0},
					intergo.Res{RankingIDx: 1, ItemIDx: 2},
				},
				{
					intergo.Res{RankingIDx: 0, ItemIDx: 0},
					intergo.Res{RankingIDx: 0, ItemIDx: 1},
				},
				{
					intergo.Res{RankingIDx: 1, ItemIDx: 0},
					intergo.Res{RankingIDx: 0, ItemIDx: 1},
				},
				{
					intergo.Res{RankingIDx: 1, ItemIDx: 0},
					intergo.Res{RankingIDx: 1, ItemIDx: 2},
				},
			},
		},
	}

	for n, tc := range cases {
		tcc := tc
		t.Run(fmt.Sprintf("%d-th unit test", n), func(t *testing.T) {
			actual, _ := bMultileaving.GetInterleavedRanking(tcc.num, tcc.inputRks...)
			t.Log("actual: ", actual)
			assert.Equal(t, true, len(actual) <= tcc.num)

			var isExpected = false
			for _, expected := range tcc.expectedPatterns {

				var isExpectedPattern = true
				for i := 0; i < tcc.num; i++ {
					if actual[i] != expected[i] {
						isExpectedPattern = false
					}
				}

				if isExpectedPattern {
					isExpected = true
					break
				}
			}
			assert.Equal(t, true, isExpected)
		})
	}
}
