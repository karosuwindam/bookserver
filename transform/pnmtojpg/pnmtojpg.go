package pnmtojpg

import (
	"image/jpeg"
	"os"

	"github.com/lmittmann/ppm"
	"github.com/spakin/netpbm"
)

const (
	QUALITY = 90
)

// netpbmを使用してpbm画像をjpgに変換する関数
func Pbm2jpg(pbmPath string, jpgPath string) error {
	// pbmファイルのパスからファイルを読み込む
	fr, err := os.Open(pbmPath)
	if err != nil {
	}
	defer fr.Close()
	pbm, err := netpbm.Decode(fr, &netpbm.DecodeOptions{
		Target:      netpbm.PPM, // Want to wind up with color
		Exact:       false,      // Can accept grayscale or B&W too
		PBMMaxValue: 42,         // B&W white --> (42, 42, 42)
	})
	if err != nil {
		return err
	}
	// jpgファイルのパスにファイルを書き込む
	fw, err := os.Create(jpgPath)
	if err != nil {
	}
	defer fw.Close()
	if err := jpeg.Encode(fw, pbm, &jpeg.Options{Quality: QUALITY}); err != nil {
		return err
	}
	return nil
}

// ppm画像をjpgに変換する関数
func Ppm2jpg(ppmPath string, jpgPath string) error {
	// ppmファイルのパスからファイルを読み込む
	fr, err := os.Open(ppmPath)
	if err != nil {
	}
	defer fr.Close()
	img, err := ppm.Decode(fr)
	if err != nil {
		return err
	}
	// jpgファイルのパスにファイルを書き込む
	fw, err := os.Create(jpgPath)
	if err != nil {
	}
	defer fw.Close()
	if err := jpeg.Encode(fw, img, &jpeg.Options{Quality: QUALITY}); err != nil {
		return err
	}
	return nil
}
