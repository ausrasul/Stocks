package controllers

import (
	"app/providers"
	"github.com/astaxie/beego"
	"github.com/ausrasul/jwt"
	"log"
)

var NOAUTH bool // if true, anyone can log in with no account checking.

func init() {
	NOAUTH = false
	noAuth, err := beego.AppConfig.Int("NoAuth")
	if err == nil && noAuth == 1 {
		NOAUTH = true
	}
}

type MainController struct {
	beego.Controller
}

func (c *MainController) AuthenticateGet() bool {
	if NOAUTH {
		c.Data["Admin"] = true
		return true
	}
	res := c.Ctx.ResponseWriter
	req := c.Ctx.Request
	u, err := jwt.ParseToken(res, req)
	if err != nil {
		log.Print("Invalid token")
		c.Ctx.Redirect(301, "/")
		return false
	}
	email, _ := u["Email"].(string)
	if !providers.AuthenticateUser(email) {
		log.Print("Invalid user")
		c.Data["Content"] = "Send your email and name to ausrasul@gmail.com to get access"
		c.TplName = "access.tpl"
		return false
	}
	if email == "ausrasul@gmail.com" {
		c.Data["Admin"] = true
	}
	return true
}

func (c *MainController) AuthenticatePost() bool {
	if NOAUTH {
		return true
	}

	res := c.Ctx.ResponseWriter
	req := c.Ctx.Request
	u, err := jwt.ParseToken(res, req)
	if err != nil {
		log.Print("Invalid token")
		mystruct := "{\"status\": \"err\", \"msg\": \"Session timed out, log in again\"}"
		c.Data["json"] = &mystruct
		c.ServeJSON()
		return false
	}
	email, _ := u["Email"].(string)
	if !providers.AuthenticateUser(email) {
		log.Print("Invalid user")
		mystruct := "{\"status\": \"err\", \"msg\": \"Send your email and name to ausrasul@gmail.com to get access\"}"
		c.Data["json"] = &mystruct
		c.ServeJSON()
		return false
	}
	return true
}

func (c *MainController) Get() {
	// Check if the user already have a valid token
	if !c.AuthenticateGet() {
		return
	}
	var s providers.Stocks
	err := s.Get()
	if err != nil {
		log.Print("Couldn't load stocks")
		c.TplName = "index.tpl"
	}
	sortedStocks, _ := s.Sort()
	c.Data["Stocks"] = sortedStocks
	c.Layout = "layout.html"
	c.TplName = "index.tpl"
}

func (c *MainController) Degiro() {
	// Check if the user already have a valid token
	c.TplName = "degiro.tpl"
}


func (c *MainController) AddStock() {
	// Check if the user already have a valid token
	//var j goJwt.GOJWT
	if !c.AuthenticatePost() {
		return
	}
	// print our state string to the console. Ideally, you should verify
	// that it's the same string as the one you set in `setState`
	//fmt.Println("key: ", gothic.)

	name := c.Input().Get("name")
	code := c.Input().Get("code")

	if len(name) < 3 || len(code) < 3 {
		mystruct := "{\"status\": \"Error\", \"msg\": \"Invalid Input\"}"
		c.Data["json"] = &mystruct
		c.ServeJSON()
		return
	}
	log.Print("Name: ", name)
	log.Print("Code: ", code)
	//p := Pool.Get()
	//defer p.Close()

	//test, _ := p.Do("GET", "astaxie")
	//red, err := cache.NewCache("redis", `{"conn":":6379"}`)
	var stocks providers.Stocks
	err := stocks.Get()
	if err != nil {
		log.Print("error", err)
	}
	err = stocks.Add(name, code)
	if err != nil {
		log.Print("error", err)
		mystruct := "{\"status\": \"Error\", \"msg\": \"Could not add stock\"}"
		c.Data["json"] = &mystruct
		c.ServeJSON()
		return
	}
	mystruct := "{\"status\": \"OK\", \"msg\": \"OK\"}"
	c.Data["json"] = &mystruct
	c.ServeJSON()
}

func (c *MainController) RunCustom() {
	// Check if the user already have a valid token
	//var j goJwt.GOJWT
	if !c.AuthenticatePost() {
		return
	}
	// print our state string to the console. Ideally, you should verify
	// that it's the same string as the one you set in `setState`
	//fmt.Println("key: ", gothic.)

	err := providers.RunCustom()
	if err != nil {
		log.Print("error", err)
		mystruct := "{\"status\": \"Error\", \"msg\": " + err.Error()
		c.Data["json"] = &mystruct
		c.ServeJSON()
		return
	}
	mystruct := "{\"status\": \"OK\", \"msg\": \"OK\"}"
	c.Data["json"] = &mystruct
	c.ServeJSON()
}

func (c *MainController) Recommendations(){
	log.Print("API get recommendations")

	var s providers.Stocks
	err := s.Get()
	if err != nil {
		log.Print("Couldn't load stocks")
		mystruct := "{\"status\": \"err\", \"msg\": \"" + err.Error() + "\"}"
		c.Data["json"] = &mystruct
		c.ServeJSON()
	}
	sortedStocks, _ := s.Sort()
	c.Data["json"] = sortedStocks
	c.ServeJSON()
}

func (c *MainController) SyncStocks() {
	// Check if the user already have a valid token
	//var j goJwt.GOJWT
	if !c.AuthenticatePost(){
		return
	}
	
	var stocks providers.Stocks
	err := stocks.Get()
	if len(stocks.Stk) < 1 || err != nil {
		mystruct := "{\"status\": \"err\", \"msg\": \"Start by adding some stocks\"}"
		c.Data["json"] = &mystruct
		c.ServeJSON()
		return
	}

	err = stocks.Sync()
	if err != nil {
		mystruct := "{\"status\": \"err\", \"msg\": \"" + err.Error() + "\"}"
		c.Data["json"] = &mystruct
		c.ServeJSON()
		return
	}

	mystruct := "{\"status\": \"OK\", \"msg\": \"OK\"}"
	c.Data["json"] = &mystruct
	c.ServeJSON()
}

