// This file provides tests with access to internal package state.

package intern

// ForgetEverything invokes all of the package-internal forgetAll* functions.
func ForgetEverything() {
	eq.forgetAll()
	eqc.forgetAll()
	lge.forgetAll()
	lgec.forgetAll()
}
