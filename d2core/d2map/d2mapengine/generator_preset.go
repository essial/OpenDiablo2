package d2mapengine

import (
	"math/rand"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2map/d2mapstamp"
)

type MapGeneratorPreset struct {
	seed   int64
	level  *MapLevel
	engine *MapEngine
}

func (m *MapGeneratorPreset) init(s int64, l *MapLevel, e *MapEngine) {
	m.seed = s
	m.level = l
	m.engine = e
}

func (m *MapGeneratorPreset) generate() {
	rand.Seed(m.seed)

	levelTypeId := d2enum.RegionIdType(m.level.details.LevelType)
	levelPresetId := m.level.preset.DefinitionId

	stamp := d2mapstamp.LoadStamp(levelTypeId, levelPresetId, -1)
	x := m.level.details.WorldOffsetX
	y := m.level.details.WorldOffsetY
	if x < 0 {
		return
	}
	if y < 0 {
		return
	}
	m.engine.PlaceStamp(stamp, x, y)
}
