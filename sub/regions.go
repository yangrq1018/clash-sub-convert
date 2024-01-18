package sub

import (
	"github.com/biter777/countries"
)

var (
	// Make sure every code is provided in upstream sub!!
	countriesNeeded = []struct {
		countries.CountryCode
		priority
	}{
		{countries.HK, Required},            // 香港
		{countries.TW, Required},            // 台湾
		{countries.JP, Required},            // 日本
		{countries.Singapore, Required},     // 新加坡
		{countries.US, Required},            // 美国
		{countries.UnitedKingdom, Required}, // 英国
		{countries.France, Required},        // 法国
		{countries.Germany, Required},       // 德国

		// optional
		{countries.Thailand, Optional}, // 泰国
		{countries.Korea, Optional},    // 韩国
		{countries.IS, Optional},       // 冰岛
	}
)
