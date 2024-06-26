package proxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type ProxyResult struct {
	Success bool `json:"success"`
	Result  []struct {
		IP   string `json:"ip"`
		Port string `json:"port"`
	} `json:"result"`
}

func GetProxy() *http.Client {
	proxyAPI := "http://15527071908.user.xiecaiyun.com/api/proxies?action=getJSON&key=NPFB728011&count=&word=&rand=false&norepeat=false&detail=false&ltime=300&idshow=false"
	proxyusernm := "15527071908" // 代理帐号
	proxypasswd := "15527071908" // 代理密码
	targetURL := "https://myip.ipip.net/"
	log.Printf("proxyAPI: %s\n", proxyAPI)
	resp, err := http.Get(proxyAPI)
	if err != nil {
		fmt.Println("获取代理失败")
		return nil
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var proxyResult ProxyResult
	json.Unmarshal(body, &proxyResult)

	if proxyResult.Success && len(proxyResult.Result) > 0 {
		p := proxyResult.Result[0]
		proxyURL := fmt.Sprintf("http://%s:%s@%s:%s", proxyusernm, proxypasswd, p.IP, p.Port)

		proxy, _ := url.Parse(proxyURL)
		log.Printf("使用代理: %v\n", proxy)
		client := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
		}
		if client == nil {
			fmt.Println("创建client失败")
			return nil
		}

		req, _ := http.NewRequest("GET", targetURL, nil)
		req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
		req.Header.Add("Accept-Encoding", "gzip, deflate")
		req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
		req.Header.Add("Cache-Control", "max-age=0")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")

		start := time.Now()
		resp, _ := client.Do(req)
		defer resp.Body.Close()

		fmt.Println("Response time:", time.Since(start))

		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))

		return client
	} else {
		fmt.Println("获取0个代理IP")
		return nil
	}
}
