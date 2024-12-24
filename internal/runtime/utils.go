package runtime

import "golang.org/x/exp/constraints"

func add[T constraints.Ordered](v1 T, v2 T) T {
	return v1 + v2
}
