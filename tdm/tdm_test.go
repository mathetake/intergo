package tdm_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/mathetake/intergo"
	"github.com/mathetake/intergo/tdm"
	"gotest.tools/assert"
)

type tRanking []int

func (rk tRanking) GetIDByIndex(i int) intergo.ID {
	return intergo.ID(strconv.Itoa(rk[i]))
}

func (rk tRanking) Len() int {
	return len(rk)
}

var _ intergo.Ranking = tRanking{}

func TestTeamDraftMultileaving(t *testing.T) {
	td := &tdm.TeamDraftMultileaving{}

	cases := []struct {
		inputRks         []intergo.Ranking
		num              int
		expectedPatterns [][]intergo.Result
		expErr           error
	}{
		{
			inputRks: []intergo.Ranking{},
			num:      10,
			expErr:   intergo.ErrInsufficientRankingsParameters,
		},
		{
			inputRks: []intergo.Ranking{
				tRanking{1, 2, 3, 4, 5},
				tRanking{10, 20, 30, 40, 50},
			},
			num:    0,
			expErr: intergo.ErrNonPositiveSamplingNumParameters,
		},
		{
			expErr: intergo.ErrNonPositiveSamplingNumParameters,
		},
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
					intergo.Result{RankingIndex: 1, ItemIndex: 0},
					intergo.Result{RankingIndex: 0, ItemIndex: 0},
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
					intergo.Result{RankingIndex: 1, ItemIndex: 0},
					intergo.Result{RankingIndex: 0, ItemIndex: 1},
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
					intergo.Result{RankingIndex: 1, ItemIndex: 1},
					intergo.Result{RankingIndex: 0, ItemIndex: 1},
				},
				{
					intergo.Result{RankingIndex: 1, ItemIndex: 0},
					intergo.Result{RankingIndex: 0, ItemIndex: 1},
					intergo.Result{RankingIndex: 1, ItemIndex: 1},
				},
				{
					intergo.Result{RankingIndex: 1, ItemIndex: 0},
					intergo.Result{RankingIndex: 0, ItemIndex: 1},
					intergo.Result{RankingIndex: 0, ItemIndex: 2},
				},
				{
					intergo.Result{RankingIndex: 0, ItemIndex: 0},
					intergo.Result{RankingIndex: 1, ItemIndex: 1},
					intergo.Result{RankingIndex: 1, ItemIndex: 2},
				},
			},
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(fmt.Sprintf("%d-th unit test", n), func(t *testing.T) {
			actual, actualErr := td.GetInterleavedRanking(tc.num, tc.inputRks...)
			if tc.expErr != nil {
				assert.Equal(t, tc.expErr, actualErr)
				return // exit
			} else if actualErr != nil {
				t.Fatal(actualErr)
			}

			assert.Equal(t, true, len(actual) <= tc.num)

			var isExpected bool
			for _, expected := range tc.expectedPatterns {

				var isExpectedPattern = true
				for i := 0; i < tc.num; i++ {
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

func TestPopRandomIdx(t *testing.T) {
	for i, cc := range []struct {
		target []int
		expLen int
	}{
		{
			target: []int{1},
			expLen: 0,
		},
		{
			target: []int{1, 2, 3, 4},
			expLen: 3,
		},
		{
			target: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expLen: 9,
		},
	} {
		c := cc
		t.Run(fmt.Sprintf("%d-th case", i), func(t *testing.T) {
			actualS, actualP := tdm.ExportedPopRandomIdx(c.target)
			assert.Equal(t, c.expLen, len(actualP))

			isIncluded := false
			for _, actual := range actualP {
				if actual == actualS {
					isIncluded = true
				}
			}

			assert.Equal(t, false, isIncluded)
		})
	}
}

func TestTeamDraftMultileaving_RankingRatio(t *testing.T) {
	ml := &tdm.TeamDraftMultileaving{}

	for i, cc := range []struct {
		itemNum, rankingNum, returnedNum int
		threshold                        float64
	}{
		{1e2, 3, 10, 1e-1},
		{1e3, 3, 100, 1e-2},
		{1e4, 3, 1000, 1e-3},
		{1e5, 3, 10000, 1e-4},
	} {
		c := cc
		t.Run(fmt.Sprintf("%d-th case", i), func(t *testing.T) {
			rks := getRankings(c.itemNum, c.rankingNum)

			res, err := ml.GetInterleavedRanking(c.returnedNum, rks...)
			if err != nil {
				t.Fatalf("GetInterleavedRanking failed: %v", err)
			}

			counts := map[int]int{}
			for _, it := range res {
				counts[it.RankingIndex]++
			}
			fmt.Println(counts)

			for _, v := range counts {
				diff := float64(v)/float64(c.returnedNum) - float64(1)/float64(c.rankingNum)

				if diff < 0 {
					diff *= -1
				}

				assert.Equal(t, true, diff < c.threshold)
			}
		})
	}
}

func getRankings(itemNum, RankingNum int) []intergo.Ranking {
	rks := make([]intergo.Ranking, RankingNum)
	for i := 0; i < RankingNum; i++ {
		rk := tRanking{}
		for j := 0; j < itemNum; j++ {
			rk = append(rk, i*itemNum+j)
		}
		rks[i] = rk
	}
	return rks
}
