package dingtalk

import (
	"net/http"
	"time"
	"net/url"
	"fmt"
	"crypto/sha1"
	"encoding/json"
)

const (
	OAPIURL = "https://oapi.dingtalk.com/"
)

type (
	DTClient struct {
		CorpID string
		CorpSecret string
		AgentID string
		AccessToken string
		AccessTokenCache Cache
		TicketCache Cache
		HttpClient *http.Client
	}
	OAPIResponse struct {
		ErrCode int `json:"errcode"`
		ErrMsg string `json:"errmsg"`
	}
	UserIDResponse struct {
		OAPIResponse
		UserID string `json:"userid"`
		DeviceID string `json:"deviceId"`
		IsSys bool `json:"is_sys"`
		SysLevel int `json:"sys_level"`
	}
	UserInfoResponse struct {
		OAPIResponse
		UserID          string `json:"userid"`
		OpenID          string `json:"openid"`
		Name            string `json:"name"`
		Tel             string
		WorkPlace       string
		Remark          string
		Mobile          string
		Email           string
		OrgEmail        string
		Active          bool
		IsAdmin         bool
		IsBoos          bool
		DingID          string
		UnionID         string
		IsHide          bool
		Department      []int
		Position        string
		Avatar          string `json:"avatar"`
		Jobnumber       string
		IsSenior        bool
		StateCode       string
		OrderInDepts    string
		IsLeaderInDepts string
		Extattr         interface{}
		Roles           []Roles
	}
	Roles struct {
		ID        int `json:"id"`
		Name      string
		GroupName string
	}
	AccessTokenResponse struct {
		OAPIResponse
		AccessToken string `json:"access_token"`
		Expires int `json:"expires_in"`
		Created int64
	}
	TicketResponse struct {
		OAPIResponse
		Ticket string `json:"ticket"`
		Expires int `json:"expires_in"`
		Created int64
	}
)

func (a *AccessTokenResponse) CreatedAt() int64{
	return a.Created
}

func (a *AccessTokenResponse) ExpiresIn() int{
	return  a.Expires
}

func (t *TicketResponse) CreatedAt() int64{
	return t.Created
}

func (t *TicketResponse) ExpiresIn() int{
	return t.Expires
}

func New(CorpID string, CorpSecret string) *DTClient{
	return &DTClient{
		CorpID:CorpID,
		CorpSecret:CorpSecret,
		AccessTokenCache:NewFileCache(".access_token"),
		TicketCache:NewFileCache(".ticket"),
		HttpClient:&http.Client{
			Timeout:10 *time.Second,
		},
	}
}

func (d *DTClient) RefreshAccessToken() error{
	var data AccessTokenResponse
	err := d.AccessTokenCache.Get(&data)
	if err == nil{
		d.AccessToken = data.AccessToken
		return  nil
	}
	params := url.Values{}
	params.Add("corpid", d.CorpID)
	params.Add("corpsecret", d.CorpSecret)
	err = d.httpRequest("gettoken", params, nil, &data)
	if err == nil{
		d.AccessToken = data.AccessToken
		data.Expires = data.Expires | 7200
		data.Created = time.Now().Unix()
		err = d.AccessTokenCache.Set(&data)
	}
	return  err
}

func (d *DTClient) GetJSAPITicket()(string, error){
	var data TicketResponse
	err := d.TicketCache.Get(&data)
	if err == nil{
		return data.Ticket, err
	}
	err = d.httpRequest("get_jsapi_ticket", nil, nil, &data)
	if err == nil{
		ticket := data.Ticket
		d.TicketCache.Set(&data)
		return ticket, err
	}
	return "", err
}

func (d *DTClient) GetConfig(nonceStr, timestamp, iurl string ) string{
	ticket, _ := d.GetJSAPITicket()
	config := map[string]string{
		"url": iurl,
		"nonceStr": nonceStr,
		"agentId": d.AgentID,
		"timeStamp": timestamp,
		"corpId": d.CorpID,
		"ticket": ticket,
		"signature": sign(ticket, nonceStr, timestamp, iurl),
	}
	bytes, _ := json.Marshal(config)
	return string(bytes)
}

func sign(ticket, nonceStr, timeStamp, iurl string) string{
	s := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%s&url=%s", ticket, nonceStr, timeStamp, iurl)
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func (d *DTClient) UserIDByCode(code string) (UserIDResponse, error){
	var data UserIDResponse
	params := url.Values{}
	params.Add("code", code)
	err := d.httpRequest("user/getuserinfo", params, nil, &data)
	return  data, err
}

func (d *DTClient) UserInfoByUserID(userID string)(UserInfoResponse, error){
	var data UserInfoResponse
	params := url.Values{}
	params.Add("userid", userID)
	err := d.httpRequest("user/get", params, nil, &data)
	return  data, err
}