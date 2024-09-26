const {ethers, network} = require("hardhat");

async function main() {
    /**
     * # dev
     * b2dev: yarn hardhat run scripts/message/set_weight.js --network b2dev
     * asdev: yarn hardhat run scripts/message/set_weight.js --network asdev
     * # pord
     * b2: yarn hardhat run scripts/message/set_weight.js --network b2
     * as: yarn hardhat run scripts/message/set_weight.js --network as
     */
    const [owner] = await ethers.getSigners()
    console.log("Owner Address: ", owner.address);
    let messageAddress;
    if (network.name == 'b2dev') {
        messageAddress = "0xe55c8D6D7Ed466f66D136f29434bDB6714d8E3a5";
    } else if (network.name == 'asdev') {
        messageAddress = "0x2A82058E46151E337Baba56620133FC39BD5B71F";
    } else if (network.name == 'b2') {
        messageAddress = "";
    } else if (network.name == 'as') {
        messageAddress = "";
    }
    console.log("Message Address: ", messageAddress);
    // bridge
    const bridge = await ethers.getContractFactory("B2MessageBridge");
    const instance = await bridge.attach(messageAddress);

    let _weight = 1;
    let chainId = 421614;

    let weight = await instance.weights(chainId);
    console.log("weight: ", weight);
    if (_weight != weight) {
        const tx = await instance.setWeight(chainId, _weight);
        const txReceipt = await tx.wait(1);
        console.log(`tx hash: ${txReceipt.hash}`)
        weight = await instance.weights(chainId);
        console.log("weight: ", weight);
    }
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })
