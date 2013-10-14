package dbg

// Println is equivalent to Default.Println
func Println(v ...interface{}) {
	Default.Println(v...)
}

// Printf is equivalent to Default.Printf
func Printf(format string, v ...interface{}) {
	Default.Printf(format, v...)
}

// Print is equivalent to Default.Print
func Print(v ...interface{}) {
	Default.Print(v...)
}
