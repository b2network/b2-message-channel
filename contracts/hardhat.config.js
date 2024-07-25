require("@nomicfoundation/hardhat-toolbox");
require('@openzeppelin/hardhat-upgrades');

// B2_TEST_DEV
const B2_DEV = {
    RPC_URL: "https://testnet-rpc.bsquared.network",
    PRIVATE_KEY_LIST: [],
}

const AS_DEV = {
    RPC_URL: "https://arbitrum-sepolia.blockpi.network/v1/rpc/public",
    PRIVATE_KEY_LIST: [],
}

const B2 = {
    RPC_URL: "https://rpc.bsquared.network",
    PRIVATE_KEY_LIST: []
}

task("accounts", "Prints the list of accounts", async (taskArgs, hre) => {
    const accounts = await hre.ethers.getSigners()

    for (const account of accounts) {
        console.log(account.address)
    }
})

module.exports = {
    networks: {
        hardhat: {}, as: {
            blockConfirmations: 1, url: AS_DEV.RPC_URL, accounts: AS_DEV.PRIVATE_KEY_LIST,
        }, b2dev: {
            blockConfirmations: 1, url: B2_DEV.RPC_URL, accounts: B2_DEV.PRIVATE_KEY_LIST,
        }, b2: {
            blockConfirmations: 1, url: B2.RPC_URL, accounts: B2.PRIVATE_KEY_LIST, gasPrice: 352,
        }
    }, solidity: {
        version: "0.8.20", settings: {
            optimizer: {
                enabled: true, runs: 1000,
            },
        },
    }, etherscan: {
        apiKey: {
            b2test: "abc", b2: "abc"
        }, customChains: [{
            network: "b2test", chainId: 1123, urls: {
                apiURL: "https://testnet-backend.bsquared.network/api",
                browserURL: "https://testnet-explorer.bsquared.network"
            }
        }, {
            network: "b2", chainId: 223, urls: {
                apiURL: "https://mainnet-backend.bsquared.network/api",
                browserURL: "https://mainnet-blockscout.bsquared.network"
                // apiURL: "https://bsquared.l2scan.co/api", browserURL: "https://bsquared.l2scan.co"
            }
        }]
    }
}