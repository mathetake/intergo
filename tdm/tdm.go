package tdm

import (
	"github.com/mathetake/intergo"
	"math/rand"
)

type TeamDraftMultileaving struct{}

var _ intergo.Interleaving = &TeamDraftMultileaving{}

func (tdm *TeamDraftMultileaving) GetInterleavedRanking(num int, rks ...intergo.Ranking) ([]intergo.Res, error) {
	var numR = len(rks)
	res := make([]intergo.Res, 0, num)

	// sIDs stores item's ID in order to prevent duplication in the generated list.
	sIDs := map[interface{}]interface{}{}

	// minRks have rankings' index whose number of selected items is minimum
	minRks := make([]int, 0, numR)

	// lastIdx has a last index of the indexed ranking
	lastIdx := map[int]int{}
	for i := 0; i < numR; i++ {
		minRks = append(minRks, i)
		lastIdx[i] = 0
	}

	// The fact that the index stored in usedUpRks means it is already used up.
	usedUpRks := map[int]bool{}

	for len(res) < num && len(usedUpRks) != numR {

		// chose one ranking from keys of minRks
		var selected int
		selected, minRks = popRandomIdx(minRks)
		var rk = rks[selected]

		var bef = len(res)

		for j := lastIdx[selected]; j < rk.Len(); j++ {
			if _, ok := sIDs[rk.GetIDByIndex(j)]; !ok {
				res = append(res, intergo.Res{
					RankingIDx: selected,
					ItemIDx:    j,
				})

				sIDs[rk.GetIDByIndex(j)] = true
				lastIdx[selected] = j
				break
			}
		}

		if len(res) == bef {
			usedUpRks[selected] = true
		}

		if len(minRks) == 0 {
			// restore the targets
			minRks = make([]int, 0, numR)
			for i := 0; i < numR; i++ {
				if !usedUpRks[i] {
					minRks = append(minRks, i)
				}
			}
		}
	}
	return res, nil
}

func popRandomIdx(target []int) (int, []int) {
	if len(target) == 1 {
		return target[0], []int{}
	}

	selectedIDx := rand.Intn(len(target))
	selected := target[selectedIDx]

	popped := make([]int, 0, len(target) -1)

	for i, idx := range target {
		if i < selectedIDx {
			popped = append(popped, idx)
		} else if i == selectedIDx {
			continue
		} else {
			popped = append(popped, idx)
		}
	}
	return selected, popped
}
