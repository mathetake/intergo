package tdm

import (
	"github.com/lukechampine/randmap"
	"github.com/mathetake/intergo"
)

type TeamDraftMultileaving struct{}

var _ intergo.Interleaving = &TeamDraftMultileaving{}

func (tdm *TeamDraftMultileaving) GetInterleavedRanking(num int, rks ...intergo.Ranking) []intergo.Res {
	var numR = len(rks)
	res := make([]intergo.Res, 0, num)

	// sIDs stores item's ID in order to prevent duplication in the generated list.
	sIDs := map[interface{}]interface{}{}

	// minRks have rankings' index whose number of selected items is minimum
	minRks := map[int]interface{}{}
	for i := 0; i < numR; i++ {
		minRks[i] = true
	}

	// The fact that the index stored in usedUpRks means it is already used up.
	usedUpRks := map[int]bool{}

	for len(res) < num && len(usedUpRks) != numR {

		// chose one ranking from keys of minRks
		var selected = getRandomKey(minRks)
		var rk = rks[selected]

		var bef = len(res)

		for j := 0; j < rk.Len(); j++ {
			if _, ok := sIDs[rk.GetIDByIndex(j)]; !ok {
				res = append(res, intergo.Res{
					RankingIDx: selected,
					ItemIDx:    j,
				})

				sIDs[rk.GetIDByIndex(j)] = true
				break
			}
		}

		if len(res) == bef {
			usedUpRks[selected] = true
		}

		// delete the selected ranking from minRks
		delete(minRks, selected)

		if len(minRks) == 0 {
			// restore the targets
			for i := 0; i < numR; i++ {
				if !usedUpRks[i] {
					minRks[i] = true
				}
			}
		}
	}
	return res
}

func getRandomKey(m map[int]interface{}) int {
	k, _ := randmap.Key(m).(int)
	return k
}
