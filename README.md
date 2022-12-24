# extract-weather
Tool to export daily Weather data into Snowflake at 6 AM NZST.

## Build binary and Docker image
```
cd build
./create-image.sh
```

## Run in Docker
```
docker run --name extract-weather -d --restart=always \
  -e SNOWFLAKE_ACCOUNT=xyz.us-central1.gcp \
  -e SNOWFLAKE_USER=username \
  -e SNOWFLAKE_PWD=password \
  -e SNOWFLAKE_DB=my-db \
  -e SNOWFLAKE_SCHEMA=PUBLIC \
  datasail/extract-weather:1.0
```

## Run as Compose
```
version: '3.8'
services:
  extract-weather:
    image: datasail/extract-weather:1.0
    environment:      
      SNOWFLAKE_ACCOUNT: xyz.us-central1.gcp
      SNOWFLAKE_USER: username
      SNOWFLAKE_PWD: password
      SNOWFLAKE_DB: my-db
      SNOWFLAKE_SCHEMA: PUBLIC
```

## Run on Kubernetes
```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: extract-weather-deployment
  labels:
    app: extract-weather-label
spec:
  replicas: 1
  selector:
    matchLabels:
      app: extract-weather-label
  template:
    metadata:
      labels:
        app: extract-weather-label
    spec:
      containers:
        - name: extract-weather
          image: datasail/extract-weather:1.0
          imagePullPolicy: Always
          env:
            - name: SNOWFLAKE_ACCOUNT
              value: xyz.us-central1.gcp
            - name: SNOWFLAKE_USER
              value: username
            - name: SNOWFLAKE_PWD
              value: password
            - name: SNOWFLAKE_DB
              value: my-db
            - name: SNOWFLAKE_SCHEMA
              value: PUBLIC
```

## Build only binary
```
cd build
./build-binary.sh
```
