package d2netpacket

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2game/d2player"
	"github.com/OpenDiablo2/OpenDiablo2/d2networking/d2netpacket/d2netpackettype"
)

type UpdateServerInfoPacket struct {
	Seed     int64  `json:"seed"`
	PlayerId string `json:"playerId"`
	PlayerState d2player.PlayerState `json:"playerState"`
}

func CreateUpdateServerInfoPacket(seed int64, playerId string, playerState d2player.PlayerState) NetPacket {
	return NetPacket{
		PacketType: d2netpackettype.UpdateServerInfo,
		PacketData: UpdateServerInfoPacket{
			Seed:     seed,
			PlayerId: playerId,
			PlayerState: playerState,
		},
	}
}
