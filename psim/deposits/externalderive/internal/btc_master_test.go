package internal

import (
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestBTCFamilyMaster(t *testing.T) {
	cases := []struct {
		network NetworkType
		private string
		public  string
	}{
		{
			NetworkTypeBTCMainnet,
			`^xprv`,
			`^xpub`,
		},
		{
			NetworkTypeBTCTestnet,
			`^tprv`,
			`^tpub`,
		},
		{
			NetworkTypeDashMainnet,
			`^xprv`,
			`^xpub`,
		},
		{
			NetworkTypeDashTestnet,
			`^tprv`,
			`^tpub`,
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("%d", tc.network), func(t *testing.T) {
			master, err := NewBTCFamilyMaster(tc.network)
			assert.NoError(t, err)
			prv, err := master.ExtendedPrivate()
			assert.NoError(t, err)
			assert.Regexp(t, tc.private, prv)
			pub, err := master.ExtendedPublic()
			assert.NoError(t, err)
			assert.Regexp(t, tc.public, pub)
		})
	}
}
