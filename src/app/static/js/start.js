conf = {
	'maxAge': 6,
	'target': 1.0142833,
	'stoploss': 0.999411,
	'courtage': 0.0002,
	'capital': 20000,
	'workerCount': 1
};

function sim(){
	var chart = $('#container').highcharts();
	var funds = [];
	var buying = [],
		selling = [],
		expired = [],
		stoploss = [];
	// Get start and end trading dates
	var from = chart.xAxis[0].min;
	var to = chart.xAxis[0].max;
	var ok = false;
	for (var i = 1; i < dataLength; i++){
		if (ohlc[i][0] > from){
			from = i;
			ok = true;
			break;
		}
	}
	if (!ok) from = 1;
	ok = false;
	for (var i = 0; i < dataLength; i++){
		if (ohlc[i][0] >= to){
			to = i;
			ok = true;
			break;
		}
	}
	if (!ok) to = dataLength;

	// Initialize
	
	var stock = new Stock();
	var workers = new Workers();
	
	// Draw flat graph before the selected period.

	for (var i=0; i< from; i++){
		funds.push([ohlc[i][0],workers.totalValue(0)]);
		buying.push([ohlc[i][0],0]);
		selling.push([ohlc[i][0],0]);
		expired.push([ohlc[i][0],0]);
		stoploss.push([ohlc[i][0],0]);
	}
	
	
	// Simulate selected period

	
	for (var i=from; i<=to; i++){
		stock.open = ohlc[i][1];
		stock.high = ohlc[i][2];
		stock.low = ohlc[i][3];
		stock.close = ohlc[i][4];
		stock.candidate = buy[i][1];
		
/*		buying.push([ohlc[i][0],0]);
		selling.push([ohlc[i][0],0]);
		expired.push([ohlc[i][0],0]);
		stoploss.push([ohlc[i][0],0]);
*/		
		// See if you can sell anything today
		var buyCode = 0;
		if (stock.candidate == 1){
			buyCode = workers.tryBuyFor(stock.open); // return 1 if bought something, 0 if no cash to buy.
		}
		var sellCode = workers.trySellFor(stock.high, stock.close); // try to sell return value {none: 0, profit: 0, stoploss: 0, expired: 0}
		
		funds.push([ohlc[i][0],workers.totalValue(stock.close)]);
		buying.push([ohlc[i][0],buyCode]);
		selling.push([ohlc[i][0],sellCode.profit]);
		expired.push([ohlc[i][0],sellCode.expired]);
		stoploss.push([ohlc[i][0],sellCode.stoploss]);
	}
	

	// Draw the current portfolio chart for days after trading is finished (fill chart)
	for (var i=to+1; i< dataLength; i++){
		funds.push([ohlc[i][0],workers.totalValue(ohlc[i][4])]);
		buying.push([ohlc[i][0],0]);
		selling.push([ohlc[i][0],0]);
		expired.push([ohlc[i][0],0]);
		stoploss.push([ohlc[i][0],0]);
	}
	var percentProfit = ((workers.totalValue(ohlc[dataLength-1][4]) / conf.capital) - 1) * 100;
	document.getElementById('prcnt').innerHTML = 'Profit/loss: ' + Math.round(percentProfit) + ' % ';
	// Refresh the portfolio chart.
	chart.series[5].setData(funds);
	chart.series[6].setData(buying);
	chart.series[7].setData(selling);
	chart.series[8].setData(expired);
	chart.series[9].setData(stoploss);
}

function Worker(value){
	this.cash = value;
	this.stockCnt = 0;
	this.stockVal = 0;
	this.daysLeft = 0;
	this.target = 0;
	this.stoploss = 0;
	this.busy = false;
}
Worker.prototype = {
	payCourtage: function(orderValue){
		
	},
	buy: function(price){
		if (!this.busy && this.cash >= price){
			this.reset();
			this.stockVal = price;
			while (this.cash >= price){
				this.stockCnt++;
				this.cash -= price;
			}
			this.target = price * conf.target;
			this.stoploss = price * conf.stoploss;
			this.daysLeft = conf.maxAge;
			this.payCourtage(this.stockCnt * this.stockVal);
			this.busy = true;
			return true;
		}
		return false;
	},
	sell: function(price){
		if (this.busy){
			this.cash += this.stockCnt * price;
			this.payCourtage(this.stockCnt * price);
			this.reset();
			this.busy = false;
		}
	},
	reset: function(){
		this.stockCnt = 0;
		this.stockVal = 0;
		this.daysLeft = 0;
		this.target = 0;
		this.stoploss = 0;
	}

}

