package d2render

import (
	"fmt"
	"strings"

	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2render/d2sdl2renderer"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2config"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2render/ebiten"
)

func Create(config *d2config.Configuration) (d2interface.Renderer, error) {
	switch strings.ToUpper(config.Backend) {
	case "EBITEN":
		return ebiten.CreateRenderer(config)
	case "SDL2":
		return d2sdl2renderer.CreateRenderer(config)
	default:
		panic(fmt.Errorf("no renderer available for backend: %s", config.Backend))
	}
}
