package tdm_test

import (
	"testing"

	"github.com/mathetake/intergo"
	"gotest.tools/assert"

	"fmt"

	"github.com/mathetake/intergo/tdm"
)

type item struct {
	id int
}

func (i item) GetID() interface{} {
	return i.id
}

var _ intergo.Item = item{}

func TestGetInterleavedRanking(t *testing.T) {
	tdm := &tdm.TeamDraftMultileaving{}

	cases := []struct {
		inputRks         []intergo.Items
		num              int
		expectedPatterns [][]intergo.Res
	}{
		{
			inputRks: []intergo.Items{
				{
					item{1}, item{2}, item{3}, item{4}, item{5},
				},
				{
					item{10}, item{20}, item{30}, item{3}, item{3},
				},
			},
			num: 2,
			expectedPatterns: [][]intergo.Res{
				{
					intergo.Res{RankingIDx: 0, ItemIDx: 0},
					intergo.Res{RankingIDx: 1, ItemIDx: 0},
				},
			},
		},
	}

	for n, tc := range cases {
		tcc := tc
		t.Run(fmt.Sprintf("%d-th unit test", n), func(t *testing.T) {
			actual := tdm.GetInterleavedRanking(tcc.inputRks, tcc.num)
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
