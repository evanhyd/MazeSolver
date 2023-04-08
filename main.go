package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Pixel struct {
	y, x int
}

// graph node type
const (
	SPACE byte = iota
	BLOCK
	SOURCE
	DESTINATION
	PATH
)

func main() {
	if len(os.Args) != 9 {
		fmt.Printf(
			"Usage: [input file] [output file] [duration] [space color] [block color] [source color] [destination color] [path color]\n" +
				"duration: gif animation in seconds\n" +
				"color: R,G,B from 0 - 255, separated by a comma\n")
		return
	}

	inputImage, err := GetInputImage(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	duration, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("GIF duration: %v s\n", duration)

	configuredColors := []color.NRGBA{}
	for i := 4; i < 9; i++ {
		color, err := RGBStringToColor(os.Args[i])
		if err != nil {
			fmt.Println(err)
			return
		}
		configuredColors = append(configuredColors, color)
	}
	fmt.Printf("Blank       Color: %+v\nBlock       Color: %+v\nSource      Color: %+v\nDestination Color: %+v\n",
		configuredColors[SPACE], configuredColors[BLOCK], configuredColors[SOURCE], configuredColors[DESTINATION])

	log.Println("parsing image into graph...")
	graph, source, destination, err := ParseToGraph(inputImage, configuredColors)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("Source      %v\nDestination %v\n", source, destination)

	log.Println("calculating shortest path...")
	path, found := GetShortestPath(graph, source, destination)
	if !found {
		fmt.Println("Can not find the solution to the maze\nTry changing the color\nDid you mark the source and destination point?")
		return
	}

	log.Printf("generating %v...\n", os.Args[2])
	outputFile, err := os.Create(os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}
	if strings.HasSuffix(os.Args[2], ".gif") {
		GenerateAnimatedImage(outputFile, graph, path, configuredColors, duration)
	} else {
		GenerateStaticImage(outputFile, graph, path, configuredColors)
	}
	log.Println("Completed")
}

func GetInputImage(inputPath string) (*image.NRGBA, error) {
	sourceBytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		return nil, err
	}
	sourceImage, format, err := image.Decode(bytes.NewReader(sourceBytes))
	if err != nil {
		return nil, err
	}
	fmt.Printf("Input file: %v, format: %v\n", os.Args[1], format)

	inputImage := image.NewNRGBA(sourceImage.Bounds())
	draw.Draw(inputImage, inputImage.Bounds(), sourceImage, sourceImage.Bounds().Min, draw.Src)
	return inputImage, nil
}

func RGBStringToColor(rgbString string) (color.NRGBA, error) {
	rgb := strings.Split(rgbString, ",")
	if len(rgb) != 3 {
		return color.NRGBA{}, errors.New("invalid RGB format")
	}

	rgbColors := [3]uint8{}
	for i := range rgbColors {
		rgbColor, err := strconv.Atoi(rgb[i])
		if err != nil {
			return color.NRGBA{}, err
		}
		if rgbColor < 0 || rgbColor > 255 {
			return color.NRGBA{}, errors.New("invalid RGB values")
		}
		rgbColors[i] = uint8(rgbColor)
	}

	return color.NRGBA{rgbColors[0], rgbColors[1], rgbColors[2], 255}, nil
}

func ParseToGraph(inputImage *image.NRGBA, configuredColors []color.NRGBA) ([][]byte, Pixel, Pixel, error) {
	ColorDistance := func(c1, c2 color.NRGBA) int32 {
		r := int32(c1.R) - int32(c2.R)
		g := int32(c1.G) - int32(c2.G)
		b := int32(c1.B) - int32(c2.B)
		return r*r + g*g + b*b
	}

	graph := make([][]byte, inputImage.Bounds().Dy())
	for y := range graph {
		graph[y] = make([]byte, inputImage.Bounds().Dx())
	}

	source, destination := Pixel{}, Pixel{}
	shortestSourceDistance, shortestDestinationDistance := int32(math.MaxInt32), int32(math.MaxInt32)

	for y := 0; y < inputImage.Bounds().Dy(); y++ {
		for x := 0; x < inputImage.Bounds().Dx(); x++ {
			pixelColor := inputImage.NRGBAAt(x, y)
			spaceDistance := ColorDistance(pixelColor, configuredColors[SPACE])
			blockDistance := ColorDistance(pixelColor, configuredColors[BLOCK])
			sourceDistance := ColorDistance(pixelColor, configuredColors[SOURCE])
			destinationDistance := ColorDistance(pixelColor, configuredColors[DESTINATION])

			if sourceDistance < shortestSourceDistance {
				shortestSourceDistance = sourceDistance
				source = Pixel{y, x}
			}

			if destinationDistance < shortestDestinationDistance {
				shortestDestinationDistance = destinationDistance
				destination = Pixel{y, x}
			}

			if blockDistance < spaceDistance && blockDistance < sourceDistance && blockDistance < destinationDistance {
				graph[y][x] = BLOCK
			}
		}
	}

	if shortestSourceDistance == math.MaxInt32 {
		return nil, Pixel{}, Pixel{}, errors.New("can not identify the source point")
	}
	if shortestDestinationDistance == math.MaxInt32 {
		return nil, Pixel{}, Pixel{}, errors.New("can not identify the destination point")
	}

	graph[source.y][source.x] = SOURCE
	graph[destination.y][destination.x] = DESTINATION

	return graph, source, destination, nil
}

