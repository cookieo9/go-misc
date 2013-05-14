package dbg

import (
	"bytes"
	"log"
	"os"
	. "testing"
	"testing/quick"
)

func setupLogger() *bytes.Buffer {
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	log.SetFlags(0)
	return buf
}

func resetLogger() {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags)
}

var debugTests = []struct {
	Format                 string
	Args                   []interface{}
	Println, Print, Printf string
}{
	{"%v %q %s", []interface{}{"foo", "bar", "baz"}, "foo bar baz", "foobarbaz", "foo \"bar\" baz"},
	{"%x %T", []interface{}{42, false}, "42 false", "42 false", "2a bool"},
}

func TestDebug(t *T) {
	buf := setupLogger()
	defer resetLogger()

	for _, test := range debugTests {
		off := Debug(false)

		buf.Reset()
		off.Println(test.Args...)
		if buf.Len() != 0 {
			t.Errorf("off.Println(%v): expected no output, got %q", test.Args, buf.String())
		}

		buf.Reset()
		off.Print(test.Args...)
		if buf.Len() != 0 {
			t.Errorf("off.Print(%v): expected no output, got %q", test.Args, buf.String())
		}

		buf.Reset()
		off.Printf(test.Format, test.Args...)
		if buf.Len() != 0 {
			t.Errorf("off.Printf(%q, %v): expected no output, got %q", test.Format, test.Args, buf.String())
		}

		on := Debug(true)

		buf.Reset()
		on.Println(test.Args...)
		out := string(buf.Bytes()[:buf.Len()-1])
		if out != test.Println {
			t.Errorf("on.Println(%v): expected %q, got %q", test.Args, test.Println, out)
		}

		buf.Reset()
		on.Print(test.Args...)
		out = string(buf.Bytes()[:buf.Len()-1])
		if out != test.Print {
			t.Errorf("on.Print(%v): expected %q, got %q", test.Args, test.Print, out)
		}

		buf.Reset()
		on.Printf(test.Format, test.Args...)
		out = string(buf.Bytes()[:buf.Len()-1])
		if out != test.Printf {
			t.Errorf("on.Printf(%q, %v): expected %q got %q", test.Format, test.Args, test.Printf, out)
		}
	}
}

func TestPackageMethods(t *T) {
	pkg := func(a, b, c, d, e, f string) string {
		buf := setupLogger()
		defer resetLogger()

		Println(a, b, c, d, e, f)
		Print(a, b, c, d, e, f)
		Printf(a, b, c, d, e, f)
		return buf.String()
	}

	dft := func(a, b, c, d, e, f string) string {
		buf := setupLogger()
		defer resetLogger()

		Default.Println(a, b, c, d, e, f)
		Default.Print(a, b, c, d, e, f)
		Default.Printf(a, b, c, d, e, f)
		return buf.String()
	}

	if err := quick.CheckEqual(pkg, dft, nil); err != nil {
		t.Error(err)
	}
}
