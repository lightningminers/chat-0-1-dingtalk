package dingtalk

import (
	"net/http"
	"time"
	"net/url"
)

const (
	OAPIURL = "https://oapi.dingtalk.com/"
)

type DTClient struct {
	CorpID string
	CorpSecret string
	AccessToken string
	AccessTokenCache Cache
	HttpClient *http.Client
}

type OAPIResponse struct {
	ErrCode int `json:"errcode"`
	ErrMsg string `json:"errmsg"`
}

type UserIDResponse struct {
	OAPIResponse
	UserID string `json:"userid"`
	DeviceID string `json:"deviceId"`
	IsSys bool `json:"is_sys"`
	SysLevel int `json:"sys_level"`
}

type AccessTokenResponse struct {
	OAPIResponse
	AccessToken string `json:"access_token"`
	Expires int `json:"expires_in"`
	Created int64
}

func (a *AccessTokenResponse) CreatedAt() int64{
	return a.Created
}

func (a *AccessTokenResponse) ExpiresIn() int{
	return  a.Expires
}

func New(CorpID string, CorpSecret string) *DTClient{
	return &DTClient{
		CorpID:CorpID,
		CorpSecret:CorpSecret,
		AccessTokenCache:NewFileCache(".access_token"),
		HttpClient:&http.Client{
			Timeout:10 *time.Second,
		},
	}
}

func (d *DTClient) RefreshAccessToken() error{
	var data AccessTokenResponse
	err := d.AccessTokenCache.Get(&data)
	if err == nil{
		d.AccessToken = d.AccessToken
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

func (d *DTClient) UserIDByCode(code string) (UserIDResponse, error){
	var data UserIDResponse
	params := url.Values{}
	params.Add("code", code)
	err := d.httpRequest("user/getuserinfo", params, nil, &data)
	return  data, err
}

func (d *DTClient) UserInfoByUserID(userID string){

}