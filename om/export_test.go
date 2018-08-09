package om

import "github.com/mathetake/intergo"

func (om *OptimizedMultiLeaving) ExportedGetCombinedRanking(num int, rks ...intergo.Ranking) []intergo.Res {
	return om.getCombinedRanking(num, rks...)
}
