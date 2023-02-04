package main

import (
	"bytes"
	"image/png"
	"io/ioutil"
	"testing"

	"github.com/FloatTech/gg"
)

// go test -benchmem -bench .
// go test -benchmem -run=^$ -bench ^BenchmarkByte$

func BenchmarkByte(b *testing.B) {
	//by, _ := ioutil.ReadFile("../examples/james-webb.png")
	by, _ := ioutil.ReadFile("../examples/gopher.png")

	b.ResetTimer() //重置时间
	for i := 0; i < b.N; i++ {
		gg.GetPNGWH(by)
	}
	/*
	   goos: windows
	   goarch: amd64
	   cpu: Intel(R) Core(TM) i3-10100F CPU @ 3.60GHz
	   BenchmarkByte-8   	1000000000	         0.8591 ns/op	       0 B/op	       0 allocs/op
	   PASS
	   ok  	_/e_/1/github/gg/test	0.978s
	*/
}

func BenchmarkByte2(b *testing.B) {
	//by, _ := ioutil.ReadFile("../examples/james-webb.png")
	by, _ := ioutil.ReadFile("../examples/gopher.png")

	b.ResetTimer() //重置时间
	for i := 0; i < b.N; i++ {
		png.DecodeConfig(bytes.NewReader(by))
	}
	/*
	   goos: windows
	   goarch: amd64
	   cpu: Intel(R) Core(TM) i3-10100F CPU @ 3.60GHz
	   BenchmarkByte2-8   	 3977091	       298.1 ns/op	    1088 B/op	       3 allocs/op
	   PASS
	   ok  	_/e_/1/github/gg/test	1.514s
	*/
}
