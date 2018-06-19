package derive

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestETHDeriver_ChildAddress(t *testing.T) {
	cases := []struct {
		key      string
		child    uint64
		expected string
	}{
		{
			`xprv9s21ZrQH143K3oKadkPYNM5cAchBSw59qowfjCBZf5CRN6zk9q9hwtxJvWf4bDaPEjCVcsPJjtvRBXZxf5y4RaYcVSPZfDLj4MicoxAncrp`,
			2,
			`0x440119D0873EceAAB900c90a60968199C0117E96`,
		},
		{
			`xprv9s21ZrQH143K3oKadkPYNM5cAchBSw59qowfjCBZf5CRN6zk9q9hwtxJvWf4bDaPEjCVcsPJjtvRBXZxf5y4RaYcVSPZfDLj4MicoxAncrp`,
			4,
			`0x722C883b7E20CE4c793555c1Cab32f0CB19DE7A8`,
		},
		{
			`tprv8ZgxMBicQKsPdFznysigxtn2jpUiWX3UCXUPYpA3MnZbk7zATu6FwVozmtc4QAnUtnhqgSe6rS8HdpPPYK3yJvAMYkuuhxSSSYBMFN2ovAR`,
			2,
			`0x0186393e2BB6E3724027aD69bADF3468490e8A7C`,
		},
		{
			`tprv8ZgxMBicQKsPdFznysigxtn2jpUiWX3UCXUPYpA3MnZbk7zATu6FwVozmtc4QAnUtnhqgSe6rS8HdpPPYK3yJvAMYkuuhxSSSYBMFN2ovAR`,
			4,
			`0xD5856d3F388BBac3D4e5a71F3F2C3471Aed5Ae44`,
		},
		{
			`tpubD6NzVbkrYhZ4WoEMn2UxJq3pLiskSnYyfP4jbXXmBFvT6FHspuAxmT3GvBUjfngJwYoKRTUwESMgQtyHSB7vvvEzrnwQ3jW28PJ7PU1Vohv`,
			2,
			`0x547D75B40543Cd0a6f38dDfb5624e5E945e49dc2`,
		},
		{
			`tpubD6NzVbkrYhZ4WoEMn2UxJq3pLiskSnYyfP4jbXXmBFvT6FHspuAxmT3GvBUjfngJwYoKRTUwESMgQtyHSB7vvvEzrnwQ3jW28PJ7PU1Vohv`,
			4,
			`0xF4846cA746353Ab5D09BE1428b8381D80Ad683Ab`,
		},
		{
			`xpub661MyMwAqRbcEoyM2EtWoi2nPmquBxkn7f6vawQoDyMuCwXmDs1v74K7pbkMtXKmnVg3TogZ9XLVuV5zEbPT6naaNM2os7ZY85YgXkNnaJZ`,
			2,
			`0x735c6f21FE0DDbC2AB42b5de0ee14E98d982220F`,
		},
		{
			`xpub661MyMwAqRbcEoyM2EtWoi2nPmquBxkn7f6vawQoDyMuCwXmDs1v74K7pbkMtXKmnVg3TogZ9XLVuV5zEbPT6naaNM2os7ZY85YgXkNnaJZ`,
			4,
			`0xba16095562b7476B630B7741bc5cD3CcF9ACfB93`,
		},
	}
	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			deriver, err := NewETHDeriver(tc.key)
			assert.NoError(t, err)
			got, err := deriver.ChildAddress(tc.child)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestETHDeriver_ChildPrivate(t *testing.T) {
	cases := []struct {
		key      string
		child    uint32
		expected string
		err      bool
	}{
		{
			`xprv9s21ZrQH143K3oKadkPYNM5cAchBSw59qowfjCBZf5CRN6zk9q9hwtxJvWf4bDaPEjCVcsPJjtvRBXZxf5y4RaYcVSPZfDLj4MicoxAncrp`,
			2,
			`3d3dde2bc4e606d6e162c1ed3a3233a64727228334d159c1278038f51bbb2ed0`,
			false,
		},
		{
			`xprv9s21ZrQH143K3oKadkPYNM5cAchBSw59qowfjCBZf5CRN6zk9q9hwtxJvWf4bDaPEjCVcsPJjtvRBXZxf5y4RaYcVSPZfDLj4MicoxAncrp`,
			4,
			`6288692d913f9feb31d88efbf4b14c0ecc0c6e16d1bd24641fc45d783c7b452c`,
			false,
		},
		{
			`tprv8ZgxMBicQKsPdFznysigxtn2jpUiWX3UCXUPYpA3MnZbk7zATu6FwVozmtc4QAnUtnhqgSe6rS8HdpPPYK3yJvAMYkuuhxSSSYBMFN2ovAR`,
			2,
			`e45bc04822e15b3b8d8f11b8a2d35ab6a69429cd7ab79c08d8eba513e42a67d8`,
			false,
		},
		{
			`tprv8ZgxMBicQKsPdFznysigxtn2jpUiWX3UCXUPYpA3MnZbk7zATu6FwVozmtc4QAnUtnhqgSe6rS8HdpPPYK3yJvAMYkuuhxSSSYBMFN2ovAR`,
			4,
			`461205724190734ac4dbff2809c80f8c27d7e7fb41011cd00007c7a205f404a2`,
			false,
		},
		{
			`tpubD6NzVbkrYhZ4WoEMn2UxJq3pLiskSnYyfP4jbXXmBFvT6FHspuAxmT3GvBUjfngJwYoKRTUwESMgQtyHSB7vvvEzrnwQ3jW28PJ7PU1Vohv`,
			2,
			``,
			true,
		},
		{
			`tpubD6NzVbkrYhZ4WoEMn2UxJq3pLiskSnYyfP4jbXXmBFvT6FHspuAxmT3GvBUjfngJwYoKRTUwESMgQtyHSB7vvvEzrnwQ3jW28PJ7PU1Vohv`,
			4,
			``,
			true,
		},
		{
			`xpub661MyMwAqRbcEoyM2EtWoi2nPmquBxkn7f6vawQoDyMuCwXmDs1v74K7pbkMtXKmnVg3TogZ9XLVuV5zEbPT6naaNM2os7ZY85YgXkNnaJZ`,
			2,
			``,
			true,
		},
		{
			`xpub661MyMwAqRbcEoyM2EtWoi2nPmquBxkn7f6vawQoDyMuCwXmDs1v74K7pbkMtXKmnVg3TogZ9XLVuV5zEbPT6naaNM2os7ZY85YgXkNnaJZ`,
			4,
			``,
			true,
		},
	}
	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			deriver, err := NewETHDeriver(tc.key)
			assert.NoError(t, err)
			got, err := deriver.ChildPrivate(tc.child)
			if tc.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, got)
			}
		})
	}
}
