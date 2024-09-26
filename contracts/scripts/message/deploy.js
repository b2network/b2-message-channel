const {ethers, upgrades, network} = require("hardhat");

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/message/deploy.js --network b2dev
     * asdev: yarn hardhat run scripts/message/deploy.js --network asdev
     * # pord
     * b2: yarn hardhat run scripts/message/deploy.js --network b2
     * as: yarn hardhat run scripts/message/deploy.js --network as
     */
    const [owner] = await ethers.getSigners()
    console.log("Owner Address:", owner.address); // 0x2BC22b1754ff4aDea4Ef9bdF9b16A7210bC45579

    const B2MessageBridge = await ethers.getContractFactory("B2MessageBridge");
    const instance = await upgrades.deployProxy(B2MessageBridge);
    await instance.waitForDeployment();
    console.log("B2MessageBridge Address:", instance.target);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })