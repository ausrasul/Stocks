package providers

import (
	"github.com/astaxie/beego"
	"math/rand"
	"sort"
	"sync"
	"time"
	//"encoding/json"
	//"encoding/binary"
)
/*
	representation
	- take 5 days chunk
	- average first day
	- conversion factor = 100/ avg
	- 

*/
var conf = struct {
	Days        int
	Window      int
	Xor         int //crossover ratio percent
	Mur         int //mutation ratio percent
	Psize       int //Population size
	Winners     int
	Generations int
	StopDate    string
	Capital		float64
}{
	Days:        6, // first day invisible just use its avg, then 5 days data, and today. totally 7 but this starts from 0 so it last point is 6
	Window:      40,
	Xor:         50,
	Mur:         5,
	Psize:       40000,
	Winners:     2000,
	Generations: 2000,
	StopDate:    "2014-03-31",
	Capital:     20000,
}

//olhc olhc olhc olhc olhc
// Representation:
/*
	o = open
	c = close
	l = low
	h = high
	a = average
	x = last close
	old Genome: [ohlc ohlc ohlc ohlc ohlc opreviousa limit stoploss timeout]
	Genome: [0.1 0.1 1.1 0.1 1.1 0.1 1.1 0.1 1.1 0.1   again  again   again   again  0.1 0.1 1.02 0.8 40]
			 \____________________________________/                                   ^   ^    ^   ^   ^
              ^   ^   ^   ^   ^   ^   ^   ^   ^   ^                                   |   |    |   |   |__ 40 day timeout
              |   |   |   |   |   |   |   |   |   |__ +/-                             |   |    |   |__Stop loss -20%
              |   |   |   |   |   |   |   |   |__ close to avg %                      |   |    |__Sell limit 2%
              |   |   |   |   |   |   |   |__ +/-                                     |   |__ yesterday +/-
              |   |   |   |   |   |   |__low to avg %                                 |__ today open price 10% higher than avg yesterday
              |   |   |   |   |   |__ +/- percent tolerance for previous reading
              |   |   |   |   |__ high price 10% higher than avg
              |   |   |   |__ +/- percent tolerance for previous reading
              |   |   |__open price 10% higher than avg
              |   |__ +/-
              |__ open to previous day avg
*/
type StockRep struct {
	//avg ohlc
	avg float64
	aohlc [5]float64
	tolerance [5]float64
}
type Chromosom struct {
	geneGrp [5]StockRep
	today float64
	todayTol float64
	limit float64
	stoploss float64
	timeout int
	Fit float64
	Hits int
	Losses int
}

func flt() float64{
	return rand.Float64() + 0.5
}
func NewStockRep() StockRep {
	s := StockRep{}
	s.avg = rand.Float64() * 20
	for k,_ := range s.aohlc {
		s.aohlc[k] = flt()
	}
	for k,_ := range s.tolerance {
		s.tolerance[k] = rand.Float64()
	}
	return s
}
func NewChromosom() Chromosom{
	c := Chromosom{}
	for k,_ := range c.geneGrp{
		c.geneGrp[k] = NewStockRep()
	}
	c.today = flt()
	c.todayTol = rand.Float64()
	c.limit = rand.Float64() + 1
	c.stoploss = rand.Float64()
	c.timeout = rand.Intn(conf.Window)
	c.Fit = 0
	c.Hits = 0
	c.Losses = 0
	return c
}
func CrossOver(cWinner Chromosom, cLoser Chromosom) Chromosom{
	child := Chromosom{}
	if rand.Intn(100) < conf.Xor {
		child.today = cWinner.today
	} else {
		child.today = cLoser.today
	}
	if rand.Intn(100) < conf.Xor {
		child.todayTol = cWinner.todayTol
	} else {
		child.todayTol = cLoser.todayTol
	}
	if rand.Intn(100) < conf.Xor {
		child.limit = cWinner.limit
	} else {
		child.limit = cLoser.limit
	}
	if rand.Intn(100) < conf.Xor {
		child.stoploss = cWinner.stoploss
	} else {
		child.stoploss = cLoser.stoploss
	}
	if rand.Intn(100) < conf.Xor {
		child.timeout = cWinner.timeout
	} else {
		child.timeout = cLoser.timeout
	}
	for k,_ := range cLoser.geneGrp {
		if rand.Intn(100) < conf.Xor {
			child.geneGrp[k].avg = cWinner.geneGrp[k].avg
		} else {
			child.geneGrp[k].avg = cLoser.geneGrp[k].avg
		}
		for k1,_ := range cLoser.geneGrp[k].aohlc {
			if rand.Intn(100) < conf.Xor {
				child.geneGrp[k].aohlc[k1] = cWinner.geneGrp[k].aohlc[k1]
			} else {
				child.geneGrp[k].aohlc[k1] = cLoser.geneGrp[k].aohlc[k1]
			}
			if rand.Intn(100) < conf.Xor {
				child.geneGrp[k].tolerance[k1] = cWinner.geneGrp[k].tolerance[k1]
			} else {
				child.geneGrp[k].tolerance[k1] = cLoser.geneGrp[k].tolerance[k1]
			}
		}
	}
	return child
}
func (c *Chromosom) Reset(){
	c.Fit = 0
	c.Hits = 0
	c.Losses = 0
}
func (c *Chromosom) Mutate(){
	if rand.Intn(100) < conf.Mur {
		c.today = flt()
	}
	if rand.Intn(100) < conf.Mur {
		c.todayTol = rand.Float64()
	}
	if rand.Intn(100) < conf.Mur {
		c.limit = rand.Float64() + 1
	}
	if rand.Intn(100) < conf.Mur {
		c.stoploss = rand.Float64()
	}
	if rand.Intn(100) < conf.Mur {
		c.timeout = rand.Intn(conf.Window)
	}
	for k,_ := range c.geneGrp {
		if rand.Intn(100) < conf.Mur {
			c.geneGrp[k].avg = rand.Float64() * 20
		}
		for k1,_ := range c.geneGrp[k].aohlc {
			if rand.Intn(100) < conf.Mur {
				c.geneGrp[k].aohlc[k1] = flt()
			}
			if rand.Intn(100) < conf.Mur {
				c.geneGrp[k].tolerance[k1] = rand.Float64()
			}
		}
	}
}

