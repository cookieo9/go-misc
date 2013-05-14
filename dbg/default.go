package dbg

// Equivalent to Default.Println
func Println(v ...interface{}) {
	Default.Println(v...)
}

// Equivalent to Default.Printf
func Printf(format string, v ...interface{}) {
	Default.Printf(format, v...)
}

// Equivalent to Default.Print
func Print(v ...interface{}) {
	Default.Print(v...)
}
