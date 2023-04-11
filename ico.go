// Copyright 2022 Staffan Thomen <staffan@thomen.fi>
// Use of this source code is governed by a BSD-style license found in the
// adjacent LICENSE file.
package ico

import "image"

// ICO file data from https://en.wikipedia.org/wiki/ICO_(file_format)

// file type, icon or cursor
const ICO = 1
const CUR = 2

// determine subimage type
const PNG = 0
const BMP = 1

const NOALPHA = -1

type iconDir struct {
	_        uint16 // Reserved, must always be 0.
	iconType uint16 // 1 for icon (.ICO) image, 2 for cursor (.CUR) image.
	// Other values are invalid
	imageCount uint16 // Number of images in the file.
}

type iconDirEntry struct {
	width  uint8  // image width, 0 means 256 pixel width
	height uint8  // image height, 0 means 256 pixel height
	colors uint8  // Number of colors in the palette, 0 if palette isn't used
	_      uint8  // Reserved, should be 0
	planes uint16 // ICO format: Color planes, should be 0 or 1
	// CUR format: Hotspot horizontal offset in pixels
	bpp uint16 // ICO format: Bits per pixel
	// CUR format: Hotspot vertical offset in pixels
	size   uint32 // Size of the image in bytes
	offset uint32 // Offset of the BMP or PNG data from the beginning of the file
}

// NOTE:
// Setting the color planes to 0 or 1 is treated equivalently by the operating
// system, but if the color planes are set higher than 1, this value should be
// multiplied by the bits per pixel to determine the final color depth of the
// image. It is unknown if the various Windows operating system versions are
// resilient to different color plane values.

// Abstract representation of the ICO/CUR file
type icon struct {
	Type int

	Entries []struct {
		Type       int
		Image      image.Image
		AlphaIndex int
	}
}

// Create a new icon
func New(t int) *icon {
	var item = new(icon)
	item.Type = t
	return item
}

// Create a new icon object
func NewIcon() *icon {
	return New(ICO)
}

// Create a new cursor object
func NewCursor() *icon {
	return New(CUR)
}

// Add a png to the image
func (self *icon) AddPng(img image.Image) {
	self.Add(PNG, img, NOALPHA)
}

// Add a bitmap to the image
func (self *icon) AddBmp(img image.Image) {
	self.Add(BMP, img, NOALPHA)
}

// Add a bitmap to the image with an alpha
func (self *icon) AddBmpAlpha(img image.Image, alphaIndex int) {
	self.Add(BMP, img, alphaIndex)
}

// Add an image to the icon structure
func (self *icon) Add(t int, img image.Image, alphaIndex int) {
	var entry = struct {
		Type       int
		Image      image.Image
		AlphaIndex int
	}{
		t,
		img,
		alphaIndex,
	}

	self.Entries = append(self.Entries, entry)
}
