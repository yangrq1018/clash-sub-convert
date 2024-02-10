package sub

// rules
var (
	RulesNintendo = []Rule{
		DomainSuffixRule("stat.ink", switchGroup),
		DomainSuffixRule("nintendo.net", switchGroup),
		DomainSuffixRule("nintendo.com", switchGroup),
		DomainSuffixRule("s3-us-west-2.amazonaws.com", switchGroup),
	}
	RulesOpenAI = []Rule{
		DomainSuffixRule("openai.com", openAIGroup),
	}
	RulesEmby = []Rule{
		DomainSuffixRule("mb3admin.com", embyUnlock),
		DomainSuffixRule("tagemby.embylianmeng.com", embyTagNewFlavor),
	}
	RulesProxyConverterRules = []Rule{
		DomainRule("subscribe.hlasw.com", proxyConverter),
		DomainRule("subscribe.tagonline.asia", proxyConverter),
	}
	RulesSpecial = []Rule{
		DomainRule("cip.cc", ipCheck),
		DomainRule("ipinfo.io", ipCheck),
	}
	RulesMinecraft = []Rule{
		DomainKeyWordRule("minecraft", minecraft),
	}
	RulesXiaohongshu = []Rule{
		DomainSuffixRule("xiaohongshu.com", xiaohongshu),
	}
	RulesReddit = []Rule{
		DomainSuffixRule("reddit.com", reddit),
	}
	RulesZhihu = []Rule{
		DomainSuffixRule("zhihu.com", zhihu),
	}
	RulesQQ = []Rule{
		DomainSuffixRule("qq.com", qq),
	}
  RulesSpotify = []Rule{
    DomainSuffixRule("spotify.com", spotify),
    DomainSuffixRule("spotifycdn.com", spotify),
  }
)
