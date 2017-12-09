# API


Service has single public endpoint accepting requests:

```json
POST /

Content-Type: application/json

{
  "type": <request-type>,
  "token": <request-producer-token>,
  "payload": { ... }
}
```

## Payload


Currently there is no payload validation implemented. Its format and content is totally
up to request producer and worker.

For example for SMS send request payload might be something like:

```json
    {
      "destination": "+155555555",
      "text": "ohai!"
    }
```

## Responses

### 200 OK

Your request is successfully added to the queue.


### 400 Bad Request

Probably malformed request body or invalid `content-type` header.


### 429 Too Many Requests

Your request didn't pass all configured limiters. Depending on request type trying again
might help.


### 500 Internal Server Error

Something bad happened. Fill bug report and try again.
 

# Failed request handling


Taking into account that service has unidirectional message flow and pipeline doesn't
really care about request specifics we can't handle errors reliably.

To make sure failing workers won't block the queue indefinitely each failed attempt will
lower request priority, so other requests can bubble up.

 

# Adding new request type


* edit `config.yaml` and add new request type
* if needed, add new worker to `server/workers`
* update `switch` statement at `server/worker.go`
* restart the app and you should be good to go

