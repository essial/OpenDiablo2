package d2app

import (
	"errors"
	"os"
	"runtime/pprof"
	"strconv"

	"github.com/OpenDiablo2/OpenDiablo2/d2game/d2gamescreen"
)

func (a *App) initTerminalCommands() {
	terminalCommands := []struct {
		name string
		desc string
		args []string
		fn   func(args []string) error
	}{
		{"dumpheap", "dumps the heap to pprof/heap.pprof", nil, a.dumpHeap},
		{"fullscreen", "toggles fullscreen", nil, a.toggleFullScreen},
		{"capframe", "captures a still frame", []string{"filename"}, a.setupCaptureFrame},
		{"capgifstart", "captures an animation (start)", []string{"filename"}, a.startAnimationCapture},
		{"capgifstop", "captures an animation (stop)", nil, a.stopAnimationCapture},
		{"vsync", "toggles vsync", nil, a.toggleVsync},
		{"fps", "toggle fps counter", nil, a.toggleFpsCounter},
		{"timescale", "set scalar for elapsed time", []string{"float"}, a.setTimeScale},
		{"quit", "exits the game", nil, a.quitGame},
		{"screen-gui", "enters the gui playground screen", nil, a.enterGuiPlayground},
		{"js", "eval JS scripts", []string{"code"}, a.evalJS},
	}

	for _, cmd := range terminalCommands {
		if err := a.terminal.Bind(cmd.name, cmd.desc, cmd.args, cmd.fn); err != nil {
			a.Fatalf("failed to bind action %q: %v", cmd.name, err.Error())
		}
	}
}

func (a *App) dumpHeap([]string) error {
	if _, err := os.Stat("./pprof/"); os.IsNotExist(err) {
		if err := os.Mkdir("./pprof/", 0750); err != nil {
			a.Fatal(err.Error())
		}
	}

	fileOut, err := os.Create("./pprof/heap.pprof")
	if err != nil {
		a.Error(err.Error())
	}

	if fileOut == nil {
		return errors.New("could not create heap output")
	}

	if err := pprof.WriteHeapProfile(fileOut); err != nil {
		a.Fatal(err.Error())
	}

	if err := fileOut.Close(); err != nil {
		a.Fatal(err.Error())
	}

	return nil
}

func (a *App) evalJS(args []string) error {
	val, err := a.scriptEngine.Eval(args[0])
	if err != nil {
		a.terminal.Errorf(err.Error())
		return nil
	}

	a.Info("%s" + val)

	return nil
}

func (a *App) toggleFullScreen([]string) error {
	fullscreen := !a.renderer.IsFullScreen()
	a.renderer.SetFullScreen(fullscreen)
	a.terminal.Infof("fullscreen is now: %v", fullscreen)

	return nil
}

func (a *App) setupCaptureFrame(args []string) error {
	a.captureState = captureStateFrame
	a.capturePath = args[0]
	a.captureFrames = nil

	return nil
}

func (a *App) startAnimationCapture(args []string) error {
	a.captureState = captureStateGif
	a.capturePath = args[0]
	a.captureFrames = nil

	return nil
}

func (a *App) stopAnimationCapture([]string) error {
	a.captureState = captureStateNone

	return nil
}

func (a *App) toggleVsync([]string) error {
	vsync := !a.renderer.GetVSyncEnabled()
	a.renderer.SetVSyncEnabled(vsync)
	a.terminal.Infof("vsync is now: %v", vsync)

	return nil
}

func (a *App) toggleFpsCounter([]string) error {
	a.showFPS = !a.showFPS
	a.terminal.Infof("fps counter is now: %v", a.showFPS)

	return nil
}

func (a *App) setTimeScale(args []string) error {
	timeScale, err := strconv.ParseFloat(args[0], 64)
	if err != nil || timeScale <= 0 {
		a.terminal.Errorf("invalid time scale value")
		return nil
	}

	a.terminal.Infof("timescale changed from %f to %f", a.timeScale, timeScale)
	a.timeScale = timeScale

	return nil
}

func (a *App) quitGame([]string) error {
	os.Exit(0)
	return nil
}

func (a *App) enterGuiPlayground([]string) error {
	a.screen.SetNextScreen(d2gamescreen.CreateGuiTestMain(a.renderer, a.guiManager, *a.Options.LogLevel, a.asset))
	return nil
}
