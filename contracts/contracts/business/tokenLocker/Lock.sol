// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/cryptography/EIP712Upgradeable.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";


interface IB2MessageSharing {
    function call(uint256 to_chain_id, address to_business_contract, bytes calldata to_message) external returns (uint256 from_id);
}

interface IBusinessContract {
    function send(uint256 from_chain_id, uint256 from_id, address from_sender, bytes calldata data) external returns (bool success);
}


contract TokenLockerContract is IBusinessContract, Initializable, UUPSUpgradeable, EIP712Upgradeable, AccessControlUpgradeable {

    using SafeERC20 for IERC20;

    bytes32 public constant ADMIN_ROLE = keccak256("admin_role");
    bytes32 public constant UPGRADE_ROLE = keccak256("upgrade_role");
    bytes32 public constant SENDER_ROLE = keccak256("sender_role");

    IB2MessageSharing public messageSharing;

    function initialize() public initializer {
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

    function setB2MessageSharing(address sharing_address) external onlyRole(ADMIN_ROLE) {
        messageSharing = IB2MessageSharing(sharing_address);
    }

    function setLocks(uint256 chain_id,address lock_address) external onlyRole(ADMIN_ROLE) {
        locks[chain_id] = lock_address;
    }

    function setTokens(uint256 chain_id, address token_address ,address to_token_address) external onlyRole(ADMIN_ROLE) {
        tokens[chain_id][token_address] = to_token_address;
    }

    mapping (uint256 => address) public locks;
    mapping (uint256 => mapping (uint256 => bool)) public ids;
    mapping (uint256 => mapping (address => address)) public tokens;

    event Unlock(uint256 from_chain_id ,uint256 from_id, address from_address, address token_address, uint256 to_chain_id, address to_business_contract, address to_address, uint256 amount);

    function send(uint256 from_chain_id, uint256 from_id, address from_sender, bytes calldata message) external onlyRole(SENDER_ROLE) override returns (bool success) {
        require(locks[from_chain_id] == from_sender, "Invalid chain id or from_sender");
        require(!ids[from_chain_id][from_id], "Have been executed");
        (address from_token_address, address from_address, address to_address, uint256 amount) = decodeLockData(message);

        address token_address = tokens[from_chain_id][from_token_address];
        require(token_address != address(0x0), "Invalid token");

        IERC20(token_address).transfer(to_address, amount);
        emit Unlock(from_chain_id, from_id, from_address, token_address, block.chainid, address(this), to_address, amount);
        return true;
    }

    event Lock(uint256 from_chain_id ,uint256 from_id, address from_address, address token_address, uint256 to_chain_id, address to_business_contract, address to_address, uint256 amount);

    function lock(address token_address, uint256 amount, uint256 to_chain_id, address to_business_contract, address to_address) external {
        IERC20(token_address).transferFrom(msg.sender, address(this), amount);
        bytes memory to_message = encodeLockData(token_address, msg.sender, to_address, amount);

        uint256 from_id =  messageSharing.call(to_chain_id, to_business_contract, to_message);
        emit Lock(block.chainid ,from_id, msg.sender, token_address, to_chain_id, to_business_contract, to_address, amount);
    }

    function encodeLockData(address token_address, address from_address, address to_address, uint256 amount) public pure returns (bytes memory) {
        return abi.encode(token_address, from_address, to_address, amount);
    }

    function decodeLockData(bytes memory data) public pure returns (address token_address, address from_address, address to_address, uint256 amount) {
        (token_address, from_address, to_address, amount) = abi.decode(data, (address, address, address, uint256));
    }
}