function Workers(){
	this.workers = [];
	for (var i = 0; i < conf.workerCount; i++){
		var w = new Worker(conf.capital / conf.workerCount);
		this.workers.push(w);
	}
}
Workers.prototype = {
	trySellFor: function(high, close){ // must run every day to refresh age counter.
		var returnVal ={none: 0, profit: 0, stoploss: 0, expired: 0};
		$.each(this.workers, function(i, w){
			if (w.busy){
				this.workers[i].daysLeft = w.daysLeft - 1;
				if (high >= w.target){
					w.sell(w.target);
					returnVal.profit++;
				} else if (close <= this.stoploss){
					w.sell(close);
					returnVal.stoploss++;
				} else if (w.daysLeft == 0) {
					w.sell(close);
					returnVal.expired++;
				}
				returnVal.none++;
			}
		}.bind(this));
		this.refill();
		return returnVal;
	},
	tryBuyFor: function(open){
		var bought = 0;
		$.each(this.workers, function(i, w){
			if (!w.busy && w.cash >= open){
				w.buy(open);
				bought = 1;
				return false;
			}
		});
		return bought;
	},
	totalValue: function(price){
		var value = 0;
		$.each(this.workers, function(i,w){
			value += w.cash;
			if (w.busy){
				value += w.stockCnt * price;
			}
		});
		return value;
	},
	refill: function(){
		var cash = 0;
		var idleWorkers = 0;
		$.each(this.workers, function(i,w){
			if (!w.busy){
				cash += w.cash;
				idleWorkers++;
			}
		});
		
		var cashToWrkr = cash / idleWorkers;
		$.each(this.workers, function(i,w){
			if (!w.busy){
				this.workers[i].cash = cashToWrkr;
			}
		}.bind(this));
	}
};

function Stock(){
	this.high = 0;
	this.open = 0;
	this.close = 0;
	this.low = 0;
	this.candidate = 0;
}


function load(id) {
	$.getJSON('/secure/load/' + id, function (data) {
		//data = $.parseJSON(data);
		// split the data set into ohlc and volume
		ohlc = [],
		volume = [],
		open = [],
		close = [],
		high = [],
		low = [],
		lowClose = [],
		highOpen = [],
		openClose = [],
		buy = [],
		test = [],
		badDay = [],
		funds = [],
		good = [],
		mapping = [],
		action  = [],
		buying = [],
		selling = [],
		expired = [],
		stoploss = [];
		dataLength = data.length;
		// set the allowed units for data grouping
		/*groupingUnits = [[
				'week',                         // unit name
				[1]                             // allowed multiples
			], [
				'month',
				[1, 2, 3, 4, 6]
			]],*/
		groupingUnits = [[
				'week',
				[1]
			], [
				'month',
				[6]
		]];

		i = 0;

		for (var i=0; i < dataLength; i += 1) {
			//if (data[i]['buy'] === 1) console.log(1);
			//mapping[data[i][0]] = i;
			ohlc.push([
				data[i]['dt'], // the date
				data[i]['open'], // open
				data[i]['high'], // high
				data[i]['low'], // low
				data[i]['close'] // close
			]);

			volume.push([
				data[i]['dt'], // the date
				data[i]['volume'] // the volume
			]);
			
			/*lowClose.push([
				data[i]['dt'], // date
				data[i]['lowclose'] // low - close ratio
			]);
			highOpen.push([
				data[i]['dt'],
				data[i]['highopen']
			]);
			openClose.push([
				data[i]['dt'],
				data[i]['openclose']
			]);*/
			buy.push([
				//data[i][0], // date
				data[i]['dt'],
				data[i]['buy'] // low - close ratio
			]);
			test.push([
				data[i]['dt'],
				data[i]['test']
			]);
			badDay.push([
				data[i]['dt'],
				data[i]['badDay']
			]);
			funds.push([
				//data[i][0], // date
				data[i]['dt'],
				20000 // low - close ratio
			]);
			/*good.push([
				data[i]['dt'],
				data[i]['goodDay']
			]);
			open.push([
				data[i]['dt'], // date
				data[i]['open'] // low - close ratio
			]);
			high.push([
				data[i]['dt'], // date
				data[i]['high'] // low - close ratio
			]);
			low.push([
				data[i]['dt'], // date
				data[i]['low'] // low - close ratio
			]);
			close.push([
				data[i]['dt'], // date
				data[i]['close'] // low - close ratio
			]);*/
			buying.push([
				data[i]['dt'],
				0
			]);
			selling.push([
				data[i]['dt'],
				0
			]);
			expired.push([
				data[i]['dt'],
				0
			]);
			stoploss.push([
				data[i]['dt'],
				0
			]);
		}
		// create the chart
		$('#container').highcharts('StockChart', {
			rangeSelector: { selected: 0 },
			title: { text: 'Historical' },
			yAxis: [{
					labels: { align: 'right', x: -3 },
					title: { text: 'Stock' },
					height: '30%',
					lineWidth: 2
				}, {
					labels: { align: 'right', x: -3 },
					title: { text: 'Volume' },
					top: '32%',
					height: '10%',
					offset: 0,
					lineWidth: 2
				}, {
					labels: { align: 'right', x: -3 },
					title: { text: 'Indicators' },
					top: '44%',
					height: '20%',
					offset: 0,
					lineWidth: 2
				}, {
					labels: { align: 'right', x: -3 },
					title: { text: 'Portfolio' },
					top: '66%',
					height: '24%',
					offset: 0,
					lineWidth: 2
				}, {
					labels: { align: 'right', x: -3 },
					title: { text: 'Action' },
					top: '92%',
					height: '8%',
					offset: 0,
					lineWidth: 2
				}
			],
			series: [{
				type: 'candlestick',
				name: 'Stock',
				data: ohlc,
				dataGrouping: {units: groupingUnits}
			}, {
				type: 'column',
				name: 'Volume',
				data: volume,
				yAxis: 1,
				dataGrouping: {units: groupingUnits}
			}, {
				type: 'column',
				name: 'Buy!',
				data: buy,
				yAxis: 2,
				dataGrouping: {units: groupingUnits}
			}, {
				type: 'column',
				name: 'BadDay',
				data: badDay,
				yAxis: 2,
				dataGrouping: {units: groupingUnits}
			}, {
				type: 'column',
				name: 'Test',
				data: test,
				yAxis: 2,
				dataGrouping: {units: groupingUnits}
			}, {
				type: 'line',
				name: 'Protfolio',
				data: funds,
				yAxis: 3,
				dataGrouping: {units: groupingUnits}
			}, {
				type: 'column',
				name: 'Buying',
				data: buying,
				yAxis: 4,
				dataGrouping: {units: groupingUnits}
			}, {
				type: 'column',
				name: 'Selling',
				data: selling,
				yAxis: 4,
				dataGrouping: {units: groupingUnits}
			}, {
				type: 'column',
				name: 'Expired',
				data: expired,
				yAxis: 4,
				dataGrouping: {units: groupingUnits}
			}, {
				type: 'column',
				name: 'Stoploss',
				data: stoploss,
				yAxis: 4,
				dataGrouping: {units: groupingUnits}
			}]
		});
	});
}

