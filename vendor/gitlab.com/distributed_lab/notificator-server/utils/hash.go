package utils

import  (
	"strings"
	"crypto/sha1"
	"io"
	"encoding/base64"
	"log"
)

func Hash(args ...string) string {
	h := sha1.New()
	_, err := io.WriteString(h, strings.Join(args, ":"))
	if err != nil {
		// all possible reasons are app-fatal
		// so just to make linter happy
		log.Printf("something really bad happened: %s\n", err)
	}

	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
