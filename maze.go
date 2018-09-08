package main

import (
	"fmt"
	"math/rand"
	"time"
)

const notABlock = -1
const notAWall = -1

type MazeBlock struct {
	index    int
	previous *MazeBlock
	walls    []bool
}

func (mazeBlock *MazeBlock) first() *MazeBlock {
	b := mazeBlock
	for b.previous != nil {
		b = b.previous
	}
	return b
}

type IMaze interface {
	BlockCount() int
	Block(index int) *MazeBlock
	Navigate(blockIndex int, wallIndex int) (newBlock int, newWall int)
	String() string
	SCad(renderSettings RenderSettings) string
}

func (block *MazeBlock) SetAllWalls() {
	for w := 0; w < len(block.walls); w++ {
		block.walls[w] = true
	}
}

func setAllWalls(maze IMaze) {
	bc := maze.BlockCount()
	for b := 0; b < bc; b++ {
		block := maze.Block(b)
		block.index = b
		block.SetAllWalls()
	}
}

func RandomMaze(maze IMaze) {
	setAllWalls(maze)

	rand.Seed(time.Now().UnixNano())
	bc := maze.BlockCount()

	// Initialize each block to being in their own groups
	for i := 0; i < bc; i++ {
		block := maze.Block(i)
		block.previous = nil
	}

	// Join walls till all cells are joined (BlockCount-1 required)
	for i := 0; i < bc-1; i++ {
		for {
			i1 := rand.Intn(bc)
			b1 := maze.Block(i1)
			w1 := rand.Intn(len(b1.walls))
			i2, w2 := maze.Navigate(i1, w1)

			if i2 == notABlock {
				continue
			}

			b2 := maze.Block(i2)

			f1 := b1.first()
			f2 := b2.first()

			if f1 == f2 {
				continue
			}

			if i1 < i2 {
				b2.previous = f1
				f2.previous = f1
			} else {
				b1.previous = f2
				f1.previous = f2
			}

			// Actually open the wall
			b1.walls[w1] = false
			b2.walls[w2] = false

			break
		}
	}
}

type Maze struct {
	blockCount int
	blocks     []MazeBlock
}

type SquareMaze struct {
	Maze
	width  int
	height int
}

func (maze *Maze) BlockCount() int {
	return maze.blockCount
}

func (maze *Maze) Block(index int) *MazeBlock {
	return &maze.blocks[index]
}

func (maze SquareMaze) Init(width int, height int) *SquareMaze {
	maze.width = width
	maze.height = height
	maze.blockCount = width * height
	maze.blocks = make([]MazeBlock, maze.blockCount)
	for b := 0; b < maze.blockCount; b++ {
		maze.blocks[b].walls = make([]bool, 4)
	}

	setAllWalls(&maze)
	return &maze
}

func (maze *SquareMaze) Navigate(blockIndex int, wallIndex int) (newBlockIndex int, newWall int) {
	if blockIndex >= len(maze.blocks) {
		return notABlock, notAWall
	}

	mazeBlock := maze.blocks[blockIndex]

	if len(mazeBlock.walls) <= wallIndex {
		return notABlock, notAWall
	}

	x := blockIndex % maze.width
	y := blockIndex / maze.width

	switch wallIndex {
	case 0: //Left
		if x <= 0 {
			return notABlock, notAWall
		}
		x--

	case 1: //Up
		if y <= 0 {
			return notABlock, notAWall
		}
		y--

	case 2: //Right
		if x >= maze.width-1 {
			return notABlock, notAWall
		}
		x++

	case 3: //Down
		if y >= maze.height-1 {
			return notABlock, notAWall
		}
		y++
	}

	return x + y*maze.width, (wallIndex + 2) % 4
}

func (maze *SquareMaze) String() string {
	s := ""
	for y := 0; y < maze.height; y++ {
		//Tops
		for x := 0; x < maze.width; x++ {
			block := maze.Block(x + y*maze.width)

			if block.walls[1] {
				s = s + "*---"
			} else {
				s = s + "*   "
			}
		}

		s = s + "*\n"

		// Lefts
		for x := 0; x < maze.width; x++ {
			block := maze.Block(x + y*maze.width)

			if block.walls[0] {
				s = s + "|" + fmt.Sprintf("%2d", x+y*maze.width) + " "
			} else {
				s = s + " " + fmt.Sprintf("%2d", x+y*maze.width) + " "
			}
		}

		// Last right
		s = s + "|\n"
	}

	// Bottom
	for x := 0; x < maze.width; x++ {
		s = s + "*---"
	}
	s = s + "*\n"

	return s
}

func (maze *SquareMaze) SCad(rs RenderSettings) string {
	s := ""

	s += "// Global resolution\n"
	s += "$fs = 0.1;  // Don't generate smaller facets than 0.1 mm\n"
	s += "$fa = 10;    // Don't generate larger angles than 5 degrees\n"
	s += "\n"

	s += fmt.Sprintf("mazeWidth = %d;\n", maze.width)
	s += fmt.Sprintf("mazeHeight = %d;\n", maze.height)
	s += fmt.Sprintf("blockSize = %6f;\n", rs.blockSize)
	s += fmt.Sprintf("blockDepth = %6f;\n", rs.blockDepth)
	s += fmt.Sprintf("ballRadius = %6f;\n", rs.ballRadius)
	s += fmt.Sprintf("ballDepth = %6f;\n", rs.ballDepth)
	s += "\n"

	s += "module maze() color(\"red\") linear_extrude(height = blockDepth) square([mazeWidth * blockSize, mazeHeight*blockSize]);\n"
	s += "\n"
	s += "module ball() color(\"green\") sphere(ballRadius);\n"
	s += "\n"
	s += "module hcylinder() color(\"green\") rotate([0, 90, 0]) cylinder(h = blockSize, r = ballRadius);\n"
	s += "\n"
	s += "module vcylinder() color(\"green\") rotate([-90, 0, 0]) cylinder(h = blockSize, r = ballRadius);\n"
	s += "\n"

	s += "translate([mazeWidth * blockSize / -2, mazeHeight * blockSize / -2, 0]) {\n"
	s += "  difference() {\n"
	s += "    maze();\n"

	s += "    union() {\n"
	for y := 0; y < maze.height; y++ {
		for x := 0; x < maze.width; x++ {
			s += fmt.Sprintf("      translate([(%d + 0.5) * blockSize, (%d + 0.5) * blockSize, ballDepth]) {\n", x, y)
			s += "        ball();\n"

			i1 := x + y*maze.width
			b1 := maze.blocks[i1]

			if !b1.walls[2] {
				s += "        hcylinder();\n"
			}
			if !b1.walls[3] {
				s += "        vcylinder();\n"
			}
			s += "      }\n"
		}
	}
	s += "    }\n"
	s += "  }\n"
	s += "}\n"

	return s
}
