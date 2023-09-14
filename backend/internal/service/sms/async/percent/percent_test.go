package percent

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

const epsilon = 1e-9

func AlmostEqual(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}

func TestPercent_Add(t *testing.T) {
	testCases := []struct {
		name     string
		size     int
		addVals  []bool
		wantRets []float64
	}{
		{
			name:     "test",
			size:     4,
			addVals:  []bool{false, true, true, false, false, true, false, false},
			wantRets: []float64{0, 0.5, 2.0 / 3.0, 0.5, 0.5, 0.5, 0.25, 0.25},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewPercent[bool](tc.size, func(b bool) bool {
				return b
			})

			for i, val := range tc.addVals {
				res := p.Add(val)
				assert.InDelta(t, tc.wantRets[i], res, 1e-9)
			}
		})
	}
}
