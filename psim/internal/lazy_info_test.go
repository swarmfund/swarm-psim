package internal

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/internal/mocks"
	"gitlab.com/tokend/regources"
)

func TestLazyInfo_Info(t *testing.T) {
	infoer := &mocks.Infoer{}
	entry := logan.New().Out(ioutil.Discard)
	info := regources.Info{}

	t.Run("is lazy", func(t *testing.T) {
		infoer.AssertNotCalled(t, "Info")
		NewLazyInfo(context.TODO(), entry, infoer)
		infoer.AssertExpectations(t)
	})

	t.Run("retry", func(t *testing.T) {
		infoer.On("Info").Return(nil, errors.New("hello")).Once()
		infoer.On("Info").Return(&info, nil).Once()

		lazy := NewLazyInfo(context.TODO(), entry, infoer)
		got, err := lazy.Info()
		assert.Equal(t, &info, got)
		assert.NoError(t, err)

		infoer.AssertExpectations(t)
	})
}
