# bench: a simple http benchmark test

Bench is a simple, concurrent http benchmark tester. It can test multiple endpoints with varying amount of connections. So far the metrics returned
are average response size, average request duration, number of requests made, and number of errors. 

## Install 
```
go get github.com/andresoro/bench
```
## Usage

```
> bench config.json
Running test on /post/ with 3 connections making a request every 100ms
Running test on /get/ with 5 connections making a request every 100ms
Test completed for endpoint: /post/
	Total requests completed: 139
	Total errors: 1
	Average response size: 20.847826 bytes
	Average response time: 8.695211ms
Test completed for endpoint: /get/
	Total requests completed: 232
	Total errors: 5
	Average response size: 4.889868 bytes
	Average response time: 8.465712ms
```

Configs consist of request(s), the host, and the duration of the test in seconds. The *request rate* is rate per tester per request in milliseconds.


### Example config
```json
{
    "requests": [
        {
            "method":"GET",
            "endpoint": "/get/",
            "data": "",
            "connections": 5,
            "rate":100
        },
        {
            "method":"POST",
            "endpoint": "/post/",
            "data": "",
            "connections": 3,
            "rate":100
        }
    ],
    "host":"http://localhost:8080",
    "duration":5
}
```



