package main

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/petherin/engotest/pkg"
)

var (
	scrollSpeed float32 = 700

	worldWidth  int = 800
	worldHeight int = 800
)

type DefaultScene struct{}

func (d DefaultScene) Preload() {}

func (d DefaultScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	common.SetBackground(color.White)
	w.AddSystem(&common.RenderSystem{})

	w.AddSystem(common.NewKeyboardScroller(scrollSpeed, engo.DefaultHorizontalAxis, engo.DefaultVerticalAxis))

	// Create the background; this way we'll see when we actually scroll
	pkg.NewBackground(w, worldWidth, worldHeight, color.RGBA{192, 153, 0, 255}, color.RGBA{102, 173, 0, 255})

	// Center camera if GlobalScale is Setup
	engo.Mailbox.Dispatch(
		common.CameraMessage{
			Axis:        common.XAxis,
			Value:       float32(worldWidth) / 2,
			Incremental: false},
	)

	engo.Mailbox.Dispatch(
		common.CameraMessage{
			Axis:        common.YAxis,
			Value:       float32(worldHeight) / 2,
			Incremental: false},
	)
}

func (d DefaultScene) Type() string {
	return "Game"
}

func main() {
	opts := engo.RunOptions{
		Title:          "KeyboardScroller Demo",
		Width:          worldWidth,
		Height:         worldHeight,
		StandardInputs: true,
	}

	engo.Run(opts, &DefaultScene{})
}
