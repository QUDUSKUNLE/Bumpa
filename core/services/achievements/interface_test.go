package achievements

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAchievementService_ImplementsService(t *testing.T) {
	var service Service = &AchievementService{}

	assert.NotNil(t, service)
}
