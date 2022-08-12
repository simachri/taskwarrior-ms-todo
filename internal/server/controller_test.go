package server

import (
	"testing"

	"github.com/simachri/taskwarrior-ms-todo/internal/taskwarrior"
	"github.com/simachri/taskwarrior-ms-todo/internal/test"
	"github.com/stretchr/testify/assert"
)


func TestCheckHealth_udasExist_isOK(t *testing.T) {
    test.NewTaskwarriorEnv(t)
    err := taskwarrior.CreateIntegrationUDAs()
    assert.NoError(t, err)

    err = checkHealth()
    assert.NoError(t, err)
}

func TestCheckHealth_udasMissing_isError(t *testing.T) {
    test.NewTaskwarriorEnv(t)

    err := checkHealth()
    assert.Error(t, err)
}

