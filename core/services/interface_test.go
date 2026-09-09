package services

import (
	"testing"

	"github.com/QUDUSKUNLE/Bumpa/core/services/achievements"
	"github.com/stretchr/testify/assert"
)

func TestAchievementService_ImplementsService(t *testing.T) {
	var service Service = &achievements.AchievementService{}

	assert.NotNil(t, service)
}
