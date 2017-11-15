package uri

import "encoding/base64"

// URL:
//	 http://host/url
//	 https://host/url
// Path:
//	 AbsolutePath	(Must start with '/')
//	 Pid:RelPath	(Pid.len = 16)
//	 Id: 			(Id.len = 16)
//	 :LinkId:RelPath
//	 :LinkId
func Encode(uri string) string {

	size := len(uri)
	if size == 0 {
		return ""
	}

	encodedURI := encode(uri)
	if c := uri[0]; c == '/' || c == ':' || (size > 16 && encodedURI[16] == ':') || (size > 5 && (encodedURI[4] == ':' || encodedURI[5] == ':')) {
		return encodedURI
	}
	return "!" + encodedURI
}

func Decode(encodedURI string) (uri string, err error) {

	size := len(encodedURI)
	if size == 0 {
		return
	}

	if c := encodedURI[0]; c == '!' || c == ':' || (size > 16 && encodedURI[16] == ':') || (size > 5 && (encodedURI[4] == ':' || encodedURI[5] == ':')) {
		uri, err = decode(encodedURI)
		if err != nil {
			return
		}
		if c == '!' {
			uri = uri[1:]
		}
		return
	}

	b := make([]byte, base64.URLEncoding.DecodedLen(len(encodedURI)))
	n, err := base64.URLEncoding.Decode(b, []byte(encodedURI))
	return string(b[:n]), err
}
