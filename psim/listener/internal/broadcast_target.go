package internal

// Target defines a place, where events sent to
// SendEvent should be retried by the caller if success == false and err == nil
// If err != nil - handle error - can be ignored or fixed and retried. Success doesn't matter
// If success == true - event sent
type Target interface {
	SendEvent(event *BroadcastedEvent) (err error)
}
