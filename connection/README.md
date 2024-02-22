# Benchmark for connection

Benchmark for conn read/write.

```shell
$ go test -v -benchtime=10s  -benchmem -run=^$ -bench ^BenchmarkConn .
goos: linux
goarch: amd64
pkg: benchmark/connection
cpu: 12th Gen Intel(R) Core(TM) i7-12700
BenchmarkConnCmux
BenchmarkConnCmux/magic
BenchmarkConnCmux/magic-20                917902             12348 ns/op        10614.63 MB/s          0 B/op          0 allocs/op
BenchmarkConnCmux/token
BenchmarkConnCmux/token-20                942536             12300 ns/op        10656.67 MB/s          0 B/op          0 allocs/op
BenchmarkConnCmux/tls
BenchmarkConnCmux/tls/bare
BenchmarkConnCmux/tls/bare-20             291968             40402 ns/op        3244.23 MB/s         191 B/op          7 allocs/op
BenchmarkConnCmux/tls/magic
BenchmarkConnCmux/tls/magic-20            298700             40465 ns/op        3239.15 MB/s         191 B/op          7 allocs/op
BenchmarkConnCmux/any
BenchmarkConnCmux/any-20                 1000000             12418 ns/op        10554.98 MB/s          0 B/op          0 allocs/op
BenchmarkConnSmux
BenchmarkConnSmux-20                      473323             24873 ns/op        5269.67 MB/s        1396 B/op         26 allocs/op
BenchmarkConnTCP
BenchmarkConnTCP-20                       926667             11862 ns/op        11049.71 MB/s          0 B/op          0 allocs/op
BenchmarkConnTLS
BenchmarkConnTLS-20                       291843             40619 ns/op        3226.87 MB/s         192 B/op          7 allocs/op
PASS
ok      benchmark/connection    95.894s
```

Benchmark for echo.

```shell
$ go test -v -benchtime=10s  -benchmem -run=^$ -bench ^BenchmarkEcho .
goos: linux
goarch: amd64
pkg: benchmark/connection
cpu: 12th Gen Intel(R) Core(TM) i7-12700
BenchmarkEchoCmux
BenchmarkEchoCmux-20               79078            167912 ns/op         780.60 MB/s       34006 B/op         26 allocs/op
BenchmarkEchoFastHTTP
BenchmarkEchoFastHTTP-20          111536             98882 ns/op        1325.54 MB/s       12810 B/op         41 allocs/op
BenchmarkEchoNetHTTP
BenchmarkEchoNetHTTP-20            44312            248886 ns/op         526.63 MB/s      340752 B/op         58 allocs/op
BenchmarkEchoSmuxTCP
BenchmarkEchoSmuxTCP-20           159150             69924 ns/op        1874.50 MB/s       79812 B/op         85 allocs/op
BenchmarkEchoSmuxTLS
BenchmarkEchoSmuxTLS-20           103567            113524 ns/op        1154.57 MB/s       80941 B/op        102 allocs/op
BenchmarkEchoTCP
BenchmarkEchoTCP-20               304132             40183 ns/op        3261.86 MB/s        1079 B/op         23 allocs/op
BenchmarkEchoTLS
BenchmarkEchoTLS-20                10000           2021083 ns/op          64.85 MB/s      242498 B/op       1112 allocs/op
PASS
ok      benchmark/connection    108.726s
```