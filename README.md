# Test server
A simple web server that includes some endpoints useful for testing proxy server behaviour.  
Use it as a docker image: [chmodx/test-server](https://hub.docker.com/r/chmodx/test-server)

## Endpoints
### `/download`
Downloads a file with a size of the specified number of random bytes.  
It writes the bytes in chunks of 1MB and prints to stdout when finished.  
Example:  
`GET /download?size={fileSize}`  

### `/upload`
Uploads a file and prints out the filename and size to stdout.  
It returns a summary of the  uploaded data in the response.  
Example:  
```
POST /test-server/upload HTTP/1.1
Host: localhost:14140
Content-Type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW

----WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="file"; filename="somefile"
Content-Type: <Content-Type header here>

(data)
----WebKitFormBoundary7MA4YWxkTrZu0gW
```

### `/headers`
Returns the http method, path, protocol and headers in the response.  
Accepts a query parameter `?print=true` which means it additionally prints the response to stdout.  
Example:  
`GET /headers`  
`GET /headers?print=true`  

### `/service`
Calls another endpoint and returns its response. This is useful for testing distributed tracing scenarios.   
When starting the server supply the `serviceBaseUrl` and `serviceCallPath` to control what gets called.  
Example:  
Start `./test-server -port 8888 -serviceBaseUrl=http://localhost:9999 -serviceCallPath=/service`(S1) and `./test-server -port 9999 -serviceBaseUrl=http://localhost:8888 -serviceCallPath=/headers?print=true`(S2).  
Then calling `GET /service` on S1 will call `GET /service` on S2, which will call `GET /headers` on S1 again.
