# Benchmark for connection

Benchmark for conn read/write.

```shell
$ go test -v -benchtime=10s  -benchmem -run=^$ -bench ^BenchmarkConn .
goos: linux
goarch: amd64
pkg: benchmark/connection
cpu: 12th Gen Intel(R) Core(TM) i7-12700
BenchmarkConnCmux
BenchmarkConnCmux/MagicMatcher
BenchmarkConnCmux/MagicMatcher-20                     997550             11862 ns/op        11049.73 MB/s          0 B/op          0 allocs/op
BenchmarkConnCmux/TokenMatcher
BenchmarkConnCmux/TokenMatcher-20                     958461             11714 ns/op        11188.94 MB/s          0 B/op          0 allocs/op
BenchmarkConnCmux/TLSMatcher
BenchmarkConnCmux/TLSMatcher/TLS
BenchmarkConnCmux/TLSMatcher/TLS-20                   295111             40471 ns/op        3238.68 MB/s         192 B/op          7 allocs/op
BenchmarkConnCmux/TLSMatcher/MagicMatcher
BenchmarkConnCmux/TLSMatcher/MagicMatcher-20          296203             39566 ns/op        3312.75 MB/s         192 B/op          7 allocs/op
BenchmarkConnCmux/AnyMatcher
BenchmarkConnCmux/AnyMatcher-20                       932871             11870 ns/op        11041.90 MB/s          0 B/op          0 allocs/op
BenchmarkConnSmux
BenchmarkConnSmux/OverTCP
BenchmarkConnSmux/OverTCP-20                          438889             24703 ns/op        5305.97 MB/s        1380 B/op         26 allocs/op
BenchmarkConnSmux/OverTLS
BenchmarkConnSmux/OverTLS-20                          210336             57345 ns/op        2285.69 MB/s        1596 B/op         36 allocs/op
BenchmarkConnTCP
BenchmarkConnTCP-20                                   917894             12120 ns/op        10814.60 MB/s          0 B/op          0 allocs/op
BenchmarkConnTLS
BenchmarkConnTLS-20                                   292843             40310 ns/op        3251.57 MB/s         192 B/op          7 allocs/op
PASS
ok      benchmark/connection    106.287s
```

Benchmark for echo.

```shell
$ go test -v -benchtime=10s  -benchmem -run=^$ -bench ^BenchmarkEcho .
goos: linux
goarch: amd64
pkg: benchmark/connection
cpu: 12th Gen Intel(R) Core(TM) i7-12700
BenchmarkEchoCmux
BenchmarkEchoCmux-20                       83162            164356 ns/op         797.49 MB/s       34005 B/op         26 allocs/op
BenchmarkEchoFastHTTP
BenchmarkEchoFastHTTP-20                  144302             95231 ns/op        1376.36 MB/s       12941 B/op         41 allocs/op
BenchmarkEchoNetHTTP
BenchmarkEchoNetHTTP-20                    65124            239187 ns/op         547.99 MB/s      370816 B/op         59 allocs/op
BenchmarkEchoSmux
BenchmarkEchoSmux/OverTCP
BenchmarkEchoSmux/OverTCP-20              153706             70494 ns/op        1859.34 MB/s       79824 B/op         85 allocs/op
BenchmarkEchoSmux/OverTLS
BenchmarkEchoSmux/OverTLS-20              106585            112120 ns/op        1169.04 MB/s       81776 B/op        102 allocs/op
BenchmarkEchoTCP
BenchmarkEchoTCP-20                       308125             39266 ns/op        3338.05 MB/s        1078 B/op         23 allocs/op
BenchmarkEchoTLS
BenchmarkEchoTLS-20                        10000           1988704 ns/op          65.91 MB/s      241188 B/op       1112 allocs/op
PASS
ok      benchmark/connection    112.673s
```