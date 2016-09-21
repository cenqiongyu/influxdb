package escape // import "github.com/influxdata/influxdb/pkg/escape"

import (
	"bytes"
	"strings"
)

func Bytes(in []byte) []byte {
	for b, esc := range Codes {
		in = bytes.Replace(in, []byte{b}, esc, -1)
	}
	return in
}

const escapeChars = `," =`

func IsEscaped(b []byte) bool {
	for len(b) > 0 {
		i := bytes.IndexByte(b, '\\')
		if i < 0 {
			return false
		}

		if i+1 < len(b) && strings.IndexByte(escapeChars, b[i+1]) >= 0 {
			return true
		}
		b = b[i+1:]
	}
	return false
}

func AppendUnescaped(dst, src []byte) []byte {
	var off int
	for len(src) > 0 {
		next := bytes.IndexByte(src[off:], '\\')
		if next < 0 || off+next+1 >= len(src) {
			return append(dst, src...)
		}

		if off+next+1 < len(src) && strings.IndexByte(escapeChars, src[off+next+1]) >= 0 {
			if off+next > 0 {
				dst = append(dst, src[:off+next]...)
			}
			src = src[off+next+1:]
		}

		off += next + 1
	}

	return dst
}

func Unescape(in []byte) []byte {
	if len(in) == 0 {
		return nil
	}

	if bytes.IndexByte(in, '\\') == -1 {
		return in
	}

	i := 0
	inLen := len(in)
	var out []byte

	for {
		if i >= inLen {
			break
		}
		if in[i] == '\\' && i+1 < inLen {
			switch in[i+1] {
			case ',':
				out = append(out, ',')
				i += 2
				continue
			case '"':
				out = append(out, '"')
				i += 2
				continue
			case ' ':
				out = append(out, ' ')
				i += 2
				continue
			case '=':
				out = append(out, '=')
				i += 2
				continue
			}
		}
		out = append(out, in[i])
		i += 1
	}
	return out
}
