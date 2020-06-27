package d2mapengine

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2data/d2datadict"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
)

type MapLevel struct {
	details       d2datadict.LevelDetailsRecord
	preset        d2datadict.LevelPresetRecord
	warps         []d2datadict.LevelWarpRecord
	substitutions d2datadict.LevelSubstitutionRecord
	types         d2datadict.LevelTypeRecord
	//generator     MapGenerator

}

//func (level *MapLevel) Advance(elapsed float64) {
//	level.mapEngine.Advance(elapsed)
//}

func CreateMapLevel(levelDetail d2datadict.LevelDetailsRecord) MapLevel {
	return MapLevel{
		details:       d2datadict.LevelDetailsForLevel(levelDetail.Id),
		preset:        d2datadict.LevelPresetForLevel(levelDetail.Id),
		warps:         d2datadict.GetLevelWarpsByLevelId(levelDetail.Id),
		substitutions: d2datadict.LevelSubstitutions[levelDetail.SubType],
		types:         d2datadict.LevelTypes[d2enum.RegionIdType(levelDetail.LevelType)],
	}
}

//switch level.details.LevelGenerationType {
//case d2enum.LevelTypeRandomMaze:
//	level.generator = &MapGeneratorMaze{}
//case d2enum.LevelTypeWilderness:
//	level.generator = &MapGeneratorWilderness{}
//case d2enum.LevelTypePreset:
//	level.generator = &MapGeneratorPreset{}
//default:
//	panic("Unknown level type specified. Cannot construct a generator.")
//}
//
//if level.generator != nil {
//	log.Printf("Initializing Level: %s", level.details.Name)
//	level.generator.init(seed, level, level.mapEngine)
//}

//func (level *MapLevel) GenerateMap() {
//	if level.isGenerated {
//		return
//	}
//	log.Printf("Generating Level: %s", level.details.Name)
//	level.generator.generate()
//	level.mapEngine.RegenerateWalkPaths()
//	level.isGenerated = true
//}
