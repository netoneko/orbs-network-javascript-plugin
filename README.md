# Orbs JS smart contract runtime

## Example contract

```javascript
import { Address, State, Events, Verify, Uint32 } from "orbs-contract-sdk/v1";

// Define uint32 later
const TOTAL_SUPPLY = Uint32("10000000");

// doesn't do anything except validation
export function TransferEvent(from, to, amount) {
    Verify.bytes(from);
    Verify.bytes(to);
    Verify.uint32(amount);
}

export function _init() {
    const ownerAddress = Address.getSignerAddress();
    State.writeUint32(ownerAddress, TOTAL_SUPPLY);
}

export function totalSupply() {
    return TOTAL_SUPPLY;
}

export function transfer(amount, targetAddress) {
    Verify.uint32(amount);
    Verify.bytes(targetAddress);

    // sender
    const callerAddress = Address.getCallerAddress();
    const callerBalance = State.readUint32(callerAddress);
    if (callerBalance < amount) {
        throw new Error(`transfer of ${amount} failed since balance is only ${callerBalance}`);
    }
    State.writeUint32(callerAddress, callerBalance-amount);

    // recipient
    Address.validateAddress(targetAddress);
    const targetBalance = State.readUint32(targetAddress);
    State.writeUint32(targetAddress, targetBalance+amount);

    Events.emitEvent(TransferEvent, callerAddress, targetAddress, amount);
}

export function balanceOf(targetAddress) {
    Verify.bytes(targetAddress);
    Address.validateAddress(targetAddress);
    return State.readUint32(targetAddress);
}
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