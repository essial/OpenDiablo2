package d2mapengine

import (
	"log"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2data/d2datadict"
)

// MapAct is a structure that represents an act's map
type MapAct struct {
	id        int
	seed      int64
	mapEngine MapEngine
	levels    []MapLevel
}


func (act *MapAct) MapEngine() *MapEngine{
	return &act.mapEngine
}

// CreateAct create an act based on the act id
func CreateAct(seed int64, actId int) MapAct {
	actLevelRecords := d2datadict.GetLevelDetailsByActId(actId)

	act := MapAct{
		id:     actId,
		seed:   seed,
		levels: make([]MapLevel, len(actLevelRecords)),
	}

	log.Printf("Initializing Act %d", actId)

	for idx := range actLevelRecords {
		act.levels[idx] = CreateMapLevel(actLevelRecords[idx])
	}

	InitializeMapEngineForAct(&act)

	return act
}
