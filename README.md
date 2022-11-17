# Study Websocket with Golang and Gorilla

* Go: 1.16

<br/>

## Run Server
- For MacOS

#### With Docker

```bash
# Use the run_server.sh file. It contains docker commands.
$ run_server.sh
```

#### With a local binary

```bash
# 1. run redis
$ docker run -p 6379:6379 redis 

# 2. compile
$ make all

# 3. run the server. Or You can use
# the debug mod of VSCode to click the button 'Launch'.
$ bin/study-websocket-go
```

<br/>

## Chat Test with Postman
- Before test, you need to install Postman.
### Create a WebSocket Request (Beta)
- URL
    - ws://{{localhost}}:{{port}}/room/{{room_id}}/broadcast?user_name="{{nick_name}}"
      - {{localhost}}
          - for Docker : 0.0.0.0
          - for a local binary : localhost or 127.0.0.1
      - {{port}}
          - for Docker : 10101 or 10102
          - for a local binary : 10812
      - {{room_id}}
          - a integer type value whatever you want
      - {{nick_name}}
          - a string type value whatever you want
- URL Examples
  - Docker
      - ws://0.0.0.0:10101/room/123/broadcast?user_name="revue"
      - ws://0.0.0.0:10102/room/123/broadcast?user_name="giraffe"
  - a local binary  
      - ws://localhost:10812/room/123/broadcast?user_name="iKnow"
### Start chat
- Click 'Connect' button on the requests, and Enjoy it!
