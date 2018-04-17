# lofai

_Hyper_ simple log http based log viewer / searcher

lofai creates a simple web listener with several HTTP endpoints for finding information about a relevant log



## Running:
`lofai --offset <int> --logfile <string> --port <string>`  

### Arguments:  
`--offset`  
For streamer operations which start  at the end of a file, how far inward the default seek will start.  
**`Default:`** 2048

`--logfie`  
Absolute path to the logfile to monitor.  
**`Default:`** Reads the OS Env Var `LOFAI_LOG`

`--port`  
Port to listen on.  
**`Default:`** 8000

---

## HTTP Endpoints:

`/streamer` -- `[GET]`  
This endpoint will stream lines as they are appended to the relevant log file; simliar to `tail -f`  
**Example:** `curl http://127.0.0.1:8000/streamer`

`/search/<searchTerm>` -- `[GET]`  
This endpoint will allow you to search for a term. Spaces must be encoded with `%20`  
**Example:** `curl http://127.0.0.1:8000/search/test%20search`


`/get/<bytes>` -- `[GET]`  
This endpoint allows you to retrieve the final chunk of a file. This uses an offset provided in the url and by default multiplies by a factor of 64.  
**Example:** `curl http://127.0.0.1:8000/get/1024`

_Note: Currently returns blank if the offset is greater than the size of the log file._

---

## Docker Image:

Pull Image:  
`docker pull walked/lofai`

Run Image:  
`docker run -d -v /host/logs:/container/log -p 8000:8000 walked/lofai --logfile /container/log/log.txt`

---

## Building docker:

Running `build.sh` should handle the build; tagging / pushing to a docker hub should be ok.

---

## External Resources / Requirements:
- [gorilla/mux](https://github.com/gorilla/mux) - golang url router  
- [hpcloud/tail](https://github.com/hpcloud/tail) - file tail library

## Design and Implementation
Conceptually prepared to be an incredibly simple / singular binary that can provide http access on a specified port to read a specific logfile.

With additional time we'd certainly implement:
- Directory watching (all logs in a path)
- Ability to Query by datetime
- Ability to provide regex for searches
- Token or similar authentication

Design is to provide http endpoints that can be accessed with all pertinent information encoded in the URL so as to minimize the 'spinup time' to be able to view the information needed.

## Alternatives:
The logical choice would be a decision between: 
- An http based log accessor that is pre-existing
- A log forwarding methodology (e.g. logstash)

The tradeoff being that http based access will not necessarily scale well (especially given ephemeral nodes / auto-scaling infrastructure), whereas log forwarding requires an existing endpoint to direct logs to for aggregation. The latter will certainly scale better, but is not as lightweight to implement for "quick" viewing.

## Deployment Plan:
This can be deployed several ways:  
- Config management  
The implementation of this naturally lends itself to a systemd unit to run as a service.
- Docker container (see above section re: Docker) 
- Manually run on any node necessary

## Monitoring / Reporting / Alerts:
The primary two deployment concepts allow for monitoring status readily.  

If via systemd, you can use (insert monitoring tool here) to verify the unit status as running.

If via docker; the container is naturally tied to the process at hand. You can use a the `--restart always` flag to ensure operational.

Finally; it would be simple to provide a service-health path that can be queried via polling mechanism; though this particular service may not lend itself to that.

Re: Metrics - you could conceivably provide a prometheus `/metrics` endpoint if it were deemed necessary to actively report metrics of such an application. 

## Configuration options:
As provided right now; the majority of configuration is optional, with the primary required input is the logfile path.

If a `--logfile` paramater is not passed, it will default to the value of the Environment Varialbe `LOFAI_LOG`

Should more in depth configuration be necessary later, I would likely implement a config file; possibly using `toml` as I did with [route53ddns](https://github.com/walked/route53ddns)

## Testing
TBD