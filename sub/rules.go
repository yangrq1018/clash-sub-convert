package sub

import (
	"time"

	"github.com/biter777/countries"
	"github.com/enescakir/emoji"
)

var (
	selfHosted = selectGroup("自建服务器",
		selfHostedServer2SZ.Name,
	)
	uncommon    = selectGroup("小众节点")
	xiaohongshu = selectGroup("小红书",
		DIRECT,
		uncommon.Name,
		grand().Name,
	)
	reddit = selectGroup("Reddit",
		DIRECT,
		grand().Name,
		countryGroup(countries.HK),
		countryGroup(countries.TW),
		countryGroup(countries.JP),
	)
	zhihu = selectGroup("知乎",
		DIRECT,
		uncommon.Name,
		grand().Name,
	)
	qq = selectGroup("QQ",
		DIRECT,
		selfHosted.Name,
		grand().Name,
	)
	// Grand proxy group that contains all proxies
	telegram = urlTestGroup(
		emoji.Airplane.String()+"Telegram",
		300*time.Second,
		countryGroup(countries.HK),
		countryGroup(countries.SG),
		countryGroup(countries.US),
	)
	rest = selectGroup(
		emoji.Fish.String()+"漏网之鱼",
		DIRECT,
		grand().Name,
		selfHosted.Name,
	)
	apple = selectGroup(
		emoji.RedApple.String()+"Apple",
		DIRECT,
		grand().Name,
	)
	embyUnlock = selectGroup(
		emoji.PuzzlePiece.String()+"Emby Unlock",
		DIRECT,
		unlockEMBYServer.Name,
	)
	// TAG New flavor EMBY proxies
	embyTagNewFlavor = selectGroup(
		emoji.PuzzlePiece.String()+"Emby Tag New Flavor",
		DIRECT,
		countryGroup(countries.HK),
		countryGroup(countries.JP),
	)
	proxyConverter = selectGroup(
		emoji.Sparkles.String()+"订阅转换",
		DIRECT,
		proxyConverterServerLocal.Name,
	)
	microsoft = selectGroup(
		emoji.DesktopComputer.String()+"Microsoft",
		countryGroup(countries.HK),
		countryGroup(countries.SG),
		countryGroup(countries.US),
		"DIRECT",
	)
	github = selectGroup(
		emoji.CatFace.String()+"Github",
		countryGroup(countries.HK),
		countryGroup(countries.SG),
		countryGroup(countries.US),
		grand().Name,
		"DIRECT",
	)
	steam = selectGroup(
		emoji.SteamingBowl.String()+"Steam",
		countryGroup(countries.HK),
		countryGroup(countries.SG),
		countryGroup(countries.US),
		"DIRECT",
	)
	youtube = selectGroup(
		emoji.Television.String()+"YouTube",
		countryGroup(countries.HK),
		countryGroup(countries.SG),
		countryGroup(countries.US),
		grand().Name,
		"DIRECT",
	)
	amazon = selectGroup(
		emoji.Television.String()+"Amazon",
		countryGroup(countries.US),
	)
	switchGroup = selectGroup(
		emoji.VideoGame.String()+"Switch",
		DIRECT,
		countryGroup(countries.Japan),
	)
	openAIGroup = selectGroup(
		emoji.GlobeShowingEuropeAfrica.String()+"Open AI",
		countryGroup(countries.Singapore),
    grand().Name,
	)
	reject = selectGroup(
		emoji.NoEntry.String()+"垃圾拦截",
		REJECT,
		DIRECT,
	)
	ipCheck = selectGroup(
		emoji.TelephoneReceiver.String()+"IP检查",
		DIRECT,
		grand().Name,
	)
	minecraft = selectGroup(
		emoji.Rock.String()+"Minecraft",
		DIRECT,
		grand().Name,
	)
	streamMedia = []struct {
		mediaKey  string
		countries []countries.CountryCode
	}{
		{
			mediaKey:  "Hulu",
			countries: []countries.CountryCode{countries.US},
		},
		{
			mediaKey:  "Netflix",
			countries: []countries.CountryCode{countries.SG, countries.JP, countries.HK},
		},
		{
			mediaKey:  "Pornhub",
			countries: []countries.CountryCode{countries.US, countries.HK},
		},
		{
			mediaKey:  "Bilibili",
			countries: []countries.CountryCode{countries.Taiwan, countries.HongKong, countries.Thailand},
		},
	}
)
