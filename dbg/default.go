package dbg

func Println(v ...interface{}) {
	Default.Println(v...)
}

func Printf(format string, v ...interface{}) {
	Default.Printf(format, v...)
}

func Print(v ...interface{}) {
	Default.Print(v...)
}
