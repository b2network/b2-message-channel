const {ethers, upgrades, network} = require("hardhat");

let messageAddress;
let businessAddress;
let to_chain_id;

async function main() {
    /**
     * b2dev: yarn hardhat run scripts/message/call.js --network b2dev
     * as: yarn hardhat run scripts/message/call.js --network as
     * b2: yarn hardhat run scripts/message/call.js --network b2
     */

    const [owner] = await ethers.getSigners()
    console.log("Owner Address:", owner.address);

    if (network.name == 'b2dev') {
        messageAddress = "0xc7441Ac47596D1356fcc70062dA0462FcA98E14e";
        businessAddress = "0x8Ac2C830532d7203a12C4C32C0BE4d3d15917534";
        to_chain_id = 421614;
    } else if (network.name == 'as') {
        messageAddress = "0x2A82058E46151E337Baba56620133FC39BD5B71F";
        businessAddress = "0x91171cf194a4B66Bd459Ada038397c7e890FB9D4";
        to_chain_id = 1123;
    } else if (network.name == 'b2') {
        messageAddress = "";
        businessAddress = "";
    }

    const B2MessageBridge = await ethers.getContractFactory("B2MessageBridge");
    const instance = await B2MessageBridge.attach(messageAddress);

    let contract_address = businessAddress;
    let data = '0x1234';
    let callTx = await instance.call(to_chain_id, contract_address, data);
    const callTxReceipt = await callTx.wait(1);
    console.log("callTxReceipt:", callTxReceipt.hash);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })