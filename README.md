# grpc-server

Example gRPC golang application. Mongo DB used as persisted storage.
The app works from Docker container, and uses NGINX as load balancer.

The app has two functions:
- FETCH: fetches csv file from url and saves it to db. 
(File format: *name*;*price*)
- LIST: list and filter db table.




Clone the repository and type the following command to build app. Docker and docker-compose has to be installed, otherwise command will fail:
```bash
make build
```

This command will fire up app on 3 instances:
```bash
make up
```

The app works on port: **8003**