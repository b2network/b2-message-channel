# b2-message-channel

Message-channel is a blockchain cross-chain messaging mechanism.

## 1. Contracts

### 1.1 Contracts interface

#### 1.1.1 Message Bridge Interface

```
interface IB2MessageBridge {

    /**
     * Get the validator role for a specific chain
     * @param chain_id The ID of the chain for which to retrieve the validator role.
     * @return bytes32 The hash associated with the validator role for the specified chain ID.
     */
    function validatorRole(uint256 chain_id) external pure returns (bytes32);

    /**
     * Generate a message hash
     * @param from_chain_id The ID of the originating chain, used to identify the source of the message.
     * @param from_id The ID of the cross-chain message, used to uniquely identify the message.
     * @param from_sender The address of the sender on the originating chain.
     * @param to_chain_id The ID of the target chain, where the message will be sent.
     * @param contract_address The address of the target contract that will receive the cross-chain message.
     * @param data The input data for the target contract's cross-chain call.
     * @return bytes32 The generated message hash, used for subsequent verification and processing.
     */
    function SendHash(uint256 from_chain_id, uint256 from_id, address from_sender, uint256 to_chain_id, address contract_address, bytes calldata data) external view returns (bytes32);

    /**
     * Verify the legitimacy of a message
     * @param from_chain_id The ID of the originating chain, used to validate the source of the message.
     * @param from_id The ID of the cross-chain message, used to check if the message has already been processed.
     * @param from_sender The address of the sender on the originating chain, used to verify the sender's legitimacy.
     * @param to_chain_id The ID of the target chain, indicating where the message will be sent.
     * @param contract_address The address of the target contract that will receive the cross-chain message.
     * @param data The input data for the target contract's cross-chain call.
     * @param signature The signature of the message, used to verify its legitimacy and integrity.
     * @return bool Returns true if the verification succeeds, and false if it fails.
     */
    function verify(uint256 from_chain_id, uint256 from_id, address from_sender, uint256 to_chain_id, address contract_address, bytes calldata data, bytes calldata signature) external view returns (bool);

    /**
     * Set the weight for message processing
     * @param chain_id The ID of the chain.
     * @param _weight The weight value that influences the logic or priority of message processing.
     */
    function setWeight(uint256 chain_id, uint256 _weight) external;

    /**
     * Request cross-chain message data
     * @param to_chain_id The ID of the target chain, specifying where the message will be sent.
     * @param contract_address The address of the target contract that will receive the cross-chain message.
     * @param data The input data for the target contract's cross-chain call.
     * @return from_id The ID of the cross-chain message, returning a unique identifier to track the request.
     */
    function call(uint256 to_chain_id, address contract_address, bytes calldata data) external returns (uint256 from_id);

    /**
     * Confirm cross-chain message data
     * @param from_chain_id The ID of the originating chain, used to validate the source of the message.
     * @param from_id The ID of the cross-chain message, used to check if the message has already been processed.
     * @param from_sender The address of the sender on the originating chain (msg.sender), used to determine the sender's security based on business needs.
     * @param contract_address The address of the target contract, indicating where the message will be sent (can be a contract on the target chain or the current chain).
     * @param data The input data for the target contract's cross-chain call.
     * @param signatures An array of signatures used to verify the legitimacy of the message, ensuring only authorized senders can send the message.
     */
    function send(uint256 from_chain_id, uint256 from_id, address from_sender, address contract_address, bytes calldata data, bytes[] calldata signatures) external;

    /**
     * Set the validator role for a specific chain
     * @param chain_id The ID of the chain for which to set the validator role.
     * @param account The address of the validator, indicating which account to set the role for.
     * @param valid A boolean indicating the validity of the validator role, true for valid and false for invalid.
     */
    function setValidatorRole(uint256 chain_id, address account, bool valid) external;

    // Event declarations
    event SetWeight(uint256 chain_id, uint256 weight); // Event emitted when weight is set
    event SetValidatorRole(uint256 chain_id, address account, bool valid); // Event emitted when validator role is set
    event Send(uint256 from_chain_id, uint256 from_id, address from_sender, uint256 to_chain_id, address contract_address, bytes data); // Event emitted when a message is sent
    event Call(uint256 from_chain_id, uint256 from_id, address from_sender, uint256 to_chain_id, address contract_address, bytes data); // Event emitted when a message call is made
}
```

#### 1.1.2 Business Contract Interface

