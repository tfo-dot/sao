package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Shopify/go-lua"
	"github.com/google/uuid"
)

func GetLuaString(state *lua.State, str string) string {
	state.Global(str)

	value, ok := state.ToString(-1)

	if !ok {
		panic("Cannot convert to string")
	}

	state.Pop(1)
	return value
}

func GetLuaBool(state *lua.State, str string) bool {
	state.Global(str)

	value := state.ToBoolean(-1)

	state.Pop(1)
	return value
}

func GetLuaFloat(state *lua.State, str string) float64 {
	state.Global(str)

	value, ok := state.ToNumber(-1)

	if !ok {
		panic("Cannot convert to float")
	}

	state.Pop(1)
	return value
}

func GetLuaInt(state *lua.State, str string) int {
	state.Global(str)

	value, ok := state.ToInteger(-1)

	if !ok {
		panic("Cannot convert to int")
	}

	state.Pop(1)
	return value
}

func GetNullableString(state *lua.State, str string) (string, bool) {
	state.Global(str)

	value, ok := state.ToString(-1)

	if !ok {
		state.Pop(1)
		return "", false
	}

	state.Pop(1)
	return value, true
}

func GetTableAsMap(state *lua.State) (map[string]interface{}, error) {
	if !state.IsTable(-1) {
		return nil, errors.New("expected table at -1")
	}

	returnValue := make(map[string]interface{})

	state.PushNil()

	for state.Next(-2) {
		key, ok := state.ToString(-2)

		if !ok {
			state.Pop(2)
			return nil, fmt.Errorf("only strings are valid keys, got %d", state.TypeOf(-2))
		}

		valueType := state.TypeOf(-1)

		switch valueType {
		case lua.TypeNil:
			returnValue[key] = nil
			state.Pop(1)
		case lua.TypeBoolean:
			returnValue[key] = state.ToBoolean(-1)
			state.Pop(1)
		case lua.TypeNumber:
			tempValue, _ := state.ToNumber(-1)

			returnValue[key] = tempValue
			state.Pop(1)
		case lua.TypeString:
			tempValue, _ := state.ToString(-1)

			returnValue[key] = tempValue
			state.Pop(1)
		case lua.TypeTable:
			if !state.IsTable(-1) {
				return nil, errors.New("expected table at -1")
			}

			isValid, err := IsLuaArray(state)

			if err != nil {
				return nil, errors.Join(fmt.Errorf("error at '%s'", key), err)
			}

			if isValid {
				tempValue, err := GetTableAsArray(state)

				if err != nil {
					return nil, errors.Join(fmt.Errorf("error at '%s'", key), err)
				}

				returnValue[key] = tempValue
			} else {
				tempValue, err := GetTableAsMap(state)

				if err != nil {
					return nil, errors.Join(fmt.Errorf("error at '%s'", key), err)
				}

				returnValue[key] = tempValue

			}
		case lua.TypeFunction:
			fUuid := strings.Replace(uuid.New().String(), "-", "_", -1)
			state.SetGlobal(fUuid)

			returnValue[key] = LuaFunctionRef{FunctionName: fUuid}
		default:
			return nil, fmt.Errorf("unsupported type %d", valueType)
		}
	}

	state.Pop(1)

	return returnValue, nil
}

func GetTableAsArray(state *lua.State) ([]interface{}, error) {
	if !state.IsTable(-1) {
		return nil, errors.New("expected table at -1")
	}

	isValidArray, err := IsLuaArray(state)

	if err != nil {
		return nil, err
	}

	if !isValidArray {
		return []interface{}{}, nil
	}

	returnValue := make([]interface{}, 0)

	eltIdx := 1

	for {
		state.PushInteger(eltIdx)

		state.RawGet(-2)

		if state.IsNil(-1) {
			state.Pop(1)
			break
		}

		valueType := state.TypeOf(-1)

		switch valueType {
		case lua.TypeNil:
			returnValue = append(returnValue, nil)
			state.Pop(1)
		case lua.TypeBoolean:
			returnValue = append(returnValue, state.ToBoolean(-1))
			state.Pop(1)
		case lua.TypeNumber:
			tempValue, _ := state.ToNumber(-1)

			returnValue = append(returnValue, tempValue)
			state.Pop(1)
		case lua.TypeString:
			tempValue, _ := state.ToString(-1)

			returnValue = append(returnValue, tempValue)
			state.Pop(1)
		case lua.TypeTable:
			if !state.IsTable(-1) {
				return nil, errors.New("expected table at -1")
			}

			isValid, err := IsLuaArray(state)

			if err != nil {
				return nil, errors.Join(fmt.Errorf("error at %d", eltIdx), err)
			}

			if isValid {
				tempValue, err := GetTableAsArray(state)

				if err != nil {
					return nil, errors.Join(fmt.Errorf("error at %d", eltIdx), err)
				}

				returnValue = append(returnValue, tempValue)
			} else {
				tempValue, err := GetTableAsMap(state)

				if err != nil {
					return nil, errors.Join(fmt.Errorf("error at %d", eltIdx), err)
				}

				returnValue = append(returnValue, tempValue)
			}

		case lua.TypeFunction:
			fUuid := strings.Replace(uuid.New().String(), "-", "_", -1)
			state.SetGlobal(fUuid)

			returnValue = append(returnValue, LuaFunctionRef{FunctionName: fUuid})
		default:
			return nil, fmt.Errorf("unsupported type %d", valueType)
		}

		eltIdx++
	}

	state.Pop(1)

	return returnValue, nil
}

type LuaFunctionRef struct {
	FunctionName string
}

func IsLuaArray(state *lua.State) (bool, error) {
	if !state.IsTable(-1) {
		return false, errors.New("expected table at -1")
	}

	state.Length(-1)

	objLen, _ := state.ToInteger(-1)

	state.Pop(1)

	count := 0

	state.PushNil()

	for state.Next(-2) {
		keyType := state.TypeOf(-2)
		keyValue, _ := state.ToInteger(-2)

		if keyType == lua.TypeNumber && keyValue == count+1 {
			count++
		} else {
			state.Pop(2)
			return false, nil
		}
		state.Pop(1)
	}

	return count == objLen, nil
}
