package mimex

import "net/http"

func GetMimeType(buf []byte) string {
	var data []byte

	if len(buf) <= 512 {
		data = buf
	} else {
		data = buf[:512]
	}

	if len(data) < 1 {
		return ""
	}

	return http.DetectContentType(data)
}
