package env_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xakepp35/pkg/env"
)

func TestGet(t *testing.T) {
	want := "DEF_VALUE"
	have := env.Get("DEFINETELY_NON_EXISTING_KEY", want)
	assert.Equal(t, want, have)
}
