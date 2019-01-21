package routers

import (
	"app/controllers"
	"github.com/astaxie/beego"
	
)

func init() {
	beego.Info("Loading router...")
	
    //beego.Router("/", &controllers.MainController{})
	beego.Router("/", &controllers.LoginController{}, "get:ShowLoginPage")
	beego.Router("/login/:provider", &controllers.LoginController{}, "get:Authenticate")
	beego.Router("/auth/:provider/callback", &controllers.LoginController{}, "get:Validate")
	beego.Router("/secure", &controllers.MainController{})
	beego.Router("/secure/addStock", &controllers.MainController{}, "post:AddStock")
	beego.Router("/secure/syncStocks", &controllers.MainController{}, "post:SyncStocks")
	beego.Router("/secure/forecast", &controllers.MainController{}, "post:Forecast")
	beego.Router("/secure/preload/:code", &controllers.MainController{}, "get:ShowStock")
	beego.Router("/secure/load/:code", &controllers.MainController{}, "get:LoadStock")
	beego.Router("/secure/analyzeHistory", &controllers.MainController{}, "post:AnalyzeHistory")
	beego.Router("/secure/getSignal", &controllers.MainController{}, "post:GetSignal")
	beego.Router("/secure/degiro", &controllers.MainController{}, "get:Degiro")
	beego.Router("/secure/runCustom", &controllers.MainController{}, "post:RunCustom")
	beego.Router("/secure/upload", &controllers.MainController{}, "post:ParsePortfolio")
	beego.Router("/secure/getRecommendations", &controllers.MainController{}, "get:Recommendations")
	//beego.Router("/secure/parsePortfolio", &controllers.MainController{}, "post:ParsePortfolio")
	
}
