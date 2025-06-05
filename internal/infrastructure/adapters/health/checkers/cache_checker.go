package checkers

import (
	"fmt"

	"github.com/MarlonG1/api-facturacion-sv/config"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/health/models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
	"github.com/dimiro1/health/redis"
)

type redisChecker struct {
	redis *config.RedisConfig
}

func NewRedisChecker(redisConfig *config.RedisConfig) health.ComponentChecker {
	return &redisChecker{redisConfig}
}

func (c *redisChecker) Name() string {
	return "redis"
}

func (c *redisChecker) Check() models.Health {
	addr := fmt.Sprintf("%s:%s", c.redis.Host, c.redis.Port)
	checker := redis.NewChecker("tcp", addr)
	health := checker.Check()
	if health.IsDown() {
		details := utils.TranslateHealthDown(c.Name())
		if health.GetInfo("error") != nil {
			details = fmt.Sprintf("%s: %v", details, health.GetInfo("error"))
		}
		return models.Health{
			Status:  constants.StatusDown,
			Details: details,
		}
	}

	return models.Health{
		Status:  constants.StatusUp,
		Details: utils.TranslateHealthUp(c.Name()),
	}
}
