//go:build ignore

package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
)

const srcUrl = "https://www.apple.com/certificateauthority/"
const outDir = "certs/"

var certLinkPattern = regexp.MustCompile(`<a [^>]*href="([^"]+\.cer)"`)

func main() {
	resp, err := http.Get(srcUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if err := os.RemoveAll(outDir); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		panic(err)
	}
	matches := certLinkPattern.FindAllSubmatch(content, -1)
	for _, match := range matches {
		certPath := string(match[1])
		var certUrl string
		if certPath[0] == '/' {
			baseUrl, err := url.Parse(srcUrl)
			if err != nil {
				panic(err)
			}
			baseUrl.Path = certPath
			certUrl = baseUrl.String()
		} else if strings.HasPrefix(certPath, "https://www.apple.com/") || strings.HasPrefix(certPath, "https://developer.apple.com/") {
			certUrl = certPath
		} else {
			certUrl, err = url.JoinPath(srcUrl, certPath)
			if err != nil {
				panic(err)
			}
		}
		resp, err := http.Get(certUrl)
		if err != nil {
			panic(err)
		}
		fmt.Println(resp.StatusCode, certUrl)
		if resp.StatusCode != 200 {
			continue
		}
		fileName := path.Base(certPath)
		f, err := os.OpenFile(outDir+fileName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(f, resp.Body)
		if err != nil {
			panic(err)
		}
		err = resp.Body.Close()
		if err != nil {
			panic(err)
		}
	}
}
