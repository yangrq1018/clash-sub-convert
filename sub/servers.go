package sub

import (
	"strconv"

	"github.com/enescakir/emoji"
)

// servers
var (
	selfHostedServer2SZ = Node{
		Name:     "自建服务器2深圳",
		Type:     "ss",
		Server:   "39.108.10.209",
		Port:     "8388",
		Cipher:   "chacha20-ietf-poly1305",
		Password: "HX1J7MYQ7H5Y",
		UDP:      true,
		TFO:      true,
	}
	unlockEMBYServer = Node{
		Name:           emoji.Joystick.String() + "Crack Emby",
		Server:         "34.92.170.135",
		Port:           strconv.Itoa(29967),
		Type:           "http",
		TLS:            false,
		SkipCertVerify: false,
		UDP:            true,
	}
	proxyConverterServerLocal = Node{
		Name:           emoji.Sparkles.String() + "Proxy Convert (Local)",
		Server:         "127.0.0.1",
		Port:           "39923",
		Type:           "http",
		TLS:            false,
		SkipCertVerify: false,
		UDP:            false,
	}
	// add every node above here
	userDefinedNodes = []Node{unlockEMBYServer, proxyConverterServerLocal, selfHostedServer2SZ}
)
