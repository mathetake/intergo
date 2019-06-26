package bm

import (
	"math/rand"

	"github.com/mathetake/intergo"
)

type BalancedMultileaving struct{}

var _ intergo.Interleaving = &BalancedMultileaving{}

func (*BalancedMultileaving) GetInterleavedRanking(num int, rankings ...intergo.Ranking) ([]*intergo.Result, error) {
	if num < 1 {
		return nil, intergo.ErrNonPositiveSamplingNumParameters
	} else if len(rankings) < 1 {
		return nil, intergo.ErrInsufficientRankingsParameters
	}

	var numR = len(rankings)
	res := make([]*intergo.Result, 0, num)

	// sIDs stores item's ID in order to prevent duplication in the generated list.
	sIDs := make(map[intergo.ID]struct{}, num)

	// The fact that the index stored in usedUpRks means it is already used up.
	usedUpRks := make(map[int]struct{}, numR)

	counter := make(map[int]int, numR)

	for len(res) < num && len(usedUpRks) != numR {

		// chose randomly one ranking from the ones used up yet
		var selectedRkIdx = rand.Intn(numR)
		if _, ok := usedUpRks[selectedRkIdx]; ok {
			continue
		}

		// get pointer on the selected ranking
		c, _ := counter[selectedRkIdx]

		// get ID of the pointed item
		itemID := rankings[selectedRkIdx].GetIDByIndex(c)

		if _, ok := sIDs[itemID]; !ok {
			res = append(res, &intergo.Result{
				RankingIndex: selectedRkIdx,
				ItemIndex:    c,
			})
			sIDs[itemID] = struct{}{}
		}

		// increment pointer on the selected ranking
		counter[selectedRkIdx]++

		if c, _ := counter[selectedRkIdx]; c >= rankings[selectedRkIdx].Len() {
			usedUpRks[selectedRkIdx] = struct{}{}
		}
	}
	return res, nil
}
