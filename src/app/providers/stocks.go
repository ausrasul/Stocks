package providers

import(
	"github.com/astaxie/beego"
	"github.com/ausrasul/redisorm"
	"time"
	"net/http"
	"net/url"
	"io/ioutil"
	"strconv"
	"bytes"
	//"net/http/httputil"
	"strings"
	"errors"
	//"fmt"
	"sort"
	"encoding/json"
	"encoding/csv"
	"io"
	"log"
)

type Position struct{
	amount float64
	cost float64
}
type Portfolio struct{
	cash float64
	stocks map[string]Position
}

func getDeposits(file io.Reader) (cash float64, datum int64, err error){
	//r := csv.NewReader(strings.NewReader(csvString))
	r := csv.NewReader(file)
	cash = 0
	datum = 0
	first := true
	for {
		record, err := r.Read()
		if err == io.EOF {
			log.Print("End of text")
			break
		}
		if err != nil {
			log.Print("Error reading text")
			return cash, datum, err
		}
		t, err := time.Parse("2/1/06", record[1])
		if err != nil {
			continue
		}
		dt := t.Unix() * 1000
		if first {
			datum = dt
			first = false
		}
		row := strings.SplitN(record[2], " ", 2)
		cmd := row[0]
		if cmd == "Deposit" || cmd == "iDEAL"{
			stage1 := strings.Replace(record[5], "\u00a0", "", -1)
			stage2 := strings.Replace(stage1, ",", ".", -1)
			deposite, err := strconv.ParseFloat(stage2 , 64)
			if err != nil {
				log.Print("Error converting deposite", stage1, stage2, deposite, err)
				return cash, datum, err
			}
			log.Print(deposite)
			cash += deposite
		}
	}
	return cash, datum, err
}
func ParsePortfolio(file io.Reader) ([]map[string]interface{}, error){
	var timeLine []map[string]interface{}
	//timeLine := make(map[int64]float64)
	//r := csv.NewReader(strings.NewReader(csvString))
	var b bytes.Buffer
	b.ReadFrom(file)
	f1 := bytes.NewReader(b.Bytes())
	f2 := bytes.NewReader(b.Bytes())
	depCash, depDate, err := getDeposits(f1)
	if err != nil {
		return timeLine, err
	}
	r := csv.NewReader(f2)
	if depCash == 0 || depDate == 0 {
		return timeLine, err
	}
	timeLine = append(timeLine, map[string]interface{}{
		"dt": depDate,
		"cash": depCash,
		"buy" : 0,
		"sell" : 0,
	})

	portfolio := Portfolio{
		cash: depCash,
		stocks: make(map[string]Position),
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			log.Print("End of text")
			break
		}
		if err != nil {
			log.Print("Error reading text")
			return timeLine, err
		}
		t, err := time.Parse("2/1/06", record[1])
		if err != nil {
			continue
		}
		dt := t.Unix() * 1000
		row := strings.SplitN(record[2], " ", 2)
		cmd := row[0]
		if cmd == "Buy" || cmd == "Sell" {
			row = strings.SplitN(row[1], " ", 2)
			sAmount := row[0]
			row = strings.SplitN(row[1], "@", 2)
			stock := row[0]
			row = strings.SplitN(row[1], " ", 2)
			sPrice := row[0]
			price, err := strconv.ParseFloat(sPrice, 64)
			if err != nil && len(row[4]) > 0 {
				log.Print("Error converting price")
				return timeLine, err
			}
			amount, err := strconv.ParseFloat(sAmount, 64)
			if err != nil && len(row[5]) > 0 {
				log.Print("Error converting amount")
				return timeLine, err
			}
			transaction := price * amount
			log.Print(t, " ", cmd, " ", stock, " ", amount, " ", price, " ", transaction)
			log.Print("Cash before: ", portfolio.cash)
			if cmd == "Buy"{
				//portfolio.cash -= transaction
				pos, _ := portfolio.stocks[stock]
				pos.amount += amount
				pos.cost += price * amount
				portfolio.stocks[stock] = pos
			} else {
				//portfolio.cash += transaction
				pos, ok := portfolio.stocks[stock]
				if ok {
					if pos.amount == 0 {
						return timeLine, err
					}
					pos.amount -= amount
					pos.cost -= price * amount
					if pos.amount == 0 {
						portfolio.cash -= pos.cost
						pos.cost = 0
					}
					portfolio.stocks[stock] = pos
				} else {
					return timeLine, err
				}
			}
			log.Print("Cash after: ", portfolio.cash)
		}
		if cmd == "Giro" {
			stage1 := strings.Replace(record[5], "\u00a0", "", -1)
			stage2 := strings.Replace(stage1, ",", ".", -1)
			deposite, err := strconv.ParseFloat(stage2 , 64)
			if err != nil {
				log.Print("Error converting deposite", stage1, stage2, deposite, err)
				return timeLine, err
			}
			log.Print(cmd, deposite)
			portfolio.cash += deposite
		}
		totalCash := portfolio.cash
		for k,v := range portfolio.stocks{
			if v.amount > 0{
				log.Print (k, " ", v.amount, " ", v.cost)
			}
		//	totalCash += v.amount * v.price
		}
		//timeLine[dt] = totalCash
		mp := make(map[string]interface{})
		b := 0
		s := 0
		if cmd == "Buy" {
			b = 1
		}
		if cmd == "Sell"{
			s = 1
		}
		mp = map[string]interface{}{
			"dt": dt,
			"cash": totalCash,
			"buy": b,
			"sell": s,
		}
		timeLine = append(timeLine, mp)
	}
	return timeLine, err
}

func init(){
	go schedules()
}
func RunCustom() error{
	Evolve("SSE1027")
	return nil
}

