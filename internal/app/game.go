package app

type game struct {
	screenWidth, screenHeight int

	temperatureStep       float64
	temperatureLosing     float64
	heatEmitterEfficiency float64

	// game flags
	pause bool

	// debug flags
	debug     bool
	debugTemp bool
	drawTemp  bool
	bloom     bool
}
