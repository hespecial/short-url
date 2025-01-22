package bloom

import "github.com/bits-and-blooms/bloom/v3"

var Filter = bloom.NewWithEstimates(100000000, 0.01)

func Add(value []byte) {
	Filter.Add(value)
}

func Contains(value []byte) bool {
	return Filter.Test(value)
}
