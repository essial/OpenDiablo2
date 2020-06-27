package d2client

import (
	"log"

	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2map/d2mapentity"
	d2cct "github.com/OpenDiablo2/OpenDiablo2/d2networking/d2client/d2clientconnectiontype"
	"github.com/OpenDiablo2/OpenDiablo2/d2networking/d2client/d2localclient"
	"github.com/OpenDiablo2/OpenDiablo2/d2networking/d2client/d2remoteclient"
)

// Creates a connections to the server and returns a game client instance
func Create(connectionType d2cct.ClientConnectionType) *GameClient {
	result := GameClient{
		Players:        make(map[string]*d2mapentity.Player),
		connectionType: connectionType,
	}

	switch connectionType {
	case d2cct.LANClient:
		result.clientConnection = d2remoteclient.Create()
	case d2cct.LANServer:
		openSocket := true
		result.clientConnection = d2localclient.Create(openSocket)
	case d2cct.Local:
		dontOpenSocket := false
		result.clientConnection = d2localclient.Create(dontOpenSocket)
	default:
		log.Panicf("unknown client connection type specified: %d", connectionType)
	}
	result.clientConnection.SetClientListener(&result)

	return &result
}
