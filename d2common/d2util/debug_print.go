package d2util

import (
	"image"
	"image/draw"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2util/assets"
)

const (
	cw = assets.CharWidth
	ch = assets.CharHeight
)

// GlyphPrinter uses an image containing glyphs to draw text onto surfaces
type GlyphPrinter struct {
	renderer        d2interface.Renderer
	glyphImageTable d2interface.Surface
	glyphsBounds    []image.Rectangle
}

// NewDebugPrinter creates a new debug printer
func NewDebugPrinter(renderer d2interface.Renderer) *GlyphPrinter {
	texImage := assets.CreateTextImage()
	rgba := image.NewRGBA(texImage.Bounds())

	draw.Draw(rgba, texImage.Bounds(), texImage, texImage.Bounds().Min, draw.Over)

	img := renderer.NewSurface(texImage.Bounds().Size().X, texImage.Bounds().Size().Y)
	img.ReplacePixels(rgba.Pix)

	charsPerRow := texImage.Bounds().Size().X / cw
	totalChars := charsPerRow * (texImage.Bounds().Size().Y / ch)

	printer := &GlyphPrinter{
		renderer:        renderer,
		glyphImageTable: img,
		glyphsBounds:    make([]image.Rectangle, totalChars),
	}

	for idx := 0; idx < totalChars; idx++ {
		sx := (idx % charsPerRow) * cw
		sy := (idx / charsPerRow) * ch
		printer.glyphsBounds[idx] = image.Rect(sx, sy, sx+cw, sy+ch)
	}

	return printer
}

// Print draws the string str on the image on left top corner.
//
// The available runes are in U+0000 to U+00FF, which is C0 Controls and
// Basic Latin and C1 Controls and Latin-1 Supplement.
//
// DebugPrint always returns nil as of 1.5.0-alpha.
func (p *GlyphPrinter) Print(target interface{}, str string) error {
	p.PrintAt(target.(d2interface.Surface), str)
	return nil
}

// PrintAt draws the string str on the image at (x, y) position.
// The available runes are in U+0000 to U+00FF, which is C0 Controls and
// Basic Latin and C1 Controls and Latin-1 Supplement.
func (p *GlyphPrinter) PrintAt(target interface{}, str string) {
	p.drawDebugText(target.(d2interface.Surface), str, 0, 0, false)
}

func (p *GlyphPrinter) drawDebugText(target d2interface.Surface, str string, ox, oy int, shadow bool) {
	px := 0
	py := 0

	target.PushEffect(d2enum.DrawEffectNormal)

	for idx := range str {
		if str[idx] == '\n' {
			px = 0
			py += ch

			continue
		}

		target.PushTranslation(px+ox, py+oy)
		target.RenderSection(p.glyphImageTable, p.glyphsBounds[int(str[idx])])
		target.Pop()

		px += cw
	}

	target.Pop()
}
