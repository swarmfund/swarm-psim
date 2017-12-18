package figure_test

import (
	"reflect"
	"testing"

	"gitlab.com/distributed_lab/figure"
)

func TestSimpleUsage(t *testing.T) {
	type Config struct {
		SomeInt    int
		SomeString string
		Missing    int
		Default    int
	}
	c := Config{
		Default: 42,
	}
	err := figure.Out(&c).From(map[string]interface{}{
		"some_int":    1,
		"some_string": "satoshi",
	}).Please()
	if err != nil {
		t.Fatalf("expected nil error got %s", err)
	}
	reference := Config{
		SomeInt:    1,
		SomeString: "satoshi",
		Default:    42,
	}
	if !reflect.DeepEqual(c, reference) {
		t.Errorf("expected %#v got %#v", reference, c)
	}

}
