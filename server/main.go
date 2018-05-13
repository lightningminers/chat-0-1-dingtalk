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
	refreshAT := func(next echo.HandlerFunc) echo.HandlerFunc{
		return func(c echo.Context) error {
			err := d.RefreshAccessToken()
			if err != nil{
					return c.JSON(http.StatusOK, map[string]string{"error": "refresh access_token error"})
			} else {
					next(c)
			}
			return err
		}
	}
	if err != nil{
		panic("钉钉服务器错误")
	} else {
		e := echo.New()
		e.File("/", "public/index.html")
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())
		e.GET("/get_config", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{
				"CorpID": CorpID,
			})
		})
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
		}, refreshAT)
		e.Logger.Fatal(e.Start(":8080"))
	}
}
