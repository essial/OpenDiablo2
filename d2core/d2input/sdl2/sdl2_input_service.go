package sdl2

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/veandco/go-sdl2/sdl"
	"os"
)

const (
	mousePressThreshold = 25
)

// static check that we implement our input service interface
var _ d2interface.InputService = &InputService{}

type mouseEventInfo struct {
	state bool
	time  uint32
}

type InputService struct {
	mouseState []mouseEventInfo
	cursorPosX int
	cursorPosY int
}

func CreateInputService() *InputService {
	result := &InputService{
		mouseState: make([]mouseEventInfo, 16),
	}

	return result
}

func (i *InputService) Process() {
	var event sdl.Event
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			os.Exit(0)
		case *sdl.MouseMotionEvent:
			i.cursorPosX = int(t.X)
			i.cursorPosY = int(t.Y)
		case *sdl.MouseButtonEvent:
			i.mouseState[t.Button].time = sdl.GetTicks()
			i.mouseState[t.Button].state = t.Type == sdl.MOUSEBUTTONDOWN
		}
	}
}

func (i *InputService) CursorPosition() (x int, y int) {
	return i.cursorPosX, i.cursorPosY
}

func (i *InputService) InputChars() []rune {
	return []rune{}
}

func (i *InputService) IsKeyPressed(key d2enum.Key) bool {
	return false
}

func (i *InputService) IsKeyJustPressed(key d2enum.Key) bool {
	return false
}

func (i *InputService) IsKeyJustReleased(key d2enum.Key) bool {
	return false
}

func (i *InputService) IsMouseButtonPressed(button d2enum.MouseButton) bool {
	switch button {
	case d2enum.MouseButtonLeft:
		return i.mouseState[1].state
	case d2enum.MouseButtonRight:
		return i.mouseState[2].state
	case d2enum.MouseButtonMiddle:
		return i.mouseState[3].state
	default:
		return false
	}
}

func (i *InputService) IsMouseButtonJustPressed(button d2enum.MouseButton) bool {
	switch button {
	case d2enum.MouseButtonLeft:
		return i.mouseState[1].state == true && (sdl.GetTicks()-i.mouseState[1].time) < mousePressThreshold
	case d2enum.MouseButtonRight:
		return i.mouseState[2].state == true && (sdl.GetTicks()-i.mouseState[2].time) < mousePressThreshold
	case d2enum.MouseButtonMiddle:
		return i.mouseState[3].state == true && (sdl.GetTicks()-i.mouseState[3].time) < mousePressThreshold
	default:
		return false
	}
}

func (i *InputService) IsMouseButtonJustReleased(button d2enum.MouseButton) bool {
	switch button {
	case d2enum.MouseButtonLeft:
		return i.mouseState[1].state == false && (sdl.GetTicks()-i.mouseState[1].time) < mousePressThreshold
	case d2enum.MouseButtonRight:
		return i.mouseState[2].state == false && (sdl.GetTicks()-i.mouseState[2].time) < mousePressThreshold
	case d2enum.MouseButtonMiddle:
		return i.mouseState[3].state == false && (sdl.GetTicks()-i.mouseState[3].time) < mousePressThreshold
	default:
		return false
	}
}

func (i *InputService) KeyPressDuration(key d2enum.Key) int {
	return 0
}
