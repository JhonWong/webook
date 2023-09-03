package lru

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestLRUCache(t *testing.T) {
	tests := []struct {
		capacity   int
		operations []string
		args       [][]int
		expected   []int
		expectedOk []bool
	}{
		{
			2,
			[]string{"Put", "Put", "Get", "Put", "Get", "Put", "Get", "Get", "Get"},
			[][]int{{1, 1}, {2, 2}, {1}, {3, 3}, {2}, {4, 4}, {1}, {3}, {4}},
			[]int{0, 0, 1, 0, 0, 0, 0, 3, 4},
			[]bool{false, false, true, false, false, false, false, true, true},
		},
	}

	for _, tt := range tests {
		cache := NewLRUCache[int](tt.capacity)
		for i, op := range tt.operations {
			arg := tt.args[i]
			expected := tt.expected[i]
			expectedOk := tt.expectedOk[i]
			switch op {
			case "Put":
				cache.Put(arg[0], arg[1])
			case "Get":
				ok, got := cache.Get(arg[0])
				assert.Equal(t, ok, expectedOk)
				if ok {
					assert.Equal(t, got, expected)
				}
			}
		}
	}
}
