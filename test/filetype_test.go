package main

import (
	"io/ioutil"
	"net/http"
	"testing"
)

// go test -benchmem -run=^$ -bench ^BenchmarkType$
func BenchmarkType(b *testing.B) {
	byt, _ := ioutil.ReadFile("../example/james-webb.png")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		http.DetectContentType(byt)
	}
	/*
		goos: android
		goarch: arm64
		BenchmarkType-4      4747869           245.5 ns/op         0 B/op          0 allocs/op
		PASS
		ok      _/sdcard/2/golang/test  1.449s

		goos: windows
		goarch: amd64
		cpu: Intel(R) Core(TM) i3-10100F CPU @ 3.60GHz
		BenchmarkType-8   	 6011888	       205.9 ns/op	       0 B/op	       0 allocs/op
		PASS
		ok  	_/e_/1/github/gg/test	1.461
	*/

}
