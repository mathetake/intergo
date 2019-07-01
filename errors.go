package intergo

import "github.com/pkg/errors"

// these errors are intended to be returned by Interleaving.GetInterleavedRanking function.
var (
	// ErrNonPositiveSamplingNumParameters should be returned when given "num" is non-positive integer
	ErrNonPositiveSamplingNumParameters = errors.New("`num` parameter should be positive")
	// ErrInsufficientRankingsParameters should be returned when given rankings is empty.
	ErrInsufficientRankingsParameters = errors.New("the number of provided rankings should be positive")
)
