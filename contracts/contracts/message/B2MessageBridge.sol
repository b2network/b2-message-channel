// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/cryptography/EIP712Upgradeable.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

interface IB2MessageBridge {
    function call(uint256 to_chain_id, address contract_address, bytes calldata data) external returns (uint256 from_id);
    function send(uint256 from_chain_id, uint256 from_id, address from_sender, address contract_address, bytes calldata data, bytes[] calldata signatures) external;
}

interface IBusinessContract {
    function send(uint256 from_chain_id, uint256 from_id, address from_sender, bytes calldata data) external returns (bool success);
}

contract B2MessageBridge is IB2MessageBridge, Initializable, UUPSUpgradeable, EIP712Upgradeable, AccessControlUpgradeable {

    using ECDSA for bytes32;
    bytes32 public constant SEND_HASH_TYPE = keccak256('Send(uint256 from_chain_id,uint256 from_id,address from_sender,uint256 to_chain_id,address contract_address,bytes data)');
    bytes32 public constant ADMIN_ROLE = keccak256("admin_role");
    bytes32 public constant UPGRADE_ROLE = keccak256("upgrade_role");
    bytes32 public constant VALIDATOR_ROLE = keccak256("validator_role");
    uint256 public sequence;
    uint256 public weight;
    mapping (uint256 => uint256) public sequences;
    mapping (uint256 => mapping (uint256 => bool)) public ids;

    function initialize() public initializer {
        __EIP712_init("B2MessageBridge", "1");
        __AccessControl_init();
        __UUPSUpgradeable_init();
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(ADMIN_ROLE, msg.sender);
        _grantRole(UPGRADE_ROLE, msg.sender);
    }

    function _authorizeUpgrade(address newImplementation)
        internal
        onlyRole(UPGRADE_ROLE)
        override
    {

    }

    function setWeight(uint256 _weight) external onlyRole(ADMIN_ROLE) {
        weight = _weight;
        emit SetWeight(weight);
    }

    event SetWeight(uint256 weight);
    event Send(uint256 from_chain_id, uint256 from_id, address from_sender, uint256 to_chain_id, address contract_address, bytes data);
    event Call(uint256 from_chain_id, uint256 from_id, address from_sender, uint256 to_chain_id, address contract_address, bytes data);

    function send(uint256 from_chain_id, uint256 from_id, address from_sender, address contract_address, bytes calldata data, bytes[] calldata signatures) external {
        require(!ids[from_chain_id][from_id], "non-repeatable processing");
        uint256 weight_ = 0;
        for(uint256 i = 0; i < signatures.length; i++) {
           bool success = verify(from_chain_id, from_id, from_sender, block.chainid, contract_address, data, signatures[i]);
           if (success) {
                weight_ = weight_ + 1;
           }
        }
        require(weight_ >= weight, "verify signatures weight invalid");

        if (contract_address != address(0x0)) {
            bool success = IBusinessContract(contract_address).send(from_chain_id, from_id, from_sender, data);
            require(success, "Call failed");
        }
        emit Send(from_chain_id, from_id, from_sender, block.chainid, contract_address, data);
    }

    function call(uint256 to_chain_id, address contract_address, bytes calldata data) external returns (uint256) {
        sequences[to_chain_id]++;
        emit Call(block.chainid, sequences[to_chain_id], msg.sender, to_chain_id, contract_address, data);
        return sequences[to_chain_id];
    }

    function verify(uint256 from_chain_id, uint256 from_id, address from_sender, uint256 to_chain_id, address contract_address, bytes calldata data, bytes calldata signature) public view  returns (bool) {
         bytes32 digest  = SendHash(from_chain_id, from_id, from_sender, to_chain_id, contract_address, data);
        return hasRole(VALIDATOR_ROLE, digest.recover(signature));
    }

    function SendHash(uint256 from_chain_id, uint256 from_id, address from_sender, uint256 to_chain_id, address contract_address, bytes calldata data) public view returns (bytes32) {
        return _hashTypedDataV4(keccak256(abi.encode(SEND_HASH_TYPE,from_chain_id, from_id, from_sender, to_chain_id, contract_address, keccak256(data))));
    }

}