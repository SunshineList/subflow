package api

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sublink/models"
	"sublink/node"
	"sublink/utils"

	"github.com/gin-gonic/gin"
)

// md5加密
func Md5(src string) string {
	m := md5.New()
	m.Write([]byte(src))
	res := hex.EncodeToString(m.Sum(nil))
	return res
}

func setSubscriptionHeader(c *gin.Context, sub models.Subcription) {
	if sub.TotalTraffic > 0 || sub.ExpireTime > 0 {
		// 格式: upload=0; download=xxx; total=xxx; expire=xxx
		headerValue := fmt.Sprintf("upload=0; download=%d; total=%d; expire=%d", sub.UsedTraffic, sub.TotalTraffic, sub.ExpireTime)
		c.Writer.Header().Set("Subscription-Userinfo", headerValue)
	}
}

func GetClient(c *gin.Context) {
	// 获取参数
	token := c.Query("token")
	ClientIndex := c.Query("client") // 客户端标识
	if token == "" {
		log.Println("token为空")
		c.Writer.WriteString("token为空")
		return
	}

	Sub := new(models.Subcription)
	// 获取所有订阅
	list, _ := Sub.List()
	// 查找订阅是否匹配
	for _, sub := range list {
		// 查找token的md5是否匹配并且转换成小写
		if Md5(sub.Name) == strings.ToLower(token) {
			// 设置通用头部（流量信息）
			setSubscriptionHeader(c, sub)

			// 判断是否带客户端参数
			switch ClientIndex {
			case "clash":
				GetClash(c, sub)
				return
			case "surge":
				GetSurge(c, sub)
				return
			case "v2ray":
				GetV2ray(c, sub)
				return
			}
			// 自动识别客户端
			userAgent := strings.ToLower(c.GetHeader("User-Agent"))
			if strings.Contains(userAgent, "clash") {
				GetClash(c, sub)
				return
			} else if strings.Contains(userAgent, "surge") {
				GetSurge(c, sub)
				return
			}
			// 默认下发 V2ray 格式
			GetV2ray(c, sub)
			return
		}
	}
	c.Writer.WriteString("订阅不存在或Token错误")
}

func GetV2ray(c *gin.Context, sub models.Subcription) {
	err := sub.GetSub()
	if err != nil {
		c.Writer.WriteString("读取节点错误")
		return
	}
	baselist := ""
	for _, v := range sub.Nodes {
		switch {
		// 如果包含多条节点
		case strings.Contains(v.Link, ","):
			links := strings.Split(v.Link, ",")
			baselist += strings.Join(links, "\n") + "\n"
			continue
		//如果是订阅转换
		case strings.Contains(v.Link, "http://") || strings.Contains(v.Link, "https://"):
			resp, err := http.Get(v.Link)
			if err != nil {
				log.Println(err)
				continue
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			nodes := utils.Base64Decode(string(body))
			baselist += nodes + "\n"
		// 默认
		default:
			baselist += v.Link + "\n"
		}
	}
	c.Set("subname", sub.Name)
	filename := fmt.Sprintf("%s.txt", sub.Name)
	encodedFilename := url.QueryEscape(filename)
	c.Writer.Header().Set("Content-Disposition", "inline; filename*=utf-8''"+encodedFilename)
	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Writer.WriteString(utils.Base64Encode(baselist))
}

func GetClash(c *gin.Context, sub models.Subcription) {
	err := sub.GetSub()
	if err != nil {
		c.Writer.WriteString("读取节点错误")
		return
	}
	var urls []node.Urls
	for _, v := range sub.Nodes {
		switch {
		// 如果包含多条节点
		case strings.Contains(v.Link, ","):
			links := strings.Split(v.Link, ",")
			for _, link := range links {
				urls = append(urls, node.Urls{
					Url:             link,
					DialerProxyName: strings.TrimSpace(v.DialerProxyName),
				})
			}
			continue
		//如果是订阅转换
		case strings.Contains(v.Link, "http://") || strings.Contains(v.Link, "https://"):
			resp, err := http.Get(v.Link)
			if err != nil {
				log.Println(err)
				continue
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			nodes := utils.Base64Decode(string(body))
			links := strings.Split(nodes, "\n")
			for _, link := range links {
				urls = append(urls, node.Urls{
					Url:             link,
					DialerProxyName: strings.TrimSpace(v.DialerProxyName),
				})
			}
		// 默认
		default:
			urls = append(urls, node.Urls{
				Url:             v.Link,
				DialerProxyName: strings.TrimSpace(v.DialerProxyName),
			})
		}
	}

	var configs utils.SqlConfig
	err = json.Unmarshal([]byte(sub.Config), &configs)
	if err != nil {
		c.Writer.WriteString("配置读取错误")
		return
	}
	DecodeClash, err := node.EncodeClash(urls, configs)
	if err != nil {
		c.Writer.WriteString(err.Error())
		return
	}
	c.Set("subname", sub.Name)
	filename := fmt.Sprintf("%s.yaml", sub.Name)
	encodedFilename := url.QueryEscape(filename)
	c.Writer.Header().Set("Content-Disposition", "inline; filename*=utf-8''"+encodedFilename)
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.WriteString(string(DecodeClash))
}

func GetSurge(c *gin.Context, sub models.Subcription) {
	err := sub.GetSub()
	if err != nil {
		c.Writer.WriteString("读取节点错误")
		return
	}
	urls := []string{}
	for _, v := range sub.Nodes {
		switch {
		// 如果包含多条节点
		case strings.Contains(v.Link, ","):
			links := strings.Split(v.Link, ",")
			urls = append(urls, links...)
			continue
		//如果是订阅转换
		case strings.Contains(v.Link, "http://") || strings.Contains(v.Link, "https://"):
			resp, err := http.Get(v.Link)
			if err != nil {
				log.Println(err)
				continue
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			nodes := utils.Base64Decode(string(body))
			links := strings.Split(nodes, "\n")
			urls = append(urls, links...)
		// 默认
		default:
			urls = append(urls, v.Link)
		}
	}

	var configs utils.SqlConfig
	err = json.Unmarshal([]byte(sub.Config), &configs)
	if err != nil {
		c.Writer.WriteString("配置读取错误")
		return
	}
	DecodeClash, err := node.EncodeSurge(urls, configs)
	if err != nil {
		c.Writer.WriteString(err.Error())
		return
	}
	c.Set("subname", sub.Name)
	filename := fmt.Sprintf("%s.conf", sub.Name)
	encodedFilename := url.QueryEscape(filename)
	c.Writer.Header().Set("Content-Disposition", "inline; filename*=utf-8''"+encodedFilename)
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	host := c.Request.Host
	url := c.Request.URL.String()
	// 如果包含头部更新信息
	if strings.Contains(DecodeClash, "#!MANAGED-CONFIG") {
		c.Writer.WriteString(DecodeClash)
		return
	}
	// 否则就插入头部更新信息
	interval := fmt.Sprintf("#!MANAGED-CONFIG %s interval=86400 strict=false", host+url)
	c.Writer.WriteString(string(interval + "\n" + DecodeClash))
}
