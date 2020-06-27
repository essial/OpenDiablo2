package d2mapengine

import (
	"log"
	"strings"

	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2asset"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dt1"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"

	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2map/d2mapentity"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2map/d2mapstamp"

	"github.com/OpenDiablo2/OpenDiablo2/d2common"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2data/d2datadict"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2ds1"
)

// MapEngine represents the map data for a specific location
type MapEngine struct {
	act           int                        // The map's act
	seed          int64                      // The map seed
	entities      []d2mapentity.MapEntity    // Entities on the map
	tiles         map[int]*d2ds1.TileRecord   // The map tiles
	size          d2common.Size              // The size of the map, in tiles
	levelType     d2datadict.LevelTypeRecord // The level type of this map
	dt1TileData   []d2dt1.Tile               // The DT1 tile data
	startSubTileX int                        // The starting X position
	startSubTileY int                        // The starting Y position
	dt1Files      []string                   // The list of DS1 strings
}

// InitializeMapEngineForAct creates a new instance of the map engine
func InitializeMapEngineForAct(act *MapAct) {
	act.mapEngine.act = act.id // Acts are 0 relative
	act.mapEngine.seed = act.seed

	mapWidth := 0
	mapHeight := 0

	for levelIdx := range act.levels {
		details := &act.levels[levelIdx].details
		right := details.WorldOffsetX + details.SizeXNormal + 1
		top := details.WorldOffsetY + details.SizeYNormal + 1

		if mapWidth < right {
			mapWidth = right
		}

		if mapHeight < top {
			mapHeight = top
		}
	}

	act.mapEngine.ResetMap(d2enum.RegionNone, mapWidth, mapHeight)

	for levelIdx := range act.levels {
		level := &act.levels[levelIdx]

		var generator MapGenerator = nil

		switch level.details.LevelGenerationType {
		case d2enum.LevelTypeNone: // Should only be for first record in the excel doc
			continue
		case d2enum.LevelTypeRandomMaze:
			generator = &MapGeneratorMaze{}
		case d2enum.LevelTypeWilderness:
			generator = &MapGeneratorWilderness{}
		case d2enum.LevelTypePreset:
			generator = &MapGeneratorPreset{}
		default:
			panic("Unknown level type specified. Cannot construct a generator.")
		}

		generator.init(act.seed, level, &act.mapEngine)
		generator.generate()

		levelType := d2datadict.LevelTypes[d2enum.RegionIdType(level.types.Id)]
		for idx := range levelType.Files {
			act.mapEngine.addDT1(levelType.Files[idx])
		}
	}
}

func (m *MapEngine) Seed() int64 {
	return m.seed
}

func (m *MapEngine) Act() int {
	return m.act
}

// Returns the starting position on the map in sub-tiles
func (m *MapEngine) GetStartingPosition() (int, int) {
	return m.startSubTileX, m.startSubTileY
}

func (m *MapEngine) ResetMap(levelType d2enum.RegionIdType, width, height int) {
	m.entities = make([]d2mapentity.MapEntity, 0)
	m.levelType = d2datadict.LevelTypes[levelType]
	m.size = d2common.Size{Width: width, Height: height}
	m.tiles = make(map[int]*d2ds1.TileRecord)
	m.dt1TileData = make([]d2dt1.Tile, 0)
	//m.walkMesh = make([]d2common.PathTile, width*height*25)
	m.dt1Files = make([]string, 0)

	for idx := range m.levelType.Files {
		m.addDT1(m.levelType.Files[idx])
	}

}

func (m *MapEngine) addDT1(fileName string) {
	if len(fileName) == 0 || fileName == "0" {
		return
	}
	fileName = strings.ToLower(fileName)
	for i := 0; i < len(m.dt1Files); i++ {
		if m.dt1Files[i] == fileName {
			return
		}
	}

	fileData, err := d2asset.LoadFile("/data/global/tiles/" + fileName)
	if err != nil {
		log.Printf("Could not load /data/global/tiles/%s", fileName)
		return
		//panic(err)
	}
	dt1, _ := d2dt1.LoadDT1(fileData)
	m.dt1TileData = append(m.dt1TileData, dt1.Tiles...)
	m.dt1Files = append(m.dt1Files, fileName)
}

func (m *MapEngine) AddDS1(fileName string) {
	if len(fileName) == 0 || fileName == "0" {
		return
	}

	fileData, err := d2asset.LoadFile("/data/global/tiles/" + fileName)
	if err != nil {
		panic(err)
	}
	ds1, _ := d2ds1.LoadDS1(fileData)
	for idx := range ds1.Files {
		dt1File := ds1.Files[idx]
		dt1File = strings.ToLower(dt1File)
		dt1File = strings.Replace(dt1File, "c:", "", -1)       // Yes they did...
		dt1File = strings.Replace(dt1File, ".tg1", ".dt1", -1) // Yes they did...
		dt1File = strings.Replace(dt1File, "\\d2\\data\\global\\tiles\\", "", -1)
		m.addDT1(strings.Replace(dt1File, "\\", "/", -1))
	}
}

func (m *MapEngine) FindTile(style, sequence, tileType int32) d2dt1.Tile {
	for idx := range m.dt1TileData {
		if m.dt1TileData[idx].Style == style && m.dt1TileData[idx].Sequence == sequence && m.dt1TileData[idx].Type == tileType {
			return m.dt1TileData[idx]
		}
	}
	panic("Could not find the requested tile!")
}

// Returns the level type of this map
func (m *MapEngine) LevelType() d2datadict.LevelTypeRecord {
	return m.levelType
}

// Returns the size of the map (in sub-tiles)
func (m *MapEngine) Size() d2common.Size {
	return m.size
}

