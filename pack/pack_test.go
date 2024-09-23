package pack

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// makeRandomN 生成指定长度的随机字符串
func makeRandomN(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return b
}

var (
	rawDataSmall = makeRandomN(72) // SYN 的大小
	smallPacket  = CapturePacket{
		CaptureInfo: CaptureInfo{
			Timestamp:      time.UnixMicro(time.Now().UnixMicro()),
			CaptureLength:  len(rawDataSmall),
			Length:         len(rawDataSmall),
			InterfaceIndex: 7,
		},
		Id:   33,
		Data: rawDataSmall,
	}

	rawDataMiddle = makeRandomN(1024)
	middlePacket  = CapturePacket{
		CaptureInfo: CaptureInfo{
			Timestamp:      time.UnixMicro(time.Now().UnixMicro()),
			CaptureLength:  len(rawDataMiddle),
			Length:         len(rawDataMiddle),
			InterfaceIndex: 7,
		},
		Id:   33,
		Data: rawDataMiddle,
	}

	rawDataLarge = makeRandomN(1024 * 16)
	largePacket  = CapturePacket{
		CaptureInfo: CaptureInfo{
			Timestamp:      time.UnixMicro(time.Now().UnixMicro()),
			CaptureLength:  len(rawDataLarge),
			Length:         len(rawDataLarge),
			InterfaceIndex: 7,
		},
		Id:   33,
		Data: rawDataLarge,
	}

	packets = []CapturePacket{smallPacket, middlePacket, largePacket}
)

func TestBinaryPack(t *testing.T) {
	data := BinaryPack.Encode(&smallPacket)
	assert.Equal(t, len(rawDataSmall)+CapturePacketMetaLen, len(data), "encode failed")
	assert.Equal(t, data[CapturePacketMetaLen:], rawDataSmall, "invalid raw data")
	t.Logf("binary pack data length=%d\n", len(data))

	var pd CapturePacket
	BinaryPack.Decode(data, &pd)
	assert.Equal(t, pd.CaptureInfo, smallPacket.CaptureInfo, "invalid capture info")
	assert.Equal(t, pd.Id, smallPacket.Id, "invalid id")
	assert.Equal(t, pd.Data, smallPacket.Data, "invalid data")
}

func BenchmarkBinaryPack(b *testing.B) {
	b.ReportAllocs()

	for _, p := range packets {
		b.Run("encode#"+strconv.Itoa(len(p.Data)), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				BinaryPack.Encode(&p)
			}
		})
	}

	for _, p := range packets {
		b.Run("encode_to#"+strconv.Itoa(len(p.Data)), func(b *testing.B) {
			b.ResetTimer()
			buf := bytes.NewBuffer(make([]byte, 0, 1024*32))
			for i := 0; i < b.N; i++ {
				buf.Reset()
				BinaryPack.EncodeTo(&p, buf)
			}
		})
	}

	for _, p := range packets {
		b.Run("decode_meta#"+strconv.Itoa(len(p.Data)), func(b *testing.B) {
			b.ResetTimer()
			data := BinaryPack.Encode(&p)
			var p CapturePacket
			for i := 0; i < b.N; i++ {
				BinaryPack.DecodeMeta(data, &p)
			}
		})
	}

	for _, p := range packets {
		b.ResetTimer()
		data := BinaryPack.Encode(&p)
		var p CapturePacket
		for i := 0; i < b.N; i++ {
			BinaryPack.Decode(data, &p)
		}
	}
}

func TestMsgPack(t *testing.T) {
	data, err := msgpack.Marshal(smallPacket)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("msgpack data length=%d\n", len(data))

	var p CapturePacket
	err = msgpack.Unmarshal(data, &p)
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkMsgPack(b *testing.B) {
	b.ReportAllocs()

	for _, p := range packets {
		b.Run("encode#"+strconv.Itoa(len(p.Data)), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				msgpack.Marshal(p)
			}
		})
	}

	for _, p := range packets {
		b.Run("encode_with_buf#"+strconv.Itoa(len(p.Data)), func(b *testing.B) {
			b.ResetTimer()
			buf := bytes.NewBuffer(make([]byte, 0, 1024*32))
			for i := 0; i < b.N; i++ {
				buf.Reset()
				msgpack.NewEncoder(buf).Encode(p)
			}
		})
	}

	for _, p := range packets {
		b.Run("decode#"+strconv.Itoa(len(p.Data)), func(b *testing.B) {
			data, err := msgpack.Marshal(p)
			if err != nil {
				b.Fatal(err)
			}
			b.ResetTimer()
			var p CapturePacket
			for i := 0; i < b.N; i++ {
				msgpack.Unmarshal(data, &p)
			}
		})
	}
}

func TestJsonPack(t *testing.T) {
	data, err := json.Marshal(smallPacket)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("json pack data length=%d\n", len(data))
}

func BenchmarkJsonPack(b *testing.B) {
	b.ReportAllocs()

	for _, p := range packets {
		b.Run("encode#"+strconv.Itoa(len(p.Data)), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				json.Marshal(p)
			}
		})
	}

	for _, p := range packets {
		b.Run("encode_with_buf#"+strconv.Itoa(len(p.Data)), func(b *testing.B) {
			b.ResetTimer()
			buf := bytes.NewBuffer(make([]byte, 0, 1024*32))
			for i := 0; i < b.N; i++ {
				buf.Reset()
				json.NewEncoder(buf).Encode(p)
			}
		})
	}

	for _, p := range packets {
		b.Run("decode#"+strconv.Itoa(len(p.Data)), func(b *testing.B) {
			data, err := json.Marshal(p)
			if err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()
			var p CapturePacket
			for i := 0; i < b.N; i++ {
				json.Unmarshal(data, &p)
			}
		})
	}
}

func TestJsonCompressPack(t *testing.T) {
	data, err := JsonCompressPack{}.Encode(&smallPacket)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("json compress pack data length=%d\n", len(data))

	var pd CapturePacket
	JsonCompressPack{}.Decode(data, &pd)
	assert.Equal(t, pd.CaptureInfo, smallPacket.CaptureInfo, "invalid capture info")
	assert.Equal(t, pd.Id, smallPacket.Id, "invalid id")
	assert.Equal(t, pd.Data, smallPacket.Data, "invalid data")
}

func BenchmarkJsonCompressPack(b *testing.B) {
	b.ReportAllocs()

	for _, p := range packets {
		b.Run("encode#"+strconv.Itoa(len(p.Data)), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				JsonCompressPack{}.Encode(&p)
			}
		})
	}

	for _, p := range packets {
		b.Run("decode#"+strconv.Itoa(len(p.Data)), func(b *testing.B) {
			data, err := JsonCompressPack{}.Encode(&p)
			if err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()
			var p CapturePacket
			for i := 0; i < b.N; i++ {
				JsonCompressPack{}.Decode(data, &p)
			}
		})
	}
}
