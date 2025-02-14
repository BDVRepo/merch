package handlers

import (
	"testing"

	"bdv-avito-merch/libs/4_common/env_vars"
)

func TestMain(m *testing.M) {
	env_vars.LoadEnvVars()

	m.Run()
}