type Data struct {
	Chrom  []Chromosom
	Prices []DayPrice
}

type ByFitness []Chromosom

func (s ByFitness) Len() int {
	return len(s)
}
func (s ByFitness) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByFitness) Less(i, j int) bool {
	return s[i].Fit > s[j].Fit
}

func Evolve(code string) {
	beego.Debug("Evolution started..")
	
	var wg sync.WaitGroup

	h := StockHistoryPrice{
		Code:   code,
		Prices: make(map[string]DayPrice),
	}
	err := h.Get()
	if err != nil {
		beego.Critical(err)
		return
	}
	stopDate, err := time.Parse("2006-01-02", conf.StopDate)
	if err != nil {
		beego.Critical(err)
		return
	}
	err = h.SortUntil(stopDate)
	if err != nil {
		beego.Critical(err)
		return
	}

	rand.Seed(time.Now().UnixNano())
	var population = make([]Chromosom, conf.Psize)
	for k, _ := range population {
		population[k] = NewChromosom()
	}
	data := RepresentData(h.sortedPrices)
	for i := 0; i < conf.Generations; i++ {
		beego.Debug("Testing generation ", i)
		for k, _ := range population {
			wg.Add(1)
			go fitness(&population[k], data, &wg)
			//fitness(&population[k], data, &wg)
		}
		wg.Wait()

		sort.Sort(ByFitness(population))
		beego.Debug("Top 10:")
		for c := 0; c < 10; c++ {
			pop := population[c]
			beego.Debug("#", c, ": Fitness=>", int(pop.Fit), " Hits=>", pop.Hits, " Losses=> ", pop.Losses)
			beego.Debug("Genes=>", pop)
		}
		poplen := len(population)
		for k, _ := range population {
			if k > conf.Winners {
				winner := GetWinner(poplen)
				loser := GetWinner(poplen)
				population[k] = CrossOver(population[winner], population[loser])
			}
			population[k].Reset()
		}
	}
	beego.Debug("Evolution complete.")
}
// returns the index of the winner according to best chance per fitness, otherwise random.
func GetWinner(poplen int) int{
	chunks := 10
	chunkSize := poplen / chunks
	chunk := 0
	for k:=0; k<chunks; k++{
		if rand.Intn(100) < 25{
			chunk = k
			break
		}
	}
	offset := chunkSize * chunk
	return rand.Intn(chunkSize) + offset
}

