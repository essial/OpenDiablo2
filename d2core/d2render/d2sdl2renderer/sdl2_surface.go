package d2sdl2renderer

import (
	"fmt"
	"image"
	"image/color"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/veandco/go-sdl2/sdl"
)

// static check that we implement our surface interface
var _ d2interface.Surface = &surface{}

type surface struct {
	texture       *sdl.Texture
	stateStackIdx int
	stateStack    [maxStack]surfaceState
	stateCurrent  surfaceState
	width         int
	height        int
	renderer      *Renderer
	format        uint32
	isRenderer    bool
}

func (s *surface) Renderer() d2interface.Renderer {
	return s.renderer
}

func (s *surface) Clear(color color.Color) error {
	sdlMutex.Lock()
	defer sdlMutex.Unlock()
	rr, g, b, a := color.RGBA()

	if err := s.renderer.renderer.SetRenderTarget(s.texture); err != nil {
		return err
	}

	if err := s.renderer.renderer.SetDrawColor(uint8(rr>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)); err != nil {
		return err
	}

	if err := s.renderer.renderer.FillRect(nil); err != nil {
		return err
	}

	if err := s.renderer.renderer.SetRenderTarget(nil); err != nil {
		return err
	}

	return nil
}

func (s *surface) PushSkew(x, y float64) {
	s.pushCurrentState()
	s.stateCurrent.skewX = x
	s.stateCurrent.skewY = y
}

func (s *surface) PushScale(x, y float64) {
	s.pushCurrentState()
	s.stateCurrent.scaleX = x
	s.stateCurrent.scaleY = y
}

func (s *surface) PushSaturation(saturation float64) {
	s.pushCurrentState()
	s.stateCurrent.saturation = saturation
}

func (s *surface) Render(sfc d2interface.Surface) error {
	sdlMutex.Lock()
	defer sdlMutex.Unlock()

	sfcSource := sfc.(*surface)

	srcRect := sdl.Rect{
		W: int32(sfcSource.width),
		H: int32(sfcSource.height),
	}

	targetRect := sdl.Rect{
		X: int32(s.stateCurrent.x),
		Y: int32(s.stateCurrent.y),
		W: int32(sfcSource.width),
		H: int32(sfcSource.height),
	}

	switch sfcSource.stateCurrent.effect {
	case d2enum.DrawEffectPctTransparency25:
		sfcSource.texture.SetAlphaMod(192)
	case d2enum.DrawEffectPctTransparency50:
		sfcSource.texture.SetAlphaMod(128)
	case d2enum.DrawEffectPctTransparency75:
		sfcSource.texture.SetAlphaMod(64)
	case d2enum.DrawEffectModulate:
		sfcSource.texture.SetAlphaMod(255)
		sfcSource.texture.SetBlendMode(sdl.BLENDMODE_MOD)
	case d2enum.DrawEffectBurn:
	case d2enum.DrawEffectNormal:
	case d2enum.DrawEffectMod2XTrans:
	case d2enum.DrawEffectMod2X:
	case d2enum.DrawEffectNone:
		sfcSource.texture.SetAlphaMod(255)
		sfcSource.texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	}

	if !s.isRenderer {
		if err := s.renderer.renderer.SetRenderTarget(s.texture); err != nil {
			return err
		}
	}

	if err := s.renderer.renderer.Copy(sfcSource.texture, &srcRect, &targetRect); err != nil {
		return err
	}

	if !s.isRenderer {
		if err := s.renderer.renderer.SetRenderTarget(nil); err != nil {
			return err
		}
	}

	return nil
}

func createSurface(r *Renderer, width, height int, currentState ...surfaceState) (*surface, error) {
	sdlMutex.Lock()
	defer sdlMutex.Unlock()

	texture, err := r.renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, int32(width), int32(height))
	if err != nil {
		return nil, err
	}

	state := surfaceState{
		effect:     d2enum.DrawEffectNone,
		saturation: defaultSaturation,
		brightness: defaultBrightness,
		skewX:      defaultSkewX,
		skewY:      defaultSkewY,
		scaleX:     defaultScaleX,
		scaleY:     defaultScaleY,
	}
	if len(currentState) > 0 {
		state = currentState[0]
	}

	result := &surface{
		width:         width,
		height:        height,
		texture:       texture,
		stateCurrent:  state,
		stateStackIdx: 0,
		renderer:      r,
	}

	result.format, _, _, _, _ = texture.Query()

	return result, nil
}