function renderProtfolio(object) {
	groupingUnits = [[
			'week',
			[1]
		], [
			'month',
			[6]
	]];
	var dt = 0;
	var s = 0;
	var b = 0;
	
	var portfolio = [];
	var buying = [];
	var selling = [];
	var dataLength = object.length;
	for (var i=0; i < dataLength; i += 1) {
		//if (data[i]['buy'] === 1) console.log(1);
		//mapping[data[i][0]] = i;
		tmpdt = object[i]['dt'];
		
		portfolio.push([
			object[i]['dt'], // the date
			object[i]['cash'] // the cash
		]);
		if (tmpdt === dt){
			s = s | object[i]['sell'];
			b = b | object[i]['buy'];
		} else {
			dt = tmpdt;
			s = object[i]['sell'];
			b = object[i]['buy'];
		}
		buying.push([
			tmpdt,
			b
		]);
		selling.push([
			tmpdt,
			s
		]);
	}
	// create the chart
	$('#portfolio-chart').highcharts('StockChart', {
		rangeSelector: { selected: 0 },
		title: { text: 'Your Degiro Portfolio' },
		xAxis: [{
            type: 'datetime',
			crosshair: true
        }],
		tooltip: {
			shared: false
		},
		yAxis: [{
				labels: { align: 'right', x: -3 },
				title: { text: 'Portfolio' },
				height: '90%',
				lineWidth: 2
			}, {
				labels: { align: 'right', x: -3 },
				title: { text: 'Actions' },
				top: '90%',
				height: '10%',
				offset: 0,
				lineWidth: 1
			}
		],
		series: [{
			type: 'spline',
			name: 'Protfolio',
			data: portfolio,
			yAxis: 0,
			dataGrouping: {units: groupingUnits}
		},{
			type: 'column',
			name: 'Buy',
			data: buying,
			yAxis: 1,
			dataGrouping: {units: groupingUnits}
		}, {
			type: 'column',
			name: 'Sell',
			data: selling,
			yAxis: 1,
			dataGrouping: {units: groupingUnits}
		}]
	});
}
