package conf

import (
	"bytes"
	"testing"

	"gitlab.com/distributed_lab/discovery-go"
	"github.com/spf13/viper"
)

func TestViperConfig_Discovery(t *testing.T) {
	def := func(c *discovery.ClientConfig) *discovery.ClientConfig {
		return c.WithDefaults()
	}

	cases := []struct {
		name, raw string
		err       bool
		exp       *discovery.ClientConfig
	}{
		{"default", "", false, def(&discovery.ClientConfig{})},
		{"host", `discovery: {host: foobar}`, false, def(&discovery.ClientConfig{Host: "foobar"})},
		{"port", `discovery: {port: 1234}`, false, def(&discovery.ClientConfig{Port: 1234})},
		// TODO Turn this off back, when vipers getting of value will be fixed
		//{"invalid port", `discovery: {port: not-a-port}`, true, def(&discovery.ClientConfig{})},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				discoveryClient = nil
			}()
			r := bytes.NewReader([]byte(tc.raw))
			v := viper.New()
			v.SetConfigType("yaml")
			err := v.ReadConfig(r)
			if err != nil {
				t.Fatal(err)
			}
			config := ViperConfig{
				viper: v,
			}
			discovery, err := config.Discovery()
			if tc.err && err == nil {
				t.Fatal("expected to throw error")
			}
			if !tc.err && err != nil {
				t.Fatalf("got %v expected nil", err)
			}
			if err != nil {
				return
			}
			gotConfig := discovery.Config()
			if *gotConfig != *tc.exp {
				t.Fatalf("got %v expected %v", gotConfig, tc.exp)

			}
		})
	}

}
