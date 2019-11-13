# Orbs JS smart contract runtime

## Example contract

```javascript
import { Address, State, Events, Verify, Uint64 } from "orbs-contract-sdk/v1";

// Define uint64 later
const TOTAL_SUPPLY = Uint64("10000000000");

// doesn't do anything except validation
export function TransferEvent(from, to, amount) {
    Verify.bytes(from);
    Verify.bytes(to);
    Verify.uint64(amount);
}

export function _init() {
    const ownerAddress = Address.getSignerAddress();
    State.writeUint64(ownerAddress, TOTAL_SUPPLY);
}

export function totalSupply() {
    return TOTAL_SUPPLY;
}

export function transfer(amount, targetAddress) {
    Verify.uint64(amount);
    Verify.bytes(targetAddress);

    // sender
    const callerAddress = Address.getCallerAddress();
    const callerBalance = State.readUint64(callerAddress);
    if (callerBalance < amount) {
        throw new Error(`transfer of ${amount} failed since balance is only ${callerBalance}`);
    }
    State.writeUint64(callerAddress, callerBalance-amount);

    // recipient
    Address.validateAddress(targetAddress);
    const targetBalance = State.readUint64(targetAddress);
    State.writeUint64(targetAddress, targetBalance+amount);

    Events.emitEvent(TransferEvent, callerAddress, targetAddress, amount);
}

export function balanceOf(targetAddress) {
    Verify.bytes(targetAddress);
    Address.validateAddress(targetAddress);
    return State.readUint64(targetAddress);
}
```

## Running locally

```bash
docker pull orbsnetwork/gamma:experimental-js && docker tag orbsnetwork/gamma:experimental-js orbsnetwork/experimental

gamma-cli start-local -env experimental -override-config '{"experimental-external-processor-plugin-path": "/opt/orbs/plugins/orbs-javascript-plugin"}'
```

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

## Testing

### Local

```bash
./build-binaries.sh
```

### Local E2E
**Make sure** that the image `orbs-network/gamma:experimental` contains `/opt/orbs/plugins/orbs-javascript-plugin`.

Start gamma:

```bash
gamma-cli start-local -env experimental -override-config '{"processor-plugin-path": "/opt/orbs/plugins/orbs-javascript-plugin"}'
```

In orbs-network-go project:

```bash
JS_ENABLED=true API_ENDPOINT=http://localhost:8080 go test ./test/e2e/... -run TestDeploymentOfJavascriptContract -v -count 1
```

### Remote E2E

```bash
# in orbs-network-go

JS_ENABLED=true API_ENDPOINT=http://35.156.69.19/vchains/1003 VCHAIN=1003 go test ./services/processor/javascript/test/e2e/... -run TestDeploymentOfJavascriptContractInteroperableWithGo -tags javascript -v
```

## Release process

```bash
# in orbs-network-go

export BUILD_FLAG=javascript
./docker/build/build.sh

# in orbs-network-javascript-plugin

./docker-build.sh

./release/build.sh

# start gamma

gamma-cli start-local -env experimental -override-config '{"experimental-external-processor-plugin-path": "/opt/orbs/plugins/orbs-javascript-plugin"}'
```

