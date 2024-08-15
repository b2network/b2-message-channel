const {ethers, upgrades, network} = require("hardhat");

async function main() {
    /**
     * b2dev: yarn hardhat run scripts/message/deploy.js --network b2dev
     * as: yarn hardhat run scripts/message/deploy.js --network as
     * b2: yarn hardhat run scripts/message/deploy.js --network b2
     */
    // 2445771839062500
    //  991038400000000
    const [owner] = await ethers.getSigners()
    console.log("Owner Address:", owner.address);

    // const B2MessageBridge = await ethers.getContractFactory("B2MessageBridge");
    // const instance = await upgrades.deployProxy(B2MessageBridge);
    // await instance.waitForDeployment();
    // console.log("B2MessageBridge Address:", instance.target);


    // Upgrading
    const simpleBridgeV4 = await ethers.getContractFactory("B2MessageBridge");
    const upgraded = await upgrades.upgradeProxy("0xc7441Ac47596D1356fcc70062dA0462FcA98E14e", simpleBridgeV4, {
        gasPrice: ethers.parseUnits('352', 'wei')
    });
    console.log("SimpleBridge upgraded:", upgraded.target);

}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })