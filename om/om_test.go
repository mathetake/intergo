package om_test

import (
	"fmt"
	"testing"

	"github.com/mathetake/intergo"
	"github.com/mathetake/intergo/om"
	"gotest.tools/assert"
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
	o := &om.OptimizedMultiLeaving{
		NumSampling: 100,
		CreditLabel: 0,
		Alpha:       0,
	}

	cases := []struct {
		num           int
		inputRankings []intergo.Ranking
		expected      []intergo.Result
	}{
		{
			inputRankings: []intergo.Ranking{
				tRanking{1, 2, 3, 4, 5},
				tRanking{10, 20, 30, 40, 50},
			},
			expected: []intergo.Result{
				{RankingIndex: 0, ItemIndex: 0},
				{RankingIndex: 1, ItemIndex: 0},
			},
			num: 2,
		},
		{
			inputRankings: []intergo.Ranking{
				tRanking{1, 2, 3},
				tRanking{10, 20, 30},
			},
			expected: []intergo.Result{
				{RankingIndex: 0, ItemIndex: 0},
				{RankingIndex: 1, ItemIndex: 0},
				{RankingIndex: 1, ItemIndex: 1},
			},
			num: 10,
		},
		{
			inputRankings: []intergo.Ranking{
				tRanking{1, 2, 3, 10, 10, 30},
				tRanking{10, 20, 30},
				tRanking{100, 200, 300},
			},
			expected: []intergo.Result{
				{RankingIndex: 0, ItemIndex: 0},
				{RankingIndex: 1, ItemIndex: 0},
				{RankingIndex: 2, ItemIndex: 0},
			},
			num: 2,
		},
	}

	for n, tc := range cases {
		tcc := tc
		t.Run(fmt.Sprintf("%d-th unit test", n), func(t *testing.T) {
			actual, err := o.GetInterleavedRanking(tcc.num, tcc.inputRankings...)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println("actual: ", actual)
		})
	}
}

func TestGetCredit(t *testing.T) {

	cases := []struct {
		RankingIndex       int
		itemId           interface{}
		idToPlacements   []map[interface{}]int
		creditLabel      int
		isSameRankingIndex bool
		expected         float64
	}{
		{
			RankingIndex: 1,
			itemId:     "item1",
			idToPlacements: []map[interface{}]int{
				{"item1": 1, "item2": 2, "item3": 3},
				{"item1": 3, "item2": 1, "item3": 2},
				{"item1": 2, "item2": 1, "item3": 3},
			},
			creditLabel:      0,
			isSameRankingIndex: false,
			expected:         0.3333333333333333,
		},
		{
			RankingIndex: 1,
			itemId:     "item1",
			idToPlacements: []map[interface{}]int{
				{"item1": 1, "item2": 2, "item3": 3},
				{"item1": 3, "item2": 1, "item3": 2},
				{"item1": 2, "item2": 1, "item3": 3},
			},
			creditLabel:      1,
			isSameRankingIndex: false,
			expected:         -2.0,
		},
		{
			RankingIndex: 0,
			itemId:     "item2",
			idToPlacements: []map[interface{}]int{
				{"item1": 1, "item3": 3},
				{"item1": 3, "item2": 1, "item3": 2},
				{"item1": 2, "item2": 1, "item3": 3},
			},
			creditLabel:      1,
			isSameRankingIndex: false,
			expected:         -3.0,
		},
		{
			RankingIndex: 1,
			itemId:     "item2",
			idToPlacements: []map[interface{}]int{
				{"item1": 1, "item3": 3},
				{"item1": 3, "item2": 1, "item3": 2},
				{"item1": 2, "item2": 1, "item3": 3},
			},
			creditLabel:      1,
			isSameRankingIndex: false,
			expected:         0.0,
		},
		{
			RankingIndex: 0,
			itemId:     "item2",
			idToPlacements: []map[interface{}]int{
				{"item1": 1, "item2": 2, "item3": 3},
				{"item1": 3, "item2": 1, "item3": 2},
				{"item1": 2, "item2": 1, "item3": 3},
			},
			creditLabel:      3,
			isSameRankingIndex: false,
			expected:         0,
		},
		{
			RankingIndex: 0,
			itemId:     "item2",
			idToPlacements: []map[interface{}]int{
				{"item1": 1, "item2": 2, "item3": 3},
				{"item1": 3, "item2": 1, "item3": 2},
				{"item1": 2, "item2": 1, "item3": 3},
			},
			creditLabel:      3,
			isSameRankingIndex: true,
			expected:         1,
		},
	}

	for n, tc := range cases {
		tcc := tc
		o := &om.OptimizedMultiLeaving{
			CreditLabel: tc.creditLabel,
		}
		t.Run(fmt.Sprintf("%d-th unit test", n), func(t *testing.T) {
			actual := o.GetCredit(tcc.RankingIndex, tcc.itemId, tcc.idToPlacements, tcc.creditLabel, tcc.isSameRankingIndex)
			assert.Equal(t, tcc.expected, actual)
		})
	}
}

func TestPrefixConstraintSampling(t *testing.T) {
	o := &om.OptimizedMultiLeaving{}

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
				tRanking{1, 20, 30, 40, 50},
			},
			num: 3,
			expectedPatterns: [][]intergo.Result{
				{
					intergo.Result{RankingIndex: 0, ItemIndex: 0},
					intergo.Result{RankingIDx: 1, ItemIndex: 1},
					intergo.Result{RankingIDx: 0, ItemIndex: 1},
				},
				{
					intergo.Result{RankingIDx: 0, ItemIndex: 0},
					intergo.Result{RankingIDx: 0, ItemIndex: 1},
					intergo.Result{RankingIDx: 1, ItemIndex: 1},
				},
				{
					intergo.Result{RankingIDx: 0, ItemIndex: 0},
					intergo.Result{RankingIDx: 0, ItemIndex: 1},
					intergo.Result{RankingIDx: 0, ItemIndex: 2},
				},
				{
					intergo.Result{RankingIDx: 0, ItemIndex: 0},
					intergo.Res{RankingIDx: 1, ItemIndex: 1},
					intergo.Res{RankingIDx: 1, ItemIndex: 2},
				},
				{
					intergo.Res{RankingIDx: 1, ItemIndex: 0},
					intergo.Res{RankingIDx: 1, ItemIndex: 1},
					intergo.Res{RankingIDx: 0, ItemIndex: 1},
				},
				{
					intergo.Res{RankingIDx: 1, ItemIndex: 0},
					intergo.Res{RankingIDx: 1, ItemIndex: 1},
					intergo.Res{RankingIDx: 1, ItemIndex: 2},
				},
				{
					intergo.Res{RankingIDx: 1, ItemIndex: 0},
					intergo.Res{RankingIDx: 0, ItemIndex: 1},
					intergo.Res{RankingIDx: 1, ItemIndex: 1},
				},
				{
					intergo.Res{RankingIDx: 1, ItemIndex: 0},
					intergo.Res{RankingIDx: 0, ItemIndex: 1},
					intergo.Res{RankingIDx: 0, ItemIndex: 2},
				},
			},
		},
	}

	for n, tc := range cases {
		tcc := tc
		t.Run(fmt.Sprintf("%d-th unit test", n), func(t *testing.T) {
			actual := o.ExportedPrefixConstraintSampling(tcc.num, tcc.inputRks...)
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
	o := &om.OptimizedMultiLeaving{Alpha: 0, CreditLabel: 0}

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
					intergo.Res{RankingIDx: 0, ItemIndex: 0},
					intergo.Res{RankingIDx: 1, ItemIndex: 0},
				},
				{
					intergo.Res{RankingIDx: 0, ItemIndex: 0},
					intergo.Res{RankingIDx: 0, ItemIDx: 1},
				},
			},
			expected:  []float64{0.1133786848, 0.8888888889},
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
			expected:  []float64{0.0376778162, 0.4923955480},
			threshold: 10e-8,
		},
		{
			inputRankings: []intergo.Ranking{
				tRanking{1, 2, 3},
				tRanking{10, 20, 30},
				tRanking{100, 200, 300},
			},
			combinedRankings: [][]intergo.Res{
				{
					intergo.Res{RankingIDx: 0, ItemIDx: 0},
					intergo.Res{RankingIDx: 1, ItemIDx: 0},
					intergo.Res{RankingIDx: 2, ItemIDx: 0},
				},
				{
					intergo.Res{RankingIDx: 0, ItemIDx: 0},
					intergo.Res{RankingIDx: 0, ItemIDx: 1},
					intergo.Res{RankingIDx: 2, ItemIDx: 0},
				},
				{
					intergo.Res{RankingIDx: 1, ItemIDx: 0},
					intergo.Res{RankingIDx: 1, ItemIDx: 1},
					intergo.Res{RankingIDx: 0, ItemIDx: 0},
				},
			},
			expected:  []float64{0.1611570248, 0.5850000000, 0.5850000000},
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
				if true != (diff < tcc.threshold) {
					t.Logf("unexpected difference at %d-th element: actual:%.10f != expected:%.10f", i, actual[i], tcc.expected[i])
				}
				assert.Equal(t, true, diff < tcc.threshold)
			}
		})
	}
}
