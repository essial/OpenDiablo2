package d2datadict

import (
	"log"

	"github.com/OpenDiablo2/OpenDiablo2/d2common"
)

type LevelWarpRecord struct {
	Id         int32
	SelectX    int32
	SelectY    int32
	SelectDX   int32
	SelectDY   int32
	ExitWalkX  int32
	ExitWalkY  int32
	OffsetX    int32
	OffsetY    int32
	LitVersion bool
	Tiles      int32
	Direction  string
}

var LevelWarps map[int]LevelWarpRecord

func LoadLevelWarps(levelWarpData []byte) {
	LevelWarps = make(map[int]LevelWarpRecord)
	streamReader := d2common.CreateStreamReader(levelWarpData)
	numRecords := int(streamReader.GetInt32())
	for i := 0; i < numRecords; i++ {
		id := int(streamReader.GetInt32())
		record := LevelWarpRecord{}
		record.Id = int32(id)
		record.SelectX = streamReader.GetInt32()
		record.SelectY = streamReader.GetInt32()
		record.SelectDX = streamReader.GetInt32()
		record.SelectDY = streamReader.GetInt32()
		record.ExitWalkX = streamReader.GetInt32()
		record.ExitWalkY = streamReader.GetInt32()
		record.OffsetX = streamReader.GetInt32()
		record.OffsetY = streamReader.GetInt32()
		record.LitVersion = streamReader.GetInt32() == 1
		record.Tiles = streamReader.GetInt32()
		record.Direction = string(streamReader.GetByte())
		streamReader.SkipBytes(3)
		LevelWarps[id] = record
	}
	log.Printf("Loaded %d level warps", len(LevelWarps))
}
