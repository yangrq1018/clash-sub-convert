package sub

import (
	"io"
	"net/http"
	"testing"

	"github.com/biter777/countries"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	clashLink = "https://stc-anycast.com/link/ABRHuzP2QGPxHHZm?sub=2&client=clash"
	ssrLink   = "https://stc-anycast.com/link/ABRHuzP2QGPxHHZm?sub=2"
	ssLink    = "https://newsubscribe.hlasw.com/api/v1/client/subscribe?token=93cb6ac0990ea0c78ee3af284a4c9c0a"
)

func TestRewrite(t *testing.T) {
	res, err := http.Get(clashLink)
	require.NoError(t, err)
	remote, err := decodeClashConfig(res)
	require.NoError(t, err)
	assert.NoError(t, Rewrite(remote, io.Discard, "", false))
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
	res, err := http.Get(ssrLink)
	assert.NoError(t, err)
	items, err := DecodeSSR(res)
	assert.NoError(t, err)
	for _, line := range items {
		t.Logf("%+v", line)
	}
}

func TestDecodeSS(t *testing.T) {
	res, err := http.Get(ssLink)
	assert.NoError(t, err)
	items, err := DecodeSS(res)
	assert.NoError(t, err)
	for _, line := range items {
		t.Logf("%+v", line)
	}
}
