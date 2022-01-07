package ico

import (
	"bytes"
	"bufio"
	"image"
	"image/png"
//	"image/bmp"
	"encoding/binary"
)

// ICO file data from https://en.wikipedia.org/wiki/ICO_(file_format)

const ICO = 1
const CUR = 2

type iconDir struct {
	_          uint16 // Reserved, must always be 0.
	iconType   uint16 // 1 for icon (.ICO) image, 2 for cursor (.CUR) image.
	                  // Other values are invalid
	imageCount uint16 // Number of images in the file.
}

type iconDirEntry struct {
	width    uint8  // image width, 0 means 256 pixel width
	height   uint8  // image height, 0 means 256 pixel height
	colors   uint8  // Number of colors in the palette, 0 if palette isn't used
	_        uint8  // Reserved, should be 0
	planes   uint16 // ICO format: Color planes, should be 0 or 1
	                // CUR format: Hotspot horizontal offset in pixels
	bpp      uint16 // ICO format: Bits per pixel
	                // CUR format: Hotspot vertical offset in pixels
	size     uint32 // Size of the image in bytes
	offset   uint32 // Offset of the BMP or PNG data from the beginning of the file
}

// NOTE:
// Setting the color planes to 0 or 1 is treated equivalently by the operating
// system, but if the color planes are set higher than 1, this value should be
// multiplied by the bits per pixel to determine the final color depth of the
// image. It is unknown if the various Windows operating system versions are
// resilient to different color plane values.

type Icon struct {
	entries []image.Image
}

func (self *Icon) Encode() ([]byte, error) {
	header := new(bytes.Buffer)
	bitmaps := new(bytes.Buffer)

	var count = len(self.entries)

	err := binary.Write(header, binary.LittleEndian, iconDir{0, ICO, uint16(count)})

	if err != nil {
		return nil, err
	}

	writer := bufio.NewWriter(bitmaps)

	var offset = header.Len() + (binary.Size(iconDirEntry{}) * count)

	for index := range self.entries {
		bounds := self.entries[index].Bounds()

		var size = bitmaps.Len()

		// FIXME support BMP as well
		err := png.Encode(writer, self.entries[index])
		writer.Flush()

		size = bitmaps.Len() - size

		if err != nil {
			return nil, err
		}

		entry := iconDirEntry{
			uint8(bounds.Dx()),
			uint8(bounds.Dy()),
			0,                // FIXME Colors for paletted images
			0,
			1,                // XXX 0 or 1 is expected, anything else is unknown
			32,               // XXX no clue if this is true, or if it matters for PNG based icons
			uint32(size),
			uint32(offset),
		}

		binary.Write(header, binary.LittleEndian, entry)
		offset += bitmaps.Len()
	}

	return bytes.Join([][]byte{header.Bytes(),bitmaps.Bytes()}, []byte{}), nil
}

func (self *Icon) AddImage(image image.Image) {
	self.entries = append(self.entries, image)
}
