package sub

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/biter777/countries"
	"github.com/enescakir/emoji"
	emojiflag "github.com/jayco/go-emoji-flag"
	"gopkg.in/yaml.v3"
)

const (
	DIRECT = "DIRECT"
	REJECT = "REJECT"
)

var (
	// Make sure every code is provided in upstream sub!!
	countriesNeeded = []countries.CountryCode{
		countries.HK,        // 香港
		countries.TW,        // 台湾
		countries.JP,        // 日本
		countries.Singapore, // 新加坡
		countries.US,        // 美国
		countries.Germany,   // 德国
		countries.BE,        // 比利时
		countries.KG,        // 吉尔吉斯斯坦
		countries.IS,        // 冰岛
		countries.LT,        // 立陶宛
		countries.VN,        // 越南
		countries.MN,        // 蒙古
		countries.AE,        // 阿联酋
		countries.CZ,        // 捷克
		countries.AD,        // 安道尔
		countries.BG,        // 保加利亚
		countries.MD,        // 摩尔多瓦
		countries.RE,        // 法属留尼汪, not available since 20220921
		countries.PA,        // 巴拿马
		countries.MU,        // 毛里求斯, not available since 20220921
	}
)

// servers
var (
	selfHostedServer1HK = Node{
		Name:     "自建服务器1香港",
		Type:     "ss",
		Server:   "34.92.170.135",
		Port:     "8388",
		Cipher:   "chacha20-ietf-poly1305",
		Password: "MP9e4JJqdkdx",
		UDP:      true,
		TFO:      true,
	}
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
	userDefinedNodes = []Node{unlockEMBYServer, proxyConverterServerLocal, selfHostedServer1HK, selfHostedServer2SZ}
)

func urlTestGroup(name string, interval time.Duration, proxies ...string) ProxyGroup {
	return ProxyGroup{
		Name:     name,
		Type:     "url-test",
		URL:      "http://www.gstatic.com/generate_204",
		Interval: int(interval.Seconds()),
		//Tolerance: 0,
		Proxies: proxies,
	}
}

func selectGroup(name string, proxies ...string) ProxyGroup {
	return ProxyGroup{
		Name:    name,
		Type:    "select",
		Proxies: proxies,
	}
}

func grand() ProxyGroup {
	return ProxyGroup{
		Name: emoji.Rocket.String() + "节点选择",
		Type: "select",
	}
}

func allNodes(remote ClashSub) ProxyGroup {
	an := ProxyGroup{
		Name: emoji.GlobeShowingAsiaAustralia.String() + "全部节点",
		Type: "select",
	}
	for _, node := range remote.Proxies {
		an.Proxies = append(an.Proxies, node.Name)
	}
	return an
}

// groups
var (
	selfHosted = selectGroup("Self host servers",
		selfHostedServer1HK.Name,
		selfHostedServer2SZ.Name,
	)
	uncommon = selectGroup("小众节点",
		countryGroup(countries.KG),
		countryGroup(countries.IS),
		countryGroup(countries.LT),
		countryGroup(countries.VN),
		countryGroup(countries.MN),
		countryGroup(countries.AE),
		countryGroup(countries.CZ),
		countryGroup(countries.AD),
		countryGroup(countries.BG),
		countryGroup(countries.MD),
		countryGroup(countries.RE),
		countryGroup(countries.PA),
		countryGroup(countries.MU),
	)
	xiaohongshu = selectGroup("小红书",
		DIRECT,
		uncommon.Name,
		grand().Name,
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
			countries: []countries.CountryCode{countries.TW, countries.HK},
		},
	}
)

