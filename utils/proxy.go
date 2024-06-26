package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ProxyResult struct {
	Success bool `json:"success"`
	Result  []struct {
		IP   string `json:"ip"`
		Port int64  `json:"port"`
	} `json:"result"`
}

func GetProxy() (*http.Client, error) {
	proxyAPI := "http://15527071908.user.xiecaiyun.com/api/proxies?action=getJSON&key=NPFB728011&count=&word=&rand=false&norepeat=false&detail=false&ltime=300&idshow=false"
	proxyusernm := "15527071908" // 代理帐号
	proxypasswd := "15527071908" // 代理密码

	resp, err := http.Get(proxyAPI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var proxyResult ProxyResult
	json.Unmarshal(body, &proxyResult)
	if proxyResult.Success && len(proxyResult.Result) > 0 {
		p := proxyResult.Result[0]
		proxyURL := fmt.Sprintf("http://%s:%s@%s:%s", proxyusernm, proxypasswd, p.IP, p.Port)
		proxy, _ := url.Parse(proxyURL)
		client := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
		}
		return client, nil
	} else {
		return nil, fmt.Errorf("Failed to get proxy")
	}
}
