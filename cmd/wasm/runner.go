//go:build js && wasm

package main

import (
	"fmt"
	"syscall/js"

	"github.com/loicsikidi/wif-go/pkg/compiler"
)

var Version string

type Runner struct {
	c *compiler.Compiler
}

// convertJsObjectValueToMap converts a js.Value object to a map[string]any
// in order to be used as a payload for the compiler
func convertJsObjectValueToMap(payload js.Value) map[string]string {
	payloadMap := map[string]string{}
	payloadKeys := js.Global().Get("Object").Call("keys", payload)
	payloadKeysLength := payloadKeys.Length()

	for i := 0; i < payloadKeysLength; i++ {
		key := payloadKeys.Index(i).String()
		payloadMap[key] = payload.Get(key).String()
	}
	return payloadMap
}

func (r *Runner) Run(this js.Value, args []js.Value) (any, error) {
	if argLength := len(args); argLength < 2 || argLength > 3 {
		return nil, fmt.Errorf("run function expect 2 or 3 args, got %d", argLength)
	}

	payload := args[0].String()
	attrMapping := convertJsObjectValueToMap(args[1])

	var attrCondition string

	if len(args) == 3 {
		attrCondition = args[2].String()
	}

	r.c = &compiler.Compiler{
		Input: &compiler.Input{
			Payload:            payload,
			AttributeMapping:   attrMapping,
			AttributeCondition: attrCondition,
		},
	}

	res, err := r.c.Run()

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Runner) Version(this js.Value, args []js.Value) (any, error) {
	return Version, nil
}
