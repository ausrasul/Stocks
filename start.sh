#!/bin/bash -x
p=`ps -ef  |grep stocks |grep -v grep|grep -v start|wc -l`
if [[ $p -eq 0 ]]
then
	cd /opt/stocks/src/app;./stocks >> /tmp/stocks.log &
fi

