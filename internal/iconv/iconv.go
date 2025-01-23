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
	"log/slog"
	"unsafe"
)

type Encoding string

const (
	ENCODING_UNDECIDED = "undecided"
	ENCODING_UTF8      = "utf-8"
	ENCODING_EUCJP     = "euc-jisx0213"
)

type Iconv struct {
	handle C.iconv_t
}

var EUCJPConverter = func() *Iconv {
	iv, err := open(ENCODING_EUCJP, ENCODING_UTF8)
	if err != nil {
		slog.Error("Failed to open iconv", "err", err)
	}

	return iv
}()

func open(from string, to string) (*Iconv, error) {
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

func (iv *Iconv) ConvertLine(s string) (string, error) {
	if len(s) == 0 {
		return "", nil
	}

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

func (iv *Iconv) close() {
	C.iconv_close(iv.handle)
}
