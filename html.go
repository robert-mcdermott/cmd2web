package main

var cmdhtml = `
	<!DOCTYPE html>
	<html>
		<head>
			<style>
				body {background-color: black;}
				h3   {color: #fc9a85;}
				pre  {color: #b1fc85;}
			</style>
		</head>
		<body>
			<h3>Command: %s </h3>
			<h3>Time: %s </h3>
			<p><pre>%s</pre></p>
		</body>
	</html>`

var four04html = `
	<!DOCTYPE html>
	<html>
		<head>
			<style>
				body {background-color: black;}
				h3   {color: #fc9a85;}
			</style>
		</head>
		<body>
			<h3>%s</h3>
		</body>
	</html>`
