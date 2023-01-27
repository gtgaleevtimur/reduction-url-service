// Package main(pkg1) datafile предназначенный исключительно для теста ExitAnalyzer.
package main

import "os"

// errExitFunc - функция хелпер для analysistest.
func errExitFunc() {
	os.Exit(1) // want "os.Exit call in main package"
}
