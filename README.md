# Price monitor


1. Clone repo.
2. Make docker image

3. Fill in config.env
```
APP_HOST=0.0.0.0
APP_PORT=9050
APP_LOG_LEVEL=error
APP_NUMCPUS=-1      //-1 for max available
APP_CHECK_PERIOD=6  //hours
```

4. Run the container
```sh
docker run -d --name <container-name> --env-file config.env --network <your-docker-network> <image-name>
```
