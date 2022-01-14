ico
===

This is a simple module library for building windows ICO files that may some
day be complete.

The module features an abstraction type for icon files so they can be decoded
as well as encoded.

Features that I intend to include are:

 - Encoding of PNG images (like most golang ico modules)
 - Encoding of BMP images
 - Encoding of paletted images
 - Encoding of paletted images with alpha
 - Decoding of ICO files

Example usage
-------------

```
img := image.NewRGBA(image.Rect(0,0,32,32))

eudenil := color.RGBA{164, 184, 135}

draw.Draw(img, img.Bounds(), &image.Uniform{eudenil}, image.ZP, draw.Src)

icon := ico.NewIcon()
icon.AddPng(img)

bytes, err := icon.Encode()
```
