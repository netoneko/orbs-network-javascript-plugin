import { Address, State, Verify, Uint32 } from "orbs-contract-sdk/v1";

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

    Events.emitEvent(TransferEvent, from, to, amount);
}

export function getBalance() {
    Verify.Bytes(targetAddress);
    Address.validateAddress(targetAddress);
    return State.readUint32(targetAddress);
}