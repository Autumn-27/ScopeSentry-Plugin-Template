// Code generated by 'yaegi extract github.com/Autumn-27/ScopeSentry-Scan/internal/redis'. DO NOT EDIT.

package symbols

import (
	"github.com/Autumn-27/ScopeSentry-Scan/internal/redis"
	"reflect"
)

func init() {
	Symbols["github.com/Autumn-27/ScopeSentry-Scan/internal/redis/redis"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"Initialize":  reflect.ValueOf(redis.Initialize),
		"RedisClient": reflect.ValueOf(&redis.RedisClient).Elem(),

		// type definitions
		"Client": reflect.ValueOf((*redis.Client)(nil)),
	}
}