func schedules(){
	time.Sleep(time.Second)
/*	var stocks Stocks
	err := stocks.Get()
	if len(stocks.Stk) < 1 || err != nil{
		beego.Error("No stocks to sync", err)
		return
	}
	err = stocks.AnalyzeHistory()
	if err != nil {
		beego.Error("No stocks to analyze ", err)
		return
	}
	beego.Error("analyze complete", err)

*/
	//Evolve("SSE1027")
	for {
		tzero := time.Now().Truncate(24 * time.Hour).Add(time.Hour * -2)
		beego.Debug("time today ", tzero)
		//tTarget:= 
		//beego.Debug("time target ", tTarget)
		//timeToTest := tzero.Add(20 * time.Hour).Add(33 * time.Minute).Sub(time.Now())
		timeToSync := tzero.Add(18 * time.Hour).Add(15 * time.Minute).Sub(time.Now())
		timeToSignal1 := tzero.Add(9 * time.Hour).Add(31 * time.Minute).Sub(time.Now())
		timeToSignal2 := tzero.Add(9 * time.Hour).Add(33 * time.Minute).Sub(time.Now())
		timeToSignal3 := tzero.Add(9 * time.Hour).Add(35 * time.Minute).Sub(time.Now())
		//if timeToTest < 0 {
		//	timeToTest = timeToTest + (24 * time.Hour)
		//}
		if timeToSync < 0 {
			timeToSync = timeToSync + (24 * time.Hour)
		}
		if timeToSignal1 < 0 {
			timeToSignal1 = timeToSignal1 + (24 * time.Hour)
		}
		if timeToSignal2 < 0 {
			timeToSignal2 = timeToSignal2 + (24 * time.Hour)
		}
		if timeToSignal3 < 0 {
			timeToSignal3 = timeToSignal3 + (24 * time.Hour)
		}
		//beego.Debug("Next test     ", timeToTest, timeToTest < 0)
		beego.Debug("Next sync     ", timeToSync, timeToSync < 0)
		beego.Debug("Next signal1  ", timeToSignal1, timeToSignal1 < 0)
		beego.Debug("Next signal2  ", timeToSignal2, timeToSignal2 < 0)
		beego.Debug("Next signal3  ", timeToSignal3, timeToSignal3 < 0)

		select {
			/*case <-time.After(timeToTest):
			beego.Debug("Test started")
			var stocks Stocks
			err := stocks.Get()
			if len(stocks.Stk) < 1 || err != nil{
				beego.Error("No stocks to test", err)
				break
			}
			nl := stocks.NewsLetter()
			if len(nl) < 1 {
				beego.Error("No stocks to send", err)
				break
			}
			err = SendNewsLetter(nl)
			if err != nil{
				beego.Error(err)
			}
			beego.Debug("Test successful")
			break
			*/
			case <-time.After(timeToSync):
			beego.Debug("Sycn started")
			var stocks Stocks
			err := stocks.Get()
			if len(stocks.Stk) < 1 || err != nil{
				beego.Error("No stocks to sync", err)
				break
			}
			err = stocks.Sync()
			if err != nil {
				beego.Error("No stocks to sync ", err)
				break
			}
			beego.Debug("Analyze started")
			err = stocks.Get()
			if len(stocks.Stk) < 1 || err != nil{
				beego.Error("No stocks to analyze", err)
				break
			}
			err = stocks.AnalyzeHistory()
			if err != nil {
				beego.Error("No stocks to analyze ", err)
				break
			}
			beego.Debug("Forecast started")
			err = stocks.Get()
			if len(stocks.Stk) < 1 || err != nil{
				beego.Error("No stocks to forecast", err)
				break
			}
			err = stocks.Forecast()
			if err != nil {
				beego.Error("No stocks to forecast ", err)
				break
			}
			beego.Debug("Sync operations successful")
			break
			case <-time.After(timeToSignal1):
			beego.Debug("Get siganl started")
			var stocks Stocks
			err := stocks.Get()
			if len(stocks.Stk) < 1 || err != nil{
				beego.Error("No stocks to forecast", err)
				break
			}
			err = stocks.GetSignal()
			if err != nil {
				beego.Error("No stocks to analyze ", err)
				break
			}
			beego.Debug("Get signal successful")
			break
			case <-time.After(timeToSignal2):
			beego.Debug("Get siganl started")
			var stocks Stocks
			err := stocks.Get()
			if len(stocks.Stk) < 1 || err != nil{
				beego.Error("No stocks to forecast", err)
				break
			}
			err = stocks.GetSignal()
			if err != nil {
				beego.Error("No stocks to analyze ", err)
				break
			}
			beego.Debug("Get signal successful")
			break
			case <-time.After(timeToSignal3):
			beego.Debug("Get siganl started")
			var stocks Stocks
			err := stocks.Get()
			if len(stocks.Stk) < 1 || err != nil{
				beego.Error("No stocks to forecast", err)
				break
			}
			err = stocks.GetSignal()
			if err != nil {
				beego.Error("No stocks to analyze ", err)
				break
			}
			//beego.Debug("Sending news letter")
			//nl := stocks.NewsLetter()
			//if len(nl) < 1 {
			//	beego.Error("No stocks to send", err)
			//	break
			//}
			//err = SendNewsLetter(nl)
			//if err != nil{
			//	beego.Error(err)
			//}
			beego.Debug("Get signal successful")
			//beego.Debug("letter sent successful")
		}
	}
}

var stockList string = "stocks"
type DayPrice struct{
	Date string
	Dt int64
	Open float64
	High float64
	Low float64
	Close float64
	Volume int
	Buy int
	GoodDay int
	BadDay int
	Test int
}
type StockHistoryPrice struct{
	Code string
	Name string
	Prices map[string]DayPrice
	sortedPrices []DayPrice
	accuracy string
	rating int
}

type Stock struct{
	Name string
	Code string
	Updated string
	Exists int
	Status int
	Candidate float64
	Signal string
	Speed string
	Accuracy string
	Rating int
}

