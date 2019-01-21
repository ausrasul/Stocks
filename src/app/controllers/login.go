package controllers

import (
	"github.com/astaxie/beego"
	"github.com/ausrasul/jwt"
	"github.com/markbates/goth/gothic"
	"log"
	"app/providers"
)

type LoginController struct {
	beego.Controller
}

type SecureContent struct {
	beego.Controller
}

func (c *LoginController) mapUrl() {
	q := c.Ctx.Request.URL.Query()
	q.Set("provider", c.Ctx.Input.Param(":provider"))
	c.Ctx.Request.URL.RawQuery = q.Encode()
	return
}

func (c *LoginController) isLoggedIn() bool {
	// Check if the user already have a valid token
	//var j goJwt.GOJWT

	res := c.Ctx.ResponseWriter
	req := c.Ctx.Request
	u, err := jwt.ParseToken(res, req)
	log.Print("UUser:", u)
	log.Print("Checking if a token already exist")
	if err == nil {
		log.Print("Yes a token exist")
		return true
	}
	return false
}

func (c *LoginController) ShowLoginPage() {
	c.mapUrl()

	// Check if the user already have a valid token
	if c.isLoggedIn() {
		log.Print("Yes a token exist")
		c.Ctx.Redirect(301, "/secure")
		return
	}
	c.Data["LoginProvider"] = "Google+"
	c.TplName = "login.tpl"
}

func (c *SecureContent) Get() {
	// Check if the user already have a valid token
	//var j goJwt.GOJWT
	res := c.Ctx.ResponseWriter
	req := c.Ctx.Request
	u, err := jwt.ParseToken(res, req)
	if err != nil {
		log.Print("Invalid token")
		c.Ctx.Redirect(301, "/")
		return
	}
	// print our state string to the console. Ideally, you should verify
	// that it's the same string as the one you set in `setState`
	//fmt.Println("key: ", gothic.)
	log.Print("User: ", u)

	//p := Pool.Get()
	//defer p.Close()

	//test, _ := p.Do("GET", "astaxie")
	//red, err := cache.NewCache("redis", `{"conn":":6379"}`)
	var p providers.Probe
	p.Get()
	beego.Debug(p.Name)
	
	c.Data["Email"] = p.Name
	c.Data["Name"] = p.Id
	p.Id = p.Id + 1
	p.Set()
	c.TplName = "index.tpl"
}
/*
func (c *LoginController) TimAuthenticate() {
	// Check if the user already have a valid token
	
	if c.isLoggedIn() {
		log.Print("Yes a token exist")
		c.Ctx.Redirect(301, "/secure")
		return
	}

	
	res := c.Ctx.ResponseWriter
	req := c.Ctx.Request

	user, err := tim.GetUser(c.Input().Get("username"), c.Input().Get("password"))
	//	gothic.CompleteUserAuth(res, req)
	beego.Debug("Tim User: ", user)
	log.Print("Tim Error: ", err)
	if err != nil {
		log.Print("Cannot authenticate user")
		c.Ctx.Redirect(301, "/")
		return
	}

	userAttributes := make(map[string]interface{})
	userAttributes["Name"] = user["cn"]
	userAttributes["Email"] = user["mail"]

	// Creating a user cookie
	_, err = jwt.CreateToken(userAttributes, res, req)
	if err != nil {
		log.Print(res, "Failed to create token", err)
		c.Ctx.Redirect(301, "/")
		return
	}

	log.Print("Authentication completed")
	c.Ctx.Redirect(301, "/secure")
}
*/
func (c *LoginController) Authenticate() {
	c.mapUrl()
	if c.isLoggedIn() {
		log.Print("Yes a token exist")
		c.Ctx.Redirect(301, "/secure")
		return
	}
	gothic.BeginAuthHandler(c.Ctx.ResponseWriter, c.Ctx.Request)
}

func (c *LoginController) Validate() {
	c.mapUrl()
	if c.isLoggedIn() {
		log.Print("Yes a token exist")
		c.Ctx.Redirect(301, "/secure")
		return
	}

	// Check if the user already have a valid token
	//var j goJwt.GOJWT

	res := c.Ctx.ResponseWriter
	req := c.Ctx.Request

	// print our state string to the console. Ideally, you should verify
	// that it's the same string as the one you set in `setState`
	//fmt.Println("key: ", gothic.)

	user, err := gothic.CompleteUserAuth(res, req)
	log.Print("GothUser: ", user)
	log.Print("GothError: ", err)
	userAttributes := make(map[string]interface{})
	userAttributes["Name"] = user.Name
	userAttributes["Email"] = user.Email
	userAttributes["AccessToken"] = user.AccessToken

	token, err := jwt.CreateToken(userAttributes, res, req)
	if err != nil {
		log.Print(res, err)
		return
	}
	token = token

	log.Print("Authentication completed")
	c.Ctx.Redirect(301, "/secure")

}
