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
**Example:** `http://127.0.0.1:8000/streamer`

`/search/<searchTerm>` -- `[GET]`  
This endpoint will allow you to search for a term. Spaces must be encoded with `%20`  
**Example:** `http://127.0.0.1:8000/search/test%20search`


`/get/<bytes>` -- `[GET]`  
This endpoint allows you to retrieve the final chunk of a file. This uses an offset provided in the url and by default multiplies by a factor of 64.  
**Example:** `http://127.0.0.1:8000/get/1024`

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