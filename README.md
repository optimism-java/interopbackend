# Interopbackend

this project is an interop-explorer backend.

Mainly listens to the events from L2ToL2CrossDomainMessenger and CrossL2Inbox, and then writes the data to the database.
And determine whether a block includes cross-chain transactions based on the ExecutingMessage event of the CrossL2Outbox
contract.

contract:
```azure
 - L2ToL2CrossDomainMessenger: 0x4200000000000000000000000000000000000023
```     
events:
```azure
event SentMessage(uint256 indexed destination, address indexed target, 
     uint256 indexed messageNonce, address sender, bytes message);
 
event RelayedMessage(uint256 indexed source, uint256 indexed messageNonce, 
    bytes32 indexed messageHash);
```

contract:
```
 - CrossL2Inbox:               0x4200000000000000000000000000000000000022
```
events:
```azure
event ExecutingMessage(bytes32 indexed msgHash, Identifier id);
```

## Install preperation: Supersim

Follow [this guide](https://github.com/ethereum-optimism/supersim/blob/main/README.md) or [supersim](https://github.com/ethereum-optimism/supersim) to install supersim.

## run the project

```azure
./deploy/start.sh
```

## stop this project

```azure
./deploy/stop.sh
```
