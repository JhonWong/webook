package percent

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPercent_Add(t *testing.T) {
	timeout := context.DeadlineExceeded

	testCases := []struct {
		name     string
		size     int
		addVals  []error
		wantRets []bool
	}{
		{
			name:     "test",
			size:     4,
			addVals:  []error{timeout, nil, nil, timeout, timeout, nil, timeout, timeout},
			wantRets: []bool{true, false, false, false, false, false, true, true},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewPercent(tc.size, func(err error) bool {
				return err == context.DeadlineExceeded
			}, 0.5)

			for i, val := range tc.addVals {
				res := p.Add(nil, val)
				assert.Equal(t, tc.wantRets[i], res)
			}
		})
	}
}
