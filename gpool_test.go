package zgo

import (
	"testing"
)

func BenchmarkMakeGPool(b *testing.B) {
	gp := MakeGPool(1000)
	work := 1024
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if ok := gp.Get(work); !ok {
			go func() {
				// var work int
				//fmt.Println(work)
				_, ok := gp.Put()
				if !ok {
					return
				}
				// work = i.(int)
			}()
		}
	}
}
