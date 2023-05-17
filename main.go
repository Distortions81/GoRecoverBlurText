package main

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/fogleman/gg"
	"github.com/matsuyoshi30/song2"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var exampleString = "SuperSecretTextHere"
var charSet string = "!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~ "
var numLetters int = len(exampleString) - 1
var bAmount int = 6
var iSizeY int = 32
var iSizeX int = 315

var highScore uint64 = 1<<63 - 1
var tempHighScore uint64
var highScoreChar string
var highScoreString string
var testStr string
var xos int = 0
var yos int = 0

var strImg *gg.Context
var strImgBlur image.Image
var src image.Image
var face font.Face

func main() {
	fdata, err := ioutil.ReadFile("SF-Mono-Font-master/SFMono-Regular.otf")
	if err != nil {
		log.Fatal(err)
	}
	f, err := opentype.Parse(fdata)
	if err != nil {
		log.Fatal(err)
	}

	face, _ = opentype.NewFace(f, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingNone,
	})

	strImg = gg.NewContext(iSizeX, iSizeY)
	strImg.SetFontFace(face)
	makeExampleImage()

	inputImgData, _ := os.Open("input.png")
	defer inputImgData.Close()
	src, _, err = image.Decode(inputImgData)

	if err != nil {
		log.Fatal(err)
	}

	iSizeX = src.Bounds().Size().X
	iSizeY = src.Bounds().Size().Y

	fmt.Println("Scan up.")
	for i := 0; i < numLetters; i++ {
		tempHighScore = 1<<63 - 1
		highScoreChar = ""
		for _, testChar := range charSet {
			testImage(testStr+string(testChar), string(testChar))
		}
		if highScore == 0 {
			break
		}

		if highScoreChar != "" {
			testStr += highScoreChar
		}
	}

	var oldHighScore uint64 = 1<<63 - 1
	for {
		for {
			fmt.Println("Rescan up.")
			for i := 0; i < numLetters; i++ {
				tempHighScore = 1<<63 - 1
				highScoreChar = ""
				for _, testChar := range charSet {
					testStr = testStr[:i] + string(testChar) + testStr[i+1:]
					testImage(testStr, string(testChar))
				}
				if highScoreChar != "" {
					testStr = testStr[:i] + string(highScoreChar) + testStr[i+1:]
				}
			}
			fmt.Println("Done", highScoreString, highScore)
			if highScore == 0 {
				return
			}

			fmt.Println("Rescan down.")
			for i := numLetters - 1; i > 0; i-- {
				tempHighScore = 1<<63 - 1
				highScoreChar = ""
				for _, testChar := range charSet {
					testStr = testStr[:i] + string(testChar) + testStr[i+1:]
					testImage(testStr, string(testChar))
				}
				if highScoreChar != "" {
					testStr = testStr[:i] + string(highScoreChar) + testStr[i+1:]
				}
			}
			fmt.Println("Done", highScoreString, highScore)
			if highScore == 0 {
				return
			}

			if oldHighScore == highScore {
				break
			}
			oldHighScore = highScore
		}

		for {
			fmt.Println("Rescan up 2 char.")
			for i := numLetters - 2; i > 0; i-- {
				tempHighScore = 1<<63 - 1
				highScoreChar = ""
				for _, testCharA := range charSet {
					for _, testCharB := range charSet {
						testStr = testStr[:i] + string(testCharA) + string(testCharB) + testStr[i+2:]
						testImage(testStr, string(testCharA)+string(testCharB))
					}
				}
				if highScoreChar != "" {
					testStr = testStr[:i] + string(highScoreChar) + testStr[i+2:]
				}
			}
			fmt.Println("Done", highScoreString, highScore)
			if highScore == 0 {
				break
			}

			fmt.Println("Rescan down 2 char.")
			for i := 0; i < numLetters-3; i++ {
				tempHighScore = 1<<63 - 1
				highScoreChar = ""
				for _, testCharA := range charSet {
					for _, testCharB := range charSet {
						testStr = testStr[:i] + string(testCharA) + string(testCharB) + testStr[i+2:]
						testImage(testStr, string(testCharA)+string(testCharB))
					}
				}
				if highScoreChar != "" {
					testStr = testStr[:i] + string(highScoreChar) + testStr[i+2:]
				}
			}
			fmt.Println("Done", highScoreString, highScore)
			if highScore == 0 {
				return
			}
			if oldHighScore == highScore {
				break
			}
			oldHighScore = highScore
		}

		for {
			fmt.Println("Rescan up 3 char.")
			for i := numLetters - 3; i > 0; i-- {
				tempHighScore = 1<<63 - 1
				highScoreChar = ""
				for _, testCharA := range charSet {
					for _, testCharB := range charSet {
						for _, testCharC := range charSet {
							testStr = testStr[:i] + string(testCharA) + string(testCharB) + string(testCharC) + testStr[i+3:]
							testImage(testStr, string(testCharA)+string(testCharB)+string(testCharB))
						}
					}
				}
				if highScoreChar != "" {
					testStr = testStr[:i] + string(highScoreChar) + testStr[i+3:]
				}
			}
			fmt.Println("Done", highScoreString, highScore)
			if highScore == 0 {
				return
			}

			fmt.Println("Rescan down 3 char.")
			for i := 0; i < numLetters-4; i++ {
				tempHighScore = 1<<63 - 1
				highScoreChar = ""
				for _, testCharA := range charSet {
					for _, testCharB := range charSet {
						for _, testCharC := range charSet {
							testStr = testStr[:i] + string(testCharA) + string(testCharB) + string(testCharC) + testStr[i+3:]
							testImage(testStr, string(testCharA)+string(testCharB)+string(testCharB))
						}
					}
				}
				if highScoreChar != "" {
					testStr = testStr[:i] + string(highScoreChar) + testStr[i+3:]
				}
			}
			fmt.Println("Done", highScoreString, highScore)
			if highScore == 0 {
				return
			}

			if oldHighScore == highScore {
				break
			}
			oldHighScore = highScore
		}
		for {
			fmt.Println("Rescan up 4 char.")
			for i := numLetters - 4; i > 0; i-- {
				tempHighScore = 1<<63 - 1
				highScoreChar = ""
				for _, testCharA := range charSet {
					for _, testCharB := range charSet {
						for _, testCharC := range charSet {
							testStr = testStr[:i] + string(testCharA) + string(testCharB) + string(testCharC) + testStr[i+3:]
							testImage(testStr, string(testCharA)+string(testCharB)+string(testCharB))
						}
					}
				}
				if highScoreChar != "" {
					testStr = testStr[:i] + string(highScoreChar) + testStr[i+3:]
				}
			}
			fmt.Println("Done", highScoreString, highScore)
			if highScore == 0 {
				return
			}

			fmt.Println("Rescan down 4 char.")
			for i := 0; i < numLetters-5; i++ {
				tempHighScore = 1<<63 - 1
				highScoreChar = ""
				for _, testCharA := range charSet {
					for _, testCharB := range charSet {
						for _, testCharC := range charSet {
							testStr = testStr[:i] + string(testCharA) + string(testCharB) + string(testCharC) + testStr[i+3:]
							testImage(testStr, string(testCharA)+string(testCharB)+string(testCharB))
						}
					}
				}
				if highScoreChar != "" {
					testStr = testStr[:i] + string(highScoreChar) + testStr[i+3:]
				}
			}
			fmt.Println("Done", highScoreString, highScore)
			if highScore == 0 {
				return
			}

			if oldHighScore == highScore {
				break
			}
			oldHighScore = highScore
		}

	}

	fmt.Println("Done", highScoreString, highScore)
}

