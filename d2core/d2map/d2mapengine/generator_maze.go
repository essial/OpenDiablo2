package d2mapengine

type MapGeneratorMaze struct {
	seed   int64
	level  *MapLevel
	engine *MapEngine
}

func (m *MapGeneratorMaze) init(s int64, l *MapLevel, e *MapEngine) {
	m.seed = s
	m.level = l
	m.engine = e
}

func (m *MapGeneratorMaze) generate() {
	if m.level.details.WorldOffsetX <0 {
		return
	}
	// TODO: This is temporary code that doesn't really do anything...
	//record := d2ds1.FloorShadowRecord{Prop1: 1, Style: 0, Sequence: 0}
	//
	//for y :=0; y < m.level.details.SizeYNormal; y++ {
	//	for x :=0; x < m.level.details.SizeXNormal; x++ {
	//		tile := m.engine.Tile(x + m.level.details.WorldOffsetX, y + m.level.details.WorldOffsetY)
	//		tile.Floors = []d2ds1.FloorShadowRecord{record}
	//	}
	//}
}
