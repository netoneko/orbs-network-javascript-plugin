# Orbs JS smart contract runtime

## Building

The layout of the projects:
```
orbs-network/
    - orbs-network-go
    - orbs-network-javascript-plugin
```

### Mac
```bash
./git-submodule-checkout.sh

cd ./vendor/github.com/ry/v8worker2/ && ./build.py && cd -

go get
./build-binaries.sh
```

### Linux

`./docker-build.sh`

## Configuration

The build will produce `_bin/orbs-javascript-plugin` which should be inserted into the Docker image of the node or gamma.

To enable the plugin in the node/gamma, update the configuration file:

```json
{
    "processor-plugin-path": "/opt/orbs/plugins/orbs-javascript-plugin"
}
```

## E2E testing

**Make sure** that the image `orbs-network/gamma:experimental` contains `/opt/orbs/plugins/orbs-javascript-plugin`.

Start gamma:

```bash
gamma-cli start-local -env experimental -override-config '{"processor-plugin-path": "/opt/orbs/plugins/orbs-javascript-plugin"}'
```

In orbs-network-go project:

```bash
API_ENDPOINT=http://localhost:8080 go test ./test/e2e/... -run TestDeploymentOfJavascriptContract -v -count 1
```