type Stocks struct{
	Stk map[string]*Stock
	sortedStk []Stock
}



func (h *StockHistoryPrice) Get() error{
	h.Prices = make(map[string]DayPrice)
	err := redisorm.Get(h.Code, h)
	//beego.Debug(err)
	/*if (err == nil){
		err = h.Sort()
	}*/
	return err
}
func (h *StockHistoryPrice) Set() error{
	err := redisorm.Set(h.Code, h)
	return err
}
func (h *StockHistoryPrice) FixTime(){
	beego.Debug("fixing time")
	for k,v := range h.Prices{
		t, _ := time.Parse("2006-01-02", v.Date)
		//h.Prices[k].Dt = t.Unix() * 1000
		dt, _ := h.Prices[k]
		dt.Dt = t.Unix() * 1000
		h.Prices[k] = dt
	}
	return
}
func (h *StockHistoryPrice) Add(dayPrice DayPrice) {
	h.Prices[dayPrice.Date] = dayPrice
	//err := h.Set()
	return
}
func (h *StockHistoryPrice) Sort() error {
	if len(h.Prices) < 1 {
		return errors.New("No history prices to sort")
	}
	h.sortedPrices = []DayPrice{}
	//i := 0
	for _,v := range h.Prices{
		h.sortedPrices = append(h.sortedPrices, v)
	}
	sort.Sort(PriceByDate(h.sortedPrices))
	return nil
}

func (h *StockHistoryPrice) SortUntil(t time.Time) error {
	dt := t.Unix() * 1000
	if len(h.Prices) < 1 {
		return errors.New("No history prices to sort")
	}
	h.sortedPrices = []DayPrice{}
	//i := 0
	for _,v := range h.Prices{
		if v.Dt < dt {
			h.sortedPrices = append(h.sortedPrices, v)
		}
	}
	sort.Sort(PriceByDate(h.sortedPrices))
	return nil
}

func (h StockHistoryPrice) GetJSReadablePrices() []map[string]interface{} {
	var jsPrices []map[string]interface{}
	h.Sort()
	for _, v := range h.sortedPrices{
		mp := make(map[string]interface{})
		mp = map[string]interface{}{
			"dt": v.Dt,
			"open": v.Open,
			"high": v.High,
			"low": v.Low,
			"close": v.Close,
			"volume": v.Volume,
			"buy": v.Buy,
			"goodDay": v.GoodDay,
			"badDay": v.BadDay,
			"test": v.Test,
		}
		jsPrices = append(jsPrices, mp)
	}
	//getNasdaqToday(h.Code)
	return jsPrices
}

// sorting stock prices for processing
type PriceByDate []DayPrice
func (s PriceByDate) Len() int{
	return len(s)
}
func (s PriceByDate) Swap(i, j int){
	s[i], s[j] = s[j], s[i]
}
func (s PriceByDate) Less(i,j int) bool{
	return s[i].Dt < s[j].Dt
}

// sorting stocks by rating
type StockByRating []Stock
func (s StockByRating) Len() int{
	return len(s)
}
func (s StockByRating) Swap(i, j int){
	s[i], s[j] = s[j], s[i]
}
func (s StockByRating) Less(i,j int) bool{
	return s[i].Rating > s[j].Rating
}


func (s *Stocks) Get() error{
	s.Stk = make(map[string]*Stock)
	err := redisorm.Get(stockList, s)
	return err
}

func (s *Stocks) Add(name string, code string) error{
	var st Stock
	st.Name = name
	st.Code = code
	st.Updated = "never"
	s.Stk[code] = &st
	err := s.Set()
	return err
}

func (s *Stocks) Set() error{
	err := redisorm.Set(stockList, s)
	return err
}

func (s *Stocks) Sort() ([]Stock,error) {
	if len(s.Stk) < 1 {
		return []Stock{},errors.New("No stocks to sort")
	}
	s.sortedStk = []Stock{}
	//i := 0
	for _,v := range s.Stk{
		s.sortedStk = append(s.sortedStk, *v)
	}
	sort.Sort(StockByRating(s.sortedStk))
	return s.sortedStk, nil
}

func (s *Stocks) Sync() error{
	err := s.Get()
	if err != nil || len(s.Stk) < 1{
		beego.Critical(err)
		return errors.New("No stocks available or error in database")
	}
	for k, v := range s.Stk{
		var start string
		var end string
		firstTime := false
		tUpdated, err := time.Parse("2006-01-02", v.Updated)
		if err != nil{
			t,_ := time.Parse("2006-01-02", "2008-01-01")
			start = t.Format("2006-01-02")
			firstTime = true
		} else {
			start = tUpdated.AddDate(0, 0, 1).Format("2006-01-02")
		}
		end = time.Now().Format("2006-01-02")
		prices, err := getNasdaq(k, start, end)
		if err != nil {
			beego.Critical(err)
			//continue
			return err
		}
		h := &StockHistoryPrice{
			Code: k,
			Name: v.Name,
			Prices: make(map[string]DayPrice),
		}
		if !firstTime{
			err = h.Get()
			if err != nil{
				return errors.New("Cannot retrieve stock prices from DB")
			}
		}
		lastUpdate, _ := time.Parse("2006-01-02", "1000-01-01")
		updates := false
		for date, price := range prices{
			updates = true
			h.Add(price)
			update, _ := time.Parse("2006-01-02", date)
			//beego.Debug("update: ", update)
			if !update.Before(lastUpdate) {
				lastUpdate = update
				//beego.Debug("lastUpdate: ", lastUpdate)
			}
		}
		if updates{
			beego.Debug("finalUpdate: ", lastUpdate)
			err = h.Set()
			if err != nil{
				return errors.New("Cannot save stock prices to DB")
			}
			s.Stk[k].Updated = lastUpdate.Format("2006-01-02")
			beego.Debug(s.Stk[k].Updated)
			err = s.Set()
			if err != nil {
				return errors.New("Cannot update stock status in DB")
			}
		}
	}
	return nil
}


