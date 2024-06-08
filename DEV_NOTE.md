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
```sh
# chmod +x ./deploy.sh
# ./deploy.sh

# docker-compose
docker-compose up -d
# Stop services only
docker-compose stop
# Stop and remove containers, networks..
docker-compose down
# Down and remove volumes(for postgres to initialize docker-entrypoint-initdb.d)
docker-compose down --volumes
# Down and remove images
docker-compose down --rmi <all|local>
```

# mino
webUI: http://127.0.0.1:9001

# postgresql
```
docker exec -it <CONTAINER ID>
psql -U <username> -d <database>
```


## References
### Problem fix
- [cmd/go, x/mod: "invalid go version" message mentions old version format](https://github.com/golang/go/issues/61888)

### notes
- [General best practices for writing Dockerfiles]https://docs.docker.com/develop/develop-images/guidelines/
- [Bash脚本中的 set -euxo pipefail](https://www.cnblogs.com/wjoyxt/p/14734502.html)
- [PostgreSQL error: Fatal: role “username” does not exist](https://openbasesystems.com/2023/06/20/postgresql-error-fatal-role-username-does-not-exist/)
docker cmd
```sh
docker build -t foo . && docker run -it foo
docker build -f /path/to/folder/account.Dockerfile /path/to/folder -t your_image_name:tag
docker build -f ./tools/test_be.Dockerfile ./tools -t be_test
```