func intAbs(input int64) uint64 {
	if input < 0 {
		return uint64(-input)
	}
	return uint64(input)
}

func testImage(str string, c string) {
	strImg.SetRGB(1, 1, 1)
	strImg.Clear()
	strImg.SetRGB(0, 0, 0)
	strImg.DrawStringAnchored(str, float64(xos), float64(iSizeY)/2+float64(yos), 0, 0.3)
	//strImgBlur = blur.Gaussian(strImg.Image(), float64(bAmount))
	strImgBlur = song2.GaussianBlur(strImg.Image(), float64(bAmount))

	var tscore uint64 = 0
	for x := 0; x < iSizeX; x++ {
		for y := 0; y < iSizeY; y++ {
			_, ag, _, _ := strImgBlur.At(x, y).RGBA()
			_, bg, _, _ := src.At(x, y).RGBA()

			//rdiff := intAbs(int64(ar) - int64(br))
			gdiff := intAbs(int64(ag) - int64(bg))
			//bdiff := intAbs(int64(ab) - int64(bb))
			//pscore := rdiff + gdiff + bdiff
			tscore += gdiff
		}
	}
	if tscore < tempHighScore {
		tempHighScore = tscore
		highScoreChar = c
	}

	if tscore < highScore {
		highScore = tscore
		fmt.Println("New high score: '", str, "'", tempHighScore)

		//Write out the image
		name := "high-score.png"
		blah, _ := os.Create(name)
		png.Encode(blah, strImgBlur)
		blah.Close()
	}
}

func makeExampleImage() {
	strImg.SetRGB(1, 1, 1)
	strImg.Clear()
	strImg.SetRGB(0, 0, 0)
	strImg.DrawStringAnchored(exampleString, float64(xos), float64(iSizeY)/2+float64(yos), 0, 0.3)
	numLetters = len(exampleString)
	//strImg.Scale(0.1, 0.1)
	//strImgBlur = blur.Gaussian(strImg.Image(), float64(bAmount))
	strImgBlur = song2.GaussianBlur(strImg.Image(), float64(bAmount))

	name := "input.png"
	blah, _ := os.Create(name)
	png.Encode(blah, strImgBlur)
	blah.Close()
}
