// This file provides tests with access to internal package state.

package intern

import (
	"runtime"
)

// ForgetAllEq provides access to the package-internal forgetAllEq function.
func ForgetAllEq() {
	lge.forgetAll()
	runtime.GC()
}

// ForgetAllEqC provides access to the package-internal forgetAllEqC function.
func ForgetAllEqC() {
	lgec.forgetAll()
	runtime.GC()
}

// ForgetAllLGE provides access to the package-internal forgetAllLGE function.
func ForgetAllLGE() {
	lge.forgetAll()
	runtime.GC()
}

// ForgetAllLGEC provides access to the package-internal forgetAllLGEC function.
func ForgetAllLGEC() {
	lgec.forgetAll()
	runtime.GC()
}
