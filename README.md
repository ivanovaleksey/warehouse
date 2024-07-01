# Warehouse

### Build
The application can be built from source or in a Docker container.
```shell
docker build -t warehouse .
```

### Run
You can run compiled binary or use docker-compose to run both the app and DB server.

1. Run DB first
```shell
docker-compose up -d db
```

2. Then start the app. It will be listening to port 8000
```shell
docker-compose up -d server
```

3. Now you can use `grpc_cli` to make requests
```shell
grpc_cli call 127.0.0.1:8000 warehouse.WarehouseService/GetProducts ''
grpc_cli call 127.0.0.1:8000 warehouse.WarehouseService/RemoveProduct 'id: 1'
```
You will see empty response because there is no data in the database.

4. Run seeds to fill the database and send one request
```shell
docker-compose up -d db-seed
```

### Test
The test suite can be run locally or using docker-compose.

1. Build test image
```shell
docker build -f tests/Dockerfile -t warehouse-tests .
```

2. Run test DB
```shell
docker-compose -f tests/docker-compose.yaml up -d db
```

3. Run the tests
```shell
docker-compose -f tests/docker-compose.yaml up tests
```
