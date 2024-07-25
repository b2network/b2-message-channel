const {ethers, upgrades, network} = require("hardhat");

async function main() {
    /**
     * b2dev: yarn hardhat run scripts/business/deploy.js --network b2dev
     * as: yarn hardhat run scripts/business/deploy.js --network as
     * b2: yarn hardhat run scripts/business/deploy.js --network b2
     */
    // 0x1c66cBEE6d4660459Fda5aa936e727398175E981
    const [owner] = await ethers.getSigners()
    console.log("Owner Address:", owner.address);

    // deploy
    const BusinessContractExample = await ethers.getContractFactory("BusinessContractExample");
    const instance = await upgrades.deployProxy(BusinessContractExample);
    await instance.waitForDeployment();
    console.log("BusinessContractExample Address:", instance.target);

    // Upgrading
    // const BusinessContractExampleV2 = await ethers.getContractFactory("BusinessContractExample");
    // const upgraded = await upgrades.upgradeProxy("0x67292d3b05C9fc8391372C23826236983C23b9a2", BusinessContractExampleV2);
    // console.log("BusinessContractExampleV2 upgraded:", upgraded.target);

}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })