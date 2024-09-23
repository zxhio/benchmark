# Benchmark for message pack

By comparing the Benchmark results below, if you are communicating data very frequently, such as sending pcap data to the server, you can use a custom binary protocol to do so, not only in terms of speed, but also in terms of memory allocation with more room for optimization.


```shell
goos: linux
goarch: amd64
pkg: benchmark/pack
cpu: 12th Gen Intel(R) Core(TM) i7-12700
BenchmarkBinaryPack/encode#72-20                  72744334        187.1 ns/op         144 B/op      2 allocs/op
BenchmarkBinaryPack/encode#1024-20                17048832        660.6 ns/op        1200 B/op      2 allocs/op
BenchmarkBinaryPack/encode#16384-20                2085050         6280 ns/op       18495 B/op      2 allocs/op
BenchmarkBinaryPack/encode_with_pool#72-20        34700313        109.2 ns/op          64 B/op      2 allocs/op
BenchmarkBinaryPack/encode_with_pool#1024-20      39370662        101.1 ns/op          64 B/op      2 allocs/op
BenchmarkBinaryPack/encode_with_pool#16384-20     18445262        177.2 ns/op          64 B/op      2 allocs/op
BenchmarkBinaryPack/encode_to#72-20              705428736        16.96 ns/op           0 B/op      0 allocs/op
BenchmarkBinaryPack/encode_to#1024-20            575312358        20.78 ns/op           0 B/op      0 allocs/op
BenchmarkBinaryPack/encode_to#16384-20           100000000        113.4 ns/op           0 B/op      0 allocs/op
BenchmarkBinaryPack/decode_meta#72-20           1000000000        2.887 ns/op           0 B/op      0 allocs/op
BenchmarkBinaryPack/decode_meta#1024-20         1000000000        2.882 ns/op           0 B/op      0 allocs/op
BenchmarkBinaryPack/decode_meta#16384-20        1000000000        2.876 ns/op           0 B/op      0 allocs/op
BenchmarkBinaryPack/decode#72-20                 100000000        85.63 ns/op          80 B/op      1 allocs/op
BenchmarkBinaryPack/decode#1024-20                 7252350        445.4 ns/op        1024 B/op      1 allocs/op
BenchmarkBinaryPack/decode#16384-20                 554329         5499 ns/op       16384 B/op      1 allocs/op
BenchmarkBinaryPack/decode_with_pool#72-20       109352595        33.97 ns/op          16 B/op      1 allocs/op
BenchmarkBinaryPack/decode_with_pool#1024-20      85589674        36.27 ns/op          16 B/op      1 allocs/op
BenchmarkBinaryPack/decode_with_pool#16384-20     26163607        140.4 ns/op          16 B/op      1 allocs/op

BenchmarkMsgPack/encode#72-20                     10466427         1199 ns/op         688 B/op      8 allocs/op
BenchmarkMsgPack/encode#1024-20                    6599528         2132 ns/op        1585 B/op      8 allocs/op
BenchmarkMsgPack/encode#16384-20                   1478127         8806 ns/op       18879 B/op      8 allocs/op
BenchmarkMsgPack/encode_with_buf#72-20            26677507        388.2 ns/op         192 B/op      4 allocs/op
BenchmarkMsgPack/encode_with_buf#1024-20          31426809        400.2 ns/op         192 B/op      4 allocs/op
BenchmarkMsgPack/encode_with_buf#16384-20         22588560        494.5 ns/op         192 B/op      4 allocs/op
BenchmarkMsgPack/decode#72-20                     19894509        654.2 ns/op         280 B/op     10 allocs/op
BenchmarkMsgPack/decode#1024-20                   18211321        664.0 ns/op         280 B/op     10 allocs/op
BenchmarkMsgPack/decode#16384-20                  13755824        769.1 ns/op         280 B/op     10 allocs/op

BenchmarkJsonPack/encode#72-20                    17315143        647.1 ns/op         320 B/op      3 allocs/op
BenchmarkJsonPack/encode#1024-20                   4616841         2835 ns/op        1666 B/op      3 allocs/op
BenchmarkJsonPack/encode#16384-20                   365313        34289 ns/op       24754 B/op      3 allocs/op
BenchmarkJsonPack/encode_with_buf#72-20           24820188        447.4 ns/op         128 B/op      2 allocs/op
BenchmarkJsonPack/encode_with_buf#1024-20         13139395        910.6 ns/op         128 B/op      2 allocs/op
BenchmarkJsonPack/encode_with_buf#16384-20         1414260         8472 ns/op         128 B/op      2 allocs/op
BenchmarkJsonPack/decode#72-20                     8699952         1364 ns/op         304 B/op      8 allocs/op
BenchmarkJsonPack/decode#1024-20                   2103712         5605 ns/op        1384 B/op      8 allocs/op
BenchmarkJsonPack/decode#16384-20                   159140        73101 ns/op       18664 B/op      8 allocs/op

BenchmarkJsonCompressPack/encode#72-20               19934       709224 ns/op     1208429 B/op     26 allocs/op
BenchmarkJsonCompressPack/encode#1024-20             17577       766349 ns/op     1212782 B/op     26 allocs/op
BenchmarkJsonCompressPack/encode#16384-20            11757       860371 ns/op     1253975 B/op     25 allocs/op
BenchmarkJsonCompressPack/decode#72-20              490164        28972 ns/op       42048 B/op     15 allocs/op
BenchmarkJsonCompressPack/decode#1024-20            187113        71612 ns/op       47640 B/op     23 allocs/op
BenchmarkJsonCompressPack/decode#16384-20            35790       346580 ns/op      173352 B/op     30 allocs/op
PASS
ok      benchmark/pack  478.757s
```