```
interface IBusinessContract {
    /**
     * Process cross-chain information in the business contract
     * @param from_chain_id The ID of the originating chain, used to validate the source of the message.
     * @param from_id The ID of the cross-chain message, used to check if the message has already been processed to prevent duplication.
     * @param from_sender The address of the sender on the originating chain, used to verify the sender's legitimacy (business needs may dictate whether verification is necessary).
     * @param data The input data for processing the cross-chain message, which may need to be decoded based on byte encoding rules.
     * @return success Indicates whether the message processing was successful, returning true for success and false for failure.
     */
    function send(uint256 from_chain_id, uint256 from_id, address from_sender, bytes calldata data) external returns (bool success);
}

```

### 1.2 Contracts code

#### 1.2.1 Message Bridge

[MessageBridge.sol](./contracts/contracts/message/MessageBridge.sol)

#### 1.2.2 Business Contract Example

[BusinessContractExample.sol](./contracts/contracts/business/BusinessContractExample.sol)

### 1.3 Deploy

#### 1.3.1 Setting up a deployment account

```
cp ./contracts/.env.test ./contracts/.env
```

[.env.test](./contracts/.env.test)

```
.env
# dev
AS_DEV_RPC_URL=https://arbitrum-sepolia.blockpi.network/v1/rpc/public
AS_DEV_PRIVATE_KEY_0=
B2_DEV_RPC_URL=https://b2-testnet.alt.technology
B2_DEV_PRIVATE_KEY_0=
# pord
AS_RPC_URL=https://arbitrum.rpc.subquery.network/public
AS_PRIVATE_KEY_0=
B2_RPC_URL=https://rpc.bsquared.network
B2_PRIVATE_KEY_0=      
```

#### 1.3.2 Deploy contract

##### 1.3.2.1 安装环境

```
cd ./contracts
npm i
```

##### 1.3.2.2 B2MessageBridge command

```
// deploy
yarn hardhat run scripts/message/deploy.js --network b2dev
// upgrade
yarn hardhat run scripts/message/upgrade.js --network b2dev
// grant_role
yarn hardhat run scripts/message/grant_role.js --network b2dev
// revoke_role
yarn hardhat run scripts/message/revoke_role.js --network b2dev
// set weight
yarn hardhat run scripts/message/set_weight.js --network b2dev
// call
yarn hardhat run scripts/message/call.js --network b2dev
// send
yarn hardhat run scripts/message/send.js --network b2dev
```

##### 1.3.2.3 BusinessContractExample command

```
// deploy
yarn hardhat run scripts/business/deploy.js --network b2dev
// upgrade
yarn hardhat run scripts/business/upgrade.js --network b2dev
// grant_role
yarn hardhat run scripts/business/grant_role.js --network b2dev
// revoke_role
yarn hardhat run scripts/business/revoke_role.js --network b2dev
```

### 1.4 Contracts instances

#### 1.4.1 Bsquared testnet

```
B2MessageBridge: 0xe55c8D6D7Ed466f66D136f29434bDB6714d8E3a5
BusinessContract: 0x804641e29f5F63a037022f0eE90A493541cCb869
```

#### 1.4.2 Arbitrum sepolia

```
B2MessageBridge: 0x2A82058E46151E337Baba56620133FC39BD5B71F
BusinessContract: 0x8Ac2C830532d7203a12C4C32C0BE4d3d15917534
```

## 2. applications

### 2.1 Listener

#### 2.1.1 Config

[listener.yaml](./applications/config/listener.yaml)

```
log:
  level: 6

particle:
  Url: https://rpc.particle.network/evm-chain
  ChainId: 1123
  ProjectUuid: 0000000000000000000000000000000000000000
  ProjectKey: 0000000000000000000000000000000000000000
  AAPubKeyAPI: https://bridge-aa-dev.bsquared.network

database:
  username: root
  password: 123456
  host: 127.0.0.1
  port: 3306
  dbname: b2_message
  loglevel: 4  # 1: Silent 2: Error 3: Warn 4: Info

bitcoin:
  status: false
  name: bitcoin
  chaintype: 2
  chainid: 0
  mainnet: false
  rpcurl: 127.0.0.1:8085
  safeblocknumber: 3
  ListenAddress: muGFcyjuyURJJsXaLXHCm43jLBmGPPU7ME
  BlockInterval: 6000
  ToChainId: 1123
  ToContractAddress: 0x0000000000000000000000000000000000000000
  BtcUser: test
  BtcPass: test
  DisableTLS: false

bsquared:
  status: false
  name: bsquared
  chaintype: 1
  mainnet: false
  chainid: 1123
  rpcurl: 127.0.0.1:8084
  safeblocknumber: 1
  ListenAddress: 0x0000000000000000000000000000000000000000
  BlockInterval: 2000
  Builders: [ "0x0000000000000000000000000000000000000000000000000000000000000000" ]

arbitrum:
  status: false
  name: arbitrum
  chaintype: 1
  chainid: 421614
  mainnet: false
  rpcurl: 127.0.0.1:8083
  safeblocknumber: 1
  ListenAddress: 0x0000000000000000000000000000000000000000
  BlockInterval: 100
  Builders: [ "0x0000000000000000000000000000000000000000000000000000000000000000" ]
```

