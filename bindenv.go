package bindenv

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type Binder interface {
	Get(key string) string
	Set(key, value string) error
	Sets(envs map[string]string) error

	Bind(raw any) error
	SetDefault(key, value string)
}

type env struct {
	mx   *sync.RWMutex
	maps map[string]string
}

func New() Binder {
	return &env{
		mx:   new(sync.RWMutex),
		maps: map[string]string{},
	}
}

func (e *env) Set(key, value string) error {
	e.mx.Lock()
	defer e.mx.Unlock()

	return os.Setenv(key, value)
}

func (e *env) Sets(envs map[string]string) (err error) {
	e.mx.Lock()
	defer e.mx.Unlock()

	for key, value := range envs {
		if err = os.Setenv(key, value); err != nil {
			return
		}
	}
	return
}

func (e *env) Get(key string) string {
	e.mx.RLock()
	defer e.mx.RUnlock()

	return os.Getenv(key)
}

func (e *env) SetDefault(key, value string) {
	e.mx.Lock()
	e.maps[key] = value
	e.mx.Unlock()
}

func (e *env) Bind(raw any) error {
	e.mx.Lock()
	defer e.mx.Unlock()

	to := reflect.TypeOf(raw)
	if raw == nil || to.Kind() != reflect.Ptr {
		return errors.New("data type sent must be pointer struct")
	}

	to = to.Elem()
	if to.Kind() != reflect.Struct {
		return errors.New("can only process struct data type")
	}

	return e.parseStruct(raw)
}

func (e *env) parseStruct(raw any) (err error) {
	to := reflect.TypeOf(raw).Elem()

	vo := reflect.ValueOf(raw)
	vo = reflect.Indirect(vo)

	for i := 0; i < vo.NumField(); i++ {
		if vo.Field(i).Type().Kind() == reflect.Struct {
			e.parseStruct(vo.Field(i).Addr().Interface())
			continue
		}

		tagEnv := to.Field(i).Tag.Get("env")
		if tagEnv == "" {
			tagEnv = strings.ToUpper(to.Field(i).Name)
		}

		envValue := strings.TrimSpace(os.Getenv(tagEnv))
		if envValue == "" {
			envValue = e.maps[tagEnv]
		}

		var res interface{}

		switch vo.Field(i).Type().Kind() {
		default:
			return errors.New("data type is not supported")

		case reflect.String:
			vo.Field(i).SetString(envValue)

		case reflect.Bool:
			if envValue == "" {
				vo.Field(i).SetBool(false)

			} else {
				res, err = strconv.ParseBool(envValue)
				vo.Field(i).SetBool(res.(bool))
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if envValue == "" {
				vo.Field(i).SetUint(0)

			} else {
				res, err = strconv.ParseUint(envValue, 10, 64)
				vo.Field(i).SetUint(res.(uint64))
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if envValue == "" {
				vo.Field(i).SetInt(0)

			} else {
				res, err = strconv.ParseInt(envValue, 10, 64)
				vo.Field(i).SetInt(res.(int64))
			}

		case reflect.Float32, reflect.Float64:
			if envValue == "" {
				vo.Field(i).SetFloat(0)

			} else {
				res, err = strconv.ParseFloat(envValue, 64)
				vo.Field(i).SetFloat(res.(float64))
			}
		}
	}
	return
}
