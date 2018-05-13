package main

import (
	"os"
	"./dingtalk"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

func main(){
	CorpID := os.Getenv("CorpId")
	CorpSecret := os.Getenv("CorpSecret")
	d := dingtalk.New(CorpID, CorpSecret)
	err := d.RefreshAccessToken()
	if err != nil{
		panic("钉钉服务器错误")
	} else {
		e := echo.New()
		e.File("/", "public/index.html")
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())
		e.GET("/user", func(c echo.Context) error {
			code := c.QueryParam("code")
			userIdRes, err := d.UserIDByCode(code)
			if err == nil{
				if userIdRes.ErrCode != 0 {
					return c.JSON(http.StatusOK, userIdRes)
				} else {
					userInfoRes, err := d.UserInfoByUserID(userIdRes.UserID)
					if err == nil{
						if userInfoRes.ErrCode != 0{
							return c.JSON(http.StatusOK, userInfoRes)
						}
						return c.JSON(http.StatusOK, userInfoRes)
					}
					return err
				}
			}
			return err
		})
		e.Logger.Fatal(e.Start(":8080"))
	}
}
