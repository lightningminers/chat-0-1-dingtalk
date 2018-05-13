package dingtalk

import (
	"net/url"
	"net/http"
	"errors"
	"io/ioutil"
	"encoding/json"
)

func (d *DTClient) httpRequest(path string, params url.Values, requestData interface{}, responseData interface{}) error{
	var request *http.Request
	if d.AccessToken != ""{
		if params == nil{
			params = url.Values{}
		}
		if params.Get("access_token") == ""{
			params.Set("access_token", d.AccessToken)
		}
	}
	requestURL := OAPIURL + path + "?" + params.Encode()
	client := d.HttpClient
	if requestData == nil{
		request, _ = http.NewRequest("GET", requestURL, nil)
	}
	res, err := client.Do(request)
	if err != nil{
		return  err
	}
	if res.StatusCode != 200{
		return  errors.New("服务器错误：" + res.Status)
	}
	defer res.Body.Close()
	if content, err := ioutil.ReadAll(res.Body); err == nil{
		json.Unmarshal(content, responseData)
	}
	return err
}
