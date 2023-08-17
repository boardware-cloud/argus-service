# BoardWare cloud

## Configuration

Example

```yaml
database:
  host: 127.0.0.1
  user: boardware
  port: 3306
  password: boardware
  database: boardware_cloud_dev
server:
  port: 8081
boardware:
  core:
    url: http://127.0.0.1:8080
jwt:
  secret: boardwaresecret
```

## Generate model from openapi

```bash
openapi-generator generate -i openapi.yaml -g go-gin-server \
  --additional-properties=packageName=model \
  --additional-properties=apiPath=model \
  -o ./controllers
```
