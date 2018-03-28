package resources

type Blob struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes BlobAttributes `json:"attributes"`
}

func (b Blob) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"id": b.ID,
		"type": b.Type,
		"attributes": b.Attributes,
	}
}

type BlobAttributes struct {
	Value string `json:"value"`
}

func (a BlobAttributes) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"value": a.Value,
	}
}
