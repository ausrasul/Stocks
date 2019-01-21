$(document).ready(function(){
	// the "href" attribute of .modal-trigger must specify the modal ID that wants to be triggered
	$('.modal-trigger').leanModal({
		dismissible: false, // Modal can be dismissed by clicking outside of the modal
		opacity: .5 //, // Opacity of modal background
		//in_duration: 300, // Transition in duration
		//out_duration: 200, // Transition out duration
		//ready: function() { alert('Ready'); }, // Callback for Modal open
		//complete: function() { alert('Closed'); } // Callback for Modal close
	});
	$(".button-collapse").sideNav();
	$('#portfolioFile').change(function(){
		var file = this.files[0];
		var name = file.name;
		var size = file.size;
		var type = file.type;
		//Your validation
		console.log(file);
		console.log(name);
		console.log(size);
		console.log(type);
	});
	$('#portfolioBtn').click(function(){
		var formData = new FormData($('#uploadPortfolio')[0]);
		$.ajax({
			url: '/secure/upload',  //Server script to process data
			type: 'POST',
			xhr: function() {  // Custom XMLHttpRequest
				var myXhr = $.ajaxSettings.xhr();
				if(myXhr.upload){ // Check if upload property exists
					myXhr.upload.addEventListener('progress',progressHandlingFunction, false); // For handling the progress of the upload
				}
				return myXhr;
			},
			//Ajax events
			//beforeSend: beforeSendHandler,
			success: completeHandler,
			error: errorHandler,
			// Form data
			data: formData,
			//Options to tell jQuery not to process data or worry about content-type.
			cache: false,
			contentType: false,
			processData: false
		});
	});
});


function completeHandler(data){
	if (typeof(data) !== "object"){
		printMsg('Unknown error');
		return;
	}
	$('#modal7').closeModal();
	$('#contentHere').children().remove();
	$('<div\>', {id: "portfolio-chart"}).css("width", "100%").css("height", "90%").appendTo($("#contentHere"));
	renderProtfolio(data);
	return;
}

function errorHandler(){
	printMsg('Cannot upload the file');
	return;
}

function progressHandlingFunction(e){
	if(e.lengthComputable){
		$('progress').attr({value:e.loaded,max:e.total});
	}
}
function runCustom(){
	var data = {};
        $.post('/secure/runCustom', '', function(dat){
                if (dat.length > 200){
                        printMsg('Unknown error');
                        return;
                }

                var result = $.parseJSON(dat);
                if (typeof(result) === 'object'){
                        if (result.status === 'OK'){
                                printMsg('Command sent.');
                        } else if (result.status === 'err'){
                                printMsg(result.msg);
                        } else {
                                printMsg('Unknown error');
                        }
                }
        });
}

function addStock(){
	var data = {};
	data.name = $('#name').prop('value');
	data.code = $('#code').prop('value');
	$.post('/secure/addStock', data, function(dat){
		if (dat.length > 200){
			printMsg('Unknown error');
			return;
		}
			
		var result = $.parseJSON(dat);
		if (typeof(result) === 'object'){
			if (result.status === 'OK'){
				printMsg('Stock was added successfully');
			} else if (result.status === 'err'){
				printMsg(result.msg);
			} else {	
				printMsg('Unknown error');
			}
		}
	});
}


function parsePortfolio(){
	var data = {};
	data.csvstring = $('#csvstring').prop('value');
	$.post('/secure/parsePortfolio', data, function(dat){
		if (typeof(dat) === 'object'){
			renderProtfolio(dat);
		} else {
			var result = $.parseJSON(dat);
			if (typeof(result) === 'object'){
				if (result.status === 'OK'){
					printMsg('Stock was added successfully');
				} else if (result.status === 'err'){
					printMsg(result.msg);
				} else {	
					printMsg('Unknown error');
				}
			}
		}
	});
}

function syncStocks(){
	$.post('/secure/syncStocks', '', function(dat){
		if (dat.length > 500){
			printMsg('Unknown error');
			$('#modal3').closeModal();
			return;
		}
			
		var result = $.parseJSON(dat);
		if (typeof(result) === 'object'){
			if (result.status === 'OK'){
				printMsg('Stocks was added successfully synchronized');
			} else if (result.status === 'err'){
				printMsg(result.msg);
			} else {	
				printMsg('Unknown error');
			}
		}
		$('#modal3').closeModal();
	});
}

function forecast(){
	$.post('/secure/forecast', '', function(dat){
		if (dat.length > 500){
			printMsg('Unknown error');
			$('#modal4').closeModal();
			return;
		}
			
		var result = $.parseJSON(dat);
		if (typeof(result) === 'object'){
			if (result.status === 'OK'){
				printMsg('Candidates are marked');
			} else if (result.status === 'err'){
				printMsg(result.msg);
			} else {	
				printMsg('Unknown error');
			}
		}
		$('#modal4').closeModal();
	});
}
function getSignal(){
	$.post('/secure/getSignal', '', function(dat){
		if (dat.length > 500){
			printMsg('Unknown error');
			$('#modal5').closeModal();
			return;
		}
			
		var result = $.parseJSON(dat);
		if (typeof(result) === 'object'){
			if (result.status === 'OK'){
				printMsg('Stocks are picked');
			} else if (result.status === 'err'){
				printMsg(result.msg);
			} else {	
				printMsg('Unknown error');
			}
		}
		$('#modal5').closeModal();
	});
}

function analyzeHistory(){
	$.post('/secure/analyzeHistory', '', function(dat){
		if (dat.length > 500){
			printMsg('Unknown error');
			$('#modal6').closeModal();
			return;
		}
			
		var result = $.parseJSON(dat);
		if (typeof(result) === 'object'){
			if (result.status === 'OK'){
				printMsg('All hisotry data analyzed.');
			} else if (result.status === 'err'){
				printMsg(result.msg);
			} else {	
				printMsg('Unknown error');
			}
		}
		$('#modal6').closeModal();
	});
}

function showStock(code){
	window.open("/secure/preload/" + code);
	//window.location.href = "/secure/preload/" + code;
}
function printMsg(msg){
	Materialize.toast(msg, 5000, 'rounded');
}
