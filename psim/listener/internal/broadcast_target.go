package internal

// Target defines a place, where events sent to
type Target interface {
	SendEvent(event BroadcastedEvent) error
}

/*
type BufferedTarget interface {
	SendEventsInBatches(events []BroadcastedEvent) error
}
*/
