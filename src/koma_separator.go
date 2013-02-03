package main

import (
	"fmt"		// Ptintnl
	"strings"	// Replace
	"strconv"	// Itoa
	"path"		// Dir, Ext
	"syscall"	// O_CREAT
	"image"		// Decode
	"os"		// Exit, Open, Close
	"image/draw"    // Draw
	"image/png"     // Encode
)

func usage() {
	fmt.Println("usage: koma_split [filename(.png)]")
}

func checkSize(src image.Image) bool {
	const srcSizeX = 858
	const srcSizeY = 1200
	pt := src.Bounds().Size()
	if pt.X != srcSizeX || pt.Y != srcSizeY {
		return false
	}

	return true
}

func saveTitle(outFileName string, fileName string, src image.Image) {
	const titleSizeX = 684
	const titleSizeY = 73
	outfile, _ := os.OpenFile(outFileName + "_title" + path.Ext(fileName), syscall.O_CREAT, 0777)
	koma_title := image.NewRGBA(image.Rect(0, 0, titleSizeX, titleSizeY))
	draw.Draw(koma_title, koma_title.Bounds(), src, image.Pt(87,98), draw.Src)
	png.Encode(outfile, koma_title)
	outfile.Close()
}

func saveKoma(outFileName string, fileName string, src image.Image) {
	const XOffset = 87
	const YOffset = 185
	const Offset  = 233
	const komaSizeX = 684 
	const komaSizeY = 218
	for i := 0; i < 4; i++ {
		outfile, _ := os.OpenFile(outFileName + "_" + strconv.Itoa(i) + path.Ext(fileName), syscall.O_CREAT, 0777)
		koma := image.NewRGBA(image.Rect(0, 0, komaSizeX, komaSizeY))
		draw.Draw(koma, koma.Bounds(), src, image.Pt(XOffset, YOffset + Offset * i), draw.Src)
		png.Encode(outfile, koma)
		outfile.Close()
	}
}

func main() {

	if len(os.Args) != 2 {
		usage()
		os.Exit(-1)
	}
	
	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		usage()
		os.Exit(-1)
	}

	// Get image data
	src, _, err := image.Decode(file)
	if checkSize(src) == false {
		file.Close()
		fmt.Println("Image size error!")
		fmt.Println("Image size must be 858x1200.")
		usage()
		os.Exit(-1)
	}
	file.Close()

	// Output Image file
	outFileName := strings.Replace(fileName, path.Ext(fileName), "", -1)
	saveTitle(outFileName, fileName, src)
	saveKoma(outFileName, fileName, src)
}
