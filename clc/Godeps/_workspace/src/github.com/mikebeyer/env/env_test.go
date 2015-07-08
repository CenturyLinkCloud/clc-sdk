package env_test

import (
	"os"
	"testing"

	"github.com/mikebeyer/env"
)

func TestString(t *testing.T) {
	err := os.Setenv("TEST_ENV_KEY", "OK")

	if err != nil {
		t.Errorf("no error expected setting environment var, err: %s", err)
	}

	s := env.String("TEST_ENV_KEY", "NOT OK")

	if s != "OK" {
		t.Errorf("expected env to be OK, was %s", s)
	}
}

func TestStringDefault(t *testing.T) {
	s := env.String("TEST_ENV_NO_KEY", "OK")

	if s != "OK" {
		t.Errorf("expected environment variable to be set")
	}
}

func TestInt(t *testing.T) {
	err := os.Setenv("TEST_INT_KEY", "1")

	if err != nil {
		t.Errorf("no error expected setting environment var, err: %s", err)
	}

	i := env.Int("TEST_INT_KEY", 0)
	if i != 1 {
		t.Errorf("expected env to be 1, was %s", i)
	}
}

func TestIntDefault(t *testing.T) {
	i := env.Int("TEST_INT_NO_KEY", 0)

	if i != 0 {
		t.Errorf("expected default env to be 0, was %v", i)
	}
}

func TestIntFailToDefault(t *testing.T) {
	err := os.Setenv("TEST_INT_INVALID_KEY", "cake")

	if err != nil {
		t.Errorf("no error expected setting environment var, err: %s", err)
	}

	i := env.Int("TEST_INT_INVALID_KEY", 0)

	if i != 0 {
		t.Errorf("expected default env to be 0, was %v", i)
	}
}

func TestBool(t *testing.T) {
	err := os.Setenv("TEST_BOOL_KEY", "true")

	if err != nil {
		t.Errorf("no error expected setting environment var, err: %s", err)
	}

	b := env.Bool("TEST_BOOL_KEY", false)
	if !b {
		t.Errorf("expected env to be true, was %s", b)
	}
}

func TestBoolDefault(t *testing.T) {
	b := env.Bool("TEST_BOOL_NO_KEY", true)

	if !b {
		t.Errorf("expected env to be true, was %s", b)
	}
}

func TestBoolFailToDefault(t *testing.T) {
	err := os.Setenv("TEST_BOOL_INVALID_KEY", "cake")

	if err != nil {
		t.Errorf("no error expected setting environment var, err: %s", err)
	}

	b := env.Bool("TEST_BOOL_INVALID_KEY", true)
	if !b {
		t.Errorf("expected env to be true, was %s", b)
	}
}

func TestFloat(t *testing.T) {
	err := os.Setenv("TEST_FLOAT_KEY", "1.5")

	if err != nil {
		t.Errorf("no error expected setting environment var, err: %s", err)
	}

	f := env.Float("TEST_FLOAT_KEY", 2.4, 64)

	if f != 1.5 {
		t.Errorf("expected env to be 1.5, but was %v", f)
	}
}

func TestFloatDefault(t *testing.T) {
	f := env.Float("TEST_FLOAT_NO_KEY", 1.5, 64)

	if f != 1.5 {
		t.Errorf("expected env to be 1.5, but was %v", f)
	}
}

func TestFloatFailToDefault(t *testing.T) {
	err := os.Setenv("TEST_FLOAT_INVALID_KEY", "cake")

	if err != nil {
		t.Errorf("no error expected setting environment var, err: %s", err)
	}

	f := env.Float("TEST_FLOAT_INVALID_KEY", 1.5, 64)
	if f != 1.5 {
		t.Errorf("expected env to be 1.5, but was %v", f)
	}
}
