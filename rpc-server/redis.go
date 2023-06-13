package main

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	cli *redis.Client
}

func (client *RedisClient) InitClient(ctx context.Context, address string, password string) error {
	rdb := redis.NewClient(&redis.Options{
		Addr: address,
		Password: password,
		DB: 0,
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		return err
	}

	client.cli = rdb
	return nil
}

type Message struct {
	Sender string `json:"sender"`
	Message string `json:"message"`
	Timestamp int64 `json:"timestamp"`
}

func (client *RedisClient) SaveMessage(ctx context.Context, roomID string, message *Message) (error) {
	text, err := json.Marshal(message)
	if err != nil {
		return err
	}

	member := &redis.Z{
		Score: float64(message.Timestamp),
		Member: text,
	}

	_, err = client.cli.ZAdd(ctx, roomID, *member).Result()
	if err != nil {
		return err
	}

	return nil
}

func (client *RedisClient) GetMessagesByRoomID(ctx context.Context, roomID string, start int64, end int64, reverse bool) ([]*Message, error){
	var (
		rawMessages []string
		messages []*Message
		err error
	)

	if reverse {
		//Get all messages from a chat room - Latest messages first
		rawMessages, err = client.cli.ZRevRange(ctx, roomID, start, end).Result()
		if err != nil {
			return nil, err
		}
	} else {
		//Get all messages from a chat room - Latest messages last
		rawMessages, err = client.cli.ZRange(ctx, roomID, start, end).Result()
		if err != nil {
			return nil, err
		}
	}

	for _, msg := range rawMessages {
		temp := &Message{}
		err := json.Unmarshal([]byte(msg), temp)
		if err != nil {
			return nil, err
		}
		messages = append(messages, temp)
	}

	return messages, nil
}

