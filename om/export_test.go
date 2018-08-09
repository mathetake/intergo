package om

import (
	"github.com/mathetake/intergo"
	"gonum.org/v1/gonum/mat"
)

func (o *OptimizedMultiLeaving) ExportedPrefixConstraintSampling(num int, rks ...intergo.Ranking) []intergo.Res {
	return o.prefixConstraintSampling(num, rks...)
}

func (o *OptimizedMultiLeaving) ExportedCalcInsensitivity(rks []intergo.Ranking, cRks [][]intergo.Res) []float64 {
	return o.calcInsensitivity(rks, cRks)
}

func (o *OptimizedMultiLeaving) ExportedGetConstraintMatrix(rks []intergo.Ranking, cRks [][]intergo.Res) *mat.Dense {
	return o.getConstraintMatrix(rks, cRks)
}