#### 2.1.2 Quick start

```
cd applications
go build -o listener cmd/listener/main.go
./listener -f=listener
```

#### 2.1.3 Set Env && Quick start

```
log:
  level: 6
=>
APP_LOG_LEVEL=6 

APP_LOG_LEVEL=6 ./listener -f=listener
```

#### 2.1.4 Env list

```
APP_LOG_LEVEL=6

APP_PARTICLE_URL=https://rpc.particle.network/evm-chain
APP_PARTICLE_CHAINID=1123
APP_PARTICLE_PROJECTUUID=0000000000000000000000000000000000000000
APP_PARTICLE_PROJECTKEY=0000000000000000000000000000000000000000
APP_PARTICLE_AAPUBKEYAPI=https://bridge-aa-dev.bsquared.network

APP_DATABASE_USERNAME=root
APP_DATABASE_PASSWORD=123456
APP_DATABASE_HOST=127.0.0.1
APP_DATABASE_PORT=3306
APP_DATABASE_DBNAME=b2_message
APP_DATABASE_LOGLEVEL=4

APP_BITCOIN_NAME=bitcoin
APP_BITCOIN_STATUS=true
APP_BITCOIN_CHAINTYPE=2
APP_BITCOIN_CHAINID=0
APP_BITCOIN_MAINNET=false
APP_BITCOIN_RPCURL=127.0.0.1:8085
APP_BITCOIN_SAFEBLOCKNUMBER=3
APP_BITCOIN_LISTENADDRESS=muGFcyjuyURJJsXaLXHCm43jLBmGPPU7ME
APP_BITCOIN_BLOCKINTERVAL=6000
APP_BITCOIN_TOCHAINID=1123
APP_BITCOIN_TOCONTRACTADDRESS=0x0000000000000000000000000000000000000000
APP_BITCOIN_BTCUSER=test
APP_BITCOIN_BTCPASS=test
APP_BITCOIN_DISABLETLS=false

APP_BSQUARED_NAME=bsquared
APP_BSQUARED_STATUS=true
APP_BSQUARED_CHAINTYPE=1
APP_BSQUARED_MAINNET=false
APP_BSQUARED_CHAINID=1123
APP_BSQUARED_RPCURL=127.0.0.1:8084
APP_BSQUARED_SAFEBLOCKNUMBER=1
APP_BSQUARED_LISTENADDRESS=0x0000000000000000000000000000000000000000
APP_BSQUARED_BLOCKINTERVAL=2000
APP_BSQUARED_BUILDERS=0x0000000000000000000000000000000000000000000000000000000000000000

APP_ARBITRUM_NAME=arbitrum
APP_ARBITRUM_STATUS=true
APP_ARBITRUM_CHAINTYPE=1
APP_ARBITRUM_CHAINID=421614
APP_ARBITRUM_MAINNET=false
APP_ARBITRUM_RPCURL=127.0.0.1:8083
APP_ARBITRUM_SAFEBLOCKNUMBER=1
APP_ARBITRUM_LISTENADDRESS=0x0000000000000000000000000000000000000000
APP_ARBITRUM_BLOCKINTERVAL=100
APP_ARBITRUM_BUILDERS=0x0000000000000000000000000000000000000000000000000000000000000000
```

### 2.2 Proposer

#### 2.2.1 Config

[proposer.yaml](./applications/config/proposer.yaml)