// rules
var (
	RulesNintendo = []Rule{
		DomainSuffixRule("stat.ink", switchGroup),
		DomainSuffixRule("nintendo.net", switchGroup),
		DomainSuffixRule("nintendo.com", switchGroup),
		DomainSuffixRule("s3-us-west-2.amazonaws.com", switchGroup),
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
	RulesZhihu = []Rule{
		DomainSuffixRule("zhihu.com", zhihu),
	}
	RulesQQ = []Rule{
		DomainSuffixRule("qq.com", qq),
	}
)

func DomainRule(domain string, group ProxyGroup) Rule {
	return Rule("DOMAIN," + domain + "," + group.Name)
}

func DomainSuffixRule(suffix string, group ProxyGroup) Rule {
	return Rule("DOMAIN-SUFFIX," + suffix + "," + group.Name)
}

func DomainKeyWordRule(keyword string, group ProxyGroup) Rule {
	return Rule("DOMAIN-KEYWORD," + keyword + "," + group.Name)
}

type Node struct {
	Name           string            `yaml:"name,omitempty"` // field needs to be modified
	Type           string            `yaml:"type,omitempty"`
	Server         string            `yaml:"server,omitempty"`
	Port           string            `yaml:"port,omitempty"`
	Cipher         string            `yaml:"cipher,omitempty"`
	Password       string            `yaml:"password,omitempty"`
	Protocol       string            `yaml:"protocol,omitempty"`
	ProtocolParam  string            `yaml:"protocol-param,omitempty"`
	Obfs           string            `yaml:"obfs,omitempty"`
	ObfsParam      string            `yaml:"obfs-param,omitempty"`
	TLS            bool              `yaml:"tls,omitempty"`
	SkipCertVerify bool              `yaml:"skip-cert-verify,omitempty"`
	UDP            bool              `yaml:"udp,omitempty"`
	Plugin         string            `yaml:"plugin,omitempty"`
	PluginOpts     map[string]string `yaml:"plugin-opts,omitempty"`
	TFO            bool              `yaml:"tfo,omitempty"`
}

type ProxyGroup struct {
	Name      string   `yaml:"name"`
	Type      string   `yaml:"type"`
	Use       []string `yaml:"use,omitempty"`
	Proxies   []string `yaml:"proxies"`
	URL       string   `yaml:"url,omitempty"`
	Tolerance int      `yaml:"tolerance,omitempty"`
	Lazy      bool     `yaml:"lazy,omitempty"`
	Interval  int      `yaml:"interval,omitempty"`
}

type RuleProvider struct {
	Type     string `yaml:"type"`
	Behavior string `yaml:"behavior"`
	URL      string `yaml:"url"`
	Path     string `yaml:"path"`
	Interval int    `yaml:"interval"`
}

// ClashSub
// See: https://github.com/Dreamacro/clash/wiki/configuration
type ClashSub struct {
	// Port of HTTP(S) proxy server on the local end
	Port int `yaml:"port,omitempty"`
	// Port of SOCKS5 proxy server on the local end
	SocksPort int `yaml:"socks-port,omitempty"`
	// HTTP(S) and SOCKS4(A)/SOCKS5 server on the same port
	MixedPort int `yaml:"mixed-port,omitempty"`
	// Transparent proxy server port for Linux and macOS (Redirect TCP and TProxy UDP)
	RedirPort int `yaml:"redir-port,omitempty"`
	// authentication of local SOCKS5/HTTP(S) server, like "user1:pass1"
	Authentication []string `yaml:"authentication,omitempty"`
	// Set to true to allow connections to the local-end server from
	// other LAN IP addresses
	AllowLan bool `yaml:"allow-lan"`
	// This is only applicable when `allow-lan` is `true`
	// '*': bind all IP addresses
	// 192.168.122.11: bind a single IPv4 address
	// "[aaaa::a8aa:ff:fe09:57d8]": bind a single IPv6 address
	BindAddress string `yaml:"bind-address,omitempty"`
	// RESTful web API listening address
	ExternalController string `yaml:"external-controller,omitempty"`
	// Secret for the RESTful API (optional)
	// Authenticate by spedifying HTTP header `Authorization: Bearer ${secret}`
	// ALWAYS set a secret if RESTful API is listening on 0.0.0.0
	// Note in CFW: secret and external-controller should be specified in config.yaml, profile doesn't override
	Secret string `yaml:"secret,omitempty"`
	// Clash router working mode
	// rule: rule-based packet routing
	// global: all packets will be forwarded to a single endpoint
	// direct: directly forward the packets to the Internet
	Mode string `yaml:"mode,omitempty"`
	// info / warning / error / debug / silent
	LogLevel string   `yaml:"log-level,omitempty"`
	Profile  *Profile `yaml:"profile,omitempty"`

	Hosts DNSMapping  `yaml:"hosts,omitempty"`
	DNS   *DNSSetting `yaml:"dns,omitempty"`

	Proxies        []Node                 `yaml:"proxies"`
	ProxyProviders map[string]interface{} `yaml:"proxy-providers,omitempty"`
	ProxyGroups    []ProxyGroup           `yaml:"proxy-groups"`

	RuleProviders map[string]RuleProvider `yaml:"rule-providers,omitempty"`
	Rules         []Rule                  `yaml:"rules"`
}

type DNSMapping map[string]string

type Rule string

func (r Rule) String() string {
	return string(r)
}

type Profile struct {
	// Tracing opens tracing exporter API
	// See: Dreamacro/clash-tracing
	Tracing bool `yaml:"tracing"`
	// StoreSelected stores the `select` results in $HOME/.config/clash/.cache
	// set false If you don't want this behavior
	// when two different configurations have groups with the same name, the selected values are shared
	StoreSelected bool `yaml:"store-selected"`
}

type DNSSetting struct {
	Enable       bool     `yaml:"enable,omitempty"`
	EnhancedMode string   `yaml:"enhanced-mode"` // "fake-ip" or "redir-host". See https://github.com/Dreamacro/clash/wiki/configuration#dns
	FakeIPRange  string   `yaml:"fake-ip-range,omitempty"`
	Listen       string   `yaml:"listen,omitempty"` // must be :53. DNS is mostly UDP Port 53, but as time progresses, DNS will rely on TCP Port 53 more heavily.
	Nameserver   []string `yaml:"nameserver,omitempty"`
	IPv6         bool     `yaml:"ipv6"` //  when the false, response to AAAA questions will be empty
}

type RuleSetTreatment func(r RuleProvider) RuleProvider

func extractCountryFromNodeName(node string) countries.CountryCode {
	if anyTwoCharCode.MatchString(node) {
		return countries.ByName(anyTwoCharCode.FindStringSubmatch(node)[1])
	} else if code := chineseMatch(node); code != "" {
		return countries.ByName(code)
	}
	log.Printf("cannot match country code: %s", node)
	return countries.Unknown
}

// Clash规则
// https://github.com/Loyalsoldier/clash-rules
func loyalSoldierClashRules(key string) RuleProvider {
	return RuleProvider{
		Type:     "http",
		URL:      "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/" + key + ".txt",
		Path:     "./ruleset/" + key + ".yaml",
		Interval: 86400,
	}
}

func newHttpClient(proxyEnv bool) *http.Client {
	var proxyFunc func(*http.Request) (*url.URL, error)
	if proxyEnv {
		proxyFunc = http.ProxyFromEnvironment
	}
	client := http.Client{
		Transport: &http.Transport{
			Proxy: proxyFunc,
		},
	}
	return &client
}

var proxyClient = newHttpClient(true)

// Clash规则碎片
// https://github.com/ACL4SSR/ACL4SSR
// 读取ACL4SSR/ACL4SSR仓库的规则碎片，转换为匹配策略key的Rule
// 对github内容的访问本身需要代理
func acl4ssrClashRules(key string, emojiPrefix emoji.Emoji) []Rule {
	res, err := proxyClient.Get("https://raw.githubusercontent.com/ACL4SSR/ACL4SSR/master/Clash/Ruleset/" + key + ".list")
	if err != nil {
		return nil
	}
	scanner := bufio.NewScanner(res.Body)
	var rules []Rule
	defer res.Body.Close()
	groupName := emojiPrefix.String() + key
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "#"):
		case strings.Contains(line, "USER-AGENT"): // 不支持
		case strings.Contains(line, "no-resolve"):
			items := strings.Split(line, ",")
			rules = append(rules, Rule(strings.Join([]string{items[0], items[1], groupName, items[2]}, ",")))
		default:
			rules = append(rules, Rule(line+","+groupName))
		}
	}
	return rules
}

// See: https://github.com/Semporia/ClashX-Pro
// Must provide a valid key, or rules are corrupted
func semporiaClashXRules(key string) RuleProvider {
	return RuleProvider{
		Type:     "http",
		Behavior: "classical",
		Path:     "./ruleset/" + strings.ToLower(key) + ".yaml",
		URL:      "https://raw.githubusercontent.com/Semporia/Clash/master/Rule/" + key + ".yaml",
		Interval: 3600,
	}
}

func Domain(r RuleProvider) RuleProvider {
	nr := r
	nr.Behavior = "domain"
	return nr
}

func IPCIRD(r RuleProvider) RuleProvider {
	nr := r
	nr.Behavior = "ipcidr"
	return nr
}

func countryGroup(c countries.CountryCode) string {
	return c.Emoji() + c.Alpha2()
}

func mediaGroup(key string, emojiPrefix emoji.Emoji, countries ...countries.CountryCode) ProxyGroup {
	pg := ProxyGroup{
		Name: emojiPrefix.String() + key,
		Type: "select",
		Proxies: []string{
			DIRECT,
		},
	}
	for i := range countries {
		pg.Proxies = append(pg.Proxies, countryGroup(countries[i]))
	}
	return pg
}

func decodeClashConfig(res *http.Response) (ClashSub, error) {
	var remote ClashSub
	err := yaml.NewDecoder(res.Body).Decode(&remote)
	return remote, err
}

