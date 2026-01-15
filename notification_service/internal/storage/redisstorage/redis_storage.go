package redisstorage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"notification_service/internal/models"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
	ttl    time.Duration
}

const (
	reminderKeyPrefix = "reminders:"
	defaultTTL        = 31 * 24 * time.Hour
)

func NewRedisStorage(host string, port int, password string, db int, ttlDays int) (*RedisStorage, error) {
	addr := fmt.Sprintf("%s:%d", host, port)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection error: %w", err)
	}

	ttl := time.Duration(ttlDays) * 24 * time.Hour
	if ttl == 0 {
		ttl = defaultTTL
	}

	return &RedisStorage{client: client, ttl: ttl}, nil
}

func (s *RedisStorage) SaveReminder(ctx context.Context, reminder *models.Reminder) error {
	data, err := json.Marshal(reminder)
	if err != nil {
		return err
	}

	key := reminder.Key()
	return s.client.Set(ctx, key, data, s.ttl).Err()
}

func (s *RedisStorage) GetAllReminders(ctx context.Context) ([]*models.Reminder, error) {
	keys, err := s.client.Keys(ctx, reminderKeyPrefix+"*").Result()
	if err != nil {
		return nil, err
	}

	var reminders []*models.Reminder

	for _, key := range keys {
		data, err := s.client.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}

		var reminder models.Reminder
		if err := json.Unmarshal(data, &reminder); err != nil {
			continue
		}

		reminders = append(reminders, &reminder)
	}

	return reminders, nil
}

func (s *RedisStorage) DeleteAllReminders(ctx context.Context) error {
	keys, err := s.client.Keys(ctx, reminderKeyPrefix+"*").Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return s.client.Del(ctx, keys...).Err()
	}

	return nil
}

func (s *RedisStorage) DeleteRemindersByPerson(ctx context.Context, personID int64) error {
	pattern := fmt.Sprintf("%s*:%d:*", reminderKeyPrefix, personID)

	keys, err := s.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return s.client.Del(ctx, keys...).Err()
	}

	return nil
}

func (s *RedisStorage) Close() error {
	return s.client.Close()
}