func (s *surface) DrawRect(width, height int, color color.Color) {
	r, g, b, a, _ := s.renderer.renderer.GetDrawColor()
	dr, dg, db, da := color.RGBA()
	s.renderer.renderer.SetRenderTarget(s.texture)
	s.renderer.renderer.SetDrawColor(uint8(dr), uint8(dg), uint8(db), uint8(da))
	s.renderer.renderer.FillRect(&sdl.Rect{X: int32(s.stateCurrent.x), Y: int32(s.stateCurrent.y), W: int32(width), H: int32(height)})
	s.renderer.renderer.SetDrawColor(r, g, b, a)
	s.renderer.renderer.SetRenderTarget(nil)
}

func (s *surface) DrawLine(x, y int, color color.Color) {
	r, g, b, a, _ := s.renderer.renderer.GetDrawColor()
	dr, dg, db, da := color.RGBA()
	s.renderer.renderer.SetRenderTarget(s.texture)
	s.renderer.renderer.SetDrawColor(uint8(dr), uint8(dg), uint8(db), uint8(da))
	s.renderer.renderer.DrawLine(int32(s.stateCurrent.x), int32(s.stateCurrent.y), int32(x), int32(y))
	s.renderer.renderer.SetDrawColor(r, g, b, a)
	s.renderer.renderer.SetRenderTarget(nil)
}

func (s *surface) DrawTextf(format string, params ...interface{}) {
	str := fmt.Sprintf(format, params...)
	_ = s.renderer.Print(s, str)
}

func (s *surface) GetSize() (width, height int) {
	return s.width, s.height
}

func (s *surface) GetDepth() int {
	return s.stateStackIdx
}

func (s *surface) Pop() {
	if s.stateStackIdx == 0 {
		panic("empty stack")
	}

	s.stateStackIdx--
	s.stateCurrent = s.stateStack[s.stateStackIdx]
}

func (s *surface) PopN(n int) {
	for i := 0; i < n; i++ {
		s.Pop()
	}
}

func (s *surface) pushCurrentState() {
	s.stateStack[s.stateStackIdx] = s.stateCurrent
	s.stateStackIdx++
	s.stateStack[s.stateStackIdx].Clear()
}

func (s *surface) PushColor(c color.Color) {
	s.pushCurrentState()
	s.stateCurrent.color = c
}

func (s *surface) PushEffect(effect d2enum.DrawEffect) {
	s.pushCurrentState()
	s.stateCurrent.effect = effect
}

func (s *surface) PushFilter(filter d2enum.Filter) {
	s.pushCurrentState()
	// s.stateCurrent.filter = d2ToEbitenFilter(filter)
}

func (s *surface) PushTranslation(x, y int) {
	s.pushCurrentState()
	s.stateCurrent.x += x
	s.stateCurrent.y += y
}

func (s *surface) PushBrightness(brightness float64) {
	s.pushCurrentState()
	s.stateCurrent.brightness = brightness
}

func (s *surface) RenderSection(srcSurface d2interface.Surface, bound image.Rectangle) error {
	sdlMutex.Lock()
	defer sdlMutex.Unlock()

	destRect := sdl.Rect{
		X: int32(s.stateCurrent.x),
		Y: int32(s.stateCurrent.y),
		W: int32(bound.Dx()),
		H: int32(bound.Dy()),
	}

	srcRect := sdl.Rect{
		X: int32(bound.Min.X),
		Y: int32(bound.Min.Y),
		W: int32(bound.Dx()),
		H: int32(bound.Dy()),
	}

	switch s.stateCurrent.effect {
	case d2enum.DrawEffectPctTransparency25:
		s.texture.SetAlphaMod(192)
		s.texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	case d2enum.DrawEffectPctTransparency50:
		s.texture.SetAlphaMod(128)
		s.texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	case d2enum.DrawEffectPctTransparency75:
		s.texture.SetAlphaMod(64)
		s.texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	case d2enum.DrawEffectModulate:
		s.texture.SetAlphaMod(255)
		s.texture.SetBlendMode(sdl.BLENDMODE_MOD)
	case d2enum.DrawEffectBurn:
	case d2enum.DrawEffectNormal:
	case d2enum.DrawEffectMod2XTrans:
	case d2enum.DrawEffectMod2X:
	case d2enum.DrawEffectNone:
		s.texture.SetAlphaMod(255)
		s.texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	}

	if err := s.renderer.renderer.SetRenderTarget(s.texture); err != nil {
		return err
	}

	if err := s.renderer.renderer.Copy(srcSurface.(*surface).texture, &srcRect, &destRect); err != nil {
		return err
	}

	if err := s.renderer.renderer.SetRenderTarget(nil); err != nil {
		return err
	}

	return nil
}

func (s *surface) ReplacePixels(pixels *[]byte) error {
	sdlMutex.Lock()
	defer sdlMutex.Unlock()

	if err := s.texture.Update(nil, *pixels, s.width*4); err != nil {
		return err
	}

	return nil
}

func (s *surface) Screenshot() *image.RGBA {
	panic("implement me")
}