/*func (s *Stocks) NewsLetter() string{
	count := 0
	c := ""
	for _, stock := range s.Stk{
		if len(stock.Signal) > 0 {
			c += "<p>Stock: " + stock.Name + "<br/>Accuracy: " + stock.Accuracy + "<br/>Recommendation: " + stock.Signal + "<br/>========================</p>"
			count += 1
		}
	}
	if count > 0{
		text := "<div>Good morning,<br/><br/>Today's stock recommendations follows.<br/>Do not forget to read the news before buying.<br/><br/>" + c + "<br/><br/>Good luck,<br/>Aus</div>"
		return text
	}
	return ""
}
*/
func (s *Stocks) Forecast() error{
	for code, _ := range s.Stk{
		h := StockHistoryPrice{
			Code: code,
		}
		err := h.Get()
		if err != nil {
			return errors.New("Cannot retrieve stock prices from db")
		}
		h.Sort()
		candidate, err := h.isCandidate()
		if err != nil {
			continue
		}
		stk, _ := s.Stk[code]
		stk.Candidate = candidate
	}
	s.Set()
	return nil
}

func (s *Stocks) GetSignal() error{
	for code, _ := range s.Stk{
		h := StockHistoryPrice{
			Code: code,
		}
		err := h.Get()
		if err != nil {
			return errors.New("Cannot retrieve stock prices from db")
		}
		h.Sort()
		signal, err := h.isGotSignal()
		//candidate, err := h.isCandidate()
		if err != nil {
			beego.Debug("Cannot get signal ", err)
			continue
		}
		stk, _ := s.Stk[code]
		if signal > 0 {
			stk.Signal = "Buy for max " + strconv.FormatFloat(signal, 'f', -1, 64)
		} else {
			stk.Signal = ""
		}
	}
	s.Set()
	return nil
}

func (s *Stocks) AnalyzeHistory() error{
	for code, _ := range s.Stk{
		h := StockHistoryPrice{
			Code: code,
		}
		err := h.Get()
		if err != nil {
			return errors.New("Cannot retrieve stock prices from db")
		}
		h.Sort()
		h.MarkAllSignals()
		//h.MarkAllSignalsCustom()
		//h.MarkAllSignals3()
		h.MarkAllSignalsTest()
		//h.RemoveInvalidPrices()
		h.Set()
		s.Stk[code].Accuracy = h.accuracy
		s.Stk[code].Rating = h.rating
	}
	s.Set()
	return nil
}
func (h *StockHistoryPrice) isGotSignal() (float64, error){
	todayPrice, err := getNasdaqToday(h.Code)
	if err != nil{
		return 0, err
	}
	oneDayAgo := len(h.sortedPrices) - 1
	twoDayAgo := oneDayAgo - 1
	threeDayAgo := oneDayAgo - 2

	condition1 := (h.sortedPrices[oneDayAgo].Close * 0.998) > todayPrice // yesterday close is 0.002 > today open
	condition2 := (h.sortedPrices[oneDayAgo].Open - h.sortedPrices[oneDayAgo].Close) > 0 //blue
	condition3 := (h.sortedPrices[twoDayAgo].Open - h.sortedPrices[twoDayAgo].Close) > 0 //blue
	condition4 := (h.sortedPrices[threeDayAgo].Open - h.sortedPrices[threeDayAgo].Close) < 0 //white
	if condition1 && condition2 && condition3 && condition4{
		return todayPrice, nil
	}
	return 0, nil
}

func (h *StockHistoryPrice) isCandidate() (float64, error){
	today := len(h.sortedPrices) - 1
	oneDayAgo := today - 1
	twoDayAgo := today - 2
	//threeDayAgo := today - 3
	dateToday := time.Now().Truncate(time.Hour * 24).Unix() * 1000
	if h.sortedPrices[today].Dt < dateToday {
		//return 0, errors.New("Too old data to get a candidate.")
	}
	condition1 := (h.sortedPrices[today].Open - h.sortedPrices[today].Close) > 0 //blue
	condition2 := (h.sortedPrices[oneDayAgo].Open - h.sortedPrices[oneDayAgo].Close) > 0 //blue
	condition3 := (h.sortedPrices[twoDayAgo].Open - h.sortedPrices[twoDayAgo].Close) < 0 //white
	if condition1 && condition2 && condition3 {
		beego.Debug("Found a candidate if open lower than ", h.sortedPrices[today].Close * 0.998)
		return h.sortedPrices[today].Close * 0.998, nil
		//return 1, nil
	}
	return 0, nil
}


func (h *StockHistoryPrice) MarkAllSignals2(){
	for k,v := range h.sortedPrices{
		if k < 5{
			continue
		}
		today := k
		oneDayAgo := today - 1
		twoDayAgo := today - 2
		threeDayAgo := today - 3
		fourDayAgo := today - 4
		//fiveDayAgo := today - 5

		price := h.sortedPrices
		condition2 := float64(price[twoDayAgo].Volume) < (float64(price[oneDayAgo].Volume) * 0.8)
		condition3 := float64(price[threeDayAgo].Volume) < (float64(price[oneDayAgo].Volume) * 0.8)
		condition4 := float64(price[fourDayAgo].Volume) < (float64(price[oneDayAgo].Volume) * 0.8)
		if condition2 && condition3 && condition4{
			priceObj, _ := h.Prices[v.Date]
			priceObj.Buy = 1
			h.Prices[v.Date] = priceObj
		} else {
			priceObj, _ := h.Prices[v.Date]
			priceObj.Buy = 0
			h.Prices[v.Date] = priceObj
		}
	}
	h.MarkAllBadDays()
}

