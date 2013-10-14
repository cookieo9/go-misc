// Package pp provides pretty printing services.
//
// NOTE: The formatting of the pretty printer should not be
// assumed to follow any known or fixed format, and could
// be changed at any time.
package pp

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"text/tabwriter"
)

// PP pretty-prints it's argument using reflection and text/tabwriter.
func PP(x interface{}) string {
	out := bytes.NewBuffer(nil)
	tw := tabwriter.NewWriter(out, 4, 4, 1, ' ', 0)
	pp(x, "", tw)
	tw.Flush()
	return string(out.Bytes())
}

func pp(x interface{}, tabs string, w io.Writer) {
	if _, ok := x.(fmt.Stringer); ok {
		fmt.Fprintf(w, "%q", x)
		return
	}

	v := reflect.Indirect(reflect.ValueOf(x))
	switch v.Kind() {
	case reflect.Struct:
		ppStruct(x, tabs, w)
	case reflect.Array, reflect.Slice:
		ppSlice(x, tabs, w)
	case reflect.Map:
		ppMap(x, tabs, w)
	case reflect.String:
		fmt.Fprintf(w, "%q", x)
	default:
		fmt.Fprintf(w, "%v", x)
	}
}

func ppStruct(x interface{}, tabs string, w io.Writer) {
	fmt.Fprintf(w, "%T (\n", x)
	v := reflect.Indirect(reflect.ValueOf(x))
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).PkgPath == "" {
			fmt.Fprint(w, tabs, "\t", t.Field(i).Name, ":\t")
			pp(v.Field(i).Interface(), tabs+"\t\t", w)
			fmt.Fprintln(w, ",")
		}
	}
	fmt.Fprint(w, tabs, ")")
}

// BUG(cookieo9) Slices are printed with 1 entry per line, which is
// not very 'pretty' for []byte and the like.
func ppSlice(x interface{}, tabs string, w io.Writer) {
	v := reflect.Indirect(reflect.ValueOf(x))
	if v.Len() == 0 {
		fmt.Fprint(w, "[]")
		return
	}

	fmt.Fprintf(w, "[\n")
	for i := 0; i < v.Len(); i++ {
		fmt.Fprint(w, tabs, "\t")
		pp(v.Index(i).Interface(), tabs+"\t", w)
		fmt.Fprintln(w, ",")
	}
	fmt.Fprint(w, tabs, "]")
}

func ppMap(x interface{}, tabs string, w io.Writer) {
	fmt.Fprintf(w, "{\n")
	v := reflect.Indirect(reflect.ValueOf(x))
	for _, k := range v.MapKeys() {
		fmt.Fprint(w, tabs, "\t", k.String(), ":\t")
		pp(v.MapIndex(k).Interface(), tabs+"\t", w)
		fmt.Fprintln(w, ",")
	}
	fmt.Fprint(w, tabs, "}")
}
