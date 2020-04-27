package main

/*
	A package to run the 'Game of Life' simulation
	Rules:
		- Any live cell with two or three neighbors survives.
		- Any dead cell with three live neighbors becomes a live cell.
		- All other live cells die in the next generation. Similarly, all other dead cells stay dead.

	Build: go build -ldflags -H=windowsgui
*/

import (
	"fmt"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"

	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

const (
	generations = -1
	gridSizeX   = 20 // Min 8 - Max 28
	gridSizeY   = 20 // Min 8 - Max 28
)

func main() {

	if gridSizeX < 8 || gridSizeY < 8 || gridSizeX > 28 || gridSizeY > 28 { // If the grid size is not within the valid range

		panic("World grid dimensions are invalid - Minimum grid = 8x8 - Maximum grid = 28x28")

	}

	if generations != -1 && generations < 8 {

		panic("Generations must be set to -1 (Forever) or a value >= 8")

	}

	wrld := world{grid: make([][]*widget.Check, gridSizeY), generating: false, currentGeneration: 0}

	for y := 0; y < gridSizeY; y++ {

		wrld.grid[y] = make([]*widget.Check, gridSizeX)

	} // Create a jagged slice of checkboxes the same size as the 'wrld'

	app := app.New()

	window := app.NewWindow("Game Of Life")

	generationLabel := widget.NewLabel("Generation: 0")

	var startButton *widget.Button

	startButton = widget.NewButton("Start Simulation", func() {

		wrld.generating = true

		startButton.DisableableWidget.Disable()

		go wrld.beginGeneration(generationLabel) // Begin generating evolutions on a goroutine

	})

	rows := make([]*fyne.Container, gridSizeY)

	for y := 0; y < gridSizeY; y++ {

		rows[y] = fyne.NewContainerWithLayout(layout.NewHBoxLayout())

		for x := 0; x < gridSizeX; x++ {

			wrld.grid[y][x] = widget.NewCheck("", func(value bool) {})

			rows[y].AddObject(wrld.grid[y][x]) // Place the current checkbox in the respective row

		}

		switch y {

		case 0:

			rows[y].AddObject(generationLabel)

		case 1:

			rows[y].AddObject(startButton)

			break

		case 2:

			rows[y].AddObject(
				widget.NewButton("Next Generation", func() {

					wrld.nextGeneration()

					wrld.currentGeneration++

					generationLabel.Text = fmt.Sprintf("Generation: %d", wrld.currentGeneration)

					generationLabel.Refresh()

				}),
			)

			break

		case 4:

			rows[y].AddObject(
				widget.NewButton("Stop Simulation", func() {

					startButton.DisableableWidget.Enable()

					wrld.generating = false

				}),
			)

			break

		case 5:

			rows[y].AddObject(
				widget.NewButton("Reset Simulation", func() {

					startButton.DisableableWidget.Enable()

					wrld.generating = false

					wrld.currentGeneration = 0

					generationLabel.Text = "Generation: 0"

					generationLabel.Refresh()

					for y := 0; y < gridSizeY; y++ {

						for x := 0; x < gridSizeX; x++ {

							wrld.grid[y][x].Checked = false

							wrld.grid[y][x].Refresh()

						}

					} // Clear the grid

				}),
			)

			break

		}

	}

	content := widget.NewVBox()

	for _, row := range rows {

		content.Append(row)

	}

	window.SetContent(content)

	window.ShowAndRun()

}

type world struct {
	grid              [][]*widget.Check
	generating        bool
	currentGeneration int
}

func (thisWorld *world) beginGeneration(genLabel *widget.Label) {

	for thisWorld.generating {

		if thisWorld.currentGeneration < generations || generations == -1 {

			thisWorld.currentGeneration++

			genLabel.Text = fmt.Sprintf("Generation: %d", thisWorld.currentGeneration)

			genLabel.Refresh()

			thisWorld.nextGeneration()

			time.Sleep(time.Second)

		} else {

			thisWorld.generating = false

		}

	}

}

func (thisWorld *world) nextGeneration() {

	nextGen := make([][]bool, len(thisWorld.grid))

	for index := 0; index < len(thisWorld.grid); index++ {

		nextGen[index] = make([]bool, len(thisWorld.grid[index]))

	}

	for y := 0; y < len(thisWorld.grid); y++ {

		for x := 0; x < len(thisWorld.grid[0]); x++ {

			neighbors := thisWorld.getNumberOfLivingNeighbors(x, y)

			if neighbors == 3 || (neighbors == 2 && thisWorld.grid[y][x].Checked) {

				nextGen[y][x] = true

			} else {

				nextGen[y][x] = false

			}

		}

	}

	for y := 0; y < len(thisWorld.grid); y++ {

		for x := 0; x < len(thisWorld.grid[0]); x++ {

			thisWorld.grid[y][x].Checked = nextGen[y][x]

			thisWorld.grid[y][x].Refresh()

		}

	}

}

func (thisWorld *world) getNumberOfLivingNeighbors(x int, y int) int {

	var numberOfNeighbors int

	for yOffset := -1; yOffset < 2; yOffset++ {

		for xOffset := -1; xOffset < 2; xOffset++ {

			if !(xOffset == 0 && yOffset == 0) {

				numberOfNeighbors += thisWorld.checkNeighbor(x, y, xOffset, yOffset)

			}

		}

	}

	return numberOfNeighbors

}

func (thisWorld *world) checkNeighbor(x int, y int, xOffset int, yOffset int) int {

	if y+yOffset >= 0 && y+yOffset < len(thisWorld.grid) {

		if x+xOffset >= 0 && x+xOffset < len(thisWorld.grid[0]) {

			if thisWorld.grid[y+yOffset][x+xOffset].Checked {

				return 1 // Neighbor is alive

			}

		}

	}

	return 0 // Neighbor is dead

}