func (h *StockHistoryPrice) MarkAllSignals3(){
	for k,v := range h.sortedPrices{
		if k < 3{
			continue
		}
		today := k
		oneDayAgo := today - 1
		//twoDayAgo := today - 2
		//threeDayAgo := today - 3
		price := h.sortedPrices
		g3 := price[today].Open / price[oneDayAgo].Close
		g2 := price[oneDayAgo].Open < price[oneDayAgo].Close // blue 1
		//g1 := price[twoDayAgo].Open > price[twoDayAgo].Close // blue 1
		//g0_1 := price[threeDayAgo].Open < price[threeDayAgo].Close // white 0
		priceObj, _ := h.Prices[v.Date]
		if g2 && (g3 >= 0.99751) {
			priceObj.Buy = 1
		} else {
			priceObj.Buy = 0
		}
		priceObj.BadDay = 0
		h.Prices[v.Date] = priceObj
	}
	h.MarkAllBadDays()
}



func (h *StockHistoryPrice) MarkAllSignals_bak(){
	for k,v := range h.sortedPrices{
		if k < 3{
			continue
		}
		today := k
		oneDayAgo := today - 1
		twoDayAgo := today - 2
		threeDayAgo := today - 3

		price := h.sortedPrices
		condition1 := (price[oneDayAgo].Close * 0.998) > price[today].Open // yesterday close is 0.002 > today open
		condition2 := (price[oneDayAgo].Open - price[oneDayAgo].Close) > 0 //blue
		condition3 := (price[twoDayAgo].Open - price[twoDayAgo].Close) > 0 //blue
		condition4 := (price[threeDayAgo].Open - price[threeDayAgo].Close) < 0 //white
		if condition1 && condition2 && condition3 && condition4{
			priceObj, _ := h.Prices[v.Date]
			priceObj.Buy = 1
			h.Prices[v.Date] = priceObj
		} else {
			priceObj, _ := h.Prices[v.Date]
			priceObj.Buy = 0
			h.Prices[v.Date] = priceObj
		}
	}
	h.MarkAllBadDays()
}

