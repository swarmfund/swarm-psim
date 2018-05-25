package derive

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeriveBTCFamilyChildAddress(t *testing.T) {
	cases := []struct {
		network  NetworkType
		key      string
		child    uint32
		expected string
	}{
		{
			NetworkTypeBTCMainnet,
			`xpub661MyMwAqRbcEoyM2EtWoi2nPmquBxkn7f6vawQoDyMuCwXmDs1v74K7pbkMtXKmnVg3TogZ9XLVuV5zEbPT6naaNM2os7ZY85YgXkNnaJZ`,
			2,
			`15bKLQSDPDe8axyarYD4G1JeqSU5QQjkvy`,
		},
		{
			NetworkTypeBTCMainnet,
			`xpub661MyMwAqRbcEoyM2EtWoi2nPmquBxkn7f6vawQoDyMuCwXmDs1v74K7pbkMtXKmnVg3TogZ9XLVuV5zEbPT6naaNM2os7ZY85YgXkNnaJZ`,
			4,
			`16rpQdsLmMh3AUFN8JQMnvedaMpTMb9m5j`,
		},
		{
			NetworkTypeBTCTestnet,
			`tpubD6NzVbkrYhZ4WoEMn2UxJq3pLiskSnYyfP4jbXXmBFvT6FHspuAxmT3GvBUjfngJwYoKRTUwESMgQtyHSB7vvvEzrnwQ3jW28PJ7PU1Vohv`,
			2,
			`n1qP1TCAhHbrLk6wzgcCV2qukHGgUfVDf8`,
		},
		{
			NetworkTypeBTCTestnet,
			`tpubD6NzVbkrYhZ4WoEMn2UxJq3pLiskSnYyfP4jbXXmBFvT6FHspuAxmT3GvBUjfngJwYoKRTUwESMgQtyHSB7vvvEzrnwQ3jW28PJ7PU1Vohv`,
			4,
			`mze8kz7ikyxmpcpsAYPpBwtjKy8XEzudLP`,
		},
		{
			NetworkTypeDashMainnet,
			`xpub661MyMwAqRbcEoyM2EtWoi2nPmquBxkn7f6vawQoDyMuCwXmDs1v74K7pbkMtXKmnVg3TogZ9XLVuV5zEbPT6naaNM2os7ZY85YgXkNnaJZ`,
			2,
			`XfHAAf67LvrijuaAiRXH7XzSfn3mMnxAkC`,
		},
		{
			NetworkTypeDashMainnet,
			`xpub661MyMwAqRbcEoyM2EtWoi2nPmquBxkn7f6vawQoDyMuCwXmDs1v74K7pbkMtXKmnVg3TogZ9XLVuV5zEbPT6naaNM2os7ZY85YgXkNnaJZ`,
			4,
			`XgYfEtXEj4udKQqwzBiaeTLRQhQ9H8CoVM`,
		},
		{
			NetworkTypeDashTestnet,
			`tpubD6NzVbkrYhZ4WoEMn2UxJq3pLiskSnYyfP4jbXXmBFvT6FHspuAxmT3GvBUjfngJwYoKRTUwESMgQtyHSB7vvvEzrnwQ3jW28PJ7PU1Vohv`,
			2,
			`ygdsZbqXHX3G4K9ThrGSYfjizuk37sFeS4`,
		},
		{
			NetworkTypeDashTestnet,
			`tpubD6NzVbkrYhZ4WoEMn2UxJq3pLiskSnYyfP4jbXXmBFvT6FHspuAxmT3GvBUjfngJwYoKRTUwESMgQtyHSB7vvvEzrnwQ3jW28PJ7PU1Vohv`,
			4,
			`yfSdK8m5MDQBYBsNsi44FanYabbsrfUFhs`,
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%d-%d", tc.network, tc.child), func(t *testing.T) {
			deriver, err := NewBTCFamilyDeriver(tc.network, tc.key)
			assert.NoError(t, err)
			got, err := deriver.ChildAddress(tc.child)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestDeriveBTCFamilyChildPrivate(t *testing.T) {
	cases := []struct {
		network  NetworkType
		key      string
		child    uint32
		expected string
	}{
		{
			NetworkTypeBTCMainnet,
			`xprv9s21ZrQH143K3oKadkPYNM5cAchBSw59qowfjCBZf5CRN6zk9q9hwtxJvWf4bDaPEjCVcsPJjtvRBXZxf5y4RaYcVSPZfDLj4MicoxAncrp`,
			2,
			`KyGkqkV7nM5iENHFmgRbcsGksSr5snG6nL7zssXc4qzwcnX9CLdA`,
		},
		{
			NetworkTypeBTCMainnet,
			`xprv9s21ZrQH143K3oKadkPYNM5cAchBSw59qowfjCBZf5CRN6zk9q9hwtxJvWf4bDaPEjCVcsPJjtvRBXZxf5y4RaYcVSPZfDLj4MicoxAncrp`,
			4,
			`KzXFDNunNeTjrdNGU6yoPPYw1hEXpE5CPdyEPhvzRoXXznFTAg6Q`,
		},
		{
			NetworkTypeBTCTestnet,
			`tprv8ZgxMBicQKsPdFznysigxtn2jpUiWX3UCXUPYpA3MnZbk7zATu6FwVozmtc4QAnUtnhqgSe6rS8HdpPPYK3yJvAMYkuuhxSSSYBMFN2ovAR`,
			2,
			`cVEbnTCHDGhKdcAcCWfuF4Tuh3zNvG8pBiFuaz1A3spaDaBcmw5N`,
		}, {
			NetworkTypeBTCTestnet,
			`tprv8ZgxMBicQKsPdFznysigxtn2jpUiWX3UCXUPYpA3MnZbk7zATu6FwVozmtc4QAnUtnhqgSe6rS8HdpPPYK3yJvAMYkuuhxSSSYBMFN2ovAR`,
			4,
			`cPvuhFT6QBSChhwyr1QshhmYHXzXSKa3XFLZuVn8AYh7S5ERK2oL`,
		},
		{
			NetworkTypeDashMainnet,
			`xprv9s21ZrQH143K3oKadkPYNM5cAchBSw59qowfjCBZf5CRN6zk9q9hwtxJvWf4bDaPEjCVcsPJjtvRBXZxf5y4RaYcVSPZfDLj4MicoxAncrp`,
			2,
			`XDLgJ1sV62iAHhHdoSRU86TmnU7fKasM9uTuQProPDH31wdsPx7a`,
		},
		{
			NetworkTypeDashMainnet,
			`xprv9s21ZrQH143K3oKadkPYNM5cAchBSw59qowfjCBZf5CRN6zk9q9hwtxJvWf4bDaPEjCVcsPJjtvRBXZxf5y4RaYcVSPZfDLj4MicoxAncrp`,
			4,
			`XEbAfeJ9gL6BuxNeVryftcjwviW7G2gSmDK8vEGBkAodPwN6tKzD`,
		},
		{
			NetworkTypeDashTestnet,
			`tprv8ZgxMBicQKsPdFznysigxtn2jpUiWX3UCXUPYpA3MnZbk7zATu6FwVozmtc4QAnUtnhqgSe6rS8HdpPPYK3yJvAMYkuuhxSSSYBMFN2ovAR`,
			2,
			`cVEbnTCHDGhKdcAcCWfuF4Tuh3zNvG8pBiFuaz1A3spaDaBcmw5N`,
		}, {
			NetworkTypeDashTestnet,
			`tprv8ZgxMBicQKsPdFznysigxtn2jpUiWX3UCXUPYpA3MnZbk7zATu6FwVozmtc4QAnUtnhqgSe6rS8HdpPPYK3yJvAMYkuuhxSSSYBMFN2ovAR`,
			4,
			`cPvuhFT6QBSChhwyr1QshhmYHXzXSKa3XFLZuVn8AYh7S5ERK2oL`,
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("%d-%d", tc.network, tc.child), func(t *testing.T) {
			deriver, err := NewBTCFamilyDeriver(tc.network, tc.key)
			assert.NoError(t, err)
			got, err := deriver.ChildPrivate(tc.child)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}
