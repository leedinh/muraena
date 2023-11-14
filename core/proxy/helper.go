package proxy

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/dsnet/compress/brotli"
	ll "github.com/evilsocket/islazy/log"
	"github.com/evilsocket/islazy/tui"

	"github.com/muraenateam/muraena/log"
)

type Response struct {
	*http.Response
}

func (response *Response) Unpack() (buffer []byte, err error) {

	var rc io.ReadCloser

	switch response.Header.Get("Content-Encoding") {
	case "x-gzip":
		fallthrough
	case "gzip":
		rc, err = gzip.NewReader(response.Body)
		if err != io.EOF {
			buffer, _ = ioutil.ReadAll(rc)
			defer rc.Close()
		} else {
			err = nil
		}
	case "br":
		c := brotli.ReaderConfig{}
		rc, err = brotli.NewReader(response.Body, &c)
		buffer, _ = ioutil.ReadAll(rc)
		defer rc.Close()
	case "deflate":
		rc = flate.NewReader(response.Body)
		buffer, _ = ioutil.ReadAll(rc)
		defer rc.Close()
	case "compress":
		fallthrough
	default:
		rc = response.Body
		buffer, _ = ioutil.ReadAll(rc)
		defer rc.Close()
	}
	return
}

func (response *Response) Encode(buffer []byte) (err error) {

	switch response.Header.Get("Content-Encoding") {
	case "x-gzip":
		fallthrough
	case "gzip":
		buffer, err = packGzip(buffer)
	case "deflate":
		buffer, err = packDeflate(buffer)
	case "br":
		response.Header.Set("Content-Encoding", "deflate")
		buffer, err = packDeflate(buffer)
	default:
		// just don't pack
	}

	if err != nil {
		ll.Error("[Encode] Error packing with %s: %s", response.Header.Get("Content-Encoding"), err)
	}

	body := ioutil.NopCloser(bytes.NewReader(buffer))
	response.Body = body
	response.ContentLength = int64(len(buffer))
	response.Header.Set("Content-Length", strconv.Itoa(len(buffer)))
	return nil
}

func packGzip(i []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(i); err != nil {
		return nil, err
	}
	if err := gz.Flush(); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func packDeflate(i []byte) ([]byte, error) {
	var b bytes.Buffer
	zz, err := flate.NewWriter(&b, 0)

	if err != nil {
		return nil, err
	}
	if _, err = zz.Write(i); err != nil {
		return nil, err
	}
	if err := zz.Flush(); err != nil {
		return nil, err
	}
	if err := zz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// ArmorDomain filters duplicate strings in place and returns a slice with
// only unique strings.
func ArmorDomain(slice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true

			// make it lowercase
			entry = strings.ToLower(entry)

			// if string is URL encoded, decode it
			if strings.Contains(entry, "%") {
				decodedURL, err := url.QueryUnescape(entry)
				if err == nil {
					entry = decodedURL
				}
			}

			// if string contains protocol, remove it
			if strings.Contains(entry, "://") {
				entry = strings.Split(entry, "://")[1]
			}

			// remove anything after the / (if any)
			if strings.Contains(entry, "/") {
				entry = strings.Split(entry, "/")[0]
			}

			list = append(list, entry)
		}
	}
	return list
}

func isWildcard(s string) bool {
	return strings.HasPrefix(s, "*.")
}

// IsSubdomain checks if a string is a subdomain of another string.
// It returns true if the given string is a subdomain of the root string.
func IsSubdomain(root string, subdomain string) bool {

	// TODO: add support for wildcard domains
	return strings.HasSuffix(subdomain, root)
}

func base64Decode(input string, padding int32) (string, bool) {
	encoding := base64.StdEncoding.WithPadding(padding)
	d, err := encoding.DecodeString(input)
	if err != nil {
		return "", false
	}

	return string(d), true
}

func base64Encode(input string, padding int32) string {
	encoding := base64.StdEncoding.WithPadding(padding)
	e := encoding.EncodeToString([]byte(input))

	return e
}

func getPadding(p string) int32 {
	if len(p) > 1 {
		// Fallback to default padding =
		log.Debug("Invalid Base64 padding value [%s], falling back to %s", tui.Red(p), tui.Green(string(Base64Padding)))
		return Base64Padding
	}

	return []rune(p)[0]
}
