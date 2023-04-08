package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"io"
	"os"
	"strconv"
	"strings"
)

type Pixel struct {
	y, x int
}

// graph node type
const (
	SPACE = iota
	BLOCK
	SOURCE
	DESTINATION
	PATH
)

func main() {
	if len(os.Args) != 9 {
		fmt.Println("Usage: [input file] [output file] [duration] [space color] [block color] [source color] [destination color] [path color]")
		fmt.Println("duration: gif animation in seconds")
		fmt.Println("color: R,G,B from 0 - 255, separated by a comma")
		return
	}

	if !strings.HasSuffix(os.Args[2], ".gif") && !strings.HasSuffix(os.Args[2], ".png") {
		fmt.Println("Unsupported output format, must be .gif or .png")
		return
	}

	inputBytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	inputImage, format, err := image.Decode(bytes.NewReader(inputBytes))
	if err != nil {
		fmt.Println(err)
		return
	}
	rgbaImage := inputImage.(*image.NRGBA)
	fmt.Printf("Input file: %v, format: %v\n", os.Args[1], format)

	outputFile, err := os.Create(os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outputFile.Close()

	duration, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("GIF frame rate: %v/s\n", duration)

	configuredColors := []color.NRGBA{}
	for i := 4; i < 9; i++ {
		color, err := RGBStringToColor(os.Args[i])
		if err != nil {
			fmt.Println(err)
			return
		}
		configuredColors = append(configuredColors, color)
	}
	fmt.Printf("Blank       Color: %+v\nBlock       Color: %+v\nSource      Color: %+v\nDestination Color: %+v\n", configuredColors[0], configuredColors[1], configuredColors[2], configuredColors[3])

	// parse images into graph
	fmt.Println("parsing image into graph...")
	graph, source, destination := ParseToGraph(rgbaImage, configuredColors[0], configuredColors[1], configuredColors[2], configuredColors[3])
	fmt.Printf("Source     : %+v\nDestination: %+v\n", source, destination)

	// calculate the shortest path with breadth first search
	fmt.Println("calculating shortest path...")
	path, found := GetShortestPath(graph, source, destination)
	if !found {
		fmt.Println("Can not find the solution to the maze. Try changing the color values.")
		return
	}

	fmt.Println("generating solution image...")
	if strings.HasSuffix(os.Args[2], ".gif") {
		GenerateGIF(outputFile, rgbaImage, path, configuredColors[4], duration)
	} else {
		GeneratePNG(outputFile, rgbaImage, path, configuredColors[4])
	}
	fmt.Printf("%v has been created\n", os.Args[2])
}

func RGBStringToColor(rgbString string) (color.NRGBA, error) {
	rgb := strings.Split(rgbString, ",")
	if len(rgb) != 3 {
		return color.NRGBA{}, errors.New("invalid RGB format")
	}

	r, err := strconv.Atoi(rgb[0])
	if err != nil {
		return color.NRGBA{}, err
	}

	g, err := strconv.Atoi(rgb[1])
	if err != nil {
		return color.NRGBA{}, err
	}

	b, err := strconv.Atoi(rgb[2])
	if err != nil {
		return color.NRGBA{}, err
	}

	if r < 0 || g < 0 || b < 0 || r > 255 || g > 255 || b > 255 {
		return color.NRGBA{}, errors.New("invalid RGB values")
	}

	return color.NRGBA{uint8(r), uint8(g), uint8(b), 255}, nil
}

func ParseToGraph(rgbaImage *image.NRGBA, spaceColor, blockColor, sourceColor, destinationColor color.NRGBA) ([][]byte, Pixel, Pixel) {
	ColorDistance := func(c1, c2 color.NRGBA) int32 {
		r := int32(c1.R) - int32(c2.R)
		g := int32(c1.G) - int32(c2.G)
		b := int32(c1.B) - int32(c2.B)
		return r*r + g*g + b*b
	}

	GetPixelType := func(pixelColor color.NRGBA) int {
		shortest := ColorDistance(pixelColor, spaceColor)
		pixelType := SPACE

		if distance := ColorDistance(pixelColor, blockColor); distance < shortest {
			shortest = distance
			pixelType = BLOCK
		}

		if distance := ColorDistance(pixelColor, sourceColor); distance < shortest {
			shortest = distance
			pixelType = SOURCE
		}

		if distance := ColorDistance(pixelColor, destinationColor); distance < shortest {
			// shortest = distance
			pixelType = DESTINATION
		}

		return pixelType
	}

	graph := make([][]byte, rgbaImage.Bounds().Dy())
	for y := range graph {
		graph[y] = make([]byte, rgbaImage.Bounds().Dx())
	}

	source, destination := Pixel{}, Pixel{}

	for y := 0; y < rgbaImage.Bounds().Dy(); y++ {
		for x := 0; x < rgbaImage.Bounds().Dx(); x++ {
			pixelColor := rgbaImage.NRGBAAt(x, y)
			switch GetPixelType(pixelColor) {
			case BLOCK:
				graph[y][x] = BLOCK
			case SOURCE:
				source = Pixel{y, x}
			case DESTINATION:
				destination = Pixel{y, x}
			}
		}
	}

	return graph, source, destination
}

func GetShortestPath(graph [][]byte, source, destination Pixel) ([]Pixel, bool) {
	offsets := []Pixel{{0, 1}, {1, 0}, {0, -1}, {-1, 0}, {-1, -1}, {1, 1}, {-1, 1}, {1, -1}}

	parent := make([][]Pixel, len(graph))
	for y := range parent {
		parent[y] = make([]Pixel, len(graph[0]))
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
				if graph[y][x] == SPACE {
					graph[y][x] = BLOCK
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

func GenerateGIF(outputFile io.Writer, nrgbaImage *image.NRGBA, path []Pixel, pathColor color.NRGBA, duration float64) {
	GetBasePalletedImage := func(sourceImage *image.NRGBA) *image.Paletted {
		paletteMap := make(map[color.NRGBA]struct{})
		baseImage := image.NewPaletted(nrgbaImage.Rect, []color.Color{pathColor})

		for y := 0; y < sourceImage.Rect.Dy(); y++ {
			for x := 0; x < sourceImage.Rect.Dx(); x++ {
				nrgba := sourceImage.NRGBAAt(x, y)
				if _, ok := paletteMap[nrgba]; !ok {
					paletteMap[nrgba] = struct{}{}
					baseImage.Palette = append(baseImage.Palette, nrgba)
				}
				baseImage.Set(x, y, nrgba)
			}
		}

		return baseImage
	}

	ClonePalettedImage := func(sourceImage *image.Paletted) *image.Paletted {
		clonedImage := image.Paletted{Rect: sourceImage.Rect, Stride: sourceImage.Stride, Palette: sourceImage.Palette}
		clonedImage.Pix = make([]uint8, len(sourceImage.Pix))
		copy(clonedImage.Pix, sourceImage.Pix)
		return &clonedImage
	}

	gifImage := gif.GIF{}
	baseImage := GetBasePalletedImage(nrgbaImage)

	/**
	x frames / second
	y frames / 10 ms
	*/
	const framesPerSecond = 50
	const delayPerFrame = 1000 / framesPerSecond / 10
	stepsPerFrame := float64(len(path)) / (framesPerSecond * duration)
	stepsRemainThisFrame := stepsPerFrame

	for i := 0; i < len(path); i++ {
		if stepsRemainThisFrame < 1.0 || i == len(path)-1 {
			stepsRemainThisFrame += stepsPerFrame
			momentImage := ClonePalettedImage(baseImage)
			gifImage.Image = append(gifImage.Image, momentImage)
			gifImage.Delay = append(gifImage.Delay, delayPerFrame)
		}

		baseImage.Set(path[i].x, path[i].y, pathColor)
		stepsRemainThisFrame -= 1.0
	}

	gif.EncodeAll(outputFile, &gifImage)
}

func GeneratePNG(outputFile io.Writer, nrgbaImage *image.NRGBA, path []Pixel, pathColor color.NRGBA) {
	for i := range path {
		nrgbaImage.SetNRGBA(path[i].x, path[i].y, pathColor)
	}
	png.Encode(outputFile, nrgbaImage)
}
