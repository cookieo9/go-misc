// +build ignore

package main

import (
	"bytes"
	"fmt"
	"github.com/cookieo9/go-misc/cgo"
	"log"
	"os"
	"unsafe"
)

// #include <stdio.h>
import "C"

func test_read() {
	log.Println("test_read")
	rdr := bytes.NewReader([]byte("Hello, World!"))

	wrapped, err := cgo.WrapReader(rdr, false)
	if err != nil {
		log.Fatal(err)
	}
	wf := (*C.FILE)(wrapped)
	defer C.fclose(wf)

	var mem [80]byte
	buf := unsafe.Pointer(&mem[0])
	siz := C.size_t(len(mem))

	n, err := C.fread(buf, 1, siz, wf)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(n, "bytes read")

	r, err := C.puts((*C.char)(buf))
	if r == C.EOF || err != nil {
		log.Fatal(err)
	}
}

func test_write() {
	log.Println("test_write")
	var buf bytes.Buffer

	wrap, err := cgo.WrapReadWriter(&buf, false)
	if err != nil {
		log.Fatal(err)
	}
	wf := (*C.FILE)(wrap)
	defer C.fclose(wf)

	if x, err := C.fputs(C.CString("Hello, world! I would like to write a number, but printf isn't supported in CGO."), wf); x == C.EOF {
		log.Fatal(err)
	}

	if x, err := C.fflush(wf); x == C.EOF {
		log.Fatal(err)
	}

	log.Printf("received: %d %q", buf.Len(), string(buf.Bytes()))
}

func test_no_close() {
	log.Println("test_no_close")
	f, err := os.Create("foo.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	wrap, err := cgo.WrapWriter(f, false)
	if err != nil {
		log.Fatal(err)
	}
	wf := (*C.FILE)(wrap)

	if x, err := C.fputs(C.CString("Foo!\n"), wf); x == C.EOF {
		log.Fatal(err)
	}
	if x, err := C.fclose(wf); x == C.EOF {
		log.Fatal(err)
	}

	fmt.Fprintln(f, "Bar!")
}

func main() {
	test_read()
	test_write()
	test_no_close()
}
