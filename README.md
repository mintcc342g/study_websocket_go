# study_websocket_go

* Go: 1.16

## run server
- for MacOS
```
## case1: using the .sh file which contains docker commands.
$ run_server.sh


## case2: run each commands, not use docker for the app server
# case2-1) run redis
$ docker run -p 6379:6379 redis 

# case2-2) compile the ws server
$ make all

# case2-3) run the ws server
$ bin/study_websocket_go
```

---
## Chat Test with Postman

* install Postman
* create New WebSoket Requests
* set url like ... ws://{{localhost}}:{{port}}/room/{{room_id}}/broadcast?user_name="{{nick_name}}"
  - local docker server
    - ws://0.0.0.0:10101/room/123/broadcast?user_name="revue"
    - ws://0.0.0.0:10102/room/123/broadcast?user_name="giraffe"
  - local server
    - ws://localhost:10812/room/123/broadcast?user_name="iKnow"
* click 'Connect' button
* enjoy it!
