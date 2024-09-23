package pack

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"sync"
	"time"
)

// from github.com/google/gopacket
type CaptureInfo struct {
	// Timestamp is the time the packet was captured, if that is known.
	Timestamp time.Time `json:"ts" msgpack:"ts"`
	// CaptureLength is the total number of bytes read off of the wire.
	CaptureLength int `json:"cap_len" msgpack:"cap_len"`
	// Length is the size of the original packet.  Should always be >=
	// CaptureLength.
	Length int `json:"len" msgpack:"len"`
	// InterfaceIndex
	InterfaceIndex int `json:"iface_idx" msgpack:"iface_idx"`
}

type CapturePacket struct {
	CaptureInfo
	Id   uint32 `json:"id" msgpack:"id"`
	Data []byte `json:"data" msgpack:"data"`
}

const (
	CapturePacketMetaLen = 22
)

// Reduce packet meta memory allocation.
var (
	metaBufPool   = sync.Pool{New: func() interface{} { return new([CapturePacketMetaLen]byte) }}
	smallBufPool  = sync.Pool{New: func() interface{} { return new([CapturePacketMetaLen + 128]byte) }}
	midBufPool    = sync.Pool{New: func() interface{} { return new([CapturePacketMetaLen + 1024]byte) }}
	largeBufPool  = sync.Pool{New: func() interface{} { return new([CapturePacketMetaLen + 8192]byte) }}
	xlargeBufPool = sync.Pool{New: func() interface{} { return new([CapturePacketMetaLen + 65536]byte) }}
)

func acquirePacketBuf(n int) ([]byte, func()) {
	var (
		buf   []byte
		putfn func()
	)
	if n <= CapturePacketMetaLen+128 {
		smallBuf := smallBufPool.Get().(*[CapturePacketMetaLen + 128]byte)
		buf = smallBuf[:0]
		putfn = func() { smallBufPool.Put(smallBuf) }
	} else if n <= CapturePacketMetaLen+1024 {
		midBuf := midBufPool.Get().(*[CapturePacketMetaLen + 1024]byte)
		buf = midBuf[:0]
		putfn = func() { midBufPool.Put(midBuf) }
	} else if n <= CapturePacketMetaLen+8192 {
		largeBuf := largeBufPool.Get().(*[CapturePacketMetaLen + 8192]byte)
		buf = largeBuf[:0]
		putfn = func() { largeBufPool.Put(largeBuf) }
	} else {
		xlargeBuf := xlargeBufPool.Get().(*[CapturePacketMetaLen + 65536]byte)
		buf = xlargeBuf[:0]
		putfn = func() { xlargeBufPool.Put(xlargeBuf) }
	}
	return buf, putfn
}

// Assuming these int values does not exceed 65535，the size can be reduced by a few bytes.
type binaryPack struct{}

var BinaryPack binaryPack

func (bp binaryPack) Encode(p *CapturePacket) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, CapturePacketMetaLen+len(p.Data)))
	bp.EncodeTo(p, buf)
	return buf.Bytes()
}

func (bp binaryPack) EncodeWithPool(p *CapturePacket) ([]byte, func()) {
	b, putfn := acquirePacketBuf(CapturePacketMetaLen + len(p.Data))
	buf := bytes.NewBuffer(b)
	bp.EncodeTo(p, buf)
	return buf.Bytes(), putfn
}

// Write encoded data directly without allocating memory.
// So at the calling point, this writer can be reused.
func (binaryPack) EncodeTo(p *CapturePacket, w io.Writer) (int, error) {
	buf := metaBufPool.Get().(*[CapturePacketMetaLen]byte)
	defer metaBufPool.Put(buf)

	binary.BigEndian.PutUint64(buf[0:], uint64(p.Timestamp.UnixMicro()))
	binary.BigEndian.PutUint16(buf[8:], uint16(p.CaptureLength))
	binary.BigEndian.PutUint16(buf[12:], uint16(p.Length)) // if _CaptureLength_ not exceed 65535, use [10:]
	binary.BigEndian.PutUint16(buf[16:], uint16(p.InterfaceIndex))
	binary.BigEndian.PutUint16(buf[18:], uint16(p.Id))

	nm, err := w.Write(buf[:])
	if err != nil {
		return 0, err
	}
	nd, err := w.Write(p.Data)
	return nm + nd, err
}

func (bp binaryPack) Decode(data []byte, p *CapturePacket) error {
	err := bp.DecodeMeta(data, p)
	if err != nil {
		return err
	}
	p.Data = make([]byte, len(data)-CapturePacketMetaLen)
	copy(p.Data, data[CapturePacketMetaLen:])
	return nil
}

func (bp binaryPack) DecodeWithPool(data []byte, p *CapturePacket) (func(), error) {
	err := bp.DecodeMeta(data, p)
	if err != nil {
		return nil, err
	}

	b, putfn := acquirePacketBuf(len(data[CapturePacketMetaLen:]))
	p.Data = b[:len(data[CapturePacketMetaLen:])]
	copy(p.Data, data[CapturePacketMetaLen:])
	return putfn, nil
}

func (binaryPack) DecodeMeta(data []byte, p *CapturePacket) error {
	if len(data) < CapturePacketMetaLen {
		return errors.New("invalid packet meta data")
	}
	p.Timestamp = time.UnixMicro(int64(binary.BigEndian.Uint64(data)))
	p.CaptureLength = int(binary.BigEndian.Uint16(data[8:]))
	p.Length = int(binary.BigEndian.Uint16(data[12:]))
	p.InterfaceIndex = int(binary.BigEndian.Uint16(data[16:]))
	p.Id = uint32(binary.BigEndian.Uint16(data[18:]))
	p.Data = nil
	return nil
}

type JsonCompressPack struct{}

func (jcp JsonCompressPack) Encode(p *CapturePacket) ([]byte, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(make([]byte, 0, len(p.Data)+CapturePacketMetaLen))
	w, _ := gzip.NewWriterLevel(buf, gzip.BestSpeed)
	_, err = w.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (jcp JsonCompressPack) Decode(data []byte, p *CapturePacket) error {
	b := bytes.NewBuffer(data)
	gr, err := gzip.NewReader(b)
	if err != nil {
		return err
	}
	defer gr.Close()

	decompressedData, err := io.ReadAll(gr)
	if err != nil {
		return err
	}
	return json.Unmarshal(decompressedData, p)
}
