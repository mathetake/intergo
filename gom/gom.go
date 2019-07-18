package gom

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/mathetake/intergo"
)

type GreedyOptimizedMultiLeaving struct {
	NumSampling int
	CreditLabel int
	Alpha       float64
}

const els = 1e-20

var _ intergo.Interleaving = &GreedyOptimizedMultiLeaving{}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetInterleavedRanking ... get a Interleaved ranking sampled from a set of interleaved rankings
// generated by `prefixConstraintSampling` method.
// Note that the way of the sampling is different from the original paper [Schuth, Anne, et al.,2014]
// where they solved LP with the unbiased constraint.
// We omit the unbiased constraint and only take `sensitivity` into account. Then we sample a ranking
// according to calculated sensitivities defined by equation (1) in [Manabe, Tomohiro, et al., 2017]
func (o *GreedyOptimizedMultiLeaving) GetInterleavedRanking(num int, rankings ...intergo.Ranking) ([]*intergo.Result, error) {
	if num < 1 {
		return nil, intergo.ErrNonPositiveSamplingNumParameters
	} else if len(rankings) < 1 {
		return nil, intergo.ErrInsufficientRankingsParameters
	}

	var wg sync.WaitGroup
	cRks := make([][]*intergo.Result, o.NumSampling)
	wg.Add(o.NumSampling)
	for i := 0; i < o.NumSampling; i++ {
		go func(i int) {
			defer wg.Done()
			cRks[i] = o.prefixConstraintSampling(num, rankings...)
		}(i)
	}
	wg.Wait()

	// calc Insensitivity of sampled rankings
	ins := o.calcInsensitivities(rankings, cRks)

	// init +inf value
	min := math.Inf(0)
	var maxIdx int
	for i, v := range ins {
		if v < min {
			maxIdx, min = i, v
		}
	}
	return cRks[maxIdx], nil
}

func (o *GreedyOptimizedMultiLeaving) GetCredit(rankingIndex int, itemId intergo.ID, idToPlacements []map[intergo.ID]int, creditLabel int, isSameRankingIndex bool) float64 {
	switch creditLabel {
	case 0:
		// credit = 1 / (original rank)
		placement, ok := idToPlacements[rankingIndex][itemId]
		if ok {
			return 1 / float64(placement)
		} else {
			return 1 / float64(len(idToPlacements[rankingIndex])+1)
		}
	case 1:
		// credit = -(relative rank - 1)
		if _, ok := idToPlacements[rankingIndex][itemId]; !ok {
			return -float64(len(idToPlacements))
		}
		var numLess float64
		for i := 0; i < len(idToPlacements); i++ {
			if _, ok := idToPlacements[i][itemId]; !ok {
				continue
			}
			if idToPlacements[i][itemId] < idToPlacements[rankingIndex][itemId] {
				numLess += 1
			}
		}
		return -numLess
	default:
		// credit = 1 if output ranking idx equals input ranking idx
		// else credit = 0
		if isSameRankingIndex {
			return 1
		}
		return 0
	}
}

func (o *GreedyOptimizedMultiLeaving) GetIdToPlacementMap(rks []intergo.Ranking) []map[intergo.ID]int {
	var iRkNum = len(rks)
	itemIds := make(map[intergo.ID]bool)
	idToPlacements := make([]map[intergo.ID]int, iRkNum)
	// idToPlacements[ranking idx][item id] -> original ranking placement
	for i := 0; i < iRkNum; i++ {
		m := make(map[intergo.ID]int, rks[i].Len())
		for j := 0; j < rks[i].Len(); j++ {
			itemId := rks[i].GetIDByIndex(j)
			m[itemId] = j + 1
			itemIds[itemId] = true
		}
		idToPlacements[i] = m
	}
	return idToPlacements
}

func (o *GreedyOptimizedMultiLeaving) CalcInsensitivityAndBias(rks []intergo.Ranking, res []*intergo.Result, creditLabel int, alpha float64) (float64, float64) {
	var iRkNum = len(rks)
	var insensitivityMean float64

	idToPlacements := o.GetIdToPlacementMap(rks)
	insensitivityMap := make([]float64, iRkNum)
	biasMap := make([][]float64, iRkNum)

	for i := 0; i < iRkNum; i++ {
		biasMap[i] = make([]float64, len(res))
		var bias float64
		for j := 0; j < len(res); j++ {
			var s = 1 / float64(j+1)
			itemId := rks[res[j].RankingIndex].GetIDByIndex(res[j].ItemIndex)
			credit := o.GetCredit(i, itemId, idToPlacements, creditLabel, res[j].RankingIndex == i)
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

	var fResLen = float64(len(res))

	insensitivityMean /= float64(iRkNum)
	if math.Abs(insensitivityMean) < els {
		return math.Inf(1), biasSum / fResLen
	}
	var insensitivitySum float64
	for i := 0; i < iRkNum; i++ {
		var in = insensitivityMap[i] - insensitivityMean
		insensitivitySum += in * in
	}
	bias := biasSum / fResLen
	return (insensitivitySum + alpha*bias) / (insensitivityMean * insensitivityMean), biasSum / fResLen
}

func (o *GreedyOptimizedMultiLeaving) calcInsensitivities(rks []intergo.Ranking, cRks [][]*intergo.Result) []float64 {
	res := make([]float64, len(cRks))

	var wg sync.WaitGroup
	wg.Add(len(cRks))
	for k := 0; k < len(cRks); k++ {
		go func(k int) {
			defer wg.Done()
			res[k], _ = o.CalcInsensitivityAndBias(rks, cRks[k], o.CreditLabel, o.Alpha)
		}(k)
	}
	wg.Wait()
	return res
}

func (*GreedyOptimizedMultiLeaving) prefixConstraintSampling(num int, rks ...intergo.Ranking) []*intergo.Result {
	var numR = len(rks)
	res := make([]*intergo.Result, 0, num)

	// sIDs stores item's ID in order to prevent duplication in the generated list.
	sIDs := make(map[intergo.ID]struct{}, num)

	// The fact that the index stored in usedUpRks means it is already used up.
	usedUpRks := make(map[int]struct{}, numR)

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
				res = append(res, &intergo.Result{
					RankingIndex: selectedRkIdx,
					ItemIndex:    j,
				})
				sIDs[rk.GetIDByIndex(j)] = struct{}{}
				break
			}
		}

		if len(res) == bef {
			usedUpRks[selectedRkIdx] = struct{}{}
		}
	}
	return res
}