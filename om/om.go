package om

import (
	"math/rand"
	"time"

	"github.com/mathetake/intergo"
)

type OptimizedMultiLeaving struct{}

var _ intergo.Interleaving = &OptimizedMultiLeaving{}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (om *OptimizedMultiLeaving) GetInterleavedRanking(num int, rks ...intergo.Ranking) []intergo.Res {
	return nil
}

func (om *OptimizedMultiLeaving) getCombinedRanking(num int, rks ...intergo.Ranking) []intergo.Res {
	var numR = len(rks)
	res := make([]intergo.Res, 0, num)

	// sIDs stores item's ID in order to prevent duplication in the generated list.
	sIDs := map[interface{}]interface{}{}

	// The fact that the index stored in usedUpRks means it is already used up.
	usedUpRks := map[int]bool{}

	for len(res) < num && len(usedUpRks) != numR {

		// chose randomly one ranking from the ones used up yet
		var selectedRkIdx = rand.Intn(numR)
		if _, ok := usedUpRks[selectedRkIdx]; ok {
			continue
		}

		var rk = rks[selectedRkIdx]
		var bef = len(res)
		for j := 0; j < rk.Len(); j++ {
			if _, ok := sIDs[rk.GetIDByIndex(j)]; !ok {
				res = append(res, intergo.Res{
					RankingIDx: selectedRkIdx,
					ItemIDx:    j,
				})
				sIDs[rk.GetIDByIndex(j)] = true
				break
			}
		}

		if len(res) == bef {
			usedUpRks[selectedRkIdx] = true
		}
	}
	return res
}
