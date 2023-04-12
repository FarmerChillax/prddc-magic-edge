package gfw

import (
	"bytes"
	"encoding/base64"
	"io"
	"log"
	"net/http"

	"github.com/pmezard/adblock/adblock"
)

// 快速生成 Interface: gfw *GFWImpl GFW
type GFW interface {
	Exist(url, domain string) (exist bool)
}

type GFWImpl struct {
	matcher adblock.RuleMatcher
	GFWList []byte
	// url     string
}

func New() (GFW, error) {
	url := "https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt"

	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 对响应体进行 base64 解码
	base64Reader := base64.NewDecoder(base64.RawStdEncoding, resp.Body)
	gfwList, err := io.ReadAll(base64Reader)
	if err != nil {
		return nil, err
	}
	// 初始化 adp 匹配器
	matcher := adblock.NewMatcher()
	rules, err := adblock.ParseRules(bytes.NewReader(gfwList))
	if err != nil {
		log.Fatalf("adblock.ParseRules err: %v", err)
		return nil, err
	}

	for index, rule := range rules {
		err = matcher.AddRule(rule, index)
		if err != nil {
			log.Printf("matcher.AddRule err: %v", err)
			return nil, err
		}
	}

	return &GFWImpl{
		GFWList: gfwList,
		matcher: *matcher,
	}, nil
}

func (gfw *GFWImpl) Exist(url, domain string) (exist bool) {
	match, ruleId, err := gfw.matcher.Match(&adblock.Request{
		URL:    url,
		Domain: domain})
	if err != nil {
		log.Printf("matcher.Match err: %v; ruleId: %d", err, ruleId)
		return false
	}
	return match
}
