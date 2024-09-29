package resources

import (
	"embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "image/jpeg" // needed for ebitenutil.NewImageFromReader()
	_ "image/png"  // needed for ebitenutil.NewImageFromReader()
	"log"
)

// Imgs is the global variable that holds all image resources
var Imgs *ImgResources

// ImgResources is a collection of all images
type ImgResources struct {
	Icon        *ebiten.Image // 288 x 288
	BgOcean     *ebiten.Image // 2475 x 1532
	BgContinent *ebiten.Image // 2475 x 1392  [16:9]
	Fortress    *ebiten.Image // 1024 x 1024
	Village     *ebiten.Image // 500 x 500
	Field       *ebiten.Image // 773 x 773
}

func init() {
	Imgs = &ImgResources{
		Icon:        loadGameImg("img/icon.png"),
		BgOcean:     loadGameImg("img/bg_ocean.jpg"),
		BgContinent: loadGameImg("img/bg_continent.png"),
		Fortress:    loadGameImg("img/fortress.png"),
		Village:     loadGameImg("img/village.png"),
		Field:       loadGameImg("img/field.png"),
	}
}

//go:embed img
var gFS embed.FS

func loadGameImg(name string) *ebiten.Image {
	// open reader
	r, err := gFS.Open(name)
	if err != nil {
		log.Fatalf("err: loadGameImg: %v\n", err)
	}
	// get image
	eim, _, err := ebitenutil.NewImageFromReader(r)
	if err != nil {
		log.Fatalf("err: loadGameImg: %v\n", err)
	}
	// return
	return eim
}
