package main

import (
	"fmt"
	"github.com/EngoEngine/engo"
)

var (
	scrollSpeed float32 = 700

	worldWidth  int = 800
	worldHeight int = 800
)

func main() {
	opts := engo.RunOptions{
		Title:          "KeyboardScroller Demo",
		Width:          worldWidth,
		Height:         worldHeight,
		StandardInputs: true,
	}

	fmt.Println(opts)
}
