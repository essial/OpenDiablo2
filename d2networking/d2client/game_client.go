package d2client

import (
	"log"
	"os"

	// "github.com/OpenDiablo2/OpenDiablo2/d2core/d2map/d2mapgen"

	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2map/d2mapengine"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2map/d2mapentity"
	"github.com/OpenDiablo2/OpenDiablo2/d2game/d2player"
	d2cct "github.com/OpenDiablo2/OpenDiablo2/d2networking/d2client/d2clientconnectiontype"
	"github.com/OpenDiablo2/OpenDiablo2/d2networking/d2netpacket"
	"github.com/OpenDiablo2/OpenDiablo2/d2networking/d2netpacket/d2netpackettype"
)

type GameClient struct {
	listener         GameClientListener
	clientConnection ClientConnection
	connectionType   d2cct.ClientConnectionType
	playerState      d2player.PlayerState
	MapAct           *d2mapengine.MapAct
	PlayerId         string
	Players          map[string]*d2mapentity.Player
	Seed             int64
}

// Using the `clientConnection`, opens a connection and passes the savefile path
func (g *GameClient) Open(connectionString string, saveFilePath string) error {
	return g.clientConnection.Open(connectionString, saveFilePath)
}

// Closes the `clientConnection`
func (g *GameClient) Close() error {
	return g.clientConnection.Close()
}

// Closes the `clientConnection`
func (g *GameClient) Destroy() error {
	return g.clientConnection.Close()
}

func (g *GameClient) SetListener(listener GameClientListener) {
	g.listener = listener
}

// Routes the incoming packets to the packet handlers
func (g *GameClient) OnPacketReceived(packet d2netpacket.NetPacket) error {

	switch packet.PacketType {
	// UNSURE: should we be bubbling up errors from these handler calls?
	case d2netpackettype.UpdateServerInfo:
		g.handleUpdateServerInfo(packet)
	case d2netpackettype.AddPlayer:
		g.handleAddPlayer(packet)
	case d2netpackettype.GenerateMap:
		act := d2mapengine.CreateAct(g.Seed, g.playerState.Act)
		g.MapAct = &act
		if g.listener != nil {
			g.listener.OnMapEngineChanged()
		}
	case d2netpackettype.MovePlayer:
		g.handleMovePlayer(packet)
	case d2netpackettype.Ping:
		g.handlePong(packet)
	case d2netpackettype.ServerClosed:
		g.handleServerClosed(packet)
	default:
		log.Fatalf("Invalid packet type: %d", packet.PacketType)
	}

	return nil
}

// Using the `clientConnection`, sends a packet to the server
func (g *GameClient) SendPacketToServer(packet d2netpacket.NetPacket) error {
	return g.clientConnection.SendPacketToServer(packet)
}

func (g *GameClient) handleUpdateServerInfo(p d2netpacket.NetPacket) {
	serverInfo := p.PacketData.(d2netpacket.UpdateServerInfoPacket)
	seed := serverInfo.Seed
	playerId := serverInfo.PlayerId

	g.playerState = serverInfo.PlayerState
	g.Seed = seed
	g.PlayerId = playerId

	if g.listener != nil {
		g.listener.OnLocalPlayerId(playerId)
	}

	log.Printf("Player id set to %s", playerId)
}

func (g *GameClient) handleAddPlayer(p d2netpacket.NetPacket) {

	player := p.PacketData.(d2netpacket.AddPlayerPacket)
	pId := player.Id
	pName := player.Name
	pX := player.X
	pY := player.Y
	pDir := 0
	pHero := player.HeroType
	pStat := player.Stats
	pEquip := player.Equipment
	newPlayer := d2mapentity.CreatePlayer(pId, pName, pX, pY, pDir, pHero, pStat, pEquip)

	g.Players[newPlayer.Id] = newPlayer
	g.MapAct.MapEngine().AddEntity(newPlayer)
}

func (g *GameClient) handleMovePlayer(p d2netpacket.NetPacket) {
	//movePlayer := p.PacketData.(d2netpacket.MovePlayerPacket)
	//
	//player := g.Players[movePlayer.PlayerId]
	//x1, y1 := movePlayer.StartX, movePlayer.StartY
	//x2, y2 := movePlayer.DestX, movePlayer.DestY
	//
	//path, _, _ := g.MapAct.MapEngine().PathFind(x1, y1, x2, y2)
	//
	//if len(path) > 0 {
	//	player.SetPath(path, func() {
	//		tile := g.MapAct.MapEngine().TileAt(player.TileX, player.TileY)
	//		if tile == nil {
	//			return
	//		}
	//
	//		regionType := tile.RegionType
	//		if regionType == d2enum.RegionAct1Town {
	//			player.SetIsInTown(true)
	//		} else {
	//			player.SetIsInTown(false)
	//		}
	//		player.SetAnimationMode(player.GetAnimationMode().String())
	//	})
	//}
}

func (g *GameClient) handlePong(p d2netpacket.NetPacket) {
	pong := d2netpacket.CreatePongPacket(g.PlayerId)
	g.clientConnection.SendPacketToServer(pong)
}

func (g *GameClient) handleServerClosed(p d2netpacket.NetPacket) {
	// TODO: Need to be tied into a character save and exit
	log.Print("Server has been closed")
	os.Exit(0)
}
