# Testapp 20231020

### Prerequisites
* Go 1.21
* Port, defined in [config.yaml](config.yaml) at `http_server.port`, must be free at application start
* [optional] Make

### How to configure
Application is configured via yaml config file. Consult with [config.yaml](config.yaml) to see the structure.
To use a configuration file, pass a path to it via `--config-path` flag.

### How to run
You could use `make run-local` in the project directory to build and run the application with default config path.

Alternatively, you could run the application with `go run ./cmd/app --config-path {path_to_config.yaml}` to use different config file.

### How to access
If application started correctly, you'll be able to access OpenAPI specification UI for it via browser. With default config, you could go to http://localhost:59999/_/docs/ to access it.

Singular API endpoint could be accessed via curl:
```shell
curl -X 'POST' \
  'http://localhost:59999/v1/fibonacci' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "position": 1000
}'
```
Response:
```json
{
  "data": {
    "result": "43466557686937456435688527675040625802564660517371780402481729089536555417949051890403879840079255169295922593080322634775209689623239873322471161642996440906533187938298969649928516003704476137795166849228875"
  },
  "errors": null,
  "status": true
}
```