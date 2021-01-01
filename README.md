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
  Thread calibration: mean lat.: 2.090ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 2.142ms, rate sampling interval: 10ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.12ms    1.00ms   9.20ms   71.72%
    Req/Sec   523.88      1.44k    5.55k    88.75%
  29002 requests in 30.00s, 3.62MB read
Requests/sec:    966.70
Transfer/sec:    123.40KB
```

Recorded performance of `/next`
```
Running 30s test @ http://0.0.0.0:8080/next
  2 threads and 100 connections
  Thread calibration: mean lat.: 1.184ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 1.190ms, rate sampling interval: 10ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.14ms  571.78us   7.48ms   74.49%
    Req/Sec   525.25    125.93     1.00k    72.73%
  29922 requests in 30.01s, 3.65MB read
Requests/sec:    997.12
Transfer/sec:    124.39KB
```


Recorded performance of `/previous`
```
Running 30s test @ http://0.0.0.0:8080/previous
  2 threads and 100 connections
  Thread calibration: mean lat.: 1.117ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 1.167ms, rate sampling interval: 10ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.11ms  558.46us   8.84ms   72.48%
    Req/Sec   524.18    123.24     1.00k    75.75%
  29922 requests in 30.00s, 3.76MB read
Requests/sec:    997.36
Transfer/sec:    128.31KB
```

As a reminder, these 3 endpoints were bombarded simulataneously each with a TPS load of 1000, coming together for a combined load of 3000 TPS over the duration of 30s.

Both `/previous` and `/next` had a reported average of ~1.1ms of latency with `/current` averaging ~2.1ms of latency - I humbly consider this within acceptable parameters given the requirements asked only for 1/3 of the TPS load.
