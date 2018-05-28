package internal

import "context"

type Extractor interface {
	Extract(ctx context.Context) (<-chan TxData, error)
}
