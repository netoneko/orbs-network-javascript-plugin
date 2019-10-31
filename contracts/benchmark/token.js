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
