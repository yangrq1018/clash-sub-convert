package sub

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ssr订阅节点格式
var (
	// group[1]是firstPattern格式数据的base64 encode
	ssrPattern   = regexp.MustCompile("ssr://(.*)")
	firstPattern = regexp.MustCompile(`(.*):(.*):(.*):(.*):(.*):(.*)/\?"`)
)

// GetRemainingDataSSR 获得机场剩余流量, SSR格式配置文件
func GetRemainingDataSSR(subsLink string) (string, error) {
	res, err := http.Get(subsLink)
	if err != nil {
		return "", err
	}
	items, err := DecodeSSR(res)
	if err != nil {
		return "", err
	}
	return getRemainingData(items), nil
}

type ClashDataUsage struct {
	Upload   int
	Download int
	Total    int
	Expire   time.Time
}

// ParseResponse 获得机场剩余流量, Clash格式配置文件
func (c *ClashDataUsage) ParseResponse(r *http.Response) error {
	userInfo := r.Header.Get("subscription-userinfo")
	var upload, download, total int
	match := regexp.MustCompile(`upload=(\d+); download=(\d+); total=(\d+); expire=(\d+)`).FindStringSubmatch(userInfo)
	if len(match) != 5 {
		return fmt.Errorf("invalid subscription-userinfo header")
	}
	upload, _ = strconv.Atoi(match[1])
	download, _ = strconv.Atoi(match[2])
	total, _ = strconv.Atoi(match[3])
	expireUnix, _ := strconv.ParseInt(match[4], 10, 64)
	c.Upload = upload
	c.Download = download
	c.Total = total
	c.Expire = time.Unix(expireUnix, 0)
	return nil
}

func (c *ClashDataUsage) String() string {
	text := fmt.Sprintf(`已用：%.1fGB
配额：%.1fGB
到期：%s
`,
		bytesToGB(c.Upload)+bytesToGB(c.Download), bytesToGB(c.Total), c.Expire.Format("2006-01-02"),
	)
	return text
}

func GetRemainingDataClash(clashLink string) (*ClashDataUsage, error) {
	res, err := http.Get(clashLink)
	if err != nil {
		return nil, err
	}
	res.Body.Close()
	var c ClashDataUsage
	err = c.ParseResponse(res)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func bytesToGB(b int) float64 {
	return float64(b) / float64(1024*1024*1024)
}

func Search(items []SSRItem, keyword string) []SSRItem {
	var out []SSRItem
	for _, item := range items {
		if strings.Contains(item.Remarks, keyword) {
			out = append(out, item)
		}
	}
	return out
}

// DecodeSSR 解析订阅地址ssr，返回解析后的节点切片
// 解析需要注意的地方是两次base64 decode，一次对订阅地址发回的
// 字符串使用base64.StdEncoding解析，获得多行文本，每行格式为ssr://<base64_data>
// 然后对<base64_data>使用base64.RawURLEncoding解析
// RawURLEncoding主要用于编码URL中的数据
// 由于网络传输，在编码中将"+“和”/“进行了替换, 所以在解码的时候要将这两个字符还原回去
func DecodeSSR(res *http.Response) ([]SSRItem, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	bodyDecoded, err := decodeBase64(body, base64.StdEncoding)
	if err != nil {
		return nil, err
	}

	lines := bytes.Split(bodyDecoded, []byte{'\n'})
	var subs []SSRItem
	for _, line := range lines {
		if ssrPattern.Match(line) {
			item, err := ssrToClash(line)
			if err != nil {
				return nil, err
			}
			subs = append(subs, *item)
		}
	}
	return subs, nil
}

func ssrToClash(line []byte) (*SSRItem, error) {
	encodedURL := ssrPattern.FindSubmatch(line)[1]
	decodedURL, err := decodeBase64(encodedURL, base64.RawURLEncoding) // Raw - ignore padding
	if err != nil {
		return nil, err
	}

	match := firstPattern.FindSubmatch(decodedURL)
	// 密码也要解码
	password, err := decodeBase64(match[6], base64.RawURLEncoding)
	if err != nil {
		return nil, err
	}
	item := SSRItem{
		Server:     string(match[1]),
		ServerPort: string(match[2]),
		Protocol:   string(match[3]),
		Method:     string(match[4]),
		OBFS:       string(match[5]),
		Password:   string(password),
	}
	err = getParam(string(decodedURL), &item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

type SSRItem struct {
	Server        string
	ServerPort    string
	Protocol      string
	Method        string
	OBFS          string
	Password      string
	OBFSParam     string
	ProtocolParam string
	Remarks       string
	Group         string
}

func decodeBase64(s []byte, encoding *base64.Encoding) ([]byte, error) {
	var ds = make([]byte, encoding.DecodedLen(len(s)))
	_, err := encoding.Decode(ds, s)
	return ds, err
}

// apply decodeBase64 to string, panic on error
func mustDecodeBase64(s string, encoding *base64.Encoding) string {
	b, err := decodeBase64([]byte(s), encoding)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// 全部params都是base64编码
func getParam(u string, item *SSRItem) error {
	s, err := url.Parse(u)
	if err != nil {
		return err
	}
	q := s.Query()
	item.OBFSParam = mustDecodeBase64(q.Get("obfsparam"), base64.RawURLEncoding)
	item.ProtocolParam = mustDecodeBase64(q.Get("protoparam"), base64.RawURLEncoding)
	item.Remarks = mustDecodeBase64(q.Get("remarks"), base64.RawURLEncoding)
	item.Group = mustDecodeBase64(q.Get("group"), base64.RawURLEncoding)
	return nil
}

func getRemainingData(items []SSRItem) string {
	var text strings.Builder
	for _, item := range items {
		if item.Server == "www.google.com" && !strings.Contains(item.Remarks, "官网") {
			text.WriteString(item.Remarks + "\n")
		}
	}
	return text.String()
}
