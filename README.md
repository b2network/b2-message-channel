# b2-message-channel

Message-channel is a blockchain cross-chain messaging mechanism.

## Usage

1. Implement the BusinessContract interface

```
interface IBusinessContract {
    // TODO 1. Verify the validity of from_chain_id and from_sender
    // TODO 2. Verify that from_id has been executed
    // TODO 3. Parse data and execute service logic
    function send(uint256 from_chain_id, uint256 from_id, address from_sender, bytes calldata data) external returns (bool success);
}
```

Reference: [BusinessContractExample.sol](https://github.com/b2network/b2-message-channel/blob/dev/contracts/contracts/business/BusinessContractExample.sol)


2. Call the call method of the B2MessageBridge contract to achieve cross-chain message passing

```
// to_chain_id is the destination chain ID
// contract_address is the address of the BusinessContract on the destination chain
// data is the metadata to be passed cross-chain
function call(uint256 to_chain_id, address contract_address, bytes calldata data) external returns (uint256)
```

3. Contract instances

bsquared testnet

```
B2MessageBridge: 0xc7441Ac47596D1356fcc70062dA0462FcA98E14e
BusinessContract: 0x91171cf194a4B66Bd459Ada038397c7e890FB9D4
```

arbitrum sepolia

```
B2MessageBridge: 0x2A82058E46151E337Baba56620133FC39BD5B71F
BusinessContract: 0x8Ac2C830532d7203a12C4C32C0BE4d3d15917534
```

4. Example call

```
messageAddress = "0xc7441Ac47596D1356fcc70062dA0462FcA98E14e";
businessAddress = "0x8Ac2C830532d7203a12C4C32C0BE4d3d15917534";
to_chain_id = 421614;
data = '0x1234';
```
Command: yarn hardhat run scripts/message/call.js --network b2dev

