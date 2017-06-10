// This file provides tests with access to internal package state.

package intern

// ForgetAllEq provides access to the package-internal forgetAllEq function.
func ForgetAllEq() {
	lge.forgetAll()
}

// ForgetAllEqC provides access to the package-internal forgetAllEqC function.
func ForgetAllEqC() {
	lgec.forgetAll()
}

// ForgetAllLGE provides access to the package-internal forgetAllLGE function.
func ForgetAllLGE() {
	lge.forgetAll()
}

// ForgetAllLGEC provides access to the package-internal forgetAllLGEC function.
func ForgetAllLGEC() {
	lgec.forgetAll()
}
