package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"password-management-service/internal/models/user"
)

type RedisService interface {
	SaveData(key, clientID string, data interface{}) error
	GetData(key, clientID string, target interface{}) error
	DeleteData(key, clientID string) error
	GetToken(clientID string) (string, error)
	DeleteToken(clientID string) error
}

type redisService struct {
	Client redis.Client
	Ctx    context.Context
}

func NewRedisService(client redis.Client) RedisService {
	return redisService{Client: client, Ctx: context.Background()}
}

func (r redisService) SaveData(key, clientID string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}
	return r.Client.Set(r.Ctx, key+":"+clientID, jsonData, 0).Err()
}

func (r redisService) GetData(key, clientID string, target interface{}) error {
	jsonData, err := r.Client.Get(r.Ctx, key+":"+clientID).Result()
	if errors.Is(err, redis.Nil) {
		return fmt.Errorf("no data found for key: %s", key+":"+clientID)
	} else if err != nil {
		return fmt.Errorf("failed to get data: %v", err)
	}
	return json.Unmarshal([]byte(jsonData), target)
}

func (r redisService) DeleteData(key, clientID string) error {
	return r.Client.Del(r.Ctx, key+":"+clientID).Err()
}

func generateRedisKey(clientID string) string {
	return "token:" + clientID
}

func (r redisService) GetToken(clientID string) (string, error) {
	token, err := r.Client.Get(r.Ctx, generateRedisKey(clientID)).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return token, err
}

func (r redisService) DeleteToken(clientID string) error {
	return r.Client.Del(r.Ctx, generateRedisKey(clientID)).Err()
}

func GetUserRedis(redis RedisService, key, clientID string) (*user.UserRedis, error) {
	var u user.UserRedis
	if err := redis.GetData(key, clientID, &u); err != nil {
		return nil, err
	}
	return &u, nil
}