```
log:
  level: 6

database:
  username: root
  password: 123456
  host: 127.0.0.1
  port: 3306
  dbname: b2_message
  loglevel: 4  # 1: Silent 2: Error 3: Warn 4: Info

bsquared:
  status: false
  name: bsquared
  chaintype: 1
  mainnet: false
  chainid: 1123
  rpcurl: 127.0.0.1:8081
  safeblocknumber: 1
  ListenAddress: 0x0000000000000000000000000000000000000000
  BlockInterval: 2000
  NodeKey: 0000000000000000000000000000000000000000000000000000000000000000
  NodePort: 20000
  SignatureWeight: 1
  Validators: [ "0x0000000000000000000000000000000000000000" ]

arbitrum:
  status: false
  name: arbitrum
  chaintype: 1
  chainid: 421614
  mainnet: false
  rpcurl: 127.0.0.1:8082
  safeblocknumber: 1
  ListenAddress: 0x0000000000000000000000000000000000000000
  BlockInterval: 100
  NodeKey: 0000000000000000000000000000000000000000000000000000000000000000
  NodePort: 20001
  SignatureWeight: 1
  Validators: [ "0x0000000000000000000000000000000000000000" ]

bitcoin:
  status: false
  name: bitcoin
  chaintype: 2
  chainid: 0
  mainnet: false
  rpcurl: 127.0.0.1:8083
  safeblocknumber: 3
  ListenAddress: muGFcyjuyURJJsXaLXHCm43jLBmGPPU7ME
  BlockInterval: 6000
  ToChainId: 1123
  ToContractAddress: 0x0000000000000000000000000000000000000000
  BtcUser: 000000000000000000
  BtcPass: 000000000000000000
  DisableTLS: true
  NodeKey: 0000000000000000000000000000000000000000000000000000000000000000
  NodePort: 20002
  SignatureWeight: 1
  Validators: [ "0x0000000000000000000000000000000000000000" ]

particle:
  Url: https://rpc.particle.network/evm-chain
  ChainId: 1123
  ProjectUuid: 000000000000000000
  ProjectKey: 000000000000000000
  AAPubKeyAPI: https://bridge-aa-dev.bsquared.network
```

#### 2.2.2 Quick start

```
cd applications
go build -o proposer cmd/proposer/main.go
./proposer -f=proposer
```

#### 2.2.3 Set Env && Quick start

```
log:
  level: 6
=>
APP_LOG_LEVEL=6 

APP_LOG_LEVEL=6 ./proposer -f=proposer
```

#### 2.2.4 Env list

```
APP_LOG_LEVEL=6

APP_DATABASE_USERNAME=root
APP_DATABASE_PASSWORD=123456
APP_DATABASE_HOST=127.0.0.1
APP_DATABASE_PORT=3306
APP_DATABASE_DBNAME=b2_message
APP_DATABASE_LOGLEVEL=4

APP_BSQUARED_NAME=bsquared
APP_BSQUARED_STATUS=true
APP_BSQUARED_CHAINTYPE=1
APP_BSQUARED_MAINNET=false
APP_BSQUARED_CHAINID=1123
APP_BSQUARED_RPCURL=127.0.0.1:8081
APP_BSQUARED_SAFEBLOCKNUMBER=1
APP_BSQUARED_LISTENADDRESS=0x0000000000000000000000000000000000000000
APP_BSQUARED_BLOCKINTERVAL=2000
APP_BSQUARED_NODEKEY=0000000000000000000000000000000000000000000000000000000000000000
APP_BSQUARED_NODEPORT=20000
APP_BSQUARED_SIGNATUREWEIGHT=1
APP_BSQUARED_VALIDATORS=0x0000000000000000000000000000000000000000

APP_ARBITRUM_NAME=arbitrum
APP_ARBITRUM_STATUS=true
APP_ARBITRUM_CHAINTYPE=1
APP_ARBITRUM_CHAINID=421614
APP_ARBITRUM_MAINNET=false
APP_ARBITRUM_RPCURL=127.0.0.1:8082
APP_ARBITRUM_SAFEBLOCKNUMBER=1
APP_ARBITRUM_LISTENADDRESS=0x0000000000000000000000000000000000000000
APP_ARBITRUM_BLOCKINTERVAL=100
APP_ARBITRUM_NODEKEY=0000000000000000000000000000000000000000000000000000000000000000
APP_ARBITRUM_NODEPORT=20001
APP_ARBITRUM_SIGNATUREWEIGHT=1
APP_ARBITRUM_VALIDATORS=0x0000000000000000000000000000000000000000

APP_BITCOIN_NAME=bitcoin
APP_BITCOIN_STATUS=true
APP_BITCOIN_CHAINTYPE=2
APP_BITCOIN_CHAINID=0
APP_BITCOIN_MAINNET=false
APP_BITCOIN_RPCURL=127.0.0.1:8083
APP_BITCOIN_SAFEBLOCKNUMBER=3
APP_BITCOIN_LISTENADDRESS=muGFcyjuyURJJsXaLXHCm43jLBmGPPU7ME
APP_BITCOIN_BLOCKINTERVAL=6000
APP_BITCOIN_TOCHAINID=1123
APP_BITCOIN_TOCONTRACTADDRESS=0x0000000000000000000000000000000000000000
APP_BITCOIN_BTCUSER=000000000000000000
APP_BITCOIN_BTCPASS=000000000000000000
APP_BITCOIN_DISABLETLS=true
APP_BITCOIN_NODEKEY=0000000000000000000000000000000000000000000000000000000000000000
APP_BITCOIN_NODEPORT=20002
APP_BITCOIN_SIGNATUREWEIGHT=1
APP_BITCOIN_VALIDATORS=0x0000000000000000000000000000000000000000

APP_PARTICLE_URL=https://rpc.particle.network/evm-chain
APP_PARTICLE_CHAINID=1123
APP_PARTICLE_PROJECTUUID=000000000000000000
APP_PARTICLE_PROJECTKEY=000000000000000000
APP_PARTICLE_AAPUBKEYAPI=https://bridge-aa-dev.bsquared.network
```

