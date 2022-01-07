ico
===

This is a simple library for building windows ICO files

Usage
-----

Create a new ico.Icon and add [image.Image](https://pkg.go.dev/image#Image)s to
it by calling Icon.AddImage(), and render the image into a byte array using
Icon.Encode().

TODO
----

 - Read the image depth from the Image object and use that instead of hardcoding
   32 bpp RGBA
 - Read ICO-files
 - BMP as well as PNG contents
