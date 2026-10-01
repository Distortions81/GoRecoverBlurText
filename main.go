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

	face, err = opentype.NewFace(f, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingNone,
	})

	if err != nil {
		log.Fatal(err)
	}

	//Create new blank image, then create an example source image
	testImg = gg.NewContext(iSizeX, iSizeY)
	testImg.SetFontFace(face)
	makeExampleImage()

	//Save image
	inputImgData, err := os.Open("input.png")
	if err != nil {
		log.Fatal(err)
	}
	defer inputImgData.Close()
	sourceImg, _, err = image.Decode(inputImgData)

	if err != nil {
		log.Fatal(err)
	}

	//Get image size
	iSizeX = sourceImg.Bounds().Size().X
	iSizeY = sourceImg.Bounds().Size().Y

	// Build the initial guess one character at a time.
	fmt.Println("Scan up.")
	for i := 0; i < numLetters; i++ {
		bestScore := ^uint64(0)
		best := testStr
		forEachCandidate(charSet, 1, func(candidate string) bool {
			guess := testStr + candidate
			if score := testImage(guess); score < bestScore {
				best, bestScore = guess, score
			}
			return bestScore == 0
		})
		testStr = best
		if highScore == 0 {
			return
		}
	}

	// Try wider replacements when shorter passes stop improving the score.
	for {
		before := highScore
		for width := 1; width <= 4 && width <= len(testStr); width++ {
			for {
				previous := highScore
				for _, reverse := range []bool{false, true} {
					fmt.Printf("Rescan %d char, reverse=%t.\n", width, reverse)
					testStr = scanText(testStr, width, reverse, charSet, testImage)
					fmt.Println("Done", highScoreString, highScore)
					if highScore == 0 {
						return
					}
				}
				if previous == highScore {
					break
				}
			}
		}
		if before == highScore {
			return
		}
	}
}

// forEachCandidate visits every replacement of the requested width in alphabet
// order. Returning true stops enumeration as soon as an exact match is found.
func forEachCandidate(alphabet string, width int, visit func(string) bool) {
	if width <= 0 || alphabet == "" {
		return
	}
	letters := []rune(alphabet)
	buffer := make([]rune, width)
	var generate func(int) bool
	generate = func(index int) bool {
		if index == width {
			return visit(string(buffer))
		}
		for _, letter := range letters {
			buffer[index] = letter
			if generate(index + 1) {
				return true
			}
		}
		return false
	}
	generate(0)
}

// scanText replaces each window of ASCII text with its best-scoring candidate.
// Include both ends of the text and keep the current guess when scores tie.
func scanText(text string, width int, reverse bool, alphabet string, score func(string) uint64) string {
	if width <= 0 || width > len(text) {
		return text
	}
	bestScore := score(text)
	last := len(text) - width
	for step := 0; step <= last && bestScore != 0; step++ {
		index := step
		if reverse {
			index = last - step
		}
		base, best := text, text
		forEachCandidate(alphabet, width, func(candidate string) bool {
			guess := base[:index] + candidate + base[index+width:]
			if candidateScore := score(guess); candidateScore < bestScore {
				best, bestScore = guess, candidateScore
			}
			return bestScore == 0
		})
		text = best
	}
	return text
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
func testImage(str string) uint64 {
	outImg := renderText(str)

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
	if tscore < highScore {
		highScore = tscore
		highScoreString = str
		fmt.Println("New high score: '", str, "'", tscore)
		if err := savePNG("high-score.png", outImg); err != nil {
			log.Fatal(err)
		}
	}
	return tscore
}

func renderText(str string) image.Image {
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

	return outImg
}

func savePNG(name string, img image.Image) error {
	output, err := os.Create(name)
	if err != nil {
		return err
	}
	if err := png.Encode(output, img); err != nil {
		output.Close()
		return err
	}
	return output.Close()
}

// Make an example input image
func makeExampleImage() {
	numLetters = len(exampleString)
	if err := savePNG("input.png", renderText(exampleString)); err != nil {
		log.Fatal(err)
	}
}
