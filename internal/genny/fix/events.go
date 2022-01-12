package fix

const (
	// EvtFixStart is emitted when fixing starts
	EvtFixStart = "buffalo:fix:start"
	// EvtFixStop is emitted when fixing stops
	EvtFixStop = "buffalo:fix:stop"
	// EvtFixStopErr is emitted when fixing is stopped due to an error
	EvtFixStopErr = "buffalo:fix:stop:err"
)
