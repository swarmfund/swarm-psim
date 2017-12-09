package notificator

// method only purpose of which is compile time type check
func (payload EmailRequestPayload) legitPayload() {}

// method only purpose of which is compile time type check
func (payload SMSRequestPayload) legitPayload() {}
