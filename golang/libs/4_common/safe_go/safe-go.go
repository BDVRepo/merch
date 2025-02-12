package safe_go

import (
	"bdv-avito-merch/libs/4_common/smart_context"
	"runtime/debug"
)

func SafeGo(logger smart_context.ISmartContext, f func()) {
	go func() {
		defer func() {
			if panicMessage := recover(); panicMessage != nil {
				stack := debug.Stack()

				logger.Errorf("RECOVERED FROM UNHANDLED PANIC: %v\nSTACK: %s", panicMessage, stack)
			}
		}()

		f()
	}()
}
