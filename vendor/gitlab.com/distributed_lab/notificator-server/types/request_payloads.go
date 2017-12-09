package types


type RequestPayload struct {
	raw string
}

func (p *RequestPayload) UnmarshalJSON(data []byte) error {
	p.raw = string(data)
	return nil
}

func (p *RequestPayload) Raw() string {
	return p.raw
}

