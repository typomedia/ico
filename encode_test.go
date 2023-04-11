package ico

import (
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
	"image/draw"
	"testing"
)

func TestEncodingAnEmptyIcon(t *testing.T) {
	facit := []byte{
		0x00, 0x00, // reserved
		0x01, 0x00, // type: 1
		0x00, 0x00, // count: 0
	}
	icon := NewIcon()
	bytes, err := icon.Encode()

	assert.Nil(t, err, "Encoding a blank ico should not produce an error")
	assert.NotNil(t, bytes, "The returned bytes must not be nil")
	assert.Equal(t, facit, bytes, "Encoding a blank ico should match the expected formt")
}

func TestEncodingAnEmptyCursor(t *testing.T) {
	facit := []byte{
		0x00, 0x00, // reserved
		0x02, 0x00, // type: 1
		0x00, 0x00, // count: 0
	}
	cursor := NewCursor()
	bytes, err := cursor.Encode()

	assert.Nil(t, err, "Encoding a blank ico should not produce an error")
	assert.NotNil(t, bytes, "The returned bytes must not be nil")
	assert.Equal(t, facit, bytes, "Encoding a blank ico should match the expected formt")
}

func TestEncodingAnIcoWithASinglePNG(t *testing.T) {
	icon := NewIcon()

	img := image.NewRGBA(image.Rect(0, 0, 1, 1))

	green := color.RGBA{0, 255, 0, 255}

	draw.Draw(img, img.Bounds(), &image.Uniform{green}, image.ZP, draw.Src)

	icon.AddPng(img)

	bytes, err := icon.Encode()

	assert.Nil(t, err, "Encoding should not produce an error")
	assert.NotNil(t, bytes, "The returned bytes must not be nil")

	dir_facit := []byte{
		0x00, 0x00, // reserved
		0x01, 0x00, // type
		0x01, 0x00, // count: 1
		0x01,       // entry width
		0x01,       // entyr height
		0x00,       // colors, palette is not used
		0x00,       // reserved
		0x01, 0x00, // ico color planes
		0x20, 0x00, // 32 bpp
	}

	assert.Equal(t, dir_facit, bytes[0:14], "The header should match the expected facit")
	assert.Equal(t, uint32(22), binary.LittleEndian.Uint32(bytes[18:22]), "The offset of the image should be just after the header")
}

func TestEncodingAnIcoWithTwoPNGs(t *testing.T) {
	icon := NewIcon()

	for _, color := range []color.RGBA{{0, 255, 0, 255}, {255, 0, 0, 255}} {
		img := image.NewRGBA(image.Rect(0, 0, 100, 100))
		draw.Draw(img, img.Bounds(), &image.Uniform{color}, image.ZP, draw.Src)

		icon.AddPng(img)
	}

	bytes, err := icon.Encode()

	assert.Nil(t, err, "Encoding should not produce an error")
	assert.NotNil(t, bytes, "The returned bytes must not be nil")

	dir_facit := []byte{
		0x00, 0x00, // reserved
		0x01, 0x00, // type
		0x02, 0x00, // count: 2
	}

	entry_facit := []byte{
		0x64,       // entry width
		0x64,       // entyr height
		0x00,       // colors, palette is not used
		0x00,       // reserved
		0x01, 0x00, // ico color planes
		0x20, 0x00, // 32 bpp
	}

	assert.Equal(t, dir_facit, bytes[0:6], "The directory should match the expected facit")
	assert.Equal(t, entry_facit, bytes[6:14], "The first image's directory entry should match the facit")
	assert.Equal(t, entry_facit, bytes[22:30], "The second image's directory entry should match the facit")
	png := []byte{0x89, 0x50, 0x4E, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}

	var first_offset uint32 = binary.LittleEndian.Uint32(bytes[18:22])

	assert.Equal(t, png, bytes[first_offset:first_offset+8], "No PNG header found at the first images offset")

	var second_offset uint32 = binary.LittleEndian.Uint32(bytes[34:38])

	assert.Equal(t, png, bytes[second_offset:second_offset+8], "No PNG header found at the second images offset")
}
