package d2audio

import (
	"fmt"
	"strings"

	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2audio/sdl2"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2util"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2asset"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2audio/ebiten"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2config"
)

func Create(logLevel d2util.LogLevel, assetManager *d2asset.AssetManager, config *d2config.Configuration) d2interface.AudioProvider {
	switch strings.ToUpper(config.Backend) {
	case "EBITEN":
		return ebiten.Create(logLevel, assetManager)
	case "SDL2":
		return sdl2.Create(logLevel, assetManager)
	default:
		panic(fmt.Errorf("no audio provider available for backend %s", config.Backend))
	}
}
