package utils

import "testing"

func TestPPMImage_WriteAndSave(t *testing.T) {
	w, h := 256, 256
	ppmImage := NewPPMImage(w, h, 255)
	var pix = make([]Pixel, w*h)
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			pix[i*w+j] = Pixel{
				R: i * 256 / w,
				G: j * 256 / w,
				B: 0,
			}
		}
	}
	ppmImage.Full(pix)
	//ppmImage.WriteAndSave("")
	ppmImage.FastWriteAndSave("")
	// P3
	// 2 2
	// 255
	// 255 0 0 0 255 0
	// 0 255 0 255 0 0
}
