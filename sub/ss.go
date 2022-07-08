package sub

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var ssPattern = regexp.MustCompile("ss://(.*)")

type SSItem struct {
	Name     string
	Plugins  map[string]string
	Server   string
	Port     string
	Password string
	Method   string
}

func ssToClash(line []byte) (*SSItem, error) {
	param := ssPattern.FindSubmatch(line)[1]
	info := SSItem{
		Plugins: make(map[string]string),
	}
	if bytes.Contains(param, []byte("#")) {
		remark, _ := url.QueryUnescape(string(param[bytes.Index(param, []byte("#"))+1:]))
		info.Name = remark
		param = param[:bytes.Index(param, []byte("#"))]
	}

	if bytes.Contains(param, []byte("/?")) {
		plugin, _ := url.QueryUnescape(string(param[bytes.Index(param, []byte("/?"))+2:]))
		param = param[:bytes.Index(param, []byte("/?"))]
		for _, p := range strings.Split(plugin, ";") {
			kv := strings.Split(p, "=")
			info.Plugins[kv[0]] = kv[1]
		}
	}

	if bytes.Contains(param, []byte("@")) {
		matcher := regexp.MustCompile("(.*)@(.*):(.*)").FindSubmatch(param)
		if matcher != nil {
			param = matcher[1]
			info.Server = string(matcher[2])
			info.Port = string(matcher[3])
		} else {
			return &info, nil
		}

		// unpadded, avoid truncate password
		param, _ = decodeBase64(param, base64.RawStdEncoding)
		matcher = regexp.MustCompile("(.*?):(.*)").FindSubmatch(param)
		if matcher != nil {
			info.Method = string(matcher[1])
			info.Password = string(matcher[2])
		} else {
			return &info, nil
		}
	} else {
		// 需要再次解密一次
		param, _ = decodeBase64(param, base64.StdEncoding)
		matcher := regexp.MustCompile("(.*?):(.*)@(.*):(.*)").FindSubmatch(param)
		if matcher != nil {
			info.Method = string(matcher[1])
			info.Password = string(matcher[2])
			info.Server = string(matcher[3])
			info.Port = string(matcher[4])
		} else {
			return &info, nil
		}
	}

	return &info, nil
}

func DecodeSS(res *http.Response) ([]SSItem, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	bodyDecoded, err := decodeBase64(body, base64.StdEncoding)
	if err != nil {
		return nil, err
	}

	lines := bytes.Split(bodyDecoded, []byte{'\n'})
	var subs []SSItem
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if ssPattern.Match(line) {
			item, err := ssToClash(line)
			if err != nil {
				return nil, err
			}
			subs = append(subs, *item)
		}
	}
	return subs, nil
}
