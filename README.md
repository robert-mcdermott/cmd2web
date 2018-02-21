# Cmd2web - Command to Web Utility

![cmd2web example](images/cmd2web-browser.png)


## What is Cmd2web?

Cmd2web is a utility that allows you to execute a command on a system and then securely view the output of that command on any other system via a web browser. Each time the page is reloaded the command is re-run and the output is updated. You can set the refresh interval that the command will be automatically re-run by you web browser.  

In addition to command output you can also also optionally expose a directory or file via Cmd2web and make it accessible via a web browser on a remove system with the --expose flag. If you provide it a path to a directory, you'll get an html directory listing that lets you navigate sub-directories and view files. If a path to a file is provided path, only that file will be available. If the exposed directory contains an 'index.html' file, it will be exposed rather than showing a directory listing.

An expiration timer can be optionally set that will stop the Cmd2web server when the timer expires. 



## Usage

```
Usage: ./cmd2web [--expose <path> --expire <minutes> --refresh <seconds> --raw --noauth] <command>

Flags:

  -expire int
    	
    	[optional] terminate the cmd2web server after the provide number of
    	minutes. If an expiration is not provide the server will run indefinitely
    	until terminated manually
    	
  -expose string
    	
    	[optional] expose this directory or file at https://*/file
    	if a directory path is given it will provide an html file/dir listing
    	that you can navigate files and sub directories. if a file path is
    	provided, the file will be available at the file URL
    	
  -help
    	
    	print usage information
    	
  -noauth
    	
    	[optional] disable basic (user/pass) authentication. By default
    	authentication is enabled
    	
  -port int
    	
    	[optional] specify a tcp port to listen on. By default a random port is
    	selected for you; this flag overrides that behavior. You must be root to assign a port below 1024
    	
  -raw
    	
    	[optional] the default output is html; this flag enables raw text
    	output that is more suitable for use with curl or using as input to
    	another program or logging.
    	
  -refresh int
    	
    	[optional] page refresh interval in seconds; only works with html
    	output format with GUI web browsers (Chrome, Firefox, etc...). each
    	refresh re-runs the command.
    	

Example 1: list the systems process table and refresh the output every 30 seconds.

	./cmd2web --refresh 30 ps aux

Example 2: expose the "myproject" directory to the web for 60 minutes.

	./cmd2web --expire 60 --expose /home/rmcdermo/myproject /usr/bin/true

```


## Accessing the Cmd2web server

After starting a Cmd2web server it will provide the required connection information on the console's standard error.  Here is an example:

```
Access Information
-------------------------------------
Command output:    https://test.rigel.net:42938/2EeRbdGXIwMaYqZQ3dk611n98piGVbp7
Remote stop:       https://test.rigel.net:42938/2EeRbdGXIwMaYqZQ3dk611n98piGVbp7/stop
Exposed directory: https://test.rigel.net:42938/2EeRbdGXIwMaYqZQ3dk611n98piGVbp7/file

Credentials:

  Username: cmd2web
  Password: fLPlehdO

Easy Access URL:   https://cmd2web:fLPlehdO@test.rigel.net:42938/2EeRbdGXIwMaYqZQ3dk611n98piGVbp7
-------------------------------------
```

## Is this secure?

The connection is SSL encrypted (AES-128/TLS 1.2), a 32 byte random path is required to access the command-output/exposed-files and authentication (username/password) is required by default. In addition, the application takes no input via the exposed web site so there is no way to inject any malicious input.

The SSL certificate included in the source/binaries is a self-singed certificate. If you want to use a valid certificate, update the 'ssl.go' source file with your certificate and recompile for your desired platform. 
