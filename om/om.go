package om

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/mathetake/intergo"
	"github.com/pkg/errors"
)

type OptimizedMultiLeaving struct {
	NumSampling int
	CreditLabel int
	Alpha       float64
}

var _ intergo.Interleaving = &OptimizedMultiLeaving{}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetInterleavedRanking ... get a Interleaved ranking sampled from a set of interleaved rankings
// generated by `prefixConstraintSampling` method.
// Note that the way of the sampling is different from the original paper [Schuth, Anne, et al.,2014]
// where they solved LP with the unbiased constraint.
// We omit the unbiased constraint and only take `sensitivity` into account. Then we sample a ranking
// according to calculated sensitivities defined by equation (1) in [Manabe, Tomohiro, et al., 2017]
func (o *OptimizedMultiLeaving) GetInterleavedRanking(num int, rks ...intergo.Ranking) ([]intergo.Res, error) {
	if num < 1 {
		return nil, errors.Errorf("invalid NumSampling: %d", o.NumSampling)
	}

	var wg sync.WaitGroup
	cRks := make([][]intergo.Res, o.NumSampling)
	for i := 0; i < o.NumSampling; i++ {
		wg.Add(1)
		go func(i int) {
			cRks[i] = o.prefixConstraintSampling(num, rks...)
			wg.Done()
		}(i)
	}
	wg.Wait()

	// calc Insensitivity of sampled rankings
	ins := o.calcInsensitivities(rks, cRks)

	// init +inf value
	min := math.Inf(0)
	var maxIDx int
	for i, v := range ins {
		if v < min {
			maxIDx, min = i, v
		}
	}
	return cRks[maxIDx], nil
}

func getCredit(rankingIdx int, itemId int, idToPlacements []map[int]int, creditLabel int, isSameRankingIdx bool) float64 {
	switch creditLabel {
	case 0:
		// credit = 1 / (original rank)
		placement, ok := idToPlacements[rankingIdx][itemId]
		if ok {
			return 1 / float64(placement)
		} else {
			return 1 / float64(len(idToPlacements[rankingIdx])+1)
		}
	case 1:
		// credit = -(relative rank - 1)
		numGreater := 0.0
		for i := 0; i < len(idToPlacements); i++ {
			_, ok1 := idToPlacements[i][itemId]
			_, ok2 := idToPlacements[rankingIdx][itemId]
			if !ok2 {
				continue
			}
			if !ok1 {
				numGreater += 1
				continue
			}
			if idToPlacements[i][itemId] > idToPlacements[rankingIdx][itemId] {
				numGreater += 1
			}
		}
		return -numGreater
	default:
		// credit = 1 if output ranking idx equals input ranking idx
		// else credit = 0
		if isSameRankingIdx {
			return 1
		} else {
			return 0
		}
	}
	return 0
}

func (o *OptimizedMultiLeaving) GetIdToPlacementMap(rks []intergo.Ranking) []map[int]int {
	var iRkNum = len(rks)
	itemIds := make(map[int]bool)
	idToPlacements := make([]map[int]int, iRkNum)
	// idToPlacements[ranking idx][item id] -> original ranking placement
	for i := 0; i < iRkNum; i++ {
		idToPlacements[i] = map[int]int{}
		for j := 0; j < (rks)[i].Len(); j++ {
			itemId := (rks)[i].GetIDByIndex(j).(int)
			idToPlacements[i][itemId] = j + 1
			itemIds[itemId] = true
		}
	}
	return idToPlacements
}

func (o *OptimizedMultiLeaving) CalcInsensitivityAndBias(rks []intergo.Ranking, res []intergo.Res, creditLabel int, alpha float64) (float64, float64) {
	var iRkNum = len(rks)
	var insensitivityMean float64

	idToPlacements := o.GetIdToPlacementMap(rks)
	insensitivityMap := make([]float64, iRkNum)
	biasMap := make([][]float64, iRkNum)

	for i := 0; i < iRkNum; i++ {
		biasMap[i] = make([]float64, len(res))
		bias := 0.0
		for j := 0; j < len(res); j++ {
			var s = 1 / float64(j+1)
			itemId := rks[res[j].RankingIDx].GetIDByIndex(res[j].ItemIDx).(int)
			credit := getCredit(i, itemId, idToPlacements, creditLabel, res[j].RankingIDx == i)
			ss := s * credit
			insensitivityMap[i] += ss
			insensitivityMean += ss
			bias += credit
			biasMap[i][j] = bias
		}
	}

	var biasSum float64
	for r := 0; r < len(res); r++ {
		min := math.Inf(1)
		max := math.Inf(-1)
		for i := 0; i < iRkNum; i++ {
			v := math.Abs(biasMap[i][r])
			if min > v {
				min = v
			}
			if max < v {
				max = v
			}
		}
		if creditLabel != 0 {
			min += 1
			max += 1
		}
		biasSum += 1.0 - math.Abs(min/max)
	}

	insensitivityMean /= float64(iRkNum)
	EPS := 1e-20
	if math.Abs(insensitivityMean) < EPS {
		return math.Inf(1), biasSum / float64(len(res))
	}
	var insensitivitySum float64
	for i := 0; i < iRkNum; i++ {
		var in = insensitivityMap[i] - insensitivityMean
		insensitivitySum += in * in
	}
	bias := biasSum / float64(len(res))
	return (insensitivitySum + alpha*bias) / (insensitivityMean * insensitivityMean), biasSum / float64(len(res))
}

func (o *OptimizedMultiLeaving) calcInsensitivities(rks []intergo.Ranking, cRks [][]intergo.Res) []float64 {
	res := make([]float64, len(cRks))

	var wg sync.WaitGroup

	for k := 0; k < len(cRks); k++ {
		wg.Add(1)
		go func(k int) {
			res[k], _ = o.CalcInsensitivityAndBias(rks, cRks[k], o.CreditLabel, o.Alpha)
			wg.Done()
		}(k)
	}
	wg.Wait()
	return res
}

func (*OptimizedMultiLeaving) prefixConstraintSampling(num int, rks ...intergo.Ranking) []intergo.Res {
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