func (c *MainController) ShowStock() {
	if !c.AuthenticateGet() {
		return
	}
	// print our state string to the console. Ideally, you should verify
	// that it's the same string as the one you set in `setState`
	//fmt.Println("key: ", gothic.)
	var code = c.Ctx.Input.Param(":code")
	if len(code) < 2 || len(code) > 10 {
		log.Print("Invalid stock code")
		c.Ctx.Redirect(404, "/")
		return
	}
	c.Data["Code"] = code
	c.Layout = "layout.html"
	c.TplName = "stock.tpl"
}

func (c *MainController) ParsePortfolio() {
	if !c.AuthenticatePost() {
		return
	}
	// print our state string to the console. Ideally, you should verify
	// that it's the same string as the one you set in `setState`
	//fmt.Println("key: ", gothic.)
	file, header, err := c.GetFile("file")
	if err != nil || header.Header["Content-Type"][0] != "application/vnd.ms-excel" {
		mystruct := "{\"status\": \"err\", \"msg\": \"" + err.Error() + "\"}"
		c.Data["json"] = &mystruct
		c.ServeJSON()
		return
	}
	h, err := providers.ParsePortfolio(file)
	if err != nil {
		mystruct := "{\"status\": \"err\", \"msg\": \"" + err.Error() + "\"}"
		c.Data["json"] = &mystruct
		c.ServeJSON()
		return
	}

	c.Data["json"] = h
	c.ServeJSON()
}

func (c *MainController) LoadStock() {
	if !c.AuthenticatePost() {
		return
	} // print our state string to the console. Ideally, you should verify
	// that it's the same string as the one you set in `setState`
	//fmt.Println("key: ", gothic.)
	var code = c.Ctx.Input.Param(":code")
	if len(code) < 2 || len(code) > 10 {
		log.Print("Invalid Stock code")
		mystruct := "{\"status\": \"err\", \"msg\": \"Invalid stock code.\"}"
		c.Data["json"] = &mystruct
		c.ServeJSON()
		return
	}
	h := providers.StockHistoryPrice{
		Code:   code,
		Prices: make(map[string]providers.DayPrice),
	}
	err := h.Get()

	if len(h.Prices) < 1 || err != nil {
		log.Print("Could not retrieve this stock from database")
		mystruct := "{\"status\": \"err\", \"msg\": \"Could not retrieve this stock from database\"}"
		c.Data["json"] = &mystruct
		c.ServeJSON()
		return
	}

	c.Data["json"] = h.GetJSReadablePrices()
	c.Data["Code"] = h.Code
	c.ServeJSON()
}

func (c *MainController) Forecast() {
	// Check if the user already have a valid token
	//var j goJwt.GOJWT
	if !c.AuthenticatePost(){
		return
	}

	var stocks providers.Stocks
	err := stocks.Get()
	if len(stocks.Stk) < 1 || err != nil {
		mystruct := "{\"status\": \"err\", \"msg\": \"Start by adding some stocks\"}"
		c.Data["json"] = &mystruct
		c.ServeJSON()
		return
	}

	err = stocks.Forecast()
	if err != nil {
		mystruct := "{\"status\": \"err\", \"msg\": \"" + err.Error() + "\"}"
		c.Data["json"] = &mystruct
		c.ServeJSON()
		return
	}

	mystruct := "{\"status\": \"OK\", \"msg\": \"OK\"}"
	c.Data["json"] = &mystruct
	c.ServeJSON()
}

func (c *MainController) AnalyzeHistory() {
	// Check if the user already have a valid token
	//var j goJwt.GOJWT
	if !c.AuthenticatePost() {
		return
	} // print our state string to the console. Ideally, you should verify
	// that it's the same string as the one you set in `setState`
	//fmt.Println("key: ", gothic.)

	var stocks providers.Stocks
	err := stocks.Get()
	if len(stocks.Stk) < 1 || err != nil {
		mystruct := "{\"status\": \"err\", \"msg\": \"Start by adding some stocks\"}"
		c.Data["json"] = &mystruct
		c.ServeJSON()
		return
	}

	err = stocks.AnalyzeHistory()
	if err != nil {
		mystruct := "{\"status\": \"err\", \"msg\": \"" + err.Error() + "\"}"
		c.Data["json"] = &mystruct
		c.ServeJSON()
		return
	}

	mystruct := "{\"status\": \"OK\", \"msg\": \"OK\"}"
	c.Data["json"] = &mystruct
	c.ServeJSON()
}

func (c *MainController) GetSignal() {
	// Check if the user already have a valid token
	//var j goJwt.GOJWT
	if !c.AuthenticatePost() {
		return
	} // print our state string to the console. Ideally, you should verify
	// that it's the same string as the one you set in `setState`
	//fmt.Println("key: ", gothic.)

	var stocks providers.Stocks
	err := stocks.Get()
	if len(stocks.Stk) < 1 || err != nil {
		mystruct := "{\"status\": \"err\", \"msg\": \"Start by adding some stocks\"}"
		c.Data["json"] = &mystruct
		c.ServeJSON()
		return
	}

	err = stocks.GetSignal()
	if err != nil {
		mystruct := "{\"status\": \"err\", \"msg\": \"" + err.Error() + "\"}"
		c.Data["json"] = &mystruct
		c.ServeJSON()
		return
	}

	mystruct := "{\"status\": \"OK\", \"msg\": \"OK\"}"
	c.Data["json"] = &mystruct
	c.ServeJSON()
}
