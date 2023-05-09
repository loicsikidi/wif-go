//go:build js && wasm

package main

import (
	"syscall/js"
)

var Promise = js.Global().Get("Promise")

func main() {
	runner := &Runner{}

	js.Global().Set("wif_run", asyncFuncOf(runner.Run))
	js.Global().Set("wif_version", asyncFuncOf(runner.Version))
	select {}
}

// asyncFuncOf avoids Go's js deadlock
//
// source: https://github.com/golang/go/issues/41310#issuecomment-725809881
func asyncFuncOf(fn func(this js.Value, args []js.Value) (any, error)) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		handler := js.FuncOf(func(_ js.Value, promise []js.Value) interface{} {
			resolve := promise[0]
			reject := promise[1]

			go func() {
				res, err := fn(this, args)
				if err != nil {
					reject.Invoke(err.Error())
					return
				}
				resolve.Invoke(res)
			}()

			return nil
		})

		return Promise.New(handler)
	})
}
