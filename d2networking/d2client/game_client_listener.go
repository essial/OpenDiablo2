package d2client

type GameClientListener interface {
	OnMapEngineChanged()
	OnLocalPlayerId(playerId string)
}