func GetShortestPath(graph [][]byte, source, destination Pixel) ([]Pixel, bool) {
	offsets := [8]Pixel{{0, 1}, {1, 0}, {0, -1}, {-1, 0}, {-1, -1}, {-1, 1}, {1, -1}, {1, 1}}

	parent := make([][]Pixel, len(graph))
	for y := range parent {
		parent[y] = make([]Pixel, len(graph[0]))
		for x := range parent[y] {
			parent[y][x].y = -1
		}
	}

	que := []Pixel{source}
	parent[source.y][source.x] = source
	found := false

BFS:
	for len(que) > 0 {
		curr := que[0]
		que = que[1:]

		for _, offset := range offsets {
			y := curr.y + offset.y
			x := curr.x + offset.x

			if 0 <= y && y < len(graph) && 0 <= x && x < len(graph[0]) {
				if graph[y][x] != BLOCK && parent[y][x].y == -1 {
					parent[y][x] = curr
					if destination.y == y && destination.x == x {
						found = true
						break BFS
					}
					que = append(que, Pixel{y, x})
				}
			}
		}
	}

	if !found {
		return nil, found
	}

	path := []Pixel{}
	trace := destination
	for trace != source {
		path = append(path, trace)
		trace = parent[trace.y][trace.x]
	}

	for i, j := 0, len(path)-1; i < j; {
		path[i], path[j] = path[j], path[i]
		i++
		j--
	}

	return path, found
}

func GenerateAnimatedImage(outputFile io.Writer, graph [][]byte, path []Pixel, configuredColors []color.NRGBA, duration float64) {
	ClonePalettedImage := func(imageToCopy *image.Paletted) *image.Paletted {
		cloned := image.Paletted{Rect: imageToCopy.Rect, Stride: imageToCopy.Stride, Palette: imageToCopy.Palette}
		cloned.Pix = make([]uint8, len(imageToCopy.Pix))
		copy(cloned.Pix, imageToCopy.Pix)
		return &cloned
	}

	// convert NRGBA image into base gif image
	baseImage := image.NewPaletted(image.Rect(0, 0, len(graph[0]), len(graph)), color.Palette{})
	for i := range configuredColors {
		baseImage.Palette = append(baseImage.Palette, configuredColors[i])
	}
	for y := 0; y < baseImage.Rect.Dy(); y++ {
		for x := 0; x < baseImage.Rect.Dx(); x++ {
			baseImage.SetColorIndex(x, y, graph[y][x])
		}
	}

	const framesPerSecond = 25
	const delayPerFrame = 1000 / framesPerSecond / 10
	stepsPerFrame := float64(len(path)) / (framesPerSecond * duration)
	stepsRemainThisFrame := stepsPerFrame

	gifImage := gif.GIF{}
	for i := 0; i < len(path); i++ {
		if stepsRemainThisFrame < 1.0 || i == len(path)-1 {
			stepsRemainThisFrame += stepsPerFrame
			momentImage := ClonePalettedImage(baseImage)
			gifImage.Image = append(gifImage.Image, momentImage)
			gifImage.Delay = append(gifImage.Delay, delayPerFrame)
		}

		baseImage.SetColorIndex(path[i].x, path[i].y, PATH)
		stepsRemainThisFrame -= 1.0
	}

	gif.EncodeAll(outputFile, &gifImage)
}

func GenerateStaticImage(outputFile io.Writer, graph [][]byte, path []Pixel, configuredColors []color.NRGBA) {
	outputImage := image.NewNRGBA(image.Rect(0, 0, len(graph[0]), len(graph)))
	for i := range graph {
		for j := range graph[i] {
			outputImage.SetNRGBA(j, i, configuredColors[graph[i][j]])
		}
	}
	for i := range path {
		outputImage.SetNRGBA(path[i].x, path[i].y, configuredColors[PATH])
	}
	png.Encode(outputFile, outputImage)
}
