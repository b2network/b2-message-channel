// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/cryptography/EIP712Upgradeable.sol";

interface IBusinessContract {
    function send(uint256 from_chain_id, uint256 from_id, address from_sender, bytes calldata data) external returns (bool success);
}

contract BusinessContractExample is IBusinessContract, Initializable, UUPSUpgradeable, EIP712Upgradeable, AccessControlUpgradeable {

    bytes32 public constant ADMIN_ROLE = keccak256("admin_role");
    bytes32 public constant UPGRADE_ROLE = keccak256("upgrade_role");
    bytes32 public constant SENDER_ROLE = keccak256("sender_role");

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

    event Send(uint256 from_chain_id, uint256 from_id, address from_sender, bytes data);

    function send(uint256 from_chain_id, uint256 from_id, address from_sender, bytes calldata data) external onlyRole(SENDER_ROLE) override returns (bool success) {
        // TODO 1. 验证 from_chain_id 和 from_sender 有效性
        // TODO 2. 验证 from_id 是否执行过
        // TODO 3. 解析 data 执行业务逻辑
        emit Send(from_chain_id, from_id, from_sender, data);
        return true;
    }

}