type Processor func(sub *ClashSub)

// AddRuleIPCIDR should be before the last MATCH rule, if any
func AddRuleIPCIDR(cidr, target string) Processor {
	return func(sub *ClashSub) {
		if len(sub.Rules) > 0 && strings.HasPrefix(sub.Rules[len(sub.Rules)-1].String(), "MATCH") {
			match := sub.Rules[len(sub.Rules)-1]
			sub.Rules[len(sub.Rules)-1] = Rule("IP-CIDR," + cidr + "," + target)
			sub.Rules = append(sub.Rules, match)
		} else {
			sub.Rules = append(sub.Rules, Rule("IP-CIDR,"+cidr+","+target))
		}
	}
}

func SetExternalController(address string) Processor {
	return func(sub *ClashSub) {
		sub.ExternalController = address
	}
}

// AddHosts 增加自定义的DNS规则,可以使用通配符
// 如{"*.example.com": "127.0.0.1"}会匹配abc.example.com
func AddHosts(records DNSMapping) Processor {
	return func(sub *ClashSub) {
		if sub.Hosts == nil {
			sub.Hosts = make(DNSMapping)
		}
		for k, v := range records {
			sub.Hosts[k] = v
		}
	}
}

func Pass(remote ClashSub, out io.Writer, proc ...Processor) error {
	config := NewSub()
	config.Proxies = remote.Proxies
	config.ProxyProviders = remote.ProxyProviders
	config.ProxyGroups = remote.ProxyGroups
	config.RuleProviders = remote.RuleProviders
	config.Rules = remote.Rules
	for i := range proc {
		proc[i](&config)
	}
	encoder := yaml.NewEncoder(out)
	encoder.SetIndent(2)
	return encoder.Encode(&config)
}

// 按国家分组
func groupByCountries(remote ClashSub, gr *ProxyGroup, dealWithEmptyGroup string) (countryGroups []ProxyGroup) {
	var countryGroupMap = make(map[countries.CountryCode][]Node)
	for _, node := range remote.Proxies {
		country := extractCountryFromNodeName(node.Name)
		countryGroupMap[country] = append(countryGroupMap[country], node)
	}

	for _, c := range countriesNeeded {
		// url-test or select? Though the need is rare, sometimes I want
		// to use a specific node in the country
		group := selectGroup(countryGroup(c))
		for _, server := range countryGroupMap[c] {
			group.Proxies = append(group.Proxies, server.Name)
		}
		if len(group.Proxies) == 0 {
			if dealWithEmptyGroup == "placeholder" {
				log.Printf("no nodes matched for country %s, put a DIRECT here", c)
				group.Proxies = []string{DIRECT}
			} else if dealWithEmptyGroup == "drop" {
				continue
			}
		}
		countryGroups = append(countryGroups, group)
		gr.Proxies = append(gr.Proxies, group.Name)
	}
	return
}

