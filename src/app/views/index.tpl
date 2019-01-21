<div class="container">
	<div class="row">
		<div class="col m12">
			<h3 style="color: #ee6e73; font-weight: 300;">List Of Stocks</h3>
			<h5 style="color: #ee6e73; font-weight: 300;">Automatically updated.</h5>
		</div>
	</div>
	<div class="row">
		<div class="col s12">
			<ul class="collapsible" data-collapsible="accordion">
				{{range $key, $val := .Stocks}}
				<li>
					<div {{if $val.Signal }} style ="background: #4db6ac" {{else if $val.Candidate}} style="background: #fff9c4" {{end}} class="collapsible-header"><!--<i class="material-icons">filter_drama</i>-->{{$val.Name}} -- {{$val.Accuracy}}</div>
					<div class="collapsible-body">
						<ul style="margin: 10px">
							<li>Buy: {{if $val.Signal}}{{$val.Signal}}{{else}}No{{end}}</li>
							<li>Maybe Tomorrow?: {{if $val.Candidate}}Yes if opens lower than {{$val.Candidate}}{{else}}No{{end}}</li>
							<li>Updated: {{$val.Updated}}</li>
							<li>Code: {{$val.Code}}</li>
							<li><a class="waves-effect waves-light btn" onClick="showStock('{{$val.Code}}')">Chart</a></li>
						</ul>
					</div>
				</li>
				{{end}}
			</ul>
		</div>
	</div>
</div>
