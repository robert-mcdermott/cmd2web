# Cmd2web - Command to Web Utility

![cmd2web example](images/cmd2web-browser.png)


## What is Cmd2web?

Cmd2web is a utility that allows you to execute a command on a system and then securely view the output of that command on any other system via a web browser. Each time the page is reloaded the command is re-run and the output is updated. You can set the refresh interval that the command will be automatically re-run by you web browser.  

In addition to command output you can also also optionally expose a directory or file via Cmd2web and make it accessable via a web browser on a remove system with the --expose flag. If you provide it a path to a directory, you'll get an html directory listing that lets you navigate sub-directories and view files. If a path to a file is provided path, only that file will be availible. If the exposed directory contains an 'index.html' file, it will be exposed rather than showing a directory listing.

An expiration timer can be optionally set that will stop the Cmd2web server when the timer expires. 



## Usage

```
Usage: ./cmd2web [--expose <path> --expire <minutes> --refresh <seconds> --raw] <command>

Flags:

  --expire int
    
        [optional] terminate the cmd2web server after the provide number of
        minutes. If an expiration is not provide the server will run indefinately
        until terminated manually
    
  --expose string
    
        [optional] expose this directory or file at https://*/file
        if a directory path is given it will provide an html file/dir listing
        that you can navigate files and sub directories. if a file path is
        provided, the file will be availible at the file URL
    
  --help
    
        print usage information
    
  --raw
    
        [optional] the default output is html; this flag enables raw text
        output that is more suitable for use with curl or using as input to
        another program or logging.
    
  --refresh int
    
        [optional] page refresh interval in seconds; only works with html
        output format with GUI web browsers (Chrome, Firefox, etc...). each
        refresh re-runs the command.
    

Example 1: list they systems process table and refresh the output every 30 seconds.

        ./cmd2web --refresh 30 ps aux

Example 2: expose the "myproject" directory to the web for 60 minutes.

        ./cmd2web --expire 60 --expose /home/rmcdermo/myproject /usr/bin/true
```


## Accessing the Cmd2web server

After starting a Cmd2web server it will provide the required connection information on the console's standard error.  Here is an example:

```
Access Information
-------------------------------------
Command output:    https://test.rigel.net:53208/IWYrKyhDVmWWWFSlmQKnDP82oTSfh9Wc
Remote stop:       https://test.rigel.net:53208/IWYrKyhDVmWWWFSlmQKnDP82oTSfh9Wc/stop
Exposed directory: https://test.rigel.net:53208/IWYrKyhDVmWWWFSlmQKnDP82oTSfh9Wc/file
Username: cmd2web
Password: hYe9SdYi

Easy Access URL:   https://cmd2web:hYe9SdYi@test.rigel.net:53208/IWYrKyhDVmWWWFSlmQKnDP82oTSfh9Wc
-------------------------------------
```

## Is this secure?

The connection is SSL encrypted, a 32 byte random access key path is required to access the command output and exposed files and authenitcation (username/password) is required. In addition the application takes no input via the exposed web site so there is no way to inject any malicious commands. 
