// This file provides tests with access to internal package state.

package intern

// ForgetAllLGE provides access to the package-internal forgetAllLGE function.
func ForgetAllLGE() {
	forgetAllLGE()
}

// ForgetAllLGEC provides access to the package-internal forgetAllLGEC function.
func ForgetAllLGEC() {
	lgec.forgetAll()
}
