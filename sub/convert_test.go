package sub

import (
	"fmt"
	"github.com/biter777/countries"
	emojiflag "github.com/jayco/go-emoji-flag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestRewrite(t *testing.T) {
	link := "https://subscribe.hlasw.com/link/vCddKFHm6Tljaz0X?clash=2"
	// read nodes from subLink
	res, err := http.Get(link)
	require.NoError(t, err)
	remote, err := decodeConfig(res)
	require.NoError(t, err)
	assert.NoError(t, Rewrite(remote, io.Discard))
}

func TestOverwriteConfigByFilename(t *testing.T) {
	link := "https://subscribe.hlasw.com/link/vCddKFHm6Tljaz0X?clash=2"
	home, _ := os.UserHomeDir()
	out := filepath.Join(home, ".config/clash/profiles/config_generated.yaml")
	processors := []Processor{
		AddHosts(DNSMapping{"*.xuanlingasset.com": "10.168.1.185"}), // 内网域名寻址,指向NGINX服务器
		AddRuleIPCIDR("10.168.1.0/24", DIRECT),                      // 内网网段
	}
	assert.NoError(t, OverwriteConfigByFilename(link, out, processors...))
}

func TestGetEmojiFlag(t *testing.T) {
	for _, code := range []string{
		"KH", // 柬埔寨
	} {
		flag := emojiflag.GetFlag(code)
		fmt.Printf("%s %s\n", code, flag)
	}
}

func TestGetCountryCode(t *testing.T) {
	for _, c := range []struct {
		name        string
		countryCode countries.CountryCode
	}{
		{"Na 美国 08 底特律 US 2倍率", countries.US},
		{"As 香港 06 HK 2倍率", countries.HK},
		{"香港 01丨1x HK", countries.HongKong},
		{"xyz 123", countries.Unknown},
	} {
		cc := extractCountryFromNodeName(c.name)
		assert.Equal(t, c.countryCode, cc)
	}
}

func TestGetRemainingDataClash(t *testing.T) {
	usage, err := GetRemainingDataClash(TAGClashSubsLink)
	assert.NoError(t, err)
	fmt.Println(usage)
}

func TestDecodeSSR(t *testing.T) {
	res, err := http.Get(TAGSSRSubsLink)
	assert.NoError(t, err)
	items, err := DecodeSSR(res)
	assert.NoError(t, err)
	for _, line := range items {
		fmt.Printf("%+v\n", line)
	}
}

func TestDecodeSS(t *testing.T) {
	res, err := http.Get(TAGSSSubsLink)
	assert.NoError(t, err)
	items, err := DecodeSS(res)
	assert.NoError(t, err)
	for _, line := range items {
		fmt.Println(line.Server, " ", line.Method)
	}
}

func TestGetRemainingDataSSR(t *testing.T) {
	text, err := GetRemainingDataSSR(TAGSSRSubsLink)
	assert.NoError(t, err)
	fmt.Println(text)
}
