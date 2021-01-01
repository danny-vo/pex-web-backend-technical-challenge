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



Testing Load Handling / High Throughput (TPS)
=============================================

### Methodology
The benchmarking tool I elected to use was [wrk2](https://github.com/giltene/wrk2) due to its ability to:
* Maintain a constant throughput load
* Allow easy customizations of factors such as TPS, number of threads, connections open, and so forth
* Extremely accurate latency measurements

Using wrk2, I would attack the 3 endpoints specified in the requirements at once using 2 threads, 100 open connections, constant TPS of 1000 *EACH* for a duration of 30s.

This this is the command I ended up running:

```bash
./wrk -t2 -c100 -d30s -R1000 http://0.0.0.0:8080/previous > ~/previous.log &
./wrk -t2 -c100 -d30s -R1000 http://0.0.0.0:8080/current > ~/current.log &
./wrk -t2 -c100 -d30s -R1000 http://0.0.0.0:8080/next > ~/next.log &
```

Thus, this ended up being a test of a 3000 constant TPS load!

### Results
Recorded performance of `/current`
```
Running 30s test @ http://0.0.0.0:8080/current
  2 threads and 100 connections
  Thread calibration: mean lat.: 1.079ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 1.066ms, rate sampling interval: 10ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.08ms  518.19us  13.21ms   72.45%
    Req/Sec   525.13    124.03     1.33k    48.96%
  29922 requests in 30.00s, 3.73MB read
Requests/sec:    997.36
Transfer/sec:    127.34KB
```

Recorded performance of `/next`
```
Running 30s test @ http://0.0.0.0:8080/next
  2 threads and 100 connections
  Thread calibration: mean lat.: 1.087ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 1.089ms, rate sampling interval: 10ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.10ms  515.76us  13.61ms   72.93%
    Req/Sec   525.52    119.71     1.20k    50.34%
  29925 requests in 30.00s, 3.65MB read
Requests/sec:    997.44
Transfer/sec:    124.42KB
```


Recorded performance of `/previous`
```
Running 30s test @ http://0.0.0.0:8080/next
  2 threads and 100 connections
  Thread calibration: mean lat.: 1.087ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 1.089ms, rate sampling interval: 10ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.10ms  515.76us  13.61ms   72.93%
    Req/Sec   525.52    119.71     1.20k    50.34%
  29925 requests in 30.00s, 3.65MB read
Requests/sec:    997.44
Transfer/sec:    124.42KB
```

These 3 endpoints were bombarded simulataneously each with a TPS load of 1000, coming together for a combined load of 3000 TPS over the duration of 30s.

Since all 3 endpoints have an average latency of about ~1.1ms I estimate that performance is acceptable given the requirements asked only for 1/3 of the TPS load.
