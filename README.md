Web Backend Technical Challenge
===============================
Please design and implement a web based API that steps through the Fibonacci sequence. 

The API must expose 3 endpoints that can be called via HTTP requests:
* current - returns the current number in the sequence
* next - returns the next number in the sequence
* previous - returns the previous number in the sequence

Example:
```
current -> 0
next -> 1
next -> 1
next -> 2
previous -> 1
```

Requirements:
* The API must be able to handle high throughput (~1k requests per second).
* The API should also be able to recover and restart if it unexpectedly crashes.
* Assume that the API will be running on a small machine with 1 CPU and 512MB of RAM.
* You may use any programming language/framework of your choice.


Using the Application
---------------------
### Running the app
This app is intended to be used with Docker for ease of environment setup and management of services.

To run the app using `docker-compose`, run this command in the root of the project
```bash
docker-compose up -d
```

To run without Docker (and forgoing redis integration), simply use the Makefile command
```bash
make run
```

Unit test execution is also available via
```bash
make test
```

### Endpoints
There are four endpoints served by the application, at the root address and port: `http://0.0.0.0:8080`  

  

#### `/health` - This endpoint simply returns the status of the server with code `200`  

To request it from the cli
```bash
curl -XGET http://0.0.0.0:8080/health
```
And receive
```bash
{"status": "healthy"}
```

#### `/current` - This endpoint retrieves the current number in the Fibonacci sequence the app is currently on - the assumption is that the app will start at `0`  

To request it from the cli
```bash
curl -XGET http://0.0.0.0:8080/current
```  
And receive:
```bash
{"current": 0}
```  

#### `/next` - This endpoint retrieves the next number in the Fibonacci sequence relative to the state of the app - this **will modify the state** of the application and advance `current` to `next`  
To request it from the cli
```bash
curl -XGET http://0.0.0.0:8080/next
```
And receive
```bash
{"next": 1}
```

### `/previous` - This endpoint retrieves the previous number in the Fibonacci sequence relative to the state of the app - an assumption was made that this WILL NOT modify the state of the app and at the starting state, `0` is `previous`  
To request it from the cli
```bash
curl -XGET http://0.0.0.0:8080/previous
```
And receive
```bash
{"previous": 0}
```  


Testing Load Handling / High Throughput (TPS)
---------------------------------------------

### Methodology
The benchmarking tool I elected to use was [wrk2](https://github.com/giltene/wrk2) due to its ability to:
* Maintain a constant throughput load
* Allow easy customizations of factors such as TPS, number of threads, connections open, and so forth
* Extremely accurate latency measurements

Using wrk2, I would attack all 3 endpoints specified in the requirements at once using 3 processes in parallel.
Each process will use 2 threads, 100 open connections, and a constant TPS of 1000 for a duration of 30s.

This this is the command I ended up running:

```bash
./wrk -t2 -c100 -d30s -R1000 http://0.0.0.0:8080/previous > ~/previous.log &
./wrk -t2 -c100 -d30s -R1000 http://0.0.0.0:8080/current > ~/current.log &
./wrk -t2 -c100 -d30s -R1000 http://0.0.0.0:8080/next > ~/next.log &
```

### Results
Recorded performance of `/current`
```
Running 30s test @ http://0.0.0.0:8080/current
  2 threads and 100 connections
  Thread calibration: mean lat.: 1.068ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 1.083ms, rate sampling interval: 10ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.09ms  735.03us  29.82ms   91.23%
    Req/Sec   525.09    125.85     1.78k    73.12%
  29922 requests in 30.00s, 4.01MB read
Requests/sec:    997.39
Transfer/sec:    136.75KB
```

Recorded performance of `/next`
```
Running 30s test @ http://0.0.0.0:8080/next
  2 threads and 100 connections
  Thread calibration: mean lat.: 1.078ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 1.082ms, rate sampling interval: 10ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.12ms  705.01us  27.70ms   89.99%
    Req/Sec   525.22    123.82     1.33k    49.30%
  29922 requests in 30.00s, 3.92MB read
Requests/sec:    997.37
Transfer/sec:    133.82KB
```


Recorded performance of `/previous`
```
Running 30s test @ http://0.0.0.0:8080/previous
  2 threads and 100 connections
  Thread calibration: mean lat.: 1.056ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 1.084ms, rate sampling interval: 10ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.09ms  724.56us  30.03ms   91.13%
    Req/Sec   525.08    121.23     1.60k    50.42%
  29924 requests in 30.00s, 4.03MB read
Requests/sec:    997.31
Transfer/sec:    137.71KB
```

These 3 endpoints were bombarded simulataneously each with a TPS load of 1000, coming together for a combined load test of 3000 TPS over the duration of 30s against the app itself.

Since all 3 endpoints have an average latency of about ~1.1ms, I estimate that performance is acceptable given the requirements asked only for 1/3 of the TPS load.

**NOTE**: My code uses uint64 to store the Fibonacci numbers, but I did attempt a solution that used math/big.Int to contain the values.
Peformance degraded by a magnitude of roughly 100x, which I did not consider a worthy tradeoff.
However that build is available in the `big` branch, and handling extremely large numbers was taken into consideration.