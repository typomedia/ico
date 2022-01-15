package ico

import (
	"bytes"
	"bufio"
	"image"
	"image/png"
	"encoding/binary"
	bmp "golang.org/x/image/bmp"
)

// Encode the icon structure to a complete ICO file
//
func (self *icon) Encode() ([]byte, error) {
	header := new(bytes.Buffer)
	bitmaps := new(bytes.Buffer)

	var count = len(self.Entries)

	var dir = iconDir{
		0,
		uint16(self.Type),
		uint16(count),
	}

	err := binary.Write(header, binary.LittleEndian, dir)

	if err != nil {
		return nil, err
	}

	var offset = header.Len() + (binary.Size(iconDirEntry{}) * count)

	for _, entry := range self.Entries {
		bounds := entry.Image.Bounds()

		var size = bitmaps.Len()

		// FIXME extract colors and bits per pixels from entry.Image
		var colors = 0
		var bpp = 32

		var bitmap []byte
		var err error

		switch (entry.Type) {
			case BMP:
				bitmap, err = encodeBMP(entry.Image)
			case PNG:
				bitmap, err = encodePNG(entry.Image)
		}

		binary.Write(bitmaps, binary.LittleEndian, bitmap)

		if err != nil {
			return nil, err
		}

		size = bitmaps.Len() - size

		entry := iconDirEntry{
			uint8(bounds.Dx()),
			uint8(bounds.Dy()),
			uint8(colors),
			0,
			1,
			uint16(bpp),
			uint32(size),
			uint32(offset),
		}

		binary.Write(header, binary.LittleEndian, entry)
		offset += bitmaps.Len()
	}

	return bytes.Join([][]byte{header.Bytes(),bitmaps.Bytes()}, []byte{}), nil
}

func encodeBMP(img image.Image) ([]byte, error) {
	bitmap := new(bytes.Buffer)

	writer := bufio.NewWriter(bitmap)

	err := bmp.Encode(writer, img)

	if err != nil {
		return nil, err
	}

	writer.Flush()

	return bitmap.Bytes()[14:], nil
}

func encodePNG(img image.Image) ([]byte, error) {
	bitmap := new(bytes.Buffer)

	writer := bufio.NewWriter(bitmap)

	err := png.Encode(writer, img)

	if err != nil {
		return nil, err
	}

	writer.Flush()

	return bitmap.Bytes(), nil
}

