package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

func main() {
	img := loadImage("input.png")
	pixels := make([][]color.Color, img.Bounds().Dx())
	for i := range pixels {
		pixels[i] = make([]color.Color, img.Bounds().Dy())
	}
	for x := range pixels {
		for y := range pixels[x] {
			pixels[x][y] = img.At(x, y)
		}
	}

	writeImage := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	waterColor := img.At(1638, 981)
	const maxDistance = 200
	for x := range pixels {
		for y := range pixels[x] {
			if pixels[x][y] == waterColor {
				if x%25 == 0 && y%10 == 0 {
					fmt.Println("X: " + fmt.Sprint(x) + ", Y: " + fmt.Sprint(y))
				}
				point, distance := findNearestDifferentColor(img, image.Pt(x, y))
				if distance <= maxDistance {
					pixels[x][y] = img.At(point.X, point.Y)
				}
			}
			writeImage.Set(x, y, pixels[x][y])
		}
	}
	createImage("output.png", writeImage)
}

func findNearestDifferentColor(img image.Image, startPixel image.Point) (image.Point, int) {
	startColor := img.At(startPixel.X, startPixel.Y)
	queue := make([]image.Point, 0)
	queue = append(queue, startPixel)
	visited := make([][]bool, img.Bounds().Dx())
	for i := range visited {
		visited[i] = make([]bool, img.Bounds().Dy())
	}
	visited[startPixel.X][startPixel.Y] = true

	for {
		pop := queue[0]
		queue = queue[1:]
		x, y := pop.X, pop.Y
		if img.At(pop.X, pop.Y) != startColor && img.At(pop.X, pop.Y) != img.At(0, 0) {
			return pop, pixelDistance(startPixel, pop)
		} else {
			if !visited[x+1][y] {
				queue = append(queue, image.Pt(x+1, y))
				visited[x+1][y] = true
			}
			if !visited[x-1][y] {
				queue = append(queue, image.Pt(x-1, y))
				visited[x-1][y] = true
			}
			if !visited[x][y+1] {
				queue = append(queue, image.Pt(x, y+1))
				visited[x][y+1] = true
			}
			if !visited[x][y-1] {
				queue = append(queue, image.Pt(x, y-1))
				visited[x][y-1] = true
			}
		}
		if len(queue) == 0 {
			return image.Pt(0, 0), -1
		}
	}
}

func pixelDistance(pixel1, pixel2 image.Point) int {
	xDistance := math.Pow(float64(pixel1.X-pixel2.X), 2)
	yDistance := math.Pow(float64(pixel1.Y-pixel2.Y), 2)
	return int(math.Sqrt(xDistance + yDistance))
}

func loadImage(file string) image.Image {
	imageFile, err := os.Open(file)
	if err != nil {
		log.Fatal("Error opening file:", file, err)
	}
	defer func(imageFile *os.File) {
		err := imageFile.Close()
		if err != nil {
			log.Fatal("Could not close file, somehow.")
		}
	}(imageFile)

	img, err := png.Decode(imageFile)
	if err != nil {
		log.Fatal("Error decoding image:", file, err)
	}
	return img
}

func createImage(outputName string, img image.Image) {
	f, err := os.Create(outputName)
	if err != nil {
		log.Fatal("Could not create image.")
	}
	encodeErr := png.Encode(f, img)
	if encodeErr != nil {
		log.Fatal("Could not encode image.")
	}
	fmt.Println("Successfully created image.")
}
