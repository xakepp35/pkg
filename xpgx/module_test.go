package xpgx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xakepp35/pkg/xpgx"
	"go.uber.org/fx"
)

func TestNewModule(t *testing.T) {
	assert.NoError(t, fx.New(xpgx.NewModule("")).Err())
}
