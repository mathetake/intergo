package tdm_test

import (
	"testing"

	"github.com/mathetake/intergo"
	"github.com/mathetake/intergo/tdm"

	"gotest.tools/assert"

	"fmt"
)

type tRanking []int

func (rk tRanking) GetIDByIndex(i int) interface{} {
	return rk[i]
}

func (rk tRanking) Len() int {
	return len(rk)
}

var _ intergo.Ranking = tRanking{}

func TestGetInterleavedRanking(t *testing.T) {
	TDM := &tdm.TeamDraftMultileaving{}

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
					intergo.Res{RankingIDx: 1, ItemIDx: 0},
					intergo.Res{RankingIDx: 0, ItemIDx: 0},
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
					intergo.Res{RankingIDx: 1, ItemIDx: 0},
					intergo.Res{RankingIDx: 0, ItemIDx: 1},
				},
			},
		},
		{
			inputRks: []intergo.Ranking{
				tRanking{1, 2, 3, 4, 5},
				tRanking{1, 20, 30, 40, 50},
			},
			num: 3,
			expectedPatterns: [][]intergo.Res{
				{
					intergo.Res{RankingIDx: 0, ItemIDx: 0},
					intergo.Res{RankingIDx: 1, ItemIDx: 1},
					intergo.Res{RankingIDx: 0, ItemIDx: 1},
				},
				{
					intergo.Res{RankingIDx: 1, ItemIDx: 0},
					intergo.Res{RankingIDx: 0, ItemIDx: 1},
					intergo.Res{RankingIDx: 1, ItemIDx: 1},
				},
				{
					intergo.Res{RankingIDx: 1, ItemIDx: 0},
					intergo.Res{RankingIDx: 0, ItemIDx: 1},
					intergo.Res{RankingIDx: 0, ItemIDx: 2},
				},
				{
					intergo.Res{RankingIDx: 0, ItemIDx: 0},
					intergo.Res{RankingIDx: 1, ItemIDx: 1},
					intergo.Res{RankingIDx: 1, ItemIDx: 2},
				},
			},
		},
	}

	for n, tc := range cases {
		tcc := tc
		t.Run(fmt.Sprintf("%d-th unit test", n), func(t *testing.T) {
			actual := TDM.GetInterleavedRanking(tcc.num, tcc.inputRks...)
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

func TestGetRandomKey(t *testing.T) {
	tc := struct {
		numCandidate int
		numSelection int
		threshold    float64
	}{
		numCandidate: 10,
		numSelection: 10000000,
		threshold:    10e-4,
	}
	input := map[int]interface{}{}
	for i := 0; i < tc.numCandidate; i++ {
		input[i] = true
	}

	chosenRatio := map[int]float64{}

	for i := 0; i < tc.numSelection; i++ {
		chosenRatio[tdm.ExportedGetRandomKey(input)] += 1 / float64(tc.numSelection)
	}

	t.Log(chosenRatio)

	for _, v := range chosenRatio {
		diff := v - 1/float64(tc.numCandidate)
		if diff < 0 {
			diff = -diff
		}
		assert.Equal(t, true, diff < tc.threshold)
	}
}
