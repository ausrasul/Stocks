function sim(){
	var chart = $('#container').highcharts();
	var funds = [];
	var buying = [],
		selling = [],
		expired = [],
		stoploss = [];//var actions = []; // 0=nothing, 1=buy, 2=sell with profit, -1=force sell, -2=stoploss
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
	
	
	
	var cash = new Cash();
	var buyOrder = new Order();
	var sellOrder = new Order();
	var stock = new Stock();
	var state = 'waiting candidate';
	var skip = false;
	var res = 0;
	var cooldown = 0;
	var startCapital = cash.cash;
	for (var i=0; i< from; i++){
		funds.push([ohlc[i][0],cash.cash]);
		buying.push([ohlc[i][0],0]);
		selling.push([ohlc[i][0],0]);
		expired.push([ohlc[i][0],0]);
		stoploss.push([ohlc[i][0],0]);
	}
	
	
	
	for (var i=from; i<=to; i++){
		stock.open = ohlc[i][1];
		stock.high = ohlc[i][2];
		stock.low = ohlc[i][3];
		stock.close = ohlc[i][4];
		stock.candidate = buy[i][1];
		buying.push([ohlc[i][0],0]);
		selling.push([ohlc[i][0],0]);
		expired.push([ohlc[i][0],0]);
		stoploss.push([ohlc[i][0],0]);

		skip = false; // new day
		var sameday = false;
		if (state === 'waiting candidate' && !skip){
			if (!cash.haveShares() && stock.candidate == 1){
				// buy stocks
				buyOrder.placeBuyOrder(cash, stock);
				state = 'waiting to buy';
				sameday = false;
				//skip = true; // skip to next business day
			}
		}
		
		if (state === 'waiting to buy' && !skip){
			// if waiting to execute buy order
			if (buyOrder.placed && !cash.haveShares()){
				res = buyOrder.executeBuy(cash, stock, sameday);
				if (res === 'bought'){
					buying.push([ohlc[i][0],1]);
					sellOrder.placeSellOrder(cash);
					state = 'waiting to sell';
				} else if (res === 'expired'){
					state = 'waiting candidate';
				}
				//skip = true;
			} 
		}
		
		if (state === 'waiting to sell' && !skip){
			if (sellOrder.placed && cash.haveShares()) {
				res = sellOrder.executeSell(cash, stock);
				if (res === 'sold') selling.push([ohlc[i][0],1]);
				if (res === 'force sell') expired.push([ohlc[i][0],1]);
				if (res === 'stoploss') stoploss.push([ohlc[i][0],1]);
				if (res === 'sold' || res === 'force sell'|| res === 'stoploss'){
					state = 'waiting candidate';
				}
				//skip = true;
				cooldown = 1;
			}
		}
			
		
		// another day
		//if (days > 0) days--;
		if (cash.cash < stock.open && !cash.haveShares()){
			console.log("Bankrupt!!!");
			return;
		}
		funds.push([ohlc[i][0],cash.cash + (cash.shares * stock.open)]);
		if (cooldown > 0) cooldown--;
	}
	

	// Draw the current portfolio chart for days after trading is finished (fill chart)
	for (var i=to+1; i< dataLength; i++){
		funds.push([ohlc[i][0],cash.cash + (cash.shares * cash.boughtfor)]);
		buying.push([ohlc[i][0],0]);
		selling.push([ohlc[i][0],0]);
		expired.push([ohlc[i][0],0]);
		stoploss.push([ohlc[i][0],0]);
	}
	var percentProfit = ((cash.getNetValue() / startCapital) - 1) * 100;
	document.getElementById('prcnt').innerHTML = 'Profit/loss: ' + Math.round(percentProfit) + ' % ';
	// Refresh the portfolio chart.
	chart.series[5].setData(funds);
	chart.series[6].setData(buying);
	chart.series[7].setData(selling);
	chart.series[8].setData(expired);
	chart.series[9].setData(stoploss);
}



function Order(){
	this.placed = false;
	this.age = 0;
	this.price = 0;
	this.maxAge = 40;
}
Order.prototype = {
	placeBuyOrder: function(cash, stock){
		this.placed = true;
		this.age = 0;
		this.price = stock.open;
	},
	executeBuy: function(cash, stock, sameday){
		//if (sameday) return 'nothing';
		console.log("buying stocks for " + stock.open);
		while (cash.cash > stock.open){
			cash.shares++;
			cash.cash -= stock.open;
		}
		cash.boughtfor = stock.open;
		cash.shareAge = 0;
		cash.takeBuyCourtage();
		this.placed = false;
		this.age = 0;
		this.price = 0;
		return 'bought';
	},
	placeSellOrder: function(cash, stock){
		this.placed = true;
		this.age = 0;
		this.price = cash.goal();
		
	},
	executeSell: function(cash, stock){
		
		if (stock.high >= this.price){
			console.log("selling stocks for " + this.price);
			while (cash.shares > 0){
				cash.shares--;
				cash.cash += this.price;
			}
			cash.shareAge = 0;
			cash.takeSellCourtage();
			this.placed = false;
			this.age = 0;
			this.price = 0;
			return 'sold';
		}
		else if (this.age > this.maxAge && this.placed){
			this.placed = false;
			this.price = 0;
			
			console.log("force sell for " + stock.close);
			while (cash.shares > 0){
				cash.shares--;
				cash.cash += stock.close;
			}
			cash.shareAge = 0;
			cash.takeSellCourtage();
			this.placed = false;
			this.age = 0;
			this.price = 0;
			return 'force sell';
		} else if (stock.open <= cash.stoploss()){
			this.placed = false;
			this.price = 0;
			
			console.log("stoploss for " + stock.close);
			while (cash.shares > 0){
				cash.shares--;
				cash.cash += stock.close;
			}
			cash.shareAge = 0;
			cash.takeSellCourtage();
			this.placed = false;
			this.age = 0;
			this.price = 0;
			return 'stoploss';
		} else if (this.age <= this.maxAge){
			this.age++;
			//return 'nothing';
		}
		return false;
	},	
}

function Stock(){
	this.high = 0;
	this.low = 0;
	this.candidate = 0;
}
Stock.prototype = {
	profitable: function(goal){ // unused function
		if (this.open >= goal || this.close >= goal) return true;
		return false;
	}
};
function Cash(){
	this.cash = 20000;
	this.shares = 0;
	this.boughtfor = 0;
	this.sharesAge = 0;
}

Cash.prototype = {
	getNetValue: function(){
		return this.cash + (this.shares * this.boughtfor);
	},
	goal: function(){
		return this.boughtfor * 1.02;
	},
	stoploss: function(){
		return this.boughtfor * 0.2;
	},
	haveShares: function(){
		if (this.shares > 0) return true;
		return false;
	},
	takeBuyCourtage: function(){
		this.cash -= ((this.shares * this.boughtfor) * 0.0002)+5;
	},
	takeSellCourtage: function(){
		this.cash -= (this.cash * 0.0002)+5;
	}
};






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
