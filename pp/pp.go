// Package PP provides a function which pretty-prints it's argument using reflection and text/tabwriter.
package pp

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"text/tabwriter"
)

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
		pp_struct(x, tabs, w)
	case reflect.Array, reflect.Slice:
		pp_slice(x, tabs, w)
	case reflect.Map:
		pp_map(x, tabs, w)
	case reflect.String:
		fmt.Fprintf(w, "%q", x)
	default:
		fmt.Fprintf(w, "%v", x)
	}
}

func pp_struct(x interface{}, tabs string, w io.Writer) {
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

func pp_slice(x interface{}, tabs string, w io.Writer) {
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

func pp_map(x interface{}, tabs string, w io.Writer) {
	fmt.Fprintf(w, "{\n")
	v := reflect.Indirect(reflect.ValueOf(x))
	for _, k := range v.MapKeys() {
		fmt.Fprint(w, tabs, "\t", k.String(), ":\t")
		pp(v.MapIndex(k).Interface(), tabs+"\t", w)
		fmt.Fprintln(w, ",")
	}
	fmt.Fprint(w, tabs, "}")
}
