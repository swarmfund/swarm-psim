package listener

// SalesforceTarget represents a target with salesforce api
type SalesforceTarget struct {
	*SalesforceConnector
}

// NewSalesforceTarget constructs a target
func (sc *SalesforceConnector) GetTarget() *SalesforceTarget {
	return &SalesforceTarget{
		sc,
	}
}

func (st *SalesforceTarget) SendEvent(event BroadcastedEvent) error {
	return nil
}
