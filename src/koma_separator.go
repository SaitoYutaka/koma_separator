package main

import (
	"fmt"		// Ptintnl
	"strings"	// Replace
	"path"		// Dir, Ext
	"syscall"	// O_CREAT
	"image"		// Decode
	"os"		// Exit, Open, Close
	"image/draw"    // Draw
	"image/png"     // Encode
	"io/ioutil"     // ReadFile
	"encoding/json" // Unmarshal
)

func getKomaPos(s string)(map[string]interface{}){
	js, _ := ioutil.ReadFile("js/frame.json")

	var f interface{}
	json.Unmarshal(js, &f)

	m := f.(map[string]interface{})
	koma := m[s].(map[string]interface{})

	return koma
}

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

	komas := getKomaPos("template28")
	for suffix, v := range komas {
		var komapos [4]int
		switch vv := v.(type) {
		case []interface{}:
			for i, u := range vv {
				komapos[i] = int(u.(float64))
			}
		default:
		}
		outfile, _ := os.OpenFile(outFileName + "_" + suffix + path.Ext(fileName), syscall.O_CREAT, 0777)
		koma_title := image.NewRGBA(image.Rect(0, 0, komapos[2], komapos[3]))
		draw.Draw(koma_title, koma_title.Bounds(), src, image.Pt(komapos[0], komapos[1]), draw.Src)
		png.Encode(outfile, koma_title)
		outfile.Close()
	}
}
