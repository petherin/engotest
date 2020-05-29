package main

import (
	"image/color"
	"strings"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

var (
	scrollSpeed float32 = 700

	worldWidth  int = 500
	worldHeight int = 500
)

type GameWorld struct{}

type Character struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type ControlSystem struct {
	entity *Character
}

// Add adds the passed Character to the ControlSystem.
func (c *ControlSystem) Add(char *Character) {
	c.entity = char
}

// Remove removes the passed BasicEntity (it'll be a character since we add one) from the ControlSystem.
func (c *ControlSystem) Remove(basic ecs.BasicEntity) {
	if c.entity != nil && basic.ID() == c.entity.ID() {
		c.entity = nil
	}
}

// Update gets the currently pressed key and moves the ControlSystem's entity in the corresponding direction.
func (c *ControlSystem) Update(dt float32) {
	if engo.Input.Button("moveup").Down() {
		c.entity.SpaceComponent.Position.Y -= 5
	}
	if engo.Input.Button("movedown").Down() {
		c.entity.SpaceComponent.Position.Y += 5
	}
	if engo.Input.Button("moveleft").Down() {
		c.entity.SpaceComponent.Position.X -= 5
	}
	if engo.Input.Button("moveright").Down() {
		c.entity.SpaceComponent.Position.X += 5
	}
}

// Preload loads assets for the world.
func (g GameWorld) Preload() {
	// A tmx file can be generated from the Tiled Map Editor.
	// The engo tmx loader only accepts tmx files that are base64 encoded and compressed with zlib.
	// When you add tilesets to the Tiled Editor, the location where you added them from is where the engo loader will look for them
	// Tileset from : http://opengameart.org

	engo.Files.SetRoot("./assets")

	if err := engo.Files.Load("example.tmx", "icon.png"); err != nil {
		panic(err)
	}
}

// Setup sets up the world, e.g. adding tiles, images, characters and their control systems.
func (dg GameWorld) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	common.SetBackground(color.RGBA{0x00, 0x00, 0x00, 0x00})

	// Add an empty RenderSystem.
	w.AddSystem(&common.RenderSystem{})

	// Load the tmx which is XML which places tiles and images in the world.
	resource, err := engo.Files.Resource("example.tmx")
	if err != nil {
		panic(err)
	}
	tmxResource := resource.(common.TMXResource)
	levelData := tmxResource.Level

	// Create render and space components for each of the tiles.
	// Tiles from the tmx are arranged in layers. Grass is in one, trees another.
	tileComponents := make([]*Tile, 0)
	for _, tileLayer := range levelData.TileLayers {
		for _, tileElement := range tileLayer.Tiles {
			if tileElement.Image != nil {

				tile := &Tile{BasicEntity: ecs.NewBasic()}

				// Each tile needs a RenderComponent to draw itself and SpaceComponent to track where it is on-screen.
				tile.RenderComponent = common.RenderComponent{
					Drawable: tileElement,
					Scale:    engo.Point{1, 1},
				}
				tile.SpaceComponent = common.SpaceComponent{
					Position: tileElement.Point,
					Width:    0,
					Height:   0,
				}

				if tileLayer.Name == "grass" {
					tile.RenderComponent.SetZIndex(0)
				}

				if tileLayer.Name == "trees" {
					tile.RenderComponent.SetZIndex(2)
				}

				tileComponents = append(tileComponents, tile)
			}
		}
	}

	// Do the same for all image layers
	for _, imageLayer := range levelData.ImageLayers {
		for _, imageElement := range imageLayer.Images {
			if imageElement.Image != nil {
				tile := &Tile{BasicEntity: ecs.NewBasic()}

				// Like tiles, each image needs a RenderComponent to draw itself and SpaceComponent to track where it is on-screen.
				tile.RenderComponent = common.RenderComponent{
					Drawable: imageElement,
					Scale:    engo.Point{1, 1},
				}
				tile.SpaceComponent = common.SpaceComponent{
					Position: imageElement.Point,
					Width:    0,
					Height:   0,
				}

				if strings.Contains(imageLayer.Name, "clouds") {
					tile.RenderComponent.SetZIndex(3)
				}

				// Add images to the tileComponents, which contains tiles and images.
				tileComponents = append(tileComponents, tile)
			}
		}
	}

	// Set up our movable character. Give it an image.
	character := Character{BasicEntity: ecs.NewBasic()}
	characterTexture, err := common.LoadedSprite("icon.png")
	if err != nil {
		panic(err)
	}

	// A character needs to be rendered and its position tracked like tiles and images.
	character.RenderComponent = common.RenderComponent{
		Drawable: characterTexture,
		Scale:    engo.Point{5, 5},
	}
	character.RenderComponent.SetZIndex(1)
	character.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{engo.CanvasWidth() / 2, engo.CanvasHeight() / 2},
		Width:    characterTexture.Width() * 5,
		Height:   characterTexture.Height() * 5,
	}

	// Add each of the tiles entities and its components to the render system along with the character.
	// At this point the world has 2 systems: CameraSystem (added automatically) and RenderSystem
	// (added ourselves above).
	// A system is an interface which implements Update() and Remove().
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:

			// If we're in the RenderSystem, add the character (with its associated RenderComponent
			// and SpaceComponent).
			// Also add all the tileComponents (i.e. tiles and images) with each one's RenderComponent
			// and SpaceComponent).
			sys.Add(&character.BasicEntity, &character.RenderComponent, &character.SpaceComponent)
			for _, v := range tileComponents {
				sys.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
			}

		}
	}

	// Add a ControlSystem, which is a custom struct implementing engo's System interface.
	// ControlSystem's Update() takes keyboard input (handled by engo) to change the character's position.
	w.AddSystem(&ControlSystem{&character})

	// Add the EntityScroller system which contains the space component of the character and which is bounded to the
	// tmx level dimensions.
	w.AddSystem(&common.EntityScroller{SpaceComponent: &character.SpaceComponent, TrackingBounds: levelData.Bounds()})

	// Register control keys.
	engo.Input.RegisterButton("moveup", engo.KeyArrowUp)
	engo.Input.RegisterButton("moveleft", engo.KeyArrowLeft)
	engo.Input.RegisterButton("moveright", engo.KeyArrowRight)
	engo.Input.RegisterButton("movedown", engo.KeyArrowDown)
}

// Type return a string describing the world type.
func (g GameWorld) Type() string {
	return "GameWorld"
}

func main() {
	opts := engo.RunOptions{
		Title:         "Engo Test",
		Width:         worldWidth,
		Height:        worldHeight,
		ScaleOnResize: false,
	}

	engo.Run(opts, &GameWorld{})
}
