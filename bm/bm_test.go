package bm_test

import (
	"fmt"
	"testing"

	"github.com/mathetake/intergo"
	"github.com/mathetake/intergo/bm"
	"gotest.tools/assert"
)

type tRanking []int

func (rk tRanking) GetIDByIndex(i int) intergo.ID {
	return rk[i]
}

func (rk tRanking) Len() int {
	return len(rk)
}

var _ intergo.Ranking = tRanking{}

func TestBalancedMultileaving(t *testing.T) {
	b := &bm.BalancedMultileaving{}

	cases := []struct {
		inputRks         []intergo.Ranking
		num              int
		expectedPatterns [][]intergo.Result
	}{
		{
			inputRks: []intergo.Ranking{
				tRanking{1, 2, 3, 4, 5},
				tRanking{10, 20, 30, 40, 50},
			},
			num: 2,
			expectedPatterns: [][]intergo.Result{
				{
					intergo.Result{RankingIndex: 0, ItemIndex: 0},
					intergo.Result{RankingIndex: 1, ItemIndex: 0},
				},
				{
					intergo.Result{RankingIndex: 0, ItemIndex: 0},
					intergo.Result{RankingIndex: 0, ItemIndex: 1},
				},
				{
					intergo.Result{RankingIndex: 1, ItemIndex: 0},
					intergo.Result{RankingIndex: 0, ItemIndex: 0},
				},
				{
					intergo.Result{RankingIndex: 1, ItemIndex: 0},
					intergo.Result{RankingIndex: 1, ItemIndex: 1},
				},
			},
		},
		{
			inputRks: []intergo.Ranking{
				tRanking{1, 2, 3, 4, 5},
				tRanking{1, 20, 30, 40, 50},
			},
			num: 2,
			expectedPatterns: [][]intergo.Result{
				{
					intergo.Result{RankingIndex: 0, ItemIndex: 0},
					intergo.Result{RankingIndex: 1, ItemIndex: 1},
				},
				{
					intergo.Result{RankingIndex: 0, ItemIndex: 0},
					intergo.Result{RankingIndex: 0, ItemIndex: 1},
				},
				{
					intergo.Result{RankingIndex: 1, ItemIndex: 0},
					intergo.Result{RankingIndex: 0, ItemIndex: 1},
				},
				{
					intergo.Result{RankingIndex: 1, ItemIndex: 0},
					intergo.Result{RankingIndex: 1, ItemIndex: 1},
				},
			},
		},
		{
			inputRks: []intergo.Ranking{
				tRanking{1, 2, 3, 4, 5},
				tRanking{1, 1, 30, 40, 50},
			},
			num: 2,
			expectedPatterns: [][]intergo.Result{
				{
					intergo.Result{RankingIndex: 0, ItemIndex: 0},
					intergo.Result{RankingIndex: 1, ItemIndex: 2},
				},
				{
					intergo.Result{RankingIndex: 0, ItemIndex: 0},
					intergo.Result{RankingIndex: 0, ItemIndex: 1},
				},
				{
					intergo.Result{RankingIndex: 1, ItemIndex: 0},
					intergo.Result{RankingIndex: 0, ItemIndex: 1},
				},
				{
					intergo.Result{RankingIndex: 1, ItemIndex: 0},
					intergo.Result{RankingIndex: 1, ItemIndex: 2},
				},
			},
		},
	}

	for n, tc := range cases {
		tcc := tc
		t.Run(fmt.Sprintf("%d-th unit test", n), func(t *testing.T) {
			actual, _ := b.GetInterleavedRanking(tcc.num, tcc.inputRks...)
			t.Log("actual: ", actual)
			assert.Equal(t, true, len(actual) <= tcc.num)

			var isExpected bool
			for _, expected := range tcc.expectedPatterns {

				var isExpectedPattern = true
				for i := 0; i < tcc.num; i++ {
					if *actual[i] != expected[i] {
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
