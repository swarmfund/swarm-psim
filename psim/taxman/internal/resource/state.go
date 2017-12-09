package resource

import "gitlab.com/tokend/psim/psim/taxman/internal/state"

type StateResource struct {
	State state.State `jsonapi:"attr,state"`
}
