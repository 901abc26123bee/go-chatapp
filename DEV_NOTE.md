## Development
### direnv
[direnv](https://direnv.net/docs/hook.html)
1.  install direnv
2.  Setup direnv
    - (zsh) Add the following line at the end of the ~/.zshrc file:
        ```sh
        eval "$(direnv hook zsh)"
        ```
    - Run command to execute your zsh setting
        ```sh
        source ~/.zshrc
        ```
    - Run command to allow direnv on your machine
        ```sh
        direnv allow .
        ```
    - If you want to disable
        ```sh
        direnv deny .
        ```

### pre-commit
- If want to ignore pre-commit test
  ```sh
  git commit --no-verify
  ```

### deploy
[docker-compose up doesn't rebuild image although Dockerfile has changed](https://github.com/docker/compose/issues/1487)
```sh

## file not executable
```sh
chmod +x ./deploy.sh
./deploy.sh
```


# docker-compose up in background
docker-compose up -d
# up with rebuild image
docker compose up --build
# Stop services only
docker-compose stop
# Stop and remove containers, networks..
docker-compose down
# Down and remove volumes(for postgres to initialize docker-entrypoint-initdb.d)
docker-compose down --volumes
# Down and remove images
docker-compose down --rmi <all|local>
# restart a specific container in a Docker Compose
docker-compose restart <service_name>



[How to Find Default Nginx Document Root Directory Location?](https://www.uptimia.com/questions/find-default-nginx-document-root-directory-location)

```sh

# reload nginx
docker exec -it gsm_nginx nginx -s reload

# see nginx applied setting
docker exec -it gsm_nginx cat /etc/nginx/nginx.config
docker exec -it gsm_nginx ls -l /var/www
docker exec -it gsm_nginx nginx -t

# see the nginx indeed applied setting
docker exec -it gsm_nginx nginx -T 
# sees if root setting is updated
docker exec -it gsm_nginx nginx -T | grep root

# inspect the active NGINX configuration 
docker exec -it gsm_nginx nginx -T | grep "nginx.conf"


# 1. Verify Configuration File Location
# Double-check the exact path of the configuration file being used by NGINX inside the container. The configuration might be split across multiple files or directories.
# Look for include directives in the output to find any additional configuration files that might be overriding your changes.

docker exec -it gsm_nginx nginx -T | grep include

include  /etc/nginx/mime.types;
include /etc/nginx/conf.d/*.conf;
#    include        fastcgi_params;


# 2. Check for Multiple Configuration Files
# NGINX can include other configuration files, which might be overriding your settings. Look in common locations like /etc/nginx/conf.d/ or /etc/nginx/sites-enabled/.
docker exec -it gsm_nginx ls /etc/nginx/conf.d/
docker exec -it gsm_nginx ls /etc/nginx/sites-enabled/
```


The issue might be due to a small typo in the path where you’re mounting your NGINX configuration file. The correct path should be /etc/nginx/nginx.conf, not /etc/nginx/nginx.config. NGINX expects its main configuration file to be named nginx.conf.




### mino
webUI: http://127.0.0.1:9001

### postgresql
```
docker exec -it <CONTAINER ID>
psql -U <username> -d <database>
```

### nginx
Templates for nginx.config : [Example nginx configuration](http://nginx.org/en/docs/example.html)




## References
### Problem fix
- [cmd/go, x/mod: "invalid go version" message mentions old version format](https://github.com/golang/go/issues/61888)

### notes
- [General best practices for writing Dockerfiles]https://docs.docker.com/develop/develop-images/guidelines/
- [Bash脚本中的 set -euxo pipefail](https://www.cnblogs.com/wjoyxt/p/14734502.html)
- [PostgreSQL error: Fatal: role “username” does not exist](https://openbasesystems.com/2023/06/20/postgresql-error-fatal-role-username-does-not-exist/)
- [Socket. IO vs. WebSocket: Keys Differences](https://apidog.com/articles/socket-io-vs-websocket/)


### Vscode
[How To Re-add An Extension To The VSCode Sidebar/Activity Bar](https://stackoverflow.com/questions/71567229/how-to-re-add-an-extension-to-the-vscode-sidebar-activity-bar)
[Customizing the Outline view](https://www.ibm.com/docs/en/wdfrhcw/1.4.0?topic=editing-customizing-outline-view)
[Activity Bar](https://code.visualstudio.com/api/ux-guidelines/activity-bar)


### docker cmd
```sh
docker build -t foo . && docker run -it foo
docker build -f /path/to/folder/account.Dockerfile /path/to/folder -t your_image_name:tag
docker build -f ./tools/test_be.Dockerfile ./tools -t be_test
docker build -f ./deployment/dockerfile/account.Dockerfile . -t test_build_account
```

### docker-compose cmd
[How to rebuild docker container in docker-compose.yml?](https://stackoverflow.com/questions/36884991/how-to-rebuild-docker-container-in-docker-compose-yml)


### address already in use
```
lsof -i tcp:3000
netstat -vanp tcp | grep 3000
```

[Is it possible to install Node.js on macOS High Sierra Version 10.13.6? This is the latest version my mid 2014 Macbook Pro can support](https://stackoverflow.com/questions/74709494/is-it-possible-to-install-node-js-on-macos-high-sierra-version-10-13-6-this-is)
```sh
# list node versions
nvm ls-remote
# install node
nvm install v20.14.0
nvm install v17.9.1
```

### Golang
[Cannot do any go command anymore](https://stackoverflow.com/questions/60406755/cannot-do-any-go-command-anymore)
[compile: version "go1.16.2" does not match go tool version "go1.15.8" #107](https://github.com/actions/setup-go/issues/107)
[Golang VSCode 開發環境建置 - 手把手教學](https://myapollo.com.tw/blog/golang-vscode/)


[cmd/go: go mod download breaks on 1.21.0 due to empty GOPROXY](https://github.com/golang/go/issues/61928)
```sh
go env -w GOPROXY=https://proxy.golang.org,direct
go env | grep GOPROXY
go install github.com/swaggo/swag/cmd/swag@latest


go get github.com/swaggo/swag  
```
### swagger
cd  doc/    
swag init  -g ./api/account/controller.go  


### redis
https://stackoverflow.com/questions/46569432/does-redis-use-a-username-for-authentication
redis-cli -a redispw ping
redis-cli -h 127.0.0.1 -p 6379 -a redispw

redis-cli -h 127.0.0.1 -p 6379
127.0.0.1:6379> AUTH redispw
OK
127.0.0.1:6379>

redis-cli -u redis://redisuser:redispw@localhost:6379

###  go-redis version
https://github.com/redis/go-redis/discussions/2241

### redis stream

In the XGroupCreateMkStream command, the $ symbol is used to specify the starting point for the consumer group. 
$ Symbol in XGroupCreateMkStream
```sh
$: This indicates that the consumer group should start reading new messages added to the stream after the group is created. In other words, it will not read any historical messages that were already in the stream before the group was created. It only processes messages that are appended to the stream after the group has been set up.
```
Detailed Explanation
```sh
When you create a consumer group with XGroupCreateMkStream, you provide the stream name, group name, and a starting point (usually $ or 0). The starting point determines from where the consumer group will start processing messages:

  $ (Dollar Sign): Start from new messages added to the stream after the group is created. This is commonly used when you want to process only new messages from that point onward, ignoring any messages that were already present in the stream.

  0 (Zero): Start from the earliest message available in the stream. This is used if you want the consumer group to process all messages, including those that were already in the stream before the group was created.
```




ws://localhost:8081/api/realtime/v1/chatroom/stream?&access_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJnc20tZGV2IiwiZXhwIjoxNzIzMDI2NzM0LCJpYXQiOjE3MjMwMjMxMzQsImlzcyI6ImdzbS1kZXYiLCJqdGkiOiIwMUo0OFk2QVRWVk1TUEtKVFJGOEsxMVRUNSIsIm5iZiI6MTcyMzAyMzEzNCwic3ViIjoiYWNjZXNzX3Rva2VuIn0.xfb8QhJQyeBUcxRds2-9K7Q-q-EVmT_eT8WL6QSAIsk

{"room_id":"test1","action":"JOIN_CHAT_ROOM"}
{"room_id":"test1","action":"LEAVE_CHAT_ROOM"}
{"room_id":"test1","action":"CHAT_ROOM_MESSAGE","chat":"hello"}

### mongodb
https://stackoverflow.com/questions/42912755/how-to-create-a-db-for-mongodb-container-on-start-up


### web
https://stackoverflow.com/questions/48712923/where-to-store-a-jwt-token-properly-and-safely-in-a-web-based-application

- localStorage - data persists until explicitly deleted. Changes made are saved and available for all current and future visits to the site.

- sessionStorage - Changes made are saved and available for the current page, as well as future visits to the site on the same window. Once the window is closed, the storage is deleted.

### websocket
When upgrading an HTTP connection to a WebSocket, the initial HTTP request can include query parameters, which are typically used for authentication or passing configuration details. These query parameters are included in the HTTP GET request used to initiate the WebSocket handshake.

However, these query parameters are only available during the initial upgrade request. Once the connection is upgraded to a WebSocket, it becomes a stateful, full-duplex connection where messages are sent and received as binary or text frames. The concept of query parameters doesn't directly apply to the messages sent over this connectio


## Nginx
```sh
root@5a10dafe097b:/# curl -v http://account:8080/api/account/v1/healthz
alive!

curl -v http://gsm-dev/api/account/v1/healthz

root@5a10dafe097b:/# curl -v http://realtime:8081/api/realtime/v1/healthz
root@5a10dafe097b:/# nginx -t
nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
nginx: configuration file /etc/nginx/nginx.conf test is successful
```

https://github.com/docker/compose/issues/3412
https://stackoverflow.com/questions/35744650/docker-network-nginx-resolver
https://forums.docker.com/t/nginx-swarm-redeploy-timeouts/68904/4

https://docs.docker.com/engine/network/
