#!/bin/bash
FROM="2016-08-01"
TO="2016-08-30"
CODE="SSE999"
URL=http://www.nasdaqomxnordic.com/webproxy/DataFeedProxy.aspx

#H1='X-Requested-With: XMLHttpRequest'
H2='User-Agent: Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36'
#H3='Content-Type: application/x-www-form-urlencoded'
QUERY="xmlquery=%3Cpost%3E%0D%0A%3Cparam+name%3D%22Exchange%22+value%3D%22NMF%22%2F%3E%0D%0A%3Cparam+name%3D%22SubSystem%22+value%3D%22History%22%2F%3E%0D%0A%3Cparam+name%3D%22Action%22+value%3D%22GetDataSeries%22%2F%3E%0D%0A%3Cparam+name%3D%22AppendIntraDay%22+value%3D%22no%22%2F%3E%0D%0A%3Cparam+name%3D%22FromDate%22+value%3D%22${FROM}%22%2F%3E%0D%0A%3Cparam+name%3D%22ToDate%22+value%3D%22${TO}%22%2F%3E%0D%0A%3Cparam+name%3D%22Instrument%22+value%3D%22${CODE}%22%2F%3E%0D%0A%3Cparam+name%3D%22hi__a%22+value%3D%220%2C5%2C6%2C3%2C1%2C2%2C4%2C21%2C8%2C10%2C12%2C9%2C11%22%2F%3E%0D%0A%3Cparam+name%3D%22OmitNoTrade%22+value%3D%22true%22%2F%3E%0D%0A%3Cparam+name%3D%22ext_xslt_lang%22+value%3D%22en%22%2F%3E%0D%0A%3Cparam+name%3D%22ext_xslt%22+value%3D%22%2FnordicV3%2Fhi_csv.xsl%22%2F%3E%0D%0A%3Cparam+name%3D%22ext_xslt_options%22+value%3D%22%2Cadjusted%2C%22%2F%3E%0D%0A%3Cparam+name%3D%22ext_contenttype%22+value%3D%22application%2Fms-excel%22%2F%3E%0D%0A%3Cparam+name%3D%22ext_contenttypefilename%22+value%3D%22TEL2-B-${FROM}-${TO}.csv%22%2F%3E%0D%0A%3Cparam+name%3D%22ext_xslt_hiddenattrs%22+value%3D%22%2Civ%2Cip%2C%22%2F%3E%0D%0A%3Cparam+name%3D%22ext_xslt_tableId%22+value%3D%22historicalTable%22%2F%3E%0D%0A%3Cparam+name%3D%22DefaultDecimals%22+value%3D%22false%22%2F%3E%0D%0A%3Cparam+name%3D%22app%22+value%3D%22%2Faktier%2Fhistoriskakurser%22%2F%3E%0D%0A%3C%2Fpost%3E"
#curl ${URL} -H "$H1" -H "$H2" -H "$H3" --data "$QUERY"
curl ${URL} -H "$H2" --data "$QUERY"
