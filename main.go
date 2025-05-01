package main

import (
	"image"
	"image/draw"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

func main() {
	topOverlay, err := imaging.Open("static/top.png")
	if err != nil {
		log.Fatal(err)
	}

	inputDir := "input"
	outputDir := "output"
	os.MkdirAll(outputDir, os.ModePerm)

	files, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		if !(strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".png")) {
			continue
		}

		srcPath := filepath.Join(inputDir, name)
		dstPath := filepath.Join(outputDir, name)

		img, err := imaging.Open(srcPath)
		if err != nil {
			log.Println("skipping:", name, "-", err)
			continue
		}

		w := img.Bounds().Dx()
		h := img.Bounds().Dy()

		resizedOverlay := imaging.Resize(topOverlay, w, 0, imaging.Lanczos)

		// overlayW := resizedOverlay.Bounds().Dx()
		overlayH := resizedOverlay.Bounds().Dy()
		offset := image.Pt(0, (h-overlayH)/2)

		out := image.NewNRGBA(img.Bounds())
		draw.Draw(out, out.Bounds(), img, image.Point{}, draw.Src)
		draw.Draw(out, resizedOverlay.Bounds().Add(offset), resizedOverlay, image.Point{}, draw.Over)

		err = imaging.Save(out, dstPath)
		if err != nil {
			log.Println("failed to save:", name, "-", err)
		}
	}
}
