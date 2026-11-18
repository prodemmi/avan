package redis

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type Cache[T any] struct {
	Redis
}

// Set sets a value in Redis with a custom expiration time.
func (r Cache[T]) Set(value interface{}, expiration time.Duration, keys ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	key := strings.Join(keys, ":")
	return r.client.Set(ctx, key, value, expiration).Err()
}

// SetWithMinutes sets a value in Redis with an expiration time in minutes.
func (r Cache[T]) SetWithMinutes(value interface{}, minutes int, keys ...string) error {
	return r.Set(value, time.Duration(minutes)*time.Minute, keys...)
}

// SetWithHours sets a value in Redis with an expiration time in hours.
func (r Cache[T]) SetWithHours(value interface{}, hours int, keys ...string) error {
	return r.Set(value, time.Duration(hours)*time.Hour, keys...)
}

// SetWithDays sets a value in Redis with an expiration time in days.
func (r Cache[T]) SetWithDays(value interface{}, days int, keys ...string) error {
	return r.Set(value, time.Duration(days)*24*time.Hour, keys...)
}

// Predefined time-based setters for common durations.
func (r Cache[T]) Set5Min(value interface{}, keys ...string) error {
	return r.Set(value, 5*time.Minute, keys...)
}

func (r Cache[T]) Set10Min(value interface{}, keys ...string) error {
	return r.Set(value, 10*time.Minute, keys...)
}

func (r Cache[T]) Set30Min(value interface{}, keys ...string) error {
	return r.Set(value, 30*time.Minute, keys...)
}

func (r Cache[T]) Set1Hour(value interface{}, keys ...string) error {
	return r.Set(value, 1*time.Hour, keys...)
}

func (r Cache[T]) Set10Hour(value interface{}, keys ...string) error {
	return r.Set(value, 10*time.Hour, keys...)
}

func (r Cache[T]) Set1Day(value interface{}, keys ...string) error {
	return r.Set(value, 24*time.Hour, keys...)
}

func (r Cache[T]) Set7Days(value interface{}, keys ...string) error {
	return r.Set(value, 7*24*time.Hour, keys...)
}

// SetNoExpire sets a value in Redis without an expiration time.
func (r Cache[T]) SetNoExpire(value interface{}, keys ...string) error {
	return r.Set(value, 0, keys...)
}

func (r Cache[T]) Get(keys ...string) (string, error) {
	ctx := context.Background()
	key := strings.Join(keys, ":")
	return r.client.Get(ctx, key).Result()
}

// GetOrSet checks if the key exists; if not, it sets the value from the callback
func (r Cache[T]) GetOrSet(keys []interface{}, expiration time.Duration, callback func() (T, error)) (T, error) {
	ctx := context.Background()
	key := strings.Join(stringifyKeys(keys...), ":")

	// Try to get the value from cache
	val, err := r.client.Get(ctx, key).Result()
	if err == nil {
		var result T
		// Convert string value to the generic type
		if parsed, ok := any(val).(T); ok {
			result = parsed
			return result, nil
		}
	}

	// If not found, call the callback to get the value
	newVal, err := callback()
	if err != nil {
		var zeroVal T
		return zeroVal, err
	}

	// Save the new value to cache (convert to string for Redis)
	if err := r.client.Set(ctx, key, fmt.Sprintf("%+v", newVal), expiration).Err(); err != nil {
		var zeroVal T
		return zeroVal, err
	}

	return newVal, nil
}

// GetOrSetArray checks if the key exists; if not, it sets the array value from the callback
func (r Cache[T]) GetOrSetArray(keys []interface{}, expiration time.Duration, callback func() ([]T, error)) ([]T, error) {
	ctx := context.Background()
	key := strings.Join(stringifyKeys(keys...), ":")

	// Try to get the value from cache
	val, err := r.client.Get(ctx, key).Result()
	if err == nil {
		var result []T
		// Convert string value to the generic type
		if parsed, ok := any(val).([]T); ok {
			result = parsed
			return result, nil
		}
	}

	// If not found, call the callback to get the value
	newVal, err := callback()
	if err != nil {
		var zeroVal []T
		return zeroVal, err
	}

	// Save the new value to cache (convert to string for Redis)
	if err := r.client.Set(ctx, key, fmt.Sprintf("%v", newVal), expiration).Err(); err != nil {
		var zeroVal []T
		return zeroVal, err
	}

	return newVal, nil
}

func (r Cache[T]) GetWithFallback(fallback func() (string, time.Duration), keys ...string) (string, error) {
	value, err := r.Get(keys...)
	if err == nil {
		return value, nil
	}

	// If not found, use fallback function
	if fallback != nil {
		newValue, expiration := fallback()
		_ = r.Set(newValue, expiration, keys...)
		return newValue, nil
	}

	return "", err
}

func (r Cache[T]) Exists(keys ...string) (bool, error) {
	ctx := context.Background()
	key := strings.Join(keys, ":")
	count, err := r.client.Exists(ctx, key).Result()
	return count > 0, err
}

func (r Cache[T]) Delete(keys ...string) error {
	ctx := context.Background()
	key := strings.Join(keys, ":")
	return r.client.Del(ctx, key).Err()
}

func (r Cache[T]) GetKey(keys ...interface{}) []interface{} {
	return keys
}

func (r Cache[T]) GetStringKey(keys ...interface{}) []string {
	return stringifyKeys(keys)
}

// stringifyKeys converts any type of key to string
func stringifyKeys(keys ...interface{}) []string {
	var strKeys []string
	for _, key := range keys {
		if key == nil || key == "" {
			key = "nil"
		}

		strKeys = append(strKeys, fmt.Sprintf("%v", key))
	}
	return strKeys
}

func HashKey(token string) string {
	hash := sha256.Sum256([]byte(token))
	hashed := hex.EncodeToString(hash[:])
	if len(hashed) > 12 {
		hashed = hashed[:12]
	}
	return hashed
}

func NewCache[T any](client Redis) Cache[T] {
	return Cache[T]{client}
}
