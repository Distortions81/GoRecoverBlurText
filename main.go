package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/fogleman/gg"
	"github.com/matsuyoshi30/song2"
	"github.com/nfnt/resize"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var exampleString = "SuperSecretTextHere"
var numLetters int = len(exampleString) - 1

var charSet string = "!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~ "

var blurAmount *int
var iSizeY int = 32
var iSizeX int = 315

var highScore uint64 = 1<<63 - 1
var tempHighScore uint64
var highScoreChar string
var highScoreString string
var testStr string
var xos int = 0
var yos int = 0

var testImg *gg.Context
var sourceImg image.Image
var face font.Face
var doBlur *bool

func main() {
	//Handle flags
	doBlur = flag.Bool("blur", false, "blur instead of pixelate")
	blurAmount = flag.Int("amount", 15, "amount to pixelate or blur")
	flag.Parse()

	//Read font
	fdata, err := os.ReadFile("SF-Mono-Font-master/SFMono-Regular.otf")
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

	//Create new blank image, then create an example source image
	testImg = gg.NewContext(iSizeX, iSizeY)
	testImg.SetFontFace(face)
	makeExampleImage()

	//Save image
	inputImgData, _ := os.Open("input.png")
	defer inputImgData.Close()
	sourceImg, _, err = image.Decode(inputImgData)

	if err != nil {
		log.Fatal(err)
	}

	//Get image size
	iSizeX = sourceImg.Bounds().Size().X
	iSizeY = sourceImg.Bounds().Size().Y

	//Start testing text against it
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

	//Scan different text combos.
	//TODO: Lots of replicated code... clean me
	var oldHighScore uint64 = 1<<63 - 1
	for {
		for {
			fmt.Println("Rescan up.")
			for i := 0; i < numLetters; i++ {
				tempHighScore = 1<<63 - 1
				highScoreChar = ""
				for _, testChar := range charSet {
					if i+1 >= len(testStr) {
						continue
					}
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
}

func intAbs(input int64) uint64 {
	if input < 0 {
		return uint64(-input)
	}
	return uint64(input)
}

func pixelate(input *gg.Context, amount int) image.Image {
	shrink := resize.Resize(uint(input.Width()/amount), uint(input.Height()/amount), input.Image(), resize.NearestNeighbor)
	return resize.Resize(uint(input.Width()), uint(input.Height()), shrink, resize.NearestNeighbor)
}

// Make a image match score
func testImage(str string, c string) {
	testImg.SetRGB(1, 1, 1)
	testImg.Clear()
	testImg.SetRGB(0, 0, 0)
	testImg.DrawStringAnchored(str, float64(xos), float64(iSizeY)/2+float64(yos), 0, 0.3)
	var outImg image.Image
	if *doBlur {
		outImg = song2.GaussianBlur(testImg.Image(), float64(*blurAmount))
	} else {
		outImg = pixelate(testImg, *blurAmount)
	}

	var tscore uint64 = 0
	for x := 0; x < iSizeX; x++ {
		for y := 0; y < iSizeY; y++ {
			_, ag, _, _ := outImg.At(x, y).RGBA()
			_, bg, _, _ := sourceImg.At(x, y).RGBA()

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
		png.Encode(blah, outImg)
		blah.Close()
	}
}

// Make an example input image
func makeExampleImage() {
	testImg.SetRGB(1, 1, 1)
	testImg.Clear()
	testImg.SetRGB(0, 0, 0)
	testImg.DrawStringAnchored(exampleString, float64(xos), float64(iSizeY)/2+float64(yos), 0, 0.3)
	numLetters = len(exampleString)
	var outImg image.Image
	if *doBlur {
		outImg = song2.GaussianBlur(testImg.Image(), float64(*blurAmount))
	} else {
		outImg = pixelate(testImg, *blurAmount)
	}

	name := "input.png"
	blah, _ := os.Create(name)
	png.Encode(blah, outImg)
	blah.Close()
}
