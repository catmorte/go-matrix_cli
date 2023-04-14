package main

import (
	"flag"
	"math"
	"math/rand"
	"sync"
	"time"

	tm "github.com/buger/goterm"
	"github.com/gookit/color"
)

var symobls = []rune("!\"#$%&'()*+,-./0123456789:;<=>?@abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~")
var (
	speed   = flag.Duration("speed", 100*time.Millisecond, "speed of matrix")
	density = flag.Float64("density", 0.02, "matrix density")

	xl, yl, sourceX, sparkY int
	symbolsMatrix           [][]rune
	brightnessMatrix        [][]float32
	matrixLock              sync.RWMutex
)

func init() {
	flag.Parse()
	xl, yl = getFireSize()
	resetMatrix()
}

func getFireSize() (int, int) {
	return tm.Width(), tm.Height()
}

func resetMatrix() {
	tm.Clear()
	tm.MoveCursor(0, 0)
	symbolsMatrix = make([][]rune, yl)
	brightnessMatrix = make([][]float32, yl)
	for y := 0; y < yl; y++ {
		brightnessMatrix[y] = make([]float32, xl)
		symbolsMatrix[y] = make([]rune, xl)
		for x := 0; x < xl; x++ {
			brightnessMatrix[y][x] = 0
			symbolsMatrix[y][x] = symobls[rand.Intn(len(symobls))]
		}
	}
}

func fg(r float32, char rune) string {
	if r == 1 {
		return color.RGB(255, 255, 255, false).Sprint(string(char))
	}
	return color.RGB(0, uint8(math.Floor(float64(r))), 0, false).Sprint(string(char))
}

func main() {
	go func() {
		for {
			newXl, newYl := getFireSize()
			if xl != newXl || yl != newYl {
				matrixLock.Lock()
				xl, yl = newXl, newYl
				resetMatrix()
				matrixLock.Unlock()
			}
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			matrixLock.Lock()
			for i := 0; i < 50; i++ {
				randX := rand.Intn(xl)
				randY := rand.Intn(yl)
				symbolsMatrix[randY][randX] = symobls[rand.Intn(len(symobls))]
			}
			matrixLock.Unlock()
			time.Sleep(1 * time.Millisecond)
		}
	}()

	for {
		matrixLock.RLock()
		topMatrix()
		for y := yl - 1; y >= 0; y-- {
			for x := xl - 1; x >= 0; x-- {
				val := brightnessMatrix[y][x]
				r := (255.0) * float32(val)
				tm.MoveCursor(x, y)
				tm.Print(fg(r, symbolsMatrix[y][x]))
			}
		}
		tm.Flush()
		matrixLock.RUnlock()
		time.Sleep(*speed)
	}
}

func topMatrix() {
	for x := 0; x < xl; x++ {
		if rand.Float64() < *density {
			brightnessMatrix[0][x] = 1
		}
	}

	for y := yl - 1; y > 0; y-- {
		for x := 0; x < xl; x++ {
			val := brightnessMatrix[y-1][x]
			brightnessMatrix[y-1][x] = val * 0.9
			brightnessMatrix[y][x] = val
		}
	}
}
