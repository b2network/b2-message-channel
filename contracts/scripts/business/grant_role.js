const {ethers, run, network, upgrades} = require("hardhat")

let businessAddress;
let senderAddress;

async function main() {
    /**
     * b2dev: yarn hardhat run scripts/business/grant_role.js --network b2dev
     * as: yarn hardhat run scripts/business/grant_role.js --network as
     * b2: yarn hardhat run scripts/business/grant_role.js --network b2
     */

    if (network.name == 'b2dev') {
        businessAddress = "0x91171cf194a4B66Bd459Ada038397c7e890FB9D4";
        senderAddress = "0xc7441Ac47596D1356fcc70062dA0462FcA98E14e";
    } else if (network.name == 'as') {
        businessAddress = "0x8Ac2C830532d7203a12C4C32C0BE4d3d15917534";
        senderAddress = "0x2A82058E46151E337Baba56620133FC39BD5B71F";
    } else if (network.name == 'b2') {
        businessAddress = "";
        senderAddress = "";
    }

    const [owner] = await ethers.getSigners()

    // Launchpad
    const BusinessContractExample = await ethers.getContractFactory("BusinessContractExample");
    const instance = await BusinessContractExample.attach(businessAddress)

    const role = await instance.SENDER_ROLE();
    console.log(role);
    const tx = await instance.grantRole(role, senderAddress);
    const txReceipt = await tx.wait(1);
    console.log(`tx hash: ${txReceipt.hash}`)
    const has = await instance.hasRole(role, senderAddress)
    console.log("has role:", has)

}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })