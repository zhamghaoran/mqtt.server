package pool

import (
	"github.com/panjf2000/ants/v2"
)

var Pool *ants.Pool

func init() {
	Pool, _ = ants.NewPool(10000)
}
func Submit(fun func()) {
	err := Pool.Submit(fun)
	if err != nil {
		return
	}
}
