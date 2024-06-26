package sub

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
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

type priority int

const (
	Required priority = iota
	Optional
)

func availableOptionalCountries(remote ClashSub) (codes []countries.CountryCode) {
	var countryGroupMap = make(map[countries.CountryCode][]Node)
	for _, node := range remote.Proxies {
		country := extractCountryFromNodeName(node.Name)
		countryGroupMap[country] = append(countryGroupMap[country], node)
	}

	for _, c := range countriesNeeded {
		// url-test or select? Though the need is rare, sometimes I want
		// to use a specific node in the country
		var proxies []string
		for _, server := range countryGroupMap[c.CountryCode] {
			proxies = append(proxies, server.Name)
		}
		if len(proxies) > 0 && c.priority == Optional {
			codes = append(codes, c.CountryCode)
		}
	}
	return
}

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
		Type: "http",
		URL:  "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/" + key + ".txt",
		Path: "./ruleset/" + key + ".yaml",
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
		group := selectGroup(countryGroup(c.CountryCode))
		for _, server := range countryGroupMap[c.CountryCode] {
			group.Proxies = append(group.Proxies, server.Name)
		}
		if len(group.Proxies) == 0 {
			if dealWithEmptyGroup == "placeholder" || c.priority == Required {
				log.Printf("no nodes matched for country %s, put a DIRECT here", c)
				group.Proxies = []string{DIRECT}
			} else if dealWithEmptyGroup == "drop" && c.priority == Optional {
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
		emptyPolicy = "drop"
	}
	countryGroups := groupByCountries(remote, &gr, emptyPolicy)
	gr.Proxies = append(gr.Proxies, an.Name, selfHosted.Name)

	// The order of groups
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
		openAIGroup,
		github,
		microsoft,
		steam,
		youtube,
		amazon,
		proxyConverter,
		ipCheck,
		selfHosted,
		xiaohongshu,
		reddit,
		zhihu,
		qq,
		spotify,
	}

	/* RULES */
	var rules []Rule

	for _, x := range [][]Rule{
		RulesNintendo,
		RulesOpenAI,
		RulesEmby,
		RulesProxyConverterRules,
		RulesSpecial,
		RulesMinecraft,
		RulesXiaohongshu,
		RulesReddit,
		RulesZhihu,
		RulesQQ,
		RulesSpotify,
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

// NewSub creates a default config
// The DNS server by default listens on :7853
// rather than :53, which could be used by
// systemd-resolved.
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
			Enable:       false,
			EnhancedMode: "fake-ip", // 虚拟IP模式
			Listen:       "0.0.0.0:7853",
			Nameserver: []string{
				"8.8.8.8",         // google
			},
			IPv6:        false,
			FakeIPRange: "198.19.0.1/16", // 虚拟IP段
		},
	}
}

var anyTwoCharCode = regexp.MustCompile(".*([A-Za-z]{2})")

// brutal way to recognize geo names in chinese
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
