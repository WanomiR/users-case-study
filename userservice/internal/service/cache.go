package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/wanomir/e"
	"log"
	"userservice/internal/entity"
)

type UserServiceCacheProxy struct {
	userService *UserService
	cache       *redis.Pool
}

func NewUserServiceCacheProxy(userService *UserService, redisAddress string) *UserServiceCacheProxy {
	redisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisAddress)
		},
	}
	return &UserServiceCacheProxy{
		userService: userService,
		cache:       redisPool,
	}
}

func (p UserServiceCacheProxy) GetUsers(ctx context.Context, offset, limit int) (users []entity.User, err error) {
	conn := p.cache.Get()
	defer conn.Close()

	key := fmt.Sprintf("%d:%d", offset, limit)

	cachedData, err := redis.Bytes(conn.Do("GET", key))
	if err == nil {
		if err = json.Unmarshal(cachedData, &users); err != nil {
			return []entity.User{}, e.Wrap("error unmarshalling users from cache", err)
		}
		log.Println("loaded users from cache")
		return users, nil
	}

	if users, err = p.userService.GetUsers(ctx, offset, limit); err != nil {
		return []entity.User{}, e.Wrap("error getting users", err)
	}

	serialized, _ := json.Marshal(users)
	if _, err = conn.Do("SETEX", key, 60*15, serialized); err != nil {
		log.Println("error caching users", err)
	}
	log.Println("cached users")

	return users, nil
}

func (p UserServiceCacheProxy) GetUser(ctx context.Context, email string) (user entity.User, err error) {
	conn := p.cache.Get()
	defer conn.Close()

	cachedData, err := redis.Bytes(conn.Do("GET", email))
	if err == nil {
		if err = json.Unmarshal(cachedData, &user); err != nil {
			return entity.User{}, e.Wrap("error unmarshalling user from cache", err)
		}
		log.Println("loaded user " + email + " from cache")
		return user, nil
	}

	if user, err = p.userService.GetUser(ctx, email); err != nil {
		return entity.User{}, e.Wrap("error getting user", err)
	}

	serialized, _ := json.Marshal(user)
	if _, err = conn.Do("SETEX", email, 60*15, serialized); err != nil {
		log.Println("error caching user", err)
	}
	log.Println("cached user " + email)

	return user, nil
}

func (p UserServiceCacheProxy) CreateUser(ctx context.Context, user entity.User) (int, error) {
	return p.userService.CreateUser(ctx, user)
}
