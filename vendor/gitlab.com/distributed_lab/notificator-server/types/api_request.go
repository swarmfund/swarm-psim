package types

import "gitlab.com/distributed_lab/notificator-server/utils"

type APIRequest struct {
	Type          RequestTypeID  `json:"type" db:"type"`
	Token         string         `json:"token" db:"token"`
	PayloadString RequestPayload `json:"payload" db:"payload"`
}

func (r *APIRequest) GetHash() string {
	return utils.Hash(string(r.Type), r.Token, r.PayloadString.raw)
}
