package main

import (
	_ "app/routers"
	//"runtime"
	"github.com/astaxie/beego"
	// OAuth authentication packages
	"github.com/ausrasul/redisorm" // Redis with resource pool
	"github.com/ausrasul/jwt"        // Web token packages
	//"github.com/ausrasul/m2mserver"	 // M2M server package
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/gplus"
	"log"
)
var NOAUTH bool // if true, anyone can log in with no account checking.
func main() {
	NOAUTH = false
	//runtime.GOMAXPROCS(8)

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	beego.SetLogFuncCall(true)
	
	beego.BConfig.WebConfig.Session.SessionOn = true
	goth.UseProviders(
		gplus.New(
			beego.AppConfig.String("CLIENT_ID"),
			beego.AppConfig.String("CLIENT_SECRET"),
			beego.AppConfig.String("CLIENT_CALLBACK"),
		),
	)
	
	noAuth, err := beego.AppConfig.Int("NoAuth")
	if err == nil && noAuth == 1 {
		NOAUTH = true
	}
	SessionTimeout, err := beego.AppConfig.Int("SESSION_TIMEOUT")
	if err != nil {
		beego.Critical(err)
	}
	SessionRefreshInterval, err := beego.AppConfig.Int("SESSION_REFRESH_INTERVAL")
	if err != nil {
		beego.Critical(err)
	}

	jwt.Configure(
		map[string]interface{}{
			"privateKeyFile":         beego.AppConfig.String("PrivateKeyFile"),
			"publicKeyFile":          beego.AppConfig.String("PublicKeyFile"),
			"algorithm":              beego.AppConfig.String("Algorithm"),
			"sessionName":            beego.AppConfig.String("SESSION_NAME"),
			"sessionTimeout":         SessionTimeout,
			"sessionRefreshInterval": SessionRefreshInterval,
		},
	)

	poolMaxIdle, err := beego.AppConfig.Int("REDIS_MaxIdle")
	if err != nil {
		beego.Critical(err)
	}
	poolMaxActive, err := beego.AppConfig.Int("REDIS_MaxActive")
	if err != nil {
		beego.Critical(err)
	}

	redisorm.Configure(
		map[string]interface{}{
			"poolMaxIdle":   poolMaxIdle,
			"poolMaxActive": poolMaxActive,
			"port":          beego.AppConfig.String("REDIS_Port"),
		},
	)
	beego.SetStaticPath("/public", "static")

	beego.Run()
}
