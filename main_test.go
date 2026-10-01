package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"image/png"
	"os"
	"reflect"
	"testing"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// PNG hashes captured from the original renderer with x/image v0.46.0.
func TestRenderingMatchesBaseline(t *testing.T) {
	data, err := os.ReadFile("SF-Mono-Font-master/SFMono-Regular.otf")
	if err != nil {
		t.Fatal(err)
	}
	f, err := opentype.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	face, err = opentype.NewFace(f, &opentype.FaceOptions{Size: 24, DPI: 72, Hinting: font.HintingNone})
	if err != nil {
		t.Fatal(err)
	}
	t.Chdir(t.TempDir())
	amount := 3
	blurAmount = &amount
	for _, tc := range []struct {
		blur bool
		hash string
	}{
		{false, "3aa220763d37536d96cde1a0af6e95e56b98b5bacaa45659328170dfd0cb0812"},
		{true, "227eb8eace1220e6397e2cf14e215fcdd83d07cf1d84df57d9224fe08b3f2174"},
	} {
		blur := tc.blur
		doBlur = &blur
		testImg = gg.NewContext(iSizeX, iSizeY)
		testImg.SetFontFace(face)
		makeExampleImage()
		data, err := os.ReadFile("input.png")
		if err != nil {
			t.Fatal(err)
		}
		if got := fmt.Sprintf("%x", sha256.Sum256(data)); got != tc.hash {
			t.Fatalf("blur=%t: rendering hash %s; want %s", blur, got, tc.hash)
		}
		sourceImg, err = png.Decode(bytes.NewReader(data))
		if err != nil {
			t.Fatal(err)
		}
		highScore = ^uint64(0)
		if score := testImage(exampleString); score != 0 {
			t.Fatalf("matching example score = %d", score)
		}
		if highScoreString != exampleString {
			t.Fatalf("best string = %q", highScoreString)
		}
		best, err := os.ReadFile("high-score.png")
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(best, data) {
			t.Fatal("candidate and example rendering differ")
		}
	}
}

func TestCandidateWidths(t *testing.T) {
	for width := 1; width <= 4; width++ {
		seen := make(map[string]bool)
		forEachCandidate("ab", width, func(candidate string) bool {
			if len(candidate) != width {
				t.Fatalf("width %d generated %q", width, candidate)
			}
			if seen[candidate] {
				t.Fatalf("duplicate %q", candidate)
			}
			seen[candidate] = true
			return false
		})
		if len(seen) != 1<<width {
			t.Fatalf("width %d: %d candidates", width, len(seen))
		}
	}
	var visited []string
	forEachCandidate("ab", 3, func(candidate string) bool {
		visited = append(visited, candidate)
		return len(visited) == 3
	})
	if !reflect.DeepEqual(visited, []string{"aaa", "aab", "aba"}) {
		t.Fatalf("early stop: %v", visited)
	}
}

func TestScanText(t *testing.T) {
	const target = "abba"
	score := func(text string) uint64 {
		if len(text) != len(target) {
			t.Fatalf("candidate length changed: %q", text)
		}
		var distance uint64
		for i := range target {
			if text[i] != target[i] {
				distance++
			}
		}
		return distance
	}
	for width := 1; width <= 4; width++ {
		for _, reverse := range []bool{false, true} {
			if got := scanText("xxxx", width, reverse, "ab", score); got != target {
				t.Fatalf("width=%d reverse=%t: got %q", width, reverse, got)
			}
		}
	}
	if got := scanText("xxxx", 2, false, "ab", func(string) uint64 { return 1 }); got != "xxxx" {
		t.Fatalf("equal scores changed guess to %q", got)
	}
	for _, width := range []int{0, -1, 5} {
		if got := scanText("xxxx", width, false, "ab", func(string) uint64 { t.Fatal("invalid width evaluated"); return 0 }); got != "xxxx" {
			t.Fatalf("invalid width changed guess to %q", got)
		}
	}
}
