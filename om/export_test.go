package om

import (
	"github.com/mathetake/intergo"
)

func (o *OptimizedMultiLeaving) ExportedPrefixConstraintSampling(num int, rks ...intergo.Ranking) []intergo.Res {
	return o.prefixConstraintSampling(num, rks...)
}

func (o *OptimizedMultiLeaving) ExportedCalcInsensitivity(rks []intergo.Ranking, cRks [][]intergo.Res) []float64 {
	return o.calcInsensitivity(rks, cRks)
}
