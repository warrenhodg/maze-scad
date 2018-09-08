package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var appName = "maze-scad"
var appDescription = "Utility to generate an Open SCAD file of a maze"
var version = "1.0.0"

type RenderSettings struct {
	blockSize  float32
	blockDepth float32
	ballRadius float32
	ballDepth  float32
}

func main() {
	app := kingpin.New(appName, appDescription)
	app.Version(version)
	width := app.Flag("width", "Width (in blocks) of the maze").Short('w').Default("10").Int()
	height := app.Flag("height", "Height (in blocks) of the maze").Short('h').Default("10").Int()
	blockSize := app.Flag("blockSize", "Size of the block of the maze").Short('s').Default("10").Float32()
	blockDepth := app.Flag("blockDepth", "Depth of each block of the maze").Short('d').Default("10").Float32()
	ballRadius := app.Flag("ballRadius", "Radius of the ball that runs through the maze").Short('r').Default("4.9").Float32()
	ballDepth := app.Flag("ballDepth", "Depth at which the ball runs through the maze").Short('D').Default("5.2").Float32()

	err := func() error {
		_, err := app.Parse(os.Args[1:])
		if err != nil {
			return err
		}

		renderSettings := RenderSettings{
			*blockSize,
			*blockDepth,
			*ballRadius,
			*ballDepth,
		}

		maze := SquareMaze{}.Init(*width, *height)
		RandomMaze(maze)
		fmt.Printf("%s", maze.SCad(renderSettings))

		return nil
	}()

	if err != nil {
		fmt.Printf("%s", err.Error())
	}
}