func (m *MapEngine) Tile(x, y int) *d2ds1.TileRecord {
	result, ok := m.tiles[x+(y*m.size.Width)]
	if !ok {
		return nil
	}
	return result
}

func (m *MapEngine) SetTile(x, y int, tile *d2ds1.TileRecord)  {
	m.tiles[x+(y*m.size.Width)] = tile
}

// Places a stamp at the specified location.
// Also adds any entities from the stamp to the map engine
func (m *MapEngine) PlaceStamp(stamp *d2mapstamp.Stamp, tileOffsetX, tileOffsetY int) {
	stampSize := stamp.Size()
	stampW := stampSize.Width
	stampH := stampSize.Height

	mapW := m.size.Width
	mapH := m.size.Height

	xMin := tileOffsetX
	yMin := tileOffsetY
	xMax := xMin + stampSize.Width
	yMax := yMin + stampSize.Height

	if (xMin < 0) || (yMin < 0) || (xMax > mapW) || (yMax > mapH) {
		panic("Tried placing a stamp outside the bounds of the map")
	}

	// Copy over the map tile data
	for y := 0; y < stampH; y++ {
		for x := 0; x < stampW; x++ {
			tile := stamp.Tile(x, y)
			tile.RegionType = stamp.RegionType()

			m.SetTile(x+xMin, y+yMin, tile)


		}
	}

	// Copy over the entities
	m.entities = append(m.entities, stamp.Entities(tileOffsetX, tileOffsetY)...)
}

// converts x,y tile coordinate into index in MapEngine.tiles
func (m *MapEngine) tileCoordinateToIndex(x, y int) int {
	return x + (y * m.size.Width)
}

// converts tile index from MapEngine.tiles to x,y coordinate
func (m *MapEngine) tileIndexToCoordinate(index int) (int, int) {
	return (index % m.size.Width), (index / m.size.Width)
}

// Returns a reference to a map tile based on the tile X,Y coordinate
func (m *MapEngine) TileAt(tileX, tileY int) *d2ds1.TileRecord {
	return m.Tile(tileX, tileY)
}

func (m *MapEngine) Tiles() map[int]*d2ds1.TileRecord {
	return m.tiles
}

// Returns a reference to the map entities
func (m *MapEngine) Entities() *[]d2mapentity.MapEntity {
	return &m.entities
}

// Adds an entity to the map engine
func (m *MapEngine) AddEntity(entity d2mapentity.MapEntity) {
	m.entities = append(m.entities, entity)
}

// Removes an entity from the map engine
func (m *MapEngine) RemoveEntity(entity d2mapentity.MapEntity) {
	if entity == nil {
		return
	}
	panic("Removing entities is not currently implemented")
	//m.entities.Remove(entity)
}

func (m *MapEngine) GetTiles(style, sequence, tileType int32) []d2dt1.Tile {
	var tiles []d2dt1.Tile
	for idx := range m.dt1TileData {
		if m.dt1TileData[idx].Style != style || m.dt1TileData[idx].Sequence != sequence || m.dt1TileData[idx].Type != tileType {
			continue
		}
		tiles = append(tiles, m.dt1TileData[idx])
	}
	if len(tiles) == 0 {
		log.Printf("Unknown tile ID [%d %d %d]\n", style, sequence, tileType)
		return nil
	}
	return tiles
}

func (m *MapEngine) GetStartPosition() (float64, float64) {
	for tileY := 0; tileY < m.size.Height; tileY++ {
		for tileX := 0; tileX < m.size.Width; tileX++ {
			tile := m.tiles[tileX+(tileY*m.size.Width)]
			if tile == nil {
				continue
			}
			for idx := range tile.Walls {
				if tile.Walls[idx].Type.Special() && tile.Walls[idx].Style == 30 {
					return float64(tileX) + 0.5, float64(tileY) + 0.5
				}
			}
		}
	}

	return m.GetCenterPosition()
}

// Returns the center of the map
func (m *MapEngine) GetCenterPosition() (float64, float64) {
	return float64(m.size.Width) / 2.0, float64(m.size.Height) / 2.0
}

// Advances time on the map engine
func (m *MapEngine) Advance(tickTime float64) {
	for idx := range m.entities {
		m.entities[idx].Advance(tickTime)
	}
}

// Checks if a tile exists
func (m *MapEngine) TileExists(tileX, tileY int) bool {
	tileIndex := m.tileCoordinateToIndex(tileX, tileY)
	if valid := (tileIndex >= 0) && (tileIndex <= len(m.tiles)); valid {
		tile := m.tiles[tileIndex]
		numFeatures := len(tile.Floors)
		numFeatures += len(tile.Shadows)
		numFeatures += len(tile.Walls)
		numFeatures += len(tile.Substitutions)
		return numFeatures > 0
	}
	return false
}

//func (m *MapEngine) GenerateMap(regionType d2enum.RegionIdType, levelPreset int, fileIndex int, cacheTiles bool) {
//	region := d2mapstamp.LoadStamp(regionType, levelPreset, fileIndex)
//	regionSize := region.Size()
//	m.ResetMap(regionType, regionSize.Width, regionSize.Height)
//	m.PlaceStamp(region, 0, 0)
//}

func (m *MapEngine) GetTileData(style int32, sequence int32, tileType d2enum.TileType) *d2dt1.Tile {
	for idx := range m.dt1TileData {
		if m.dt1TileData[idx].Style == style && m.dt1TileData[idx].Sequence == sequence && m.dt1TileData[idx].Type == int32(tileType) {
			return &m.dt1TileData[idx]
		}
	}
	return nil
}
