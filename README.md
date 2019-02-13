# fake-take                                                  

![Alt text](assets/fake-take-avatar.png?raw=true "Fake Take Logo")
Simple HTTP API read only server to test Single Page Application

Immitate HTTP API for testing Single Page Application, e.g. VueJs app. 
Typycally a reply is a 'appliaction/json' format but you can use anything you want by passing the '-contentType' param.   

## Setup

API root path by default '/api/v1.0/'  to set another root path use '-apiroot' param and set any prefix you want.

### Adding paths and respose bodies

Create a file with *.resp extension. Name it folow the rule below.
Any part of the name separated by one underscore '_' treated as url part before the slash '/'.
Example:
events_recent.resp will be transalted to /api/v1.0/events/recent    

_No query params supported!_

You can crate as many files as you want.
New files can be added during server work, no need to restart.
You can even change the contents of '*.resp' file - server always read the file at request time.

Error 404 reported If '*.resp' file not exists.

## Dev env

Compile with Go 1.9.x and https://github.com/gin-gonic/gin