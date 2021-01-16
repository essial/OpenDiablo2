package sdl2

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2util"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2asset"
)

const sampleRate = 22050
const logPrefix = "SDL2 Audio Provider"

var _ d2interface.AudioProvider = &SDL2AudioProvider{} // Static check to confirm struct conforms to interface

type SDL2AudioProvider struct {
	asset  *d2asset.AssetManager
	logger *d2util.Logger
}

func Create(l d2util.LogLevel, am *d2asset.AssetManager) *SDL2AudioProvider {
	result := &SDL2AudioProvider{
		asset: am,
	}

	result.logger = d2util.NewLogger()
	result.logger.SetLevel(l)
	result.logger.SetPrefix(logPrefix)

	return result
}

func (a SDL2AudioProvider) PlayBGM(song string) {

}

func (a SDL2AudioProvider) LoadSound(sfx string, loop bool, bgm bool) (d2interface.SoundEffect, error) {
	return CreateSoundEffect()
}

func (a SDL2AudioProvider) SetVolumes(bgmVolume, sfxVolume float64) {

}
