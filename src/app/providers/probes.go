package providers

import(
	"github.com/astaxie/beego"
	"github.com/ausrasul/redisorm"
)

type Probe struct{
	Name string
	Id int
}

var probeKey string = "astaxie"

func (p *Probe) Get() error{
	err := redisorm.Get(probeKey, p)
	return err
}

func (p *Probe) Set() error{
	beego.Debug(p)
	err := redisorm.Set(probeKey, p)
	return err
}