func RepresentData(p []DayPrice) Data {
	limit := len(p) - conf.Window
	start := conf.Days
	var data Data
	data.Chrom = make([]Chromosom, len(p))
	data.Prices = make([]DayPrice, len(p))
	lenp := len(p)
	for i := 0; i < lenp; i++ {
		st := Chromosom{}
		if i >= start && i < limit {
			predayAvg := (p[i - conf.Days].Open + p[i - conf.Days].Close) / 2
			for k,_ := range st.geneGrp {
				dayIndex := i - conf.Days + k + 1
				st.geneGrp[k].avg = (p[dayIndex].Open + p[dayIndex].Close) / 2
				st.geneGrp[k].aohlc[0] = st.geneGrp[k].avg / predayAvg
				st.geneGrp[k].aohlc[1] = p[dayIndex].Open / st.geneGrp[k].avg
				st.geneGrp[k].aohlc[2] = p[dayIndex].High / st.geneGrp[k].avg
				st.geneGrp[k].aohlc[3] = p[dayIndex].Low / st.geneGrp[k].avg
				st.geneGrp[k].aohlc[4] = p[dayIndex].Close / st.geneGrp[k].avg
			}
			st.today = p[i].Open / st.geneGrp[len(st.geneGrp) - 1].avg
		}
		data.Chrom[i] = st
		data.Prices[i] = p[i]
	}
	return data
}


func fitness(a *Chromosom, st Data, wg *sync.WaitGroup) {
	// Set the test limit
	//lenp := len(st.Prices)
	//beego.Debug("Checking fitness...")
	limit := len(st.Prices) - conf.Window
	start := conf.Days
	//beego.Debug("Limit ", limit, " Start ", start)
	// and reset this chromosom score
	// Test the chromosome on each day (sliding window)
	a.Fit = 0
	for i := start; i < limit; i++ {
		//beego.Debug("a.today >= (st.Chrom[i].today - a.todayTol) && a.today <= (st.Chrom[i].today + a.todayTol)")
		//beego.Debug(a.today, " >= (", st.Chrom[i].today, " - ", a.todayTol, ") && ", a.today, " <= (", st.Chrom[i].today, " + ", a.todayTol, ")")
		//beego.Debug(a.today, " >= (", st.Chrom[i].today - a.todayTol, ") && ", a.today, " <= (", st.Chrom[i].today + a.todayTol, ")")
		//beego.Debug(a.today >= (st.Chrom[i].today - a.todayTol) && a.today <= (st.Chrom[i].today + a.todayTol))
		if !(a.today >= (st.Chrom[i].today - a.todayTol) && a.today <= (st.Chrom[i].today + a.todayTol)) {
			continue
		}
		allGenesMatch := true
		for k,v := range a.geneGrp {
			stg := st.Chrom[i].geneGrp[k]
			if !(v.aohlc[0] >= stg.aohlc[0] - v.tolerance[0] && v.aohlc[0] <= stg.aohlc[0] + v.tolerance[0]) {
				allGenesMatch = false
			}
			if !(v.aohlc[1] >= stg.aohlc[1] - v.tolerance[1] && v.aohlc[1] <= stg.aohlc[1] + v.tolerance[1]) {
				allGenesMatch = false
			}
			if !(v.aohlc[2] >= stg.aohlc[2] - v.tolerance[2] && v.aohlc[2] <= stg.aohlc[2] + v.tolerance[2]) {
				allGenesMatch = false
			}
			if !(v.aohlc[3] >= stg.aohlc[3] - v.tolerance[3] && v.aohlc[3] <= stg.aohlc[3] + v.tolerance[3]) {
				allGenesMatch = false
			}
			if !(v.aohlc[4] >= stg.aohlc[4] - v.tolerance[4] && v.aohlc[4] <= stg.aohlc[4] + v.tolerance[4]) {
				allGenesMatch = false
			}
		}
		if !allGenesMatch {
			continue
		}
		target := st.Prices[i].Open * a.limit
		sold := false
		k := 0
		for k = 0; k < a.timeout; k++ {
			if st.Prices[k+i].High >= target {
				a.Fit = a.Fit + (a.limit - 1)
				sold = true
				a.Hits = a.Hits + 1
				break
			}
			if st.Prices[k+i].Close <= st.Prices[i].Open * a.stoploss {
				a.Fit = a.Fit - (1 - (st.Prices[k+i].Close / st.Prices[i].Close))
				a.Losses = a.Losses + 1
				sold = true
				break
			}
		}
		if !sold {
			a.Fit = a.Fit + ((st.Prices[i + a.timeout].Close / st.Prices[i].Close) - 1)
			a.Losses = a.Losses + 1
		}
		i += k
	}
	wg.Done()
}
