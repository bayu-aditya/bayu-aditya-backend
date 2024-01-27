package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArrayAppendHead(t *testing.T) {
	type args[T comparable] struct {
		slice      []T
		newElement T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "insert int",
			args: args[int]{
				slice:      []int{1, 2, 3},
				newElement: 4,
			},
			want: []int{4, 1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ArrayAppendHead(tt.args.slice, tt.args.newElement)
			assert.Equal(t, tt.want, got)
		})
	}
}
