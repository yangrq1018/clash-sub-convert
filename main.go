package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yangrq1018/clash-sub-convert/sub"
	"github.com/yangrq1018/clash-sub-convert/util"
	"gopkg.in/yaml.v3"
)

var fileProcessors []sub.Processor

func copyFileProcessors() []sub.Processor {
	var cp = make([]sub.Processor, len(fileProcessors))
	copy(cp, fileProcessors)
	return cp
}

func splitKeyValue(s string, delim string) (k, v string) {
	items := strings.Split(s, delim)
	if len(items) != 2 {
		return
	}
	k, v = items[0], items[1]
	k, v = strings.Trim(k, " "), strings.Trim(v, " ")
	return
}

func initProcessorsFromFile() error {
	var processors = make([]sub.Processor, 0)
	f, err := os.Open(util.FirstString(os.Getenv("CONFIG_FILE"), "config.yaml"))
	if err != nil {
		return err
	}
	viper.SetConfigType("yaml")
	err = viper.ReadConfig(f)
	if err != nil {
		return err
	}
	for k, v := range viper.GetStringMapString("hosts") {
		log.Infof("Add processor: %s -> %s", k, v)
		processors = append(processors, sub.AddHosts(sub.DNSMapping{k: v}))
	}

	for _, item := range viper.GetStringSlice("rules.IPCIDR") {
		k, v := splitKeyValue(item, ":")
		log.Infof("Add processor: %s -> %s", k, v)
		processors = append(processors, sub.AddRuleIPCIDR(k, v))
	}
	fileProcessors = processors
	return nil
}

func fetchUpStream(link string, ctx echo.Context) (res *http.Response, err error) {
	res, err = http.Get(link)
	if err != nil {
		_ = ctx.String(http.StatusInternalServerError, err.Error())
		return nil, err
	}
	if res.StatusCode != 200 {
		// cloning the response like this better?
		ctx.Response().Writer.WriteHeader(res.StatusCode)
		_, _ = io.Copy(ctx.Response().Writer, res.Body)
		return nil, fmt.Errorf("upstream server returns error: %d", res.StatusCode)
	}
	return
}

func setHeader(res *http.Response, writer http.ResponseWriter, attachment bool) {
	// this goes before writer.Write, or the content has no effect
	for k, v := range res.Header {
		// 复制流量使用情况的请求头
		switch k {
		case "Subscription-Userinfo", "Content-Type":
			writer.Header().Set(k, v[0])
		case "Content-Disposition":
			if attachment {
				writer.Header().Set(k, "attachment;filename="+res.Request.Host)
			} else {
				writer.Header().Set(k, "inline;filename="+res.Request.Host)
			}
		}
	}
	writer.Header().Set("Content-Type", "application/x-yaml")
}

func main() {
	err := initProcessorsFromFile()
	if err != nil {
		log.Fatal(err)
	}
	e := echo.New()
	e.Use(middleware.Logger())
	e.HideBanner = true
	e.GET("/", func(c echo.Context) (err error) {
		subLink := c.QueryParam("sub")
		subType := c.QueryParam("type")
		if subLink == "" {
			return c.String(http.StatusBadRequest, "param \"sub\" missing")
		}
		res, err := fetchUpStream(subLink, c)
		if err != nil {
			return err
		}
		if res.Body != nil {
			defer func() {
				_ = res.Body.Close()
			}()
		}
		setHeader(res, c.Response().Writer, false)
		var remote sub.ClashSub
		switch subType {
		case "clash", "":
			err = yaml.NewDecoder(res.Body).Decode(&remote)
			if err != nil {
				return err
			}
		case "ss":
			var proxies []sub.SSItem
			proxies, err = sub.DecodeSS(res)
			if err != nil {
				return err
			}
			for _, proxy := range proxies {
				node := sub.Node{
					Name:     proxy.Name,
					Type:     "ss",
					Server:   proxy.Server,
					Port:     proxy.Port,
					Cipher:   proxy.Method,
					Password: proxy.Password,
					TFO:      true, // where are these coming from?
					UDP:      true,
				}
				switch proxy.Plugins["plugin"] {
				case "simple-obfs":
					node.Plugin = "obfs"
					node.PluginOpts = map[string]string{
						"mode": proxy.Plugins["obfs"],
						"host": proxy.Plugins["obfs-host"],
					}
				}
				remote.Proxies = append(remote.Proxies, node)
			}
		case "ssr":
			var proxies []sub.SSRItem
			proxies, err = sub.DecodeSSR(res)
			if err != nil {
				return err
			}
			for _, proxy := range proxies {
				remote.Proxies = append(remote.Proxies, sub.Node{
					Name:          proxy.Remarks,
					Type:          "ssr",
					Server:        proxy.Server,
					Port:          proxy.ServerPort,
					Cipher:        proxy.Method,
					Password:      proxy.Password,
					Protocol:      proxy.Protocol,
					ProtocolParam: proxy.ProtocolParam,
					Obfs:          proxy.OBFS,
					ObfsParam:     proxy.OBFSParam,
					UDP:           true,
				})
			}
		}
		body := bytes.NewBuffer(nil)

		processors := copyFileProcessors()
		if externalController := c.QueryParam("controller"); externalController != "" {
			processors = append(processors, sub.SetExternalController(externalController))
		}

		if c.QueryParam("pass") == "true" {
			err = sub.Pass(remote, body, processors...)
		} else {
			err = sub.Rewrite(remote, body, processors...)
		}
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		setHeader(res, c.Response(), false)
		if err = c.Stream(200, "application/x-yaml", body); err != nil {
			return err
		}
		log.WithFields(log.Fields{
			"sub":    subLink,
			"remote": c.Request().RemoteAddr,
		}).Infof("fetch remote sub")
		return
	})
	port := util.FirstString(os.Getenv("PORT"), "8080")
	log.Infof("binding to port %s", port)
	if err = e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}
