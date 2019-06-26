package intergo

import "github.com/pkg/errors"

var (
	ErrNonPositiveSamplingNumParameters = errors.New("`num` parameter should be positive")
	ErrInsufficientRankingsParameters   = errors.New("the number of provided rankings should be positive")
)