func (h *StockHistoryPrice) MarkAllSignals(){
	st := RepresentData(h.sortedPrices)
	//# 0 : Fitness=> 6  Hits=> 527  Losses=>  588
	//2016/09/11 22:14:29 [genetic.go:259][D] Genes=> {[{14.421641423537874 [1.3126530727750847 1.2399122963544271 0.8310266835531134 0.8057175593278929 1.0684218636689424] [0.711562486788633 0.626074113672291 0.9991328084385397 0.7108747633907797 0.9963127212842214]} {16.0914041175913 [0.8700022854509393 0.8101665396241504 0.9230299574474878 1.0324824601070954 1.0793497307010669] [0.9390821752302948 0.5975049767502956 0.4457704589194387 0.4167881623713668 0.8627220368654299]} {6.351454156324596 [0.7847659785141268 1.0845641762004044 1.2223651393835664 0.703007805513547 0.6100530809191] [0.9563896473078616 0.7472954550032166 0.22117474785585847 0.29282191495007565 0.6537305586447606]} {12.862302581872239 [0.9811701498104446 0.9517443982683742 1.3158558953681667 1.1516065481407334 1.1264424407047682] [0.8646047106226192 0.8374361923404612 0.4406040033957521 0.6037199115830559 0.43197214102092407]} {12.565194450499746 [0.5603260172469886 1.1885503326136915 1.336969129840545 0.9070650667413283 0.8120487275720101] [0.8104877025540802 0.8044372938101474 0.8697544026628621 0.6490316289088398 0.5269808236191077]}] 0.611994702844997 0.5335953238597572 1.0142833523641146 0.9994110562315196 6 6.404999062796704 527 588}
	
	/*{
		[
			{
				14.421641423537874
				[1.3126530727750847 1.2399122963544271 0.8310266835531134 0.8057175593278929 1.0684218636689424]
				[0.711562486788633 0.626074113672291 0.9991328084385397 0.7108747633907797 0.9963127212842214]
			}
			{
				16.0914041175913
				[0.8700022854509393 0.8101665396241504 0.9230299574474878 1.0324824601070954 1.0793497307010669]
				[0.9390821752302948 0.5975049767502956 0.4457704589194387 0.4167881623713668 0.8627220368654299]
			}
			{
				6.351454156324596
				[0.7847659785141268 1.0845641762004044 1.2223651393835664 0.703007805513547 0.6100530809191]
				[0.9563896473078616 0.7472954550032166 0.22117474785585847 0.29282191495007565 0.6537305586447606]
			}
			{
				12.862302581872239
				[0.9811701498104446 0.9517443982683742 1.3158558953681667 1.1516065481407334 1.1264424407047682]
				[0.8646047106226192 0.8374361923404612 0.4406040033957521 0.6037199115830559 0.43197214102092407]
			}
			{
				12.565194450499746
				[0.5603260172469886 1.1885503326136915 1.336969129840545 0.9070650667413283 0.8120487275720101]
				[0.8104877025540802 0.8044372938101474 0.8697544026628621 0.6490316289088398 0.5269808236191077]
			}
		]
		0.611994702844997
		0.5335953238597572
		1.0142833523641146
		0.9994110562315196
		6
		6.404999062796704
		527
		588
	}*/

	var a Chromosom
	a.today= 0.611994702844997
	a.todayTol = 0.5335953238597572
	a.limit = 1.0142833523641146
	a.stoploss = 0.9994110562315196
	a.timeout = 6
	a.Fit = 6.404999062796704
	a.Hits = 527
	a.Losses = 588
	a.geneGrp[0] = StockRep{
		avg: 14.421641423537874,
		aohlc: [5]float64{1.3126530727750847, 1.2399122963544271, 0.8310266835531134, 0.8057175593278929, 1.0684218636689424},
		tolerance: [5]float64{0.711562486788633, 0.626074113672291, 0.9991328084385397, 0.7108747633907797, 0.9963127212842214},
	}
	a.geneGrp[1] = StockRep{
		avg: 16.0914041175913,
		aohlc: [5]float64{0.8700022854509393, 0.8101665396241504, 0.9230299574474878, 1.0324824601070954, 1.0793497307010669},
		tolerance: [5]float64{0.9390821752302948, 0.5975049767502956, 0.4457704589194387, 0.4167881623713668, 0.8627220368654299},
	}
	a.geneGrp[2] = StockRep{
		avg: 6.351454156324596,
		aohlc: [5]float64{0.7847659785141268, 1.0845641762004044, 1.2223651393835664, 0.703007805513547, 0.6100530809191},
		tolerance: [5]float64{0.9563896473078616, 0.7472954550032166, 0.22117474785585847, 0.29282191495007565, 0.6537305586447606},
	}
	a.geneGrp[3] = StockRep{
		avg: 12.862302581872239,
		aohlc: [5]float64{0.9811701498104446, 0.9517443982683742, 1.3158558953681667, 1.1516065481407334, 1.1264424407047682},
		tolerance: [5]float64{0.8646047106226192, 0.8374361923404612, 0.4406040033957521, 0.6037199115830559, 0.43197214102092407},
	}
	a.geneGrp[4] = StockRep{
		avg: 12.565194450499746,
		aohlc: [5]float64{0.5603260172469886, 1.1885503326136915, 1.336969129840545, 0.9070650667413283, 0.8120487275720101},
		tolerance: [5]float64{0.8104877025540802, 0.8044372938101474, 0.8697544026628621, 0.6490316289088398, 0.5269808236191077},
	}
	
	
	for i := 0; i<len(st.Prices); i++{
		if !(a.today >= (st.Chrom[i].today - a.todayTol) && a.today <= (st.Chrom[i].today + a.todayTol)) {
			continue
		}
		allGenesMatch := true
		for k,v := range a.geneGrp {
			stg := st.Chrom[i].geneGrp[k]
			//beego.Debug("stg: ", stg)
			//beego.Debug("test1 ", v.aohlc[0], ">= ", stg.aohlc[0], " - ", v.tolerance[0], " && ", v.aohlc[0], " <= ", stg.aohlc[0]," + ", v.tolerance[0])
			//beego.Debug("test1 ", !(v.aohlc[0] >= stg.aohlc[0] - v.tolerance[0] && v.aohlc[0] <= stg.aohlc[0] + v.tolerance[0]))
			if !(v.aohlc[0] >= stg.aohlc[0] - v.tolerance[0] && v.aohlc[0] <= stg.aohlc[0] + v.tolerance[0]) {
				allGenesMatch = false
			}
			//beego.Debug("test1 ", v.aohlc[1], ">= ", stg.aohlc[1], " - ", v.tolerance[1], " && ", v.aohlc[1], " <= ", stg.aohlc[1]," + ", v.tolerance[1])
			//beego.Debug("test2 ", !(v.aohlc[1] >= stg.aohlc[1] - v.tolerance[1] && v.aohlc[1] <= stg.aohlc[1] + v.tolerance[1]))
			if !(v.aohlc[1] >= stg.aohlc[1] - v.tolerance[1] && v.aohlc[1] <= stg.aohlc[1] + v.tolerance[1]) {
				allGenesMatch = false
			}
			//beego.Debug("test1 ", v.aohlc[2], ">= ", stg.aohlc[2], " - ", v.tolerance[2], " && ", v.aohlc[2], " <= ", stg.aohlc[2]," + ", v.tolerance[2])
			//beego.Debug("test3 ", !(v.aohlc[2] >= stg.aohlc[2] - v.tolerance[2] && v.aohlc[2] <= stg.aohlc[2] + v.tolerance[2]))
			if !(v.aohlc[2] >= stg.aohlc[2] - v.tolerance[2] && v.aohlc[2] <= stg.aohlc[2] + v.tolerance[2]) {
				allGenesMatch = false
			}
			//beego.Debug("test1 ", v.aohlc[3], ">= ", stg.aohlc[3], " - ", v.tolerance[3], " && ", v.aohlc[3], " <= ", stg.aohlc[3]," + ", v.tolerance[3])
			//beego.Debug("test4 ", !(v.aohlc[3] >= stg.aohlc[3] - v.tolerance[3] && v.aohlc[3] <= stg.aohlc[3] + v.tolerance[3]))
			if !(v.aohlc[3] >= stg.aohlc[3] - v.tolerance[3] && v.aohlc[3] <= stg.aohlc[3] + v.tolerance[3]) {
				allGenesMatch = false
			}
			//beego.Debug("test1 ", v.aohlc[4], ">= ", stg.aohlc[4], " - ", v.tolerance[4], " && ", v.aohlc[4], " <= ", stg.aohlc[4]," + ", v.tolerance[4])
			//beego.Debug("test5 ", !(v.aohlc[4] >= stg.aohlc[4] - v.tolerance[4] && v.aohlc[4] <= stg.aohlc[4] + v.tolerance[4]))
			if !(v.aohlc[4] >= stg.aohlc[4] - v.tolerance[4] && v.aohlc[4] <= stg.aohlc[4] + v.tolerance[4]) {
				allGenesMatch = false
			}
		}
		if allGenesMatch {
			priceObj, _ := h.Prices[st.Prices[i].Date]
			priceObj.Buy = 1
			h.Prices[st.Prices[i].Date] = priceObj
		} else {
			priceObj, _ := h.Prices[st.Prices[i].Date]
			priceObj.Buy = 0
			h.Prices[st.Prices[i].Date] = priceObj
		}
	}
	h.MarkAllBadDays()
}
func (h *StockHistoryPrice) MarkAllSignalsTest(){
	s := h.sortedPrices
	rng := len(h.sortedPrices)
	for i := 0; i < rng - 6; i++{
		hit := false
		for j := i+1; j < i + 6; j++{
			if s[j].High >= s[i].Close*1.02 {
				hit = true
			}
		}
		p := h.Prices[s[i].Date]
		if hit {
			p.Test = 1
		} else {
			p.Test = 0
		}
		h.Prices[s[i].Date] = p
	}
}
/*
func (h *StockHistoryPrice) MarkAllSignalsCustom(){
	dna := [10]float64{
		10.735665465212037, -0.14269973541134284,
		10.315327688382006, 0.17753423680381886,
		-4.09634199337719, -3.931301616464415,
		11.495845245539833, -0.4318440719036918,
		-4.314543388941171, 11.761770326458134, // today color today vs tomorrow open
	}
	s := h.sortedPrices
	rng := len(h.sortedPrices)
	buys := 0
	for i := 5; i < rng -1; i++{
		match := true
		k := 0
		for j := i-5; j < i; j++{
			var nextDayUp float64
			var avgToday float64
			var opTomor float64
			var ocRatio  float64
			avgToday = (s[j].Close + s[j].Open) / 2
			opTomor = s[j+1].Open
			nextDayUp = (opTomor/avgToday - 1) * 100
			ocRatio = ((s[j].Close / s[j].Open) - 1) * 100
			if math.Abs(dna[k] - ocRatio) >= 3 && dna[k] < 10{
				match = false
				break
			}
			k++
			if math.Abs(dna[k] - nextDayUp) >= 3 && dna[k] < 10{
				match = false
				break
			}
			k++
		}
		p := h.Prices[s[i+1].Date]
		if match {
			p.Buy = 1
			buys += 1
		} else {
			p.Buy = 0
		}
		h.Prices[s[i+1].Date] = p
	}
	beego.Debug("buys ", buys)
	h.MarkAllBadDays()
}
*/

