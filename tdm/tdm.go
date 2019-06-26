package tdm

import (
	"math/rand"

	"github.com/mathetake/intergo"
)

type TeamDraftMultileaving struct{}

var _ intergo.Interleaving = &TeamDraftMultileaving{}

func (tdm *TeamDraftMultileaving) GetInterleavedRanking(num int, rks ...intergo.Ranking) ([]intergo.Result, error) {
	var numR = len(rks)
	res := make([]intergo.Result, 0, num)

	// sIDs stores item's ID in order to prevent duplication in the generated list.
	sIDs := make(map[intergo.ID]struct{}, num)

	// minRks have rankings' index whose number of selected items is minimum
	minRks := make([]int, 0, numR)

	// lastIdx has a last index of the indexed ranking
	lastIdx := map[int]int{}
	for i := 0; i < numR; i++ {
		minRks = append(minRks, i)
		lastIdx[i] = 0
	}

	// The fact that the index stored in usedUpRks means it is already used up.
	usedUpRks := make(map[int]struct{}, numR)

	for len(res) < num && len(usedUpRks) != numR {

		// chose one ranking from keys of minRks
		var selected int
		selected, minRks = popRandomIdx(minRks)
		var rk = rks[selected]

		var bef = len(res)

		for j := lastIdx[selected]; j < rk.Len(); j++ {
			if _, ok := sIDs[rk.GetIDByIndex(j)]; !ok {
				res = append(res, intergo.Result{
					RankingIndex: selected,
					ItemIndex:    j,
				})

				sIDs[rk.GetIDByIndex(j)] = struct{}{}
				lastIdx[selected] = j
				break
			}
		}

		if len(res) == bef {
			usedUpRks[selected] = struct{}{}
		}

		if len(minRks) == 0 {
			// restore the targets
			minRks = make([]int, 0, numR-len(usedUpRks))
			for i := 0; i < numR; i++ {
				if _, ok := usedUpRks[i]; !ok {
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

	popped := make([]int, 0, len(target)-1)

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
