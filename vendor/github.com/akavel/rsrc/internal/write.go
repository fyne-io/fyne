package internal

import (
	"fmt"
	"os"
	"reflect"

	"github.com/akavel/rsrc/binutil"
	"github.com/akavel/rsrc/coff"
)

// TODO(akavel): maybe promote this to coff.Coff.WriteTo(io.Writer) (int64, error)
func Write(coff *coff.Coff, fnameout string) error {
	out, err := os.Create(fnameout)
	if err != nil {
		return err
	}
	defer out.Close()
	w := binutil.Writer{W: out}

	// write the resulting file to disk
	binutil.Walk(coff, func(v reflect.Value, path string) error {
		if binutil.Plain(v.Kind()) {
			w.WriteLE(v.Interface())
			return nil
		}
		vv, ok := v.Interface().(binutil.SizedReader)
		if ok {
			w.WriteFromSized(vv)
			return binutil.WALK_SKIP
		}
		return nil
	})

	if w.Err != nil {
		return fmt.Errorf("Error writing output file: %s", w.Err)
	}

	return nil
}
