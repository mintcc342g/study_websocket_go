# Study Websocket with Golang and Gorilla

* Go: 1.16

<br/>

## Run Server
- For MacOS

#### First Way

```bash
# Use the run_server.sh file. It contains docker commands.
$ run_server.sh
```

#### Second Way

```bash
# 1. run redis
$ docker run -p 6379:6379 redis 

# 2. compile the project
$ make all

# 3. run the server (You can use the debug mod of VSCode to click the button 'Launch')
$ bin/study-websocket-go
```

<br/>

## Chat Test with Postman

- Before test, you need to install Postman.
- Create new webSoket requests on Postman.
- Set request URL and query strings on the requests you created.
  ```Makefile
  # Format
  ws://{{localhost}}:{{port}}/room/{{room_id}}/broadcast?user_name="{{nick_name}}"

  # an Example for Docker Instances
  ws://0.0.0.0:10101/room/123/broadcast?user_name="revue"
  ws://0.0.0.0:10102/room/123/broadcast?user_name="giraffe"
  
  # an Example for a Local Server
  ws://localhost:10812/room/123/broadcast?user_name="iKnow"

  ```
- Click 'Connect' button on the requests, and Enjoy it!
