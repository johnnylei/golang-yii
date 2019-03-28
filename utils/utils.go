package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func IsStringOrNumber(data interface{}) bool {
	if _, ok := data.(string); ok {
		return true
	}

	if _, ok := data.(int); ok {
		return true
	}

	if _, ok := data.(float32); ok {
		return true
	}

	if _, ok := data.(float64); ok {
		return true
	}

	return false
}

func NumberArrayToString(data interface{}, delimiter string) string  {
	if _, ok := data.([]int); ok {
		return strings.Trim(strings.Join(strings.Split(fmt.Sprint(data), " "), delimiter), "[]")
	}

	if _, ok := data.([]float32); ok {
		return strings.Trim(strings.Join(strings.Split(fmt.Sprint(data), " "), delimiter), "[]")
	}

	if _, ok := data.([]float64); ok {
		return strings.Trim(strings.Join(strings.Split(fmt.Sprint(data), " "), delimiter), "[]")
	}

	panic("invalid data")
}

func DeepCopy(value interface{}) interface{} {
	if valueMap, ok := value.(map[string]interface{}); ok {
		newMap := make(map[string]interface{})
		for k, v := range valueMap {
			newMap[k] = DeepCopy(v)
		}

		return newMap
	} else if valueSlice, ok := value.([]interface{}); ok {
		newSlice := make([]interface{}, len(valueSlice))
		for k, v := range valueSlice {
			newSlice[k] = DeepCopy(v)
		}

		return newSlice
	}

	return value
}

func ConvertInt64(data interface{}) int64  {
	ret, ok := data.(int64)
	if ok {
		return ret
	}

	_ret, ok := data.(string)
	if !ok {
		panic("invalid data")
	}

	ret, err := strconv.ParseInt(_ret, 10, 64)
	if err != nil {
		panic(err)
	}

	return ret
}

func ConvertInt(data interface{}) int  {
	ret, ok := data.(int)
	if ok {
		return ret
	}

	_ret, ok := data.(string)
	if !ok {
		panic("invalid data")
	}

	ret, err := strconv.Atoi(_ret)
	if err != nil {
		panic(err)
	}

	return ret
}

func ConvertFloat64(data interface{}) float64  {
	ret, ok := data.(float64)
	if ok {
		return ret
	}

	_ret, ok := data.(string)
	if !ok {
		panic("invalid data")
	}

	ret, err := strconv.ParseFloat(_ret, 64)
	if err != nil {
		panic(err)
	}

	return ret
}