# fake-take                                                  

![Alt text](assets/fake-take-avatar.png?raw=true "Fake Take Logo")
Simple HTTP API read only server to test Single Page Application

Simulates HTTP API for testing Single Page Application, e.g. VueJs app. 
Typically a reply is a 'application/json' format but you can use anything you want by passing the '-contentType' param.  Allows only "GET" requests. 

## Setup

API root path by default '/api/v1.0/'  to set another root path use '-apiroot' param and set any prefix you want.

### Adding paths and response bodies

Create a file with a *.resp extension. Name it follow the rule below.
Any part of the name separated by one underscore '_' treated as url part before the slash '/'.
Example:
events_recent.resp will be transalted to /api/v1.0/events/recent    

_No query params supported!_

You can crate as many '*.resp' files as you want.
New files can be added during server work, no need to restart.
You can even change the contents in '*.resp' file - server always reads the file at request time.

Error 404 reported if '*.resp' file not exists.  

## Usage 

To see usage notes use `-help` flag.

You can set:
* -addr string - to which IP bind the server (default `localhost`), to bind all interfaces place 0.0.0.0
* -apiroot string - relative to the server root path (default `/api/v1`)
* -contentType string - content type for response (default "application/json")
* -ext string - extension for data files (default "resp")
* -port int  - server port to listen at (default 3000)
 
Compile the app with `go install` or just run with `go run main.go`

By default server starts at `localhost:3000` and API available at `http://localhost:3000/api/v1`.
Root '/' path always redirects to apiRoot.
At api root `http://localhost:3000/api/v1` you can see clickable link list of available API paths.

Example:                    

    Available paths to GET:
        /api/v1/events/history
        /api/v1/events/recent



## Dev env

Compile with Go 1.9.x and https://github.com/gin-gonic/gin