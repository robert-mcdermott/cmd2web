# cmd2web - Command to Web

```
Usage: ./cmd2web [--expose <path> --expire <minutes> --refresh <seconds> --raw] <command>

Flags:

  -expire int
    
        [optional] terminate the cmd2web server after the provide number of
        minutes. If an expiration is not provide the server will run indefinately
        until terminated manually
    
  -expose string
    
        [optional] expose this directory or file at https://*/file
        if a directory path is given it will provide an html file/dir listing
        that you can navigate files and sub directories. if a file path is
        provided, the file will be availible at the file URL
    
  -help
    
        print usage information
    
  -raw
    
        [optional] the default output is html; this flag enables raw text
        output that is more suitable for use with curl or using as input to
        another program or logging.
    
  -refresh int
    
        [optional] page refresh interval in seconds; only works with html
        output format with GUI web browsers (Chrome, Firefox, etc...). each
        refresh re-runs the command.
    

Example 1: list they systems process table and refresh the output every 30 seconds.

        ./cmd2web --refresh 30 ps aux

Example 2: expose the "myproject" directory to the web for 60 minutes.

        ./cmd2web --expire 60 --expose /home/rmcdermo/myproject /usr/bin/true
```