### 2.3 Validator

#### 2.3.1 Config

[validator.yaml](./applications/config/validator.yaml)

```
log:
  level: 6

bsquared:
  status: false
  name: bsquared
  chaintype: 1
  mainnet: false
  chainid: 1123
  rpcurl: 127.0.0.1:8081
  safeblocknumber: 1
  ListenAddress: 0x0000000000000000000000000000000000000000
  BlockInterval: 2000
  NodeKey: 0000000000000000000000000000000000000000000000000000000000000000
  Endpoint: /ip4/127.0.0.1/tcp/20000/p2p/16Uiu2HAkwynt59WSsNRS9sk1aszgeQ1PXUS8ax3a3tsewaVMgvZX # /ip4/{host}/tcp/{port}/p2p/{peerId}
  SignatureWeight: 1

arbitrum:
  status: false
  name: arbitrum
  chaintype: 1
  chainid: 421614
  mainnet: false
  rpcurl: 127.0.0.1:8082
  safeblocknumber: 1
  ListenAddress: 0x0000000000000000000000000000000000000000
  BlockInterval: 100
  NodeKey: 0000000000000000000000000000000000000000000000000000000000000000
  Endpoint: /ip4/127.0.0.1/tcp/20001/p2p/16Uiu2HAkwynt59WSsNRS9sk1aszgeQ1PXUS8ax3a3tsewaVMgvZX # /ip4/{host}/tcp/{port}/p2p/{peerId}
  SignatureWeight: 1

bitcoin:
  status: false
  name: bitcoin
  chaintype: 2
  chainid: 0
  mainnet: false
  rpcurl: 127.0.0.1:8083
  safeblocknumber: 3
  ListenAddress: muGFcyjuyURJJsXaLXHCm43jLBmGPPU7ME
  BlockInterval: 6000
  ToChainId: 1123
  ToContractAddress: 0x0000000000000000000000000000000000000000
  BtcUser: 000000000000000000
  BtcPass: 000000000000000000
  DisableTLS: true
  NodeKey: 0000000000000000000000000000000000000000000000000000000000000000
  Endpoint: /ip4/127.0.0.1/tcp/20001/p2p/16Uiu2HAkwynt59WSsNRS9sk1aszgeQ1PXUS8ax3a3tsewaVMgvZX # /ip4/{host}/tcp/{port}/p2p/{peerId}
  SignatureWeight: 1

particle:
  Url: https://rpc.particle.network/evm-chain
  ChainId: 1123
  ProjectUuid: 000000000000000000
  ProjectKey: 000000000000000000
  AAPubKeyAPI: https://bridge-aa-dev.bsquared.network

```

#### 2.3.2 Quick start

```
cd applications
go build -o validator cmd/validator/main.go
./validator -f=validator
```

#### 2.3.3 Set Env && Quick start

```
log:
  level: 6
=>
APP_LOG_LEVEL=6 

APP_LOG_LEVEL=6 ./validator -f=validator
```

#### 2.3.4 Env list

```
APP_LOG_LEVEL=6

APP_BSQUARED_NAME=bsquared
APP_BSQUARED_STATUS=true
APP_BSQUARED_CHAINTYPE=1
APP_BSQUARED_MAINNET=false
APP_BSQUARED_CHAINID=1123
APP_BSQUARED_RPCURL=127.0.0.1:8081
APP_BSQUARED_SAFEBLOCKNUMBER=1
APP_BSQUARED_LISTENADDRESS=0x0000000000000000000000000000000000000000
APP_BSQUARED_BLOCKINTERVAL=2000
APP_BSQUARED_NODEKEY=0000000000000000000000000000000000000000000000000000000000000000
APP_BSQUARED_ENDPOINT=/ip4/127.0.0.1/tcp/20000/p2p/16Uiu2HAkwynt59WSsNRS9sk1aszgeQ1PXUS8ax3a3tsewaVMgvZX
APP_BSQUARED_SIGNATUREWEIGHT=1

APP_ARBITRUM_NAME=arbitrum
APP_ARBITRUM_STATUS=true
APP_ARBITRUM_CHAINTYPE=1
APP_ARBITRUM_CHAINID=421614
APP_ARBITRUM_MAINNET=false
APP_ARBITRUM_RPCURL=127.0.0.1:8082
APP_ARBITRUM_SAFEBLOCKNUMBER=1
APP_ARBITRUM_LISTENADDRESS=0x0000000000000000000000000000000000000000
APP_ARBITRUM_BLOCKINTERVAL=100
APP_ARBITRUM_NODEKEY=0000000000000000000000000000000000000000000000000000000000000000
APP_ARBITRUM_ENDPOINT=/ip4/127.0.0.1/tcp/20001/p2p/16Uiu2HAkwynt59WSsNRS9sk1aszgeQ1PXUS8ax3a3tsewaVMgvZX
APP_ARBITRUM_SIGNATUREWEIGHT=1

APP_BITCOIN_NAME=bitcoin
APP_BITCOIN_STATUS=true
APP_BITCOIN_CHAINTYPE=2
APP_BITCOIN_CHAINID=0
APP_BITCOIN_MAINNET=false
APP_BITCOIN_RPCURL=127.0.0.1:8083
APP_BITCOIN_SAFEBLOCKNUMBER=3
APP_BITCOIN_LISTENADDRESS=muGFcyjuyURJJsXaLXHCm43jLBmGPPU7ME
APP_BITCOIN_BLOCKINTERVAL=6000
APP_BITCOIN_TOCHAINID=1123
APP_BITCOIN_TOCONTRACTADDRESS=0x0000000000000000000000000000000000000000
APP_BITCOIN_BTCUSER=000000000000000000
APP_BITCOIN_BTCPASS=000000000000000000
APP_BITCOIN_DISABLETLS=true
APP_BITCOIN_NODEKEY=0000000000000000000000000000000000000000000000000000000000000000
APP_BITCOIN_ENDPOINT=/ip4/127.0.0.1/tcp/20001/p2p/16Uiu2HAkwynt59WSsNRS9sk1aszgeQ1PXUS8ax3a3tsewaVMgvZX
APP_BITCOIN_SIGNATUREWEIGHT=1

APP_PARTICLE_URL=https://rpc.particle.network/evm-chain
APP_PARTICLE_CHAINID=1123
APP_PARTICLE_PROJECTUUID=000000000000000000
APP_PARTICLE_PROJECTKEY=000000000000000000
APP_PARTICLE_AAPUBKEYAPI=https://bridge-aa-dev.bsquared.network
```

### 2.4 Builder

#### 2.4.1 Config

[builder.yaml](./applications/config/builder.yaml)

