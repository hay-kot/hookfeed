package hookfeed

import (
	"encoding/json"
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

// WebhookRequest represents a simulated webhook request
type WebhookRequest struct {
	Body map[string]interface{} `json:"body"`
}

// Transformer handles Lua script execution for JSON transformation
type Transformer struct {
	scriptPath string
}

// NewTransformer creates a new Transformer with the given Lua script path
func NewTransformer(scriptPath string) *Transformer {
	return &Transformer{
		scriptPath: scriptPath,
	}
}

// Transform executes the Lua script's transform function with the given JSON input
func (t *Transformer) Transform(input []byte) ([]byte, error) {
	L := lua.NewState()
	defer L.Close()

	// Load the Lua script
	if err := L.DoFile(t.scriptPath); err != nil {
		return nil, fmt.Errorf("failed to load lua script: %w", err)
	}

	// Parse input JSON
	var inputData map[string]interface{}
	if err := json.Unmarshal(input, &inputData); err != nil {
		return nil, fmt.Errorf("failed to parse input JSON: %w", err)
	}

	// Convert JSON to Lua table
	inputTable := jsonToLuaTable(L, inputData)

	// Get the transform function
	transformFn := L.GetGlobal("transform")
	if transformFn.Type() != lua.LTFunction {
		return nil, fmt.Errorf("transform function not found in lua script")
	}

	// Call the transform function
	L.Push(transformFn)
	L.Push(inputTable)
	if err := L.PCall(1, 1, nil); err != nil {
		return nil, fmt.Errorf("failed to execute transform function: %w", err)
	}

	// Get the result
	result := L.Get(-1)
	L.Pop(1)

	// Convert Lua result back to JSON
	resultData := luaValueToGo(result)
	output, err := json.MarshalIndent(resultData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal output JSON: %w", err)
	}

	return output, nil
}

// jsonToLuaTable converts a Go map to a Lua table
func jsonToLuaTable(L *lua.LState, data map[string]interface{}) *lua.LTable {
	table := L.NewTable()
	for key, value := range data {
		table.RawSetString(key, goToLuaValue(L, value))
	}
	return table
}

// goToLuaValue converts Go values to Lua values
func goToLuaValue(L *lua.LState, value interface{}) lua.LValue {
	switch v := value.(type) {
	case nil:
		return lua.LNil
	case bool:
		return lua.LBool(v)
	case float64:
		return lua.LNumber(v)
	case string:
		return lua.LString(v)
	case map[string]interface{}:
		return jsonToLuaTable(L, v)
	case []interface{}:
		arr := L.NewTable()
		for i, item := range v {
			arr.RawSetInt(i+1, goToLuaValue(L, item))
		}
		return arr
	default:
		return lua.LNil
	}
}

// luaValueToGo converts Lua values to Go values
func luaValueToGo(lv lua.LValue) interface{} {
	switch v := lv.(type) {
	case *lua.LNilType:
		return nil
	case lua.LBool:
		return bool(v)
	case lua.LNumber:
		return float64(v)
	case lua.LString:
		return string(v)
	case *lua.LTable:
		// Check if it's an array or a map
		maxn := v.MaxN()
		if maxn > 0 {
			// It's an array
			arr := make([]interface{}, 0, maxn)
			for i := 1; i <= maxn; i++ {
				arr = append(arr, luaValueToGo(v.RawGetInt(i)))
			}
			return arr
		}
		// It's a map
		m := make(map[string]interface{})
		v.ForEach(func(key, value lua.LValue) {
			if keyStr, ok := key.(lua.LString); ok {
				m[string(keyStr)] = luaValueToGo(value)
			}
		})
		return m
	default:
		return nil
	}
}
