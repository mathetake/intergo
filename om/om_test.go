package om_test

import (
	"testing"

	"github.com/mathetake/intergo"
	"github.com/mathetake/intergo/om"

	"gotest.tools/assert"

	"fmt"

	"gonum.org/v1/gonum/mat"
)

type tRanking []int

func (rk tRanking) GetIDByIndex(i int) interface{} {
	return rk[i]
}

func (rk tRanking) Len() int {
	return len(rk)
}

var _ intergo.Ranking = tRanking{}

func TestOptimizedMultileaving(t *testing.T) {
	oM := &om.OptimizedMultiLeaving{}

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
					intergo.Res{RankingIDx: 0, ItemIDx: 0},
					intergo.Res{RankingIDx: 0, ItemIDx: 1},
					intergo.Res{RankingIDx: 1, ItemIDx: 1},
				},
				{
					intergo.Res{RankingIDx: 0, ItemIDx: 0},
					intergo.Res{RankingIDx: 0, ItemIDx: 1},
					intergo.Res{RankingIDx: 0, ItemIDx: 2},
				},
				{
					intergo.Res{RankingIDx: 0, ItemIDx: 0},
					intergo.Res{RankingIDx: 1, ItemIDx: 1},
					intergo.Res{RankingIDx: 1, ItemIDx: 2},
				},
				{
					intergo.Res{RankingIDx: 1, ItemIDx: 0},
					intergo.Res{RankingIDx: 1, ItemIDx: 1},
					intergo.Res{RankingIDx: 0, ItemIDx: 1},
				},
				{
					intergo.Res{RankingIDx: 1, ItemIDx: 0},
					intergo.Res{RankingIDx: 1, ItemIDx: 1},
					intergo.Res{RankingIDx: 1, ItemIDx: 2},
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
			},
		},
	}

	for n, tc := range cases {
		tcc := tc
		t.Run(fmt.Sprintf("%d-th unit test", n), func(t *testing.T) {
			actual := oM.ExportedPrefixConstraintSampling(tcc.num, tcc.inputRks...)
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

func TestCalcInsensitivity(t *testing.T) {
	o := &om.OptimizedMultiLeaving{}

	cases := []struct {
		inputRankings    []intergo.Ranking
		combinedRankings [][]intergo.Res
		expected         []float64
		threshold        float64
	}{
		{
			inputRankings: []intergo.Ranking{
				tRanking{1, 2, 3, 4, 5},
				tRanking{10, 20, 30, 40, 50},
			},
			combinedRankings: [][]intergo.Res{
				{
					intergo.Res{RankingIDx: 0, ItemIDx: 0},
					intergo.Res{RankingIDx: 1, ItemIDx: 0},
				},
				{
					intergo.Res{RankingIDx: 0, ItemIDx: 0},
					intergo.Res{RankingIDx: 0, ItemIDx: 1},
				},
			},
			expected:  []float64{0.2222222222222222, 0.5},
			threshold: 10e-7,
		},
		{
			inputRankings: []intergo.Ranking{
				tRanking{1, 2, 3},
				tRanking{10, 20, 30},
			},
			combinedRankings: [][]intergo.Res{
				{
					intergo.Res{RankingIDx: 0, ItemIDx: 0},
					intergo.Res{RankingIDx: 1, ItemIDx: 0},
					intergo.Res{RankingIDx: 1, ItemIDx: 1},
				},
				{
					intergo.Res{RankingIDx: 0, ItemIDx: 0},
					intergo.Res{RankingIDx: 0, ItemIDx: 1},
					intergo.Res{RankingIDx: 0, ItemIDx: 2},
				},
			},
			expected:  []float64{0.17833719135802462, 0.4075038580246914},
			threshold: 10e-8,
		},
	}

	for n, tc := range cases {
		tcc := tc
		t.Run(fmt.Sprintf("%d-th unit test", n), func(t *testing.T) {
			actual := o.ExportedCalcInsensitivity(tcc.inputRankings, tcc.combinedRankings)
			assert.Equal(t, len(tcc.expected), len(actual))
			for i := range tcc.expected {
				diff := actual[i] - tcc.expected[i]
				if diff < 0 {
					diff = -diff
				}
				assert.Equal(t, true, diff < tcc.threshold)
			}
		})
	}
}

func TestGetConstraintMatrix(t *testing.T) {
	o := &om.OptimizedMultiLeaving{}

	cases := []struct {
		inputRankings    []intergo.Ranking
		combinedRankings [][]intergo.Res
		expected         []float64
		threshold        float64
	}{
		{
			inputRankings: []intergo.Ranking{
				tRanking{1, 2, 3, 4, 5},
				tRanking{10, 20, 30, 40, 50},
			},
			combinedRankings: [][]intergo.Res{
				{
					intergo.Res{RankingIDx: 0, ItemIDx: 0},
					intergo.Res{RankingIDx: 1, ItemIDx: 0},
				},
				{
					intergo.Res{RankingIDx: 0, ItemIDx: 0},
					intergo.Res{RankingIDx: 0, ItemIDx: 1},
				},
			},
			expected:  []float64{0.2222222222222222, 0.5},
			threshold: 10e-7,
		},
		{
			inputRankings: []intergo.Ranking{
				tRanking{1, 2, 3},
				tRanking{10, 20, 30},
			},
			combinedRankings: [][]intergo.Res{
				{
					intergo.Res{RankingIDx: 0, ItemIDx: 0},
					intergo.Res{RankingIDx: 1, ItemIDx: 0},
					intergo.Res{RankingIDx: 1, ItemIDx: 1},
				},
				{
					intergo.Res{RankingIDx: 0, ItemIDx: 0},
					intergo.Res{RankingIDx: 0, ItemIDx: 1},
					intergo.Res{RankingIDx: 0, ItemIDx: 2},
				},
			},
			expected:  []float64{0.17833719135802462, 0.4075038580246914},
			threshold: 10e-8,
		},
	}
	for n, tc := range cases {
		tcc := tc
		t.Run(fmt.Sprintf("%d-th unit test", n), func(t *testing.T) {
			actual := o.ExportedGetConstraintMatrix(tcc.inputRankings, tcc.combinedRankings)
			fmt.Println(mat.Formatted(actual))
		})
	}
}
