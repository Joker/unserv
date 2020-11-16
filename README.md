# UnServ - stub-server with reverse proxy

Simple tool for developers to mock REST API without being dependent on any backend applications or servers. 


## Features

* simple endpoint creation
* reverse proxy (ability to redirect the endpoint request to a real working server)
* HTTP server (`./build/index.html` and `./build/static` folder by default)


## Installation

```
$ go get -u github.com/Joker/unserv
```


## Example

### Make API

If you want to make an endpoint on `http://localhost:8080/api/v1/question/`  
you should create a file `./stub/api/v1/question.json` in a file structure similar to this one:

`$ tree`
```
.
├── proxy
├── src
│   ├── App.js
│   └── App.scss
└── stub
    └── api
        ├── users.json
        └── v1
            └── question.json
```
`$ cat ./stub/api/v1/question.json`
```json
{
    "answer": 42
}
```
`$ unserv`
```
server start on:  http://localhost:8080

endpoints:
  http://localhost:8080/api/users
  http://localhost:8080/api/v1/question
```

#### Response

`$ curl http://localhost:8080/api/v1/question/`
```json
{
    "answer": 42
}
```


### Proxy example

If you want redirect the endpoint request to a real working server
move endpoint file to the `./proxy` folder:

`$ mkdir -p ./proxy/api/v1/`  
`$ mv ./stub/api/v1/question.json ./proxy/api/v1/`  
  
`$ tree`
```
.
├── build
│   ├── index.html
│   └── static
│       └── bundle.js
├── proxy
│   └── api
│       └── v1
│           └── question.json
├── src
│   ├── App.js
│   └── App.scss
└── stub
    └── api
        └── users.json
```
`$ unserv -proxy https://realserver.url`
```
server start on:  http://localhost:8080

endpoints:
  http://localhost:8080/api/users


reverse on:  https://realserver.url

proxy url:
  http://localhost:8080/api/v1/question =>
      https://realserver.url/api/v1/question
```

#### Response

`$ curl http://localhost:8080/api/v1/question/`
```json
{
    "answer": 2084
}
```


## Usage

```
Usage of ./unserv:
  -log
    	shows HTTP request log
  -p int
    	server port (default 8080)
  -proxy string
    	reverse proxy url (example:  http://localhost:9000 )
  -react
    	write setupProxy.js file (for React CRA with http-proxy-middleware)
  -reread
    	disable endpoint file cache (reread file on every request)
  -root string
    	root path for index.html and ./static foder (default "./build")
```