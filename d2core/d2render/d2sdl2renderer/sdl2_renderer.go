package d2sdl2renderer

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2config"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2util"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/veandco/go-sdl2/sdl"
)

// static check that we implement our renderer interface
var (
	_ d2interface.Renderer = &Renderer{}
	_ d2interface.Surface  = &Renderer{}
)

const (
	maxStack          = 1024
	screenWidth       = 800
	screenHeight      = 600
	defaultSaturation = 1.0
	defaultBrightness = 1.0
	defaultSkewX      = 0.0
	defaultSkewY      = 0.0
	defaultScaleX     = 1.0
	defaultScaleY     = 1.0
)

type Renderer struct {
	*surface
	window     *sdl.Window
	renderer   *sdl.Renderer
	fullscreen bool
	cursorPosX int
	cursorPosY int
	*d2util.GlyphPrinter
}

func (r *Renderer) Renderer() d2interface.Renderer {
	return r
}

func (r *Renderer) ShowPanicScreen(message string) {
	panic("implement me")
}

func CreateRenderer(cfg *d2config.Configuration) (*Renderer, error) {
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_EVENTS); err != nil {
		return nil, err
	}

	window, err := sdl.CreateWindow("OpenDiablo 2", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 640, 480, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE|sdl.WINDOW_INPUT_FOCUS)
	if err != nil {
		return nil, err
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_TARGETTEXTURE|sdl.RENDERER_PRESENTVSYNC)
	// renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
	if err != nil {
		return nil, err
	}

	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	renderer.SetIntegerScale(false)
	renderer.SetLogicalSize(800, 600)
	window.SetMinimumSize(800, 600)

	result := &Renderer{
		window:   window,
		renderer: renderer,
	}

	sfcState := surfaceState{
		//filter:     ebiten.FilterNearest,
		effect:     d2enum.DrawEffectNone,
		saturation: defaultSaturation,
		brightness: defaultBrightness,
		skewX:      defaultSkewX,
		skewY:      defaultSkewY,
		scaleX:     defaultScaleX,
		scaleY:     defaultScaleY,
	}

	if result.surface, err = createSurface(result, screenWidth, screenHeight, sfcState); err != nil {
		return nil, err
	}

	result.GlyphPrinter = d2util.NewDebugPrinter(result)
	result.isRenderer = true

	sdl.ShowCursor(0)

	return result, nil
}

func (r *Renderer) GetRendererName() string {
	return "SDL2"
}

func (r *Renderer) SetWindowIcon(fileName string) {
}

func (r *Renderer) Run(render d2interface.RenderCallback, update d2interface.UpdateCallback, width, height int, title string) error {
	r.window.SetTitle(title)
	r.window.SetSize(int32(width), int32(height))

	for {
		sdlMutex.Lock()
		if err := r.renderer.SetDrawColor(0, 0, 0, 0); err != nil {
			return err
		}

		if err := r.renderer.Clear(); err != nil {
			return err
		}
		sdlMutex.Unlock()

		update()
		if err := render(r); err != nil {
			return err
		}

		sdlMutex.Lock()
		r.renderer.Present()
		sdlMutex.Unlock()
	}
}

func (r *Renderer) IsDrawingSkipped() bool {
	return false
}

func (r *Renderer) CreateSurface(surface d2interface.Surface) (d2interface.Surface, error) {
	panic("implement me")
}

func (r *Renderer) NewSurface(width, height int) d2interface.Surface {
	result, err := createSurface(r, width, height)
	if err != nil {
		panic(err)
	}

	return result
}

func (r *Renderer) IsFullScreen() bool {
	return r.fullscreen
}

func (r *Renderer) SetFullScreen(fullScreen bool) {
	if fullScreen == r.fullscreen {
		return
	}

	if fullScreen {
		r.window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
	} else {
		r.window.SetFullscreen(0)
	}
}

func (r *Renderer) SetVSyncEnabled(vsync bool) {
}

func (r *Renderer) GetVSyncEnabled() bool {
	return false
}

func (r *Renderer) GetCursorPos() (int, int) {
	return r.cursorPosX, r.cursorPosY
}

func (r *Renderer) CurrentFPS() float64 {
	return 60.0
}
