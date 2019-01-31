# bench: a simple http benchmark test

Bench is a simple, concurrent http benchmark tester. It can test multiple endpoints with varying amount of connections. So far the metrics returned
are average response size, average request duration, number of requests made, and number of errors. 

## Usage

Configs consist of request(s), the host, and the duration of the test in seconds. 

```json
{
    "requests": [
        {
            "method":"GET",
            "endpoint": "/get/",
            "data": "",
            "connections": 5
        },
        {
            "method":"POST",
            "endpoint": "/post/",
            "data": "",
            "connections": 3
            
        }
    ],
    "host":"http://localhost:8080",
    "duration":5
}
```

```
> bench config.json
Running test on /post/ with 3 connections for 5s
Running test on /get/ with 5 connections for 5s
Test completed for endpoint: /post/
	Total requests completed: 146
	Total errors: 1
	Average response size: 20.855172 bytes
	Average response time: 0.103662s
Test completed for endpoint: /get/
	Total requests completed: 492
	Total errors: 13
	Average response size: 4.906054 bytes
	Average response time: 0.052348s
```

