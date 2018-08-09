package om

import (
	"math/rand"
	"time"

	"sync"

	"github.com/mathetake/intergo"
	"github.com/pkg/errors"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize/convex/lp"
)

type OptimizedMultiLeaving struct{}

const (
	numSampling = 200
)

var _ intergo.Interleaving = &OptimizedMultiLeaving{}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (o *OptimizedMultiLeaving) GetInterleavedRanking(num int, rks ...intergo.Ranking) ([]intergo.Res, error) {

	var wg sync.WaitGroup
	cRks := make([][]intergo.Res, numSampling)
	for i := 0; i < numSampling; i++ {
		ii := i
		wg.Add(1)
		go func() {
			cRks[ii] = o.prefixConstraintSampling(num, rks...)
			wg.Done()
		}()
	}
	wg.Wait()

	// get matrix for objective
	wg.Add(1)
	c := make([]float64, 0, numSampling)
	go func() {
		c = o.calcInsensitivity(rks, cRks)
		wg.Done()
	}()

	// get LHS of constraint
	wg.Add(1)
	var cMat = &mat.Dense{}
	go func() {
		cMat = o.getConstraintMatrix(rks, cRks)
		wg.Done()
	}()

	// RHS of constraint
	b := make([]float64, 1+num*len(rks))
	b[0] = 1 // for probability constraint

	// solve linear programming
	wg.Wait()
	_, ps, err := lp.Simplex(c, cMat, b, 10e-5, nil)
	if err != nil {
		return nil, errors.Wrap(err, "lp.Simplex failed.")
	}

	var max float64
	var maxIDx int
	for i, v := range ps {
		if v > max {
			maxIDx, max = i, v
		}
	}
	return cRks[maxIDx], nil
}

func (*OptimizedMultiLeaving) getConstraintMatrix(rks []intergo.Ranking, cRks [][]intergo.Res) *mat.Dense {
	var numInputRankings = len(rks)
	var numItem = len(cRks[0])
	var numCombinedList = len(cRks)

	var wg = sync.WaitGroup{}

	// len(cRks[0]) = r, len(rks) = j
	ret := mat.NewDense(1+numInputRankings*numItem, numCombinedList, nil)

	wg.Add(1)
	go func() {
		for k := 0; k < numCombinedList; k++ {
			ret.Set(0, k, 1)
		}
		wg.Done()
	}()

	for jj := 0; jj < numInputRankings; jj++ {
		j := jj
		for rr := 0; rr < numItem; rr++ {
			r := rr
			for kk := 0; kk < numCombinedList; kk++ {
				k := kk
				wg.Add(1)
				go func() {
					var c float64
					for i := 0; i <= r; i++ {
						var s = i + 1
						if cRks[k][i].RankingIDx == j {
							s *= s
						} else {
							s *= rks[j].Len()
						}
						c += 1 / float64(s)
					}
					ret.Set(1+j*(1+r), k, c)
					wg.Done()
				}()
			}
		}
	}
	wg.Wait()
	return ret
}

func (*OptimizedMultiLeaving) calcInsensitivity(rks []intergo.Ranking, cRks [][]intergo.Res) []float64 {
	res := make([]float64, len(cRks))

	var iRkNum = len(rks)
	var wg sync.WaitGroup

	for kk := 0; kk < len(cRks); kk++ {
		k := kk
		wg.Add(1)
		go func() {
			var mean float64

			jToScoreMap := make([]float64, iRkNum)
			for j := 0; j < iRkNum; j++ {

				for i := 0; i < len(cRks[0]); i++ {
					var s = i + 1
					if cRks[k][i].RankingIDx == j {
						s *= s
					} else {
						s *= rks[j].Len() + 1
					}
					ss := 1 / float64(s)
					jToScoreMap[j] += ss
					mean += ss
				}
			}

			mean /= float64(iRkNum)

			var score float64
			for j := 0; j < iRkNum; j++ {
				var s = jToScoreMap[j] - mean
				score += s * s
			}
			res[k] = score
			wg.Done()
		}()
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