func capStopLossDays(arrayLen int, current int) int {
	if (arrayLen - current) > 40 {
		return current + 40
	}
	return arrayLen
}

func (h *StockHistoryPrice) RemoveInvalidPrices(){
	validPrices := make(map[string]DayPrice)
	for k,v := range h.Prices{
		if v.Open != 0 && v.High != 0 && v.Low != 0 && v.Close != 0 {
			validPrices[k] = v
		}
	}
	h.Prices = validPrices
}
func (h *StockHistoryPrice) MarkAllBadDays(){
	arrayLen := len(h.sortedPrices)
	signalCount := 0
	successCount := 0
	bads := 0
	for k,v := range h.sortedPrices {
		if k < 3{
			continue
		}
		priceObj, _ := h.Prices[v.Date]
		priceObj.BadDay = 0
		if v.Buy == 1 {
			signalCount ++
			remainingDays := h.sortedPrices[k: capStopLossDays(arrayLen, k)]
			ok := false
			for _, v2 := range remainingDays{
				if v2.High >= v.Open * 1.02 {
					//priceObj, _ := h.Prices[v.Date]
					//priceObj.BadDay = 0
					//h.Prices[v.Date] = priceObj
					ok = true
					successCount ++
					break
				}
			}
			if !ok {
				priceObj.BadDay = 1
				bads += 1
			}
		}
		h.Prices[v.Date] = priceObj
	}
	beego.Debug("buyss ", signalCount)
	beego.Debug("succs ", successCount)
	beego.Debug("bad ", bads)
	//beego.Debug(successCount, " ", signalCount, " ", successCount/signalCount)
	fRatio := float64(successCount)/float64(signalCount) * 100
	iRatio := int(fRatio)
	h.accuracy = strconv.Itoa(successCount) + "/" + strconv.Itoa(signalCount) + " " + strconv.Itoa(iRatio) + "%"
	h.rating = iRatio
}

