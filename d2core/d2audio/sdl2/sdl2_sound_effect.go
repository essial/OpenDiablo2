package sdl2

type SDL2SoundEffect struct {
}

func CreateSoundEffect() (*SDL2SoundEffect, error) {
	result := &SDL2SoundEffect{}

	return result, nil
}

func (S SDL2SoundEffect) Play() {

}

func (S SDL2SoundEffect) Stop() {

}

func (S SDL2SoundEffect) SetPan(pan float64) {

}

func (S SDL2SoundEffect) IsPlaying() bool {
	return true
}

func (S SDL2SoundEffect) SetVolume(volume float64) {

}
