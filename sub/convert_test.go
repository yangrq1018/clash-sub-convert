package sub

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/biter777/countries"
	emojiflag "github.com/jayco/go-emoji-flag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// TAGSSRSubsLink TAG的SSR订阅地址
	TAGSSRSubsLink = "https://subscribe.hlasw.com/link/HheFeZMnfkMide9X?sub=1&extend=1"
	// TAGSSSubsLink TAG的SS订阅地址
	TAGSSSubsLink = "https://newsubscribe.hlasw.com/api/v1/client/subscribe?token=93cb6ac0990ea0c78ee3af284a4c9c0a"
)

func TestRewrite(t *testing.T) {
	link := "https://subscribe.hlasw.com/link/vCddKFHm6Tljaz0X?clash=2"
	// read nodes from subLink
	res, err := http.Get(link)
	require.NoError(t, err)
	remote, err := decodeClashConfig(res)
	require.NoError(t, err)
	assert.NoError(t, Rewrite(remote, io.Discard))
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
		fmt.Printf("%+v\n", line)
	}
}

func TestGetRemainingDataSSR(t *testing.T) {
	text, err := GetRemainingDataSSR(TAGSSRSubsLink)
	assert.NoError(t, err)
	fmt.Println(text)
}
