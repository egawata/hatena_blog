package main

import (
	"fmt"
	"log"
	"reflect"

	"github.com/augustoroman/hexdump"
	"github.com/shamaton/msgpack/v2"
	"github.com/shamaton/msgpack/v2/ext"
)

type ProfileList struct {
	Profiles []Profile
}

type Profile struct {
	Name      string
	Languages map[string]struct{}
}

func main() {
	err := msgpack.AddExtCoder(&ProfileEncoder{}, &ProfileDecoder{})
	if err != nil {
		log.Fatal(err)
	}
	defer msgpack.RemoveExtCoder(&ProfileEncoder{}, &ProfileDecoder{})

	pl := &ProfileList{
		Profiles: []Profile{
			{
				Name:      "Alice",
				Languages: map[string]struct{}{"en": {}, "ja": {}},
			},
			{
				Name:      "Bob",
				Languages: map[string]struct{}{"en": {}, "fr": {}},
			},
		},
	}

	enc, err := msgpack.Marshal(pl)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(hexdump.Dump(enc))

	var dec ProfileList
	if err = msgpack.Unmarshal(enc, &dec); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("decoded: %#v\n", dec)
}

const extCodeProfile = 0x01

var typeProfile = reflect.TypeOf(Profile{})

type ProfileEncoder struct {
	ext.EncoderCommon
}

var _ ext.Encoder = (*ProfileEncoder)(nil)

func (e *ProfileEncoder) Code() int8 {
	return extCodeProfile
}

func (e *ProfileEncoder) Type() reflect.Type {
	return typeProfile
}

func (e *ProfileEncoder) CalcByteSize(value reflect.Value) (int, error) {
	p := value.Interface().(Profile)

	size := e.calcDataSize(p)

	// data length
	if size < (1 << 8) {
		size++
	} else if size < (1 << 16) {
		size += 2
	} else {
		size += 4
	}

	size++ // ext code

	return size, nil
}

// データ本体のサイズを計算
func (e *ProfileEncoder) calcDataSize(p Profile) int {
	size := len(p.Name) + 1
	size++                       // length of language
	size += len(p.Languages) * 2 // 2 bytes for each language

	return size
}

func (e *ProfileEncoder) WriteToBytes(value reflect.Value, offset int, bytes *[]byte) int {
	p := value.Interface().(Profile)

	dataLen := e.calcDataSize(p)
	if dataLen < (1 << 8) {
		offset = e.SetByte1Int(0xc7, offset, bytes)
		offset = e.SetByte1Int(dataLen, offset, bytes)
	} else if dataLen < (1 << 16) {
		offset = e.SetByte1Int(0xc8, offset, bytes)
		offset = e.SetByte2Int(dataLen, offset, bytes)
	} else {
		offset = e.SetByte1Int(0xc9, offset, bytes)
		offset = e.SetByte4Int(dataLen, offset, bytes)
	}
	offset = e.SetByte1Int(extCodeProfile, offset, bytes)

	// Name
	offset = e.SetByte1Int(len(p.Name), offset, bytes)
	offset = e.SetBytes([]byte(p.Name), offset, bytes)

	// Languages
	offset = e.SetByte1Int(len(p.Languages), offset, bytes)
	for l := range p.Languages {
		offset = e.SetBytes([]byte(l[:2]), offset, bytes)
	}

	return offset
}

type ProfileDecoder struct {
	ext.DecoderCommon
}

var _ ext.Decoder = (*ProfileDecoder)(nil)

func (d *ProfileDecoder) Code() int8 {
	return extCodeProfile
}

func (d *ProfileDecoder) IsType(offset int, bytes *[]byte) bool {
	code, _ := d.ReadSize1(offset, bytes)
	ext := byte(extCodeProfile)

	switch code {
	case 0xc7:
		return offset+2 < len(*bytes) && (*bytes)[offset+2] == ext
	case 0xc8:
		return offset+3 < len(*bytes) && (*bytes)[offset+3] == ext
	case 0xc9:
		return offset+5 < len(*bytes) && (*bytes)[offset+5] == ext
	}
	return false
}

func (d *ProfileDecoder) AsValue(offset int, k reflect.Kind, bytes *[]byte) (interface{}, int, error) {
	p := Profile{}

	// Data より前の部分は今回は不要なのでスキップ
	code, o := d.ReadSize1(offset, bytes)
	offset = o
	switch code {
	case 0xc7:
		_, offset = d.ReadSize1(offset, bytes) // skip data length
	case 0xc8:
		_, offset = d.ReadSize2(offset, bytes) // skip data length
	case 0xc9:
		_, offset = d.ReadSize4(offset, bytes) // skip data length
	default:
		return nil, offset, fmt.Errorf("invalid code: %x", code)
	}
	_, offset = d.ReadSize1(offset, bytes) // skip ext code

	// ここからデータ本体

	// Name
	nameLen, o := d.ReadSize1(offset, bytes)
	offset = o

	name, o := d.ReadSizeN(offset, int(nameLen), bytes)
	offset = o
	p.Name = string(name)

	// Languages
	langLen, o := d.ReadSize1(offset, bytes)
	offset = o
	p.Languages = make(map[string]struct{}, int(langLen))

	for range int(langLen) {
		lng, o := d.ReadSize2(offset, bytes)
		offset = o
		p.Languages[string(lng)] = struct{}{}
	}

	return p, offset, nil
}
