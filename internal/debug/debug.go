package debug

import "fmt"

func Assert(condition bool, message ...string) {
	if !condition {
		if len(message) > 0 {
			panic(fmt.Sprintf("assertion failed: %s", message[0]))
		}
		panic("assertion failed")
	}
}

func Assertf(condition bool, format string, args ...any) {
	if !condition {
		panic(fmt.Sprintf("assertion failed: "+format, args...))
	}
}

func Panic() {
	panic("unreachable")
}

func Panicf(format string, args ...any) {
	panic(fmt.Sprintf(format, args...))
}

func Fail(message ...string) {
	if len(message) > 0 {
		panic(fmt.Sprintf("fail: %s", message[0]))
	}
	panic("fail")
}

func FailBadSyntaxKind(node any, message ...string) {
	if len(message) > 0 {
		panic(fmt.Sprintf("unexpected syntax kind: %T - %s", node, message[0]))
	}
	panic(fmt.Sprintf("unexpected syntax kind: %T", node))
}

func AssertNever(value any, message ...string) {
	if len(message) > 0 {
		panic(fmt.Sprintf("unexpected value: %v - %s", value, message[0]))
	}
	panic(fmt.Sprintf("unexpected value: %v", value))
}

func SetCrashOnPanic(value bool) {}

func IsCrashOnPanicEnabled() bool { return false }

func ApplyDebugStackLimit() {}