```
log:
  level: 6

particle:
  Url: https://rpc.particle.network/evm-chain
  ChainId: 1123
  ProjectUuid: 0000000000000000000000000000000000000000
  ProjectKey: 0000000000000000000000000000000000000000
  AAPubKeyAPI: https://bridge-aa-dev.bsquared.network

database:
  username: root
  password: 123456
  host: 127.0.0.1
  port: 3306
  dbname: b2_message
  loglevel: 4  # 1: Silent 2: Error 3: Warn 4: Info

bitcoin:
  status: true
  name: bitcoin
  chaintype: 2
  chainid: 0
  mainnet: false
  rpcurl: 127.0.0.1:8085
  safeblocknumber: 3
  ListenAddress: muGFcyjuyURJJsXaLXHCm43jLBmGPPU7ME
  BlockInterval: 6000
  ToChainId: 1123
  ToContractAddress: 0x0000000000000000000000000000000000000000
  BtcUser: test
  BtcPass: test
  DisableTLS: false

bsquared:
  status: true
  name: bsquared
  chaintype: 1
  mainnet: false
  chainid: 1123
  rpcurl: 127.0.0.1:8084
  safeblocknumber: 1
  ListenAddress: 0x0000000000000000000000000000000000000000
  BlockInterval: 2000
  Builders: [ "0x0000000000000000000000000000000000000000000000000000000000000000" ]

arbitrum:
  status: true
  name: arbitrum
  chaintype: 1
  chainid: 421614
  mainnet: false
  rpcurl: 127.0.0.1:8083
  safeblocknumber: 1
  ListenAddress: 0x0000000000000000000000000000000000000000
  BlockInterval: 100
  Builders: [ "0x0000000000000000000000000000000000000000000000000000000000000000" ]
```

#### 2.4.2 Quick start

```
cd applications
go build -o builder cmd/builder/main.go
./builder -f=builder
```

#### 2.4.3 Set Env && Quick start

```
log:
  level: 6
=>
APP_LOG_LEVEL=6 

APP_LOG_LEVEL=6 ./builder -f=builder
```

#### 2.4.4 Env list

```
APP_LOG_LEVEL=6
APP_PARTICLE_URL=https://rpc.particle.network/evm-chain
APP_PARTICLE_CHAINID=1123
APP_PARTICLE_PROJECTUUID=0000000000000000000000000000000000000000
APP_PARTICLE_PROJECTKEY=0000000000000000000000000000000000000000
APP_PARTICLE_AAPUBKEYAPI=https://bridge-aa-dev.bsquared.network

APP_DATABASE_USERNAME=root
APP_DATABASE_PASSWORD=123456
APP_DATABASE_HOST=127.0.0.1
APP_DATABASE_PORT=3306
APP_DATABASE_DBNAME=b2_message
APP_DATABASE_LOGLEVEL=4

APP_BITCOIN_NAME=bitcoin
APP_BITCOIN_STATUS=true
APP_BITCOIN_CHAINTYPE=2
APP_BITCOIN_CHAINID=0
APP_BITCOIN_MAINNET=false
APP_BITCOIN_RPCURL=127.0.0.1:8085
APP_BITCOIN_SAFEBLOCKNUMBER=3
APP_BITCOIN_LISTENADDRESS=muGFcyjuyURJJsXaLXHCm43jLBmGPPU7ME
APP_BITCOIN_BLOCKINTERVAL=6000
APP_BITCOIN_TOCHAINID=1123
APP_BITCOIN_TOCONTRACTADDRESS=0x0000000000000000000000000000000000000000
APP_BITCOIN_BTCUSER=test
APP_BITCOIN_BTCPASS=test
APP_BITCOIN_DISABLETLS=false

APP_BSQUARED_NAME=bsquared
APP_BSQUARED_STATUS=true
APP_BSQUARED_CHAINTYPE=1
APP_BSQUARED_MAINNET=false
APP_BSQUARED_CHAINID=1123
APP_BSQUARED_RPCURL=127.0.0.1:8084
APP_BSQUARED_SAFEBLOCKNUMBER=1
APP_BSQUARED_LISTENADDRESS=0x0000000000000000000000000000000000000000
APP_BSQUARED_BLOCKINTERVAL=2000
APP_BSQUARED_BUILDERS=0x0000000000000000000000000000000000000000000000000000000000000000

APP_ARBITRUM_NAME=arbitrum
APP_ARBITRUM_STATUS=true
APP_ARBITRUM_CHAINTYPE=1
APP_ARBITRUM_CHAINID=421614
APP_ARBITRUM_MAINNET=false
APP_ARBITRUM_RPCURL=127.0.0.1:8083
APP_ARBITRUM_SAFEBLOCKNUMBER=1
APP_ARBITRUM_LISTENADDRESS=0x0000000000000000000000000000000000000000
APP_ARBITRUM_BLOCKINTERVAL=100
APP_ARBITRUM_BUILDERS=0x0000000000000000000000000000000000000000000000000000000000000000
```