func getNasdaq(code string, from string, to string) (map[string]DayPrice, error){
	var prices = make(map[string]DayPrice)
	//set POST variables
	beego.Critical(code, "---", from, "---", to)
	data := url.Values{
		"xmlquery": {
			`<post>
				<param name="Exchange" value="NMF"/>
				<param name="SubSystem" value="History"/>
				<param name="Action" value="GetDataSeries"/>
				<param name="AppendIntraDay" value="no"/>
				<param name="Instrument" value="` + code + `"/>
				<param name="FromDate" value="` + from + `"/>
				<param name="ToDate" value="` + to + `"/>
				<param name="hi__a" value="0,3,1,2,4,8"/>
				<param name="ext_xslt" value="/nordicV3/hi_csv.xsl"/>
				<param name="OmitNoTrade" value="true"/>
				<param name="ext_xslt_lang" value="en"/>
				<param name="ext_xslt_options" value=",,"/>
				<param name="ext_contenttype" value="application/ms-excel"/>
				<param name="ext_contenttypefilename" value="` + code + `-` + from + `-` + to + `.csv"/>
				<param name="ext_xslt_hiddenattrs" value=",iv,ip,"/>
				<param name="ext_xslt_tableId" value="historicalTable"/>
				<param name="app" value="/shares/microsite"/>
			</post>`,
		},
	}
	/*data = url.Values{
		"xmlquery": {
			`xmlquery=%3Cpost%3E%0D%0A%3Cparam+name%3D%22Exchange%22+value%3D%22NMF%22%2F%3E%0D%0A%3Cparam+name%3D%22SubSystem%22+value%3D%22History%22%2F%3E%0D%0A%3Cparam+name%3D%22Action%22+value%3D%22GetDataSeries%22%2F%3E%0D%0A%3Cparam+name%3D%22AppendIntraDay%22+value%3D%22no%22%2F%3E%0D%0A%3Cparam+name%3D%22FromDate%22+value%3D%22` + from + `%22%2F%3E%0D%0A%3Cparam+name%3D%22ToDate%22+value%3D%22` + to + `%22%2F%3E%0D%0A%3Cparam+name%3D%22Instrument%22+value%3D%22` + code + `%22%2F%3E%0D%0A%3Cparam+name%3D%22hi__a%22+value%3D%220%2C5%2C6%2C3%2C1%2C2%2C4%2C21%2C8%2C10%2C12%2C9%2C11%22%2F%3E%0D%0A%3Cparam+name%3D%22OmitNoTrade%22+value%3D%22true%22%2F%3E%0D%0A%3Cparam+name%3D%22ext_xslt_lang%22+value%3D%22en%22%2F%3E%0D%0A%3Cparam+name%3D%22ext_xslt%22+value%3D%22%2FnordicV3%2Fhi_csv.xsl%22%2F%3E%0D%0A%3Cparam+name%3D%22ext_xslt_options%22+value%3D%22%2Cadjusted%2C%22%2F%3E%0D%0A%3Cparam+name%3D%22ext_contenttype%22+value%3D%22application%2Fms-excel%22%2F%3E%0D%0A%3Cparam+name%3D%22ext_contenttypefilename%22+value%3D%22TEL2-B-` + from + `-` + to + `.csv%22%2F%3E%0D%0A%3Cparam+name%3D%22ext_xslt_hiddenattrs%22+value%3D%22%2Civ%2Cip%2C%22%2F%3E%0D%0A%3Cparam+name%3D%22ext_xslt_tableId%22+value%3D%22historicalTable%22%2F%3E%0D%0A%3Cparam+name%3D%22DefaultDecimals%22+value%3D%22false%22%2F%3E%0D%0A%3Cparam+name%3D%22app%22+value%3D%22%2Faktier%2Fhistoriskakurser%22%2F%3E%0D%0A%3C%2Fpost%3E`,
		},
	}*/
	req, err := http.NewRequest(
		"POST",
		"http://www.nasdaqomxnordic.com/webproxy/DataFeedProxy.aspx",
		bytes.NewBufferString(data.Encode()),
	)
	//req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	//req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36")
	client := &http.Client{}
	resp, err := client.Do(req)
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error(err)
		resp.Body.Close()
		return prices, errors.New("Nasdaq API not returning result")
	}
	resp.Body.Close()
	lines := strings.Split(string(contents), "\n")
	beego.Debug(len(lines))
	if len(lines) < 1 {
		return prices, errors.New("No data received")
	}
	beego.Debug(lines)
	for _, lineCont := range lines{
		row := strings.Split(lineCont, ";")
		if row[0] == "Date"{
			continue
		}
		if len(row) != 6 {
			beego.Critical(lineCont, " -- ", len(row))
			continue
		}
		for key,_ := range row{
			row[key] = strings.Replace(row[key], ",", ".", -1)
			row[key] = strings.TrimSpace(row[key])
		}
		Date := row[0]
		dt,_ := time.Parse("2006-01-02", Date)
		Dt := dt.Unix() * 1000
		Open, err := strconv.ParseFloat(row[1], 64)
		if err != nil && len(row[1]) > 0 {
			return prices, errors.New("Got invalid open value, unable to convert it"+ row[1]+ ".")
		}
		High, err := strconv.ParseFloat(row[2], 64)
		if err != nil && len(row[2]) > 0 {
			return prices, errors.New("Got invalid high value, unable to convert it"+ row[2]+ ".")
		}
		Low, err := strconv.ParseFloat(row[3], 64)
		if err != nil && len(row[3]) > 0 {
			return prices, errors.New("Got invalid low value, unable to convert it" + row[3]+ ".")
		}
		Close, err := strconv.ParseFloat(row[4], 64)
		if err != nil && len(row[4]) > 0 {
			return prices, errors.New("Got invalid close value, unable to convert it"+ row[4]+ ".")
		}
		Volume, err := strconv.Atoi(row[5])
		if err != nil && len(row[5]) > 0 {
			return prices, errors.New("Got invalid volume value, unable to convert it"+ row[5]+ ".")
		}
		prices[Date] = DayPrice{
			Date: Date,
			Dt: Dt,
			Open: Open,
			High: High,
			Low: Low,
			Close: Close,
			Volume: Volume,
		}
	}
	return prices, nil
}

func getNasdaqToday(code string) (float64, error){
	var price float64
	req, err := http.NewRequest(
		"GET",
		"http://www.nasdaqomxnordic.com/webproxy/DataFeedProxy.aspx?SubSystem=Prices&Action=GetInstrument&Source=OMX&Instrument=" + code + "&inst.an=op&json=1&app=/shares/microsite-ShareInformation",
		nil,
	)
	if err != nil {
		beego.Error(err)
		return price, errors.New("Cannot write an http request")
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36")
	client := &http.Client{}
	resp, err := client.Do(req)
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error(err)
		resp.Body.Close()
		return price, errors.New("Nasdaq API not returning result")
	}
	resp.Body.Close()
	jsonString := strings.Replace(string(contents), "@", "", -1)
	var object struct{
		Status string `json:"status"`
		Ts string `json:"ts"`
		S string `json:"s"`
		Inst struct{
			Id string `json:"id"`
			Op string `json:"op"`
		} `json:"inst"`
	}
	err = json.Unmarshal([]byte(jsonString), &object)
	if err != nil || len(object.Ts) < 1{
		beego.Critical(err)
		return price, err
	}
	beego.Debug("nasdaq today ", object)
	dateToday := time.Now().Truncate(time.Hour * 24).Unix() * 1000
	ts,_ := strconv.ParseInt(object.Ts[:len(object.Ts)-3], 10, 64)
	tsTruncated := time.Unix(ts,0).Truncate(time.Hour * 24).Unix() * 1000
	if dateToday != tsTruncated{
		return price, errors.New("Old data fetched from API")
	}
	openPrice, _ := strconv.ParseFloat(object.Inst.Op, 64)
	if openPrice <= 0 {
		return price, errors.New("Invalid data fetched, or market not open yet")
	}
	price = openPrice
	return price, nil
}
