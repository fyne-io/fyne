// ◄◄◄ gobmp/rle.go ►►►
// Copyright © 2012 Jason Summers
// Use of this code is governed by an MIT-style license that can
// be found in the readme.md file.
//
// BMP RLE decoder
//

package gobmp

import "io"
import "bufio"

type rleState struct {
	xpos, ypos   int // Position in the target image
	badColorFlag bool
}

func (d *decoder) rlePutPixel(rle *rleState, v byte) {
	// Make sure the position is valid.
	if rle.xpos < 0 || rle.xpos >= d.width ||
		rle.ypos < 0 || rle.ypos >= d.height {
		return
	}
	// Make sure the palette index is valid.
	if int(v) >= d.dstPalNumEntries {
		rle.badColorFlag = true
		return
	}

	// Set the pixel, and advance the current position.
	var dstRow int

	if d.isTopDown {
		// Top-down RLE-compressed images are not legal in any known BMP
		// specification, but we'll tolerate them.
		dstRow = rle.ypos
	} else {
		dstRow = d.height - rle.ypos - 1
	}

	d.img_Paletted.Pix[dstRow*d.img_Paletted.Stride+rle.xpos] = v
	rle.xpos++
}

func (d *decoder) readBitsRLE() error {
	var err error
	var b1, b2 byte
	var uncPixelsLeft int
	var deltaFlag bool
	var k int

	bufferedR := bufio.NewReader(d.r)
	rle := new(rleState)
	rle.xpos = 0
	rle.ypos = 0

	for {
		if rle.badColorFlag {
			return FormatError("palette index out of range")
		}

		if rle.ypos >= d.height || (rle.ypos == (d.height-1) && rle.xpos >= d.width) {
			break // Reached the end of the target image; may as well stop
		}

		// Read the next two bytes
		b1, err = bufferedR.ReadByte()
		if err == nil {
			b2, err = bufferedR.ReadByte()
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if uncPixelsLeft > 0 {
			if d.biCompression == bI_RLE4 {
				// The two bytes we're processing store up to 4 uncompressed pixels.
				d.rlePutPixel(rle, b1>>4)
				uncPixelsLeft--
				if uncPixelsLeft > 0 {
					d.rlePutPixel(rle, b1&0x0f)
					uncPixelsLeft--
				}
				if uncPixelsLeft > 0 {
					d.rlePutPixel(rle, b2>>4)
					uncPixelsLeft--
				}
				if uncPixelsLeft > 0 {
					d.rlePutPixel(rle, b2&0x0f)
					uncPixelsLeft--
				}
			} else { // RLE8
				// The two bytes we're processing store up to 2 uncompressed pixels.
				d.rlePutPixel(rle, b1)
				uncPixelsLeft--
				if uncPixelsLeft > 0 {
					d.rlePutPixel(rle, b2)
					uncPixelsLeft--
				}
			}
		} else if deltaFlag {
			rle.xpos += int(b1)
			rle.ypos += int(b2)
			deltaFlag = false
		} else if b1 == 0 {
			// An uncompressed run, or a special code.
			//
			// Any pixels skipped by special codes will be left at whatever
			// image.NewPaletted() initialized them to, which we assume is 0,
			// meaning palette entry 0.
			if b2 == 0 { // End of row
				rle.ypos++
				rle.xpos = 0
			} else if b2 == 1 { // End of bitmap
				break
			} else if b2 == 2 { // Delta
				deltaFlag = true
			} else {
				// An upcoming uncompressed run of b2 pixels
				uncPixelsLeft = int(b2)
			}
		} else { // A compressed run of pixels
			if d.biCompression == bI_RLE4 {
				// b1 pixels, alternating between two colors
				for k = 0; k < int(b1); k++ {
					if k%2 == 0 {
						d.rlePutPixel(rle, b2>>4)
					} else {
						d.rlePutPixel(rle, b2&0x0f)
					}
				}
			} else { // RLE8
				// b1 pixels of color b2
				for k = 0; k < int(b1); k++ {
					d.rlePutPixel(rle, b2)
				}
			}
		}
	}

	return nil
}
