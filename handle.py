import requests
import json
import time
import hashlib
import httpx
import re

def get_proxy() -> str | None:
    # API链接    后台获取链接地址
    proxyAPI = "http://15527071908.user.xiecaiyun.com/api/proxies?action=getJSON&key=NPFB728011&count=&word=&rand=false&norepeat=false&detail=false&ltime=300&idshow=false"
    proxyusernm = "15527071908"  # 代理帐号
    proxypasswd = "15527071908"  # 代理密码
    url = "https://myip.ipip.net/"
    # 获取IP
    r = requests.get(proxyAPI)
    if r.status_code == 200:
        j = json.loads(r.text)
        if j["success"] and len(j["result"]) > 0:
            p = j["result"][0]
            # name = input();
            proxyurl = f"http://{proxyusernm}:{proxypasswd}@{p['ip']}:{p['port']}"
            t1 = time.time()
            r = requests.get(
                url,
                proxies={"http": proxyurl, "https": proxyurl},
                headers={
                    "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
                    "Accept-Encoding": "gzip, deflate",
                    "Accept-Language": "zh-CN,zh;q=0.9",
                    "Cache-Control": "max-age=0",
                    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36",
                },
            )
            r.encoding = "utf-8"

            t2 = time.time()

            print(r.text)
            return proxyurl
        else:
            print("获取0个代理IP")
            return None
    else:
        print("获取代理失败")
        return None

get_proxy()