package iconv

// #cgo darwin,arm64 CFLAGS: -I/opt/homebrew/opt/libiconv/include
// #cgo darwin,!arm64 CFLAGS: -I/usr/local/opt/libiconv/include
// #cgo darwin,arm64 LDFLAGS: -L/opt/homebrew/opt/libiconv/lib -liconv
// #cgo darwin,!arm64 LDFLAGS: -L/usr/local/opt/libiconv/lib -liconv
//
// #include <stdlib.h>
// #include <iconv.h>
//
// size_t iconv_bridge(iconv_t ctx, char *in, size_t *size_in, char *out, size_t *size_out) {
//   return iconv(ctx, &in, size_in, &out, size_out);
// }
import "C"
import (
	"unsafe"
)

type Iconv struct {
	handle C.iconv_t
}

func Open(file string, from string, to string) (*Iconv, error) {
	fromcode := C.CString(from)
	defer C.free(unsafe.Pointer(fromcode))

	tocode := C.CString(to)
	defer C.free(unsafe.Pointer(tocode))

	ret, err := C.iconv_open(tocode, fromcode)
	if err != nil {
		return nil, err
	}

	return &Iconv{handle: ret}, nil
}

func (iv *Iconv) Convert(s string) (string, error) {
	buff := [4096]byte{}
	outptr := &buff[0]
	outlen := C.size_t(len(buff))

	in := []byte(s)
	inptr := &in[0]
	inlen := C.size_t(len(in))
	_, err := C.iconv_bridge(
		iv.handle,
		(*C.char)(unsafe.Pointer(inptr)),
		&inlen,
		(*C.char)(unsafe.Pointer(outptr)),
		&outlen,
	)
	return string(buff[:len(buff)-int(outlen)]), err
}

func (iv *Iconv) Close() {
	C.iconv_close(iv.handle)
}
