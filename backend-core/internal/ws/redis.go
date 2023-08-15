package ws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Wave-95/boards/backend-core/internal/config"
	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg config.RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%v:%v", cfg.Host, cfg.Port),
	})
	return rdb
}

// setUser sets a user into the redis hash store organized by board ID. This hash store is used to
// manage the list of connected users.
func setUser(rdb *redis.Client, boardID string, user models.User) error {
	userBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	_, err = rdb.HSet(context.Background(), boardID, user.ID.String(), userBytes).Result()
	if err != nil {
		return err
	}

	return nil
}

// getUsers returns a map of all the connected users for a board
func getUsers(rdb *redis.Client, boardID string) (map[string]string, error) {
	if res, err := rdb.HGetAll(context.Background(), boardID).Result(); err != nil {
		return map[string]string{}, err
	} else {
		fmt.Println(res)
		return res, nil
	}
}

// delUser deletes a user from the redis hash store.
func delUser(rdb *redis.Client, boardID string, userID string) error {
	_, err := rdb.HDel(context.Background(), boardID, userID).Result()
	if err != nil {
		return err
	}

	return nil
}
