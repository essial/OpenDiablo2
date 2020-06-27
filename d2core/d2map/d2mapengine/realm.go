package d2mapengine

import (
	"log"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2data/d2datadict"
)

/*
	A MapRealm represents the state of the maps/levels/quests for a server

	A MapRealm has MapActs
	A MapAct has:
		MapLevels
		MapEngine
	A MapLevel has:
		a MapEngine
		a MapGenerator for the level
		data records from the txt files for the level

	The MapRealm is created by the game server

	The first player to connect to the realm becomes the host
	The host determines the difficulty and which quests are completed

	The Realm, Acts, and Levels do not advance unless they are `active`
	Nothing happens in a realm unless it is active
	Levels do not generate maps until the level becomes `active`

	A Level is active if a player is within it OR in an adjacent level
	An Act is active if one of its levels is active
	The Realm is active if and only if one of its Acts is active
*/
type MapRealm struct {
	seed       int64
	difficulty d2datadict.DifficultyLevelRecord
	acts       []MapAct
	host       string
}

// Advance advances the realm, which advances the acts, which advances the levels...
func (realm *MapRealm) Advance(elapsed float64) {
	for _, act := range realm.acts {
		act.MapEngine().Advance(elapsed)
	}
}

// Init initializes the realm
func (realm *MapRealm) Init(seed int64) {
	log.Printf("Initializing Realm...")

	////////////////////////////////////////////////////////////////////// FIXME
	// We need to set the difficulty level of the realm in order to pull
	// the right data from level details. testing this for now with normal diff
	// NOTE: we would be setting difficulty level in the realm when a host
	// is connected (the first player)
	diffTestKey := "Normal"
	realm.difficulty = d2datadict.DifficultyLevels[diffTestKey] // hack
	////////////////////////////////////////////////////////////////////////////

	realm.seed = seed
	actIds := d2datadict.GetActIds()
	realm.acts = make([]MapAct, len(actIds))

	for _, actID := range actIds {
		realm.acts[actID] = CreateAct(seed, actID)
	}
}

func (realm *MapRealm) Act(id int) *MapAct {
	for idx := range realm.acts {
		if realm.acts[idx].id != id {
			continue
		}

		return &realm.acts[idx]
	}

	return nil
}

func (realm *MapRealm) GetFirstActLevelId(actId int) int {
	return d2datadict.GetFirstLevelIdByActId(actId)
}
