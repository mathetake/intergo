package gom

import (
	"github.com/mathetake/intergo"
)

func (o *GreedyOptimizedMultiLeaving) ExportedPrefixConstraintSampling(num int, rks ...intergo.Ranking) []*intergo.Result {
	return o.prefixConstraintSampling(num, rks...)
}

func (o *GreedyOptimizedMultiLeaving) ExportedCalcInsensitivity(rks []intergo.Ranking, cRks [][]*intergo.Result) []float64 {
	return o.calcInsensitivities(rks, cRks)
}
