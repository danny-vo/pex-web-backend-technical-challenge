Table of Contents
=================
* [Challenge Description](#web-backend-technical-challenge)
* [Using the Application](#using-the-application)
    + [Running the app](#running-the-app)
    + [Endpoints](#endpoints)
        - [`/health`](#health---this-endpoint-simply-returns-the-status-of-the-server-with-code-200)
        - [`/current`](##current---this-endpoint-retrieves-the-current-number-in-the-fibonacci-sequence-the-app-is-currently-on---the-assumption-is-that-the-app-will-start-at-0)
        - [`/next`](#next---this-endpoint-retrieves-the-next-number-in-the-fibonacci-sequence-relative-to-the-state-of-the-app---this-will-modify-the-state-of-the-application-and-advance-current-to-next)
        - [`/previous`](#previous---this-endpoint-retrieves-the-previous-number-in-the-fibonacci-sequence-relative-to-the-state-of-the-app---an-assumption-was-made-that-this-will-not-modify-the-state-of-the-app-and-at-the-starting-state-0-is-previous)
* [Testing Load Handling / High Throughput (TPS)](#testing-load-handling--high-throughput-tps)
    + [Methodology](#methodology)
    + [Results](#results)
* [Application Design](#application-design)
    + [Language / Framework](#language--framework)
    + [Technologies](#technologies)
    + [Fault Tolerance and Recovery](#fault-tolerance-and-recovery)
        - [Software Solution](#software-solution)
        - ["Infrastructure" Solution](#infrastructure-solution)

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

#### `/previous` - This endpoint retrieves the previous number in the Fibonacci sequence relative to the state of the app - an assumption was made that this **WILL NOT modify the state** of the app and **at the starting state, `0` is `previous`**  
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
  Thread calibration: mean lat.: 1.323ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 1.317ms, rate sampling interval: 10ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.10ms  497.46us   7.84ms   70.04%
    Req/Sec   525.32    123.56     1.00k    49.36%
  29925 requests in 30.00s, 4.00MB read
Requests/sec:    997.45
Transfer/sec:    136.69KB
```

Recorded performance of `/next`
```
Running 30s test @ http://0.0.0.0:8080/next
  2 threads and 100 connections
  Thread calibration: mean lat.: 1.337ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 2.075ms, rate sampling interval: 10ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.58ms    0.89ms   7.41ms   72.91%
    Req/Sec   524.25      1.04k    5.55k    94.31%
  29462 requests in 30.00s, 3.86MB read
Requests/sec:    982.03
Transfer/sec:    131.74KB
```


Recorded performance of `/previous`
```
Running 30s test @ http://0.0.0.0:8080/previous
  2 threads and 100 connections
  Thread calibration: mean lat.: 1.301ms, rate sampling interval: 10ms
  Thread calibration: mean lat.: 1.294ms, rate sampling interval: 10ms
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.09ms  500.79us   7.19ms   69.92%
    Req/Sec   524.86    125.49     1.00k    72.16%
  29924 requests in 30.00s, 4.03MB read
Requests/sec:    997.42
Transfer/sec:    137.67KB
```

These 3 endpoints were bombarded simultaneously each with a TPS load of 1000, coming together for a combined load test of 3000 TPS over the duration of 30s against the app itself.

Since the read only endpoints have an average latency of about ~1.1ms, and the read/write endpoint at ~1.58ms latency - we can estimate that performance is acceptable given the requirements asked only for 1/3 of the TPS load actually tested.

**NOTE**: My code uses uint64 to store the Fibonacci numbers, but I did attempt a solution that used math/big.Int to contain the values.
Performance degraded by a magnitude of roughly 100x, which I did not consider a worthy tradeoff.
However that build is available in the `big` branch, and handling extremely large numbers was taken into consideration.


Application Design
------------------
### Language / Framework
Due to the performance requirements stated out of the application, I elected to use `Golang` as the language of choice due to native asynchronous ability and general speed/efficiency. This was important as there was also a restriction on the amount of resources our application had access to.  

The nature of this application's specifications also requires us to be mindful of read/write concurrency issues with respect to the state of the Fibonacci values. `Go` helps by natively providing an extremely easy to use native `sync` library, more specifically `sync.RWMutex`. Having the `-race` argument available to both test and compile binaries that reported race conditions is also crucial in avoiding any programming concurrency mishaps.

I utilized [httprouter](https://github.com/julienschmidt/httprouter) as the routing framework of choice, as it outperforms native `net/http` and other 3rd party libraries with extremely fast times. Additionally it is fairly lightweight and since the application itself is not large scale, there is no need to use something more cumbersome.  


### Technologies
`Docker` immediately came to mind as a highly appropriate tool to use for containerizing the application. This ensures that regardless of where the execution takes place, the actual environment the application runs in is the same, mitigating any potential compatibility failures due to OS, Go version, and so forth.

I also specify a multistage build in the `Dockerfile`, to keep the image itself small. Using `Alpine` as the OS serves to further compact memory complexity.  

I decided upon using `docker-compose` as the orchestration tool to prop our services up and specify various constraints and policies. Since I did not have access to infrastructure (AWS, Azure, etc.) this would be a suitable way for me to develop, test, and deploy locally.

`docker-compose` also makes it easy to instantly start and stop all related services without having to type out multiple convoluted `docker run` arguments each time. Other noteworthy things to me using `compose` was the ability to easily limit the containers' resources to match the requirements of CPU and memory units. Additionally the restart policies would come in handy later.  

As a state caching solution I used `redis` alongside [go-redis](https://github.com/go-redis/redis) to interface with the database. We only needed to store a single value from our application to restore state so there was no need for something more complicated such as PostgreSQL or Cassandra.  

Additionally since `redis` stores using in memory, it is not limited by disk I/O which would assist with the theme of being as performant as possible.


### Fault Tolerance and Recovery
I attempted to provide maximum resiliency and recoverability from both software solutions (inside the main code) as well as from an "infrastructure" standpoint using the orchestration layer and `redis`.  

#### Software Solution
At the very top level of the application, it is driven by a simple `http.ListenAndServe` which props up the routes and handlers to the outside world. This `ListenAndServe` function will run and "block" until an error occurs in which case it exits upon returning an error value.  

To mitigate the permanent disruption and shutdown of the application within the code base itself, I thus wrapped the aforementioned function and any server initializations within a `run` function.
```go
func run() error {
	var err error = nil
	s := server.InitializeServer()

	log.Println("Server has been initialized, now serving...")
	err = http.ListenAndServe("0.0.0.0:8080", s.GetRouter())

	return err
}
```
This `run` function was then evoked with `main`, the starting point of the app
```go
func main() {
	log.Println("Starting server...")
	for {
		if err := run(); nil != err {
			log.Printf("Error occurred while serving: %v\n", err)
		}
	}
}
```
By using this loop construct, we can get around the limitations of `main` only executing once, and ensures the application will attempt to restart again, with proper initialization and serving, ignoring fringe errors that may occur during the software solution.

The next point of potential failures are the handler functions, which are mapped to the endpoint routes. To prevent a `panic` from killing the app off due to some uncaught fringe errors, I wrap the handler functions with the exclusion of `/health` with a `recoveryWrapper`
```go
func recoveryWrapper(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		defer func() {
			if r := recover(); nil != r {
				log.Printf("Error occurred: %v\n, recovered", r)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()

		h.ServeHTTP(w, r)
	})
}
```
This wrapper function works around the limitations of `Go`'s `recover` function, which can only be executed with any real meaningful impact within defered functions. The `recover` function allows a panicked `Goroutine` to essentially regain control, and resume intended execution. The panic cause is propagated upwards however for logging. This middleware essentially acts an effective "catchall"  

I decided not to capture `/health` with the middleware since if some issue happened with such basic functionality something is *probably* very wrong and a panic is more than acceptable. Additionally the `/health` endpoint is used as a litmus test for app healthiness with our `Docker` configurations.

#### "Infrastructure" Solution
Admittedly, it is difficult to ensure high tolerance and reliability with only containers and not a fully blown infrastructure / cloud service but there are still some tools and methodology that I found to be useful.  

A bit of a more obvious solution, but using `redis` to store the current value in the Fibonacci sequence the application is currently on helps to manage and recover state. Whenever state is modified, namely through the execution of the `/next` endpoint, the app fires off a `Goroutine` to set the value in `redis`.  

During application startup, the app attempts to retrieve the "current" value from `redis` while initializing its Fibonacci state. If an error occurred, `redis` is not up, or the value is not yet set, the app starts from a fresh state. Elsewise the app retrieves the value stored in `redis` and constructs the `previous` and `next` values and uses that as its starting state.

Note that a limitation exists in the case **BOTH** `app` and `redis` goes boom, there is no other option but to start from a fresh state. There is also the rare case where `redis` restarts and the app fails to write a value to the database before crashing, will result in restarting in a fresh state. A potential solution to this would to have the `redis` service `curl` the `/current` endpoint on start and attempt to set the value itself.

Another challenge to handle was what if the "machine" (i.e. container the app is on) goes unhealthy, and not in the sense that the container itself is unhealthy. Some such scenarios could be the app gets stuck in a 3rd party library in an internal loop, or it simply got overloaded with requests to the point of non-responsiveness and failure.  

Normally if we had this service deployed to AWS on ECS, we could just note a failing health check and have ECS restart the service, feeding the new host/port to the ALB. However... we are just dealing with local containers. Fortunately, there is a [restart](https://docs.docker.com/compose/compose-file/compose-file-v2/#restart) parameter that `docker-compose` provides. Alongside that there is also a very neat customizable [healthcheck](https://docs.docker.com/compose/compose-file/compose-file-v2/#healthcheck) parameter.

Rather unfortunately however, `restart` will only be triggered when the container itself is unhealthy (e.g. PID 1 is terminated) and there is no way of triggering it based off a simply failing healthcheck without using overly complicated scriptings and/or events.

Luckily for us, the `healthcheck` parameter allows for custom command execution for healthchecking. To take advantage of this we can assign the `test` subparameter with
```bash
curl -fail --retry 3 --max-time 5 --retry-delay 5 --retry-max-time 30 "http://0.0.0.0:8080/health" || bash -c 'kill -s 15 -1 && (sleep 10; kill -s 9 -1)'
```
To explain what this essentially does, it pings our `/health` endpoint with certain timeout and retry restrictions. If it ends up receiving a failing status code or times out without any response, it then pipes over to a `kill` command. The `kill` command attempts to first tell all processes to gracefully die then after a while terminates them forcibly. Since it is prefixed with `bash` rather than just `sh`, the `-1` PID means to target **ALL** running processes.  

As mentioned before, containers will then be viewed as *unhealthy* by `docker-compose`, triggering our `restart: always` policy. There is a small caveat to this though - PID 1 is protected by default in `Docker`. To circumvent this safeguard from terminating our main process, we can use the [init](https://docs.docker.com/compose/compose-file/compose-file-v2/#init): `true` parameter. This ends up wrapping all of our processes using [tini](https://github.com/krallin/tini), circumventing the need to manually define a kill signal handler for our processes.