func Rewrite(remote ClashSub, out io.Writer, emptyPolicy string, ruleStreamMedia bool, proc ...Processor) error {
	// since 20220921, TAG node names alreay have emoji prefix
	// no prefixEmoji needed

	/* PROXY GROUPS */
	gr := grand()
	an := allNodes(remote)
	if emptyPolicy == "" {
		emptyPolicy = "placeholder"
	}
	countryGroups := groupByCountries(remote, &gr, emptyPolicy)
	gr.Proxies = append(gr.Proxies, an.Name, selfHosted.Name)

	// 自定义组
	proxyGroups := []ProxyGroup{
		gr,
		an,
		reject,
		rest,
		minecraft,
		apple,
		embyUnlock,
		embyTagNewFlavor,
		telegram,
		switchGroup,
		github,
		microsoft,
		steam,
		youtube,
		amazon,
		proxyConverter,
		ipCheck,
		selfHosted,
		xiaohongshu,
		zhihu,
		qq,
		uncommon,
	}

	/* RULES */
	var rules []Rule

	for _, x := range [][]Rule{
		RulesNintendo,
		RulesEmby,
		RulesProxyConverterRules,
		RulesSpecial,
		RulesMinecraft,
		RulesXiaohongshu,
		RulesZhihu,
		RulesQQ,
	} {
		rules = append(rules, x...)
	}

	/* 流媒体 */
	// could be slow if server cannot access Github properly
	if ruleStreamMedia {
		var (
			streamingEmoji = emoji.Television
		)
		for _, x := range streamMedia {
			proxyGroups = append(proxyGroups, mediaGroup(x.mediaKey, streamingEmoji, x.countries...))
			moreRules := acl4ssrClashRules(x.mediaKey, streamingEmoji)
			rules = append(rules, moreRules...)
		}
	}

	proxyGroups = append(proxyGroups, countryGroups...)

	/* RULE PROVIDERS */
	ruleProviders := map[string]RuleProvider{
		"microsoft": semporiaClashXRules("Microsoft"),
		"github":    semporiaClashXRules("GitHub"),
		"steam":     semporiaClashXRules("Steam"),
		"youtube":   semporiaClashXRules("YouTube"),
		"amazon":    semporiaClashXRules("Amazon"),
	}
	rules = append(rules,
		Rule("RULE-SET,microsoft,"+microsoft.Name),
		Rule("RULE-SET,github,"+github.Name),
		Rule("RULE-SET,steam,"+steam.Name),
		Rule("RULE-SET,youtube,"+youtube.Name),
		Rule("RULE-SET,amazon,"+amazon.Name),
	)

	for _, x := range []struct {
		Treatment RuleSetTreatment
		RuleSet   string
		Chain     string // ignore this ruleset if left empty, add it though
	}{
		{Treatment: Domain, RuleSet: "reject", Chain: reject.Name},
		{Treatment: Domain, RuleSet: "icloud", Chain: DIRECT},
		{Treatment: Domain, RuleSet: "apple", Chain: apple.Name},
		{Treatment: Domain, RuleSet: "google", Chain: DIRECT},
		{Treatment: Domain, RuleSet: "proxy", Chain: grand().Name},
		{Treatment: Domain, RuleSet: "direct", Chain: DIRECT},
		{Treatment: Domain, RuleSet: "private", Chain: DIRECT},
		{Treatment: Domain, RuleSet: "gfw"},
		{Treatment: Domain, RuleSet: "greatfire"},
		{Treatment: Domain, RuleSet: "tld-not-cn"},

		{Treatment: IPCIRD, RuleSet: "telegramcidr", Chain: telegram.Name},
		{Treatment: IPCIRD, RuleSet: "cncidr"},
		{Treatment: IPCIRD, RuleSet: "lancidr"},
	} {
		ruleProviders[x.RuleSet] = x.Treatment(loyalSoldierClashRules(x.RuleSet))
		// 白名单模式
		if x.Chain != "" {
			rules = append(rules, Rule("RULE-SET"+","+x.RuleSet+","+x.Chain))
		}
	}
	// 白名单模式
	rules = append(rules, "GEOIP,CN,DIRECT")

	config := NewSub()
	config.ProxyGroups = proxyGroups
	config.Rules = rules
	config.Proxies = append(remote.Proxies, userDefinedNodes...)
	config.RuleProviders = ruleProviders

	for i := range proc {
		proc[i](&config)
	}

	// 漏网之鱼
	config.Rules = append(config.Rules, Rule("MATCH,"+rest.Name))
	encoder := yaml.NewEncoder(out)
	encoder.SetIndent(2)
	return encoder.Encode(&config)
}

func NewSub() ClashSub {
	return ClashSub{
		MixedPort:          7890,
		ExternalController: "0.0.0.0:9090",
		AllowLan:           true,
		Mode:               "rule",
		LogLevel:           "info",
		Profile: &Profile{
			Tracing:       true,
			StoreSelected: true,
		},
		DNS: &DNSSetting{
			Enable:       true,
			EnhancedMode: "fake-ip", // 虚拟IP模式
			Listen:       "0.0.0.0:53",
			Nameserver: []string{
				// public DNS servers
				"223.5.5.5",       // alibaba
				"8.8.8.8",         // google
				"114.114.114.114", // 114
			},
			IPv6:        true,
			FakeIPRange: "198.19.0.1/16", // 虚拟IP段
		},
	}
}

var anyTwoCharCode = regexp.MustCompile(".*([A-Za-z]{2})")
var chineseCodeMap = map[string]string{
	"台湾": "TW",
	"日本": "JP",
	"狮城": "SG",
	"美国": "US",
	"香港": "HK",
}

func chineseMatch(name string) string {
	for ch, code := range chineseCodeMap {
		if strings.Contains(name, ch) {
			return code
		}
	}
	return ""
}

func prefixEmoji(sub *ClashSub) {
	for i := range sub.Proxies {
		name := sub.Proxies[i].Name
		var countryCode string
		if anyTwoCharCode.MatchString(name) {
			countryCode = anyTwoCharCode.FindStringSubmatch(name)[1]
		} else {
			log.Printf("cannot match country code: %s", name)
			continue
		}
		sub.Proxies[i].Name = emojiflag.GetFlag(countryCode) + name
	}
}
