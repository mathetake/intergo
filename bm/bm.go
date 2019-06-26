package bm

import (
	"math/rand"

	"github.com/mathetake/intergo"
)

type BalancedMultileaving struct{}

var _ intergo.Interleaving = &BalancedMultileaving{}

func (*BalancedMultileaving) GetInterleavedRanking(num int, rks ...intergo.Ranking) ([]intergo.Result, error) {
	var numR = len(rks)
	res := make([]intergo.Result, 0, num)

	// sIDs stores item's ID in order to prevent duplication in the generated list.
	sIDs := map[interface{}]interface{}{}

	// The fact that the index stored in usedUpRks means it is already used up.
	usedUpRks := map[int]bool{}

	counter := map[int]int{}

	for len(res) < num && len(usedUpRks) != numR {

		// chose randomly one ranking from the ones used up yet
		var selectedRkIdx = rand.Intn(numR)
		if _, ok := usedUpRks[selectedRkIdx]; ok {
			continue
		}

		// get pointer on the selected ranking
		c, _ := counter[selectedRkIdx]

		// get ID of the pointed item
		itemID := rks[selectedRkIdx].GetIDByIndex(c)

		if _, ok := sIDs[itemID]; !ok {
			res = append(res, intergo.Result{
				RankingIndex: selectedRkIdx,
				ItemIDx:    c,
			})
			sIDs[itemID] = true
		}

		// increment pointer on the selected ranking
		counter[selectedRkIdx]++

		if c, _ := counter[selectedRkIdx]; c >= rks[selectedRkIdx].Len() {
			usedUpRks[selectedRkIdx] = true
		}
	}
	return res, nil
}
