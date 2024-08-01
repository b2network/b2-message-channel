const {ethers, upgrades, network} = require("hardhat");

const messageAddress = "0x2A82058E46151E337Baba56620133FC39BD5B71F";

async function main() {
    /**
     * b2dev: yarn hardhat run scripts/message/send.js --network b2dev
     * b2: yarn hardhat run scripts/message/send.js --network b2
     * as: yarn hardhat run scripts/message/send.js --network as
     */

    const [owner] = await ethers.getSigners()
    console.log("Owner Address:", owner.address);

    const B2MessageBridge = await ethers.getContractFactory("B2MessageBridge");
    const instance = await B2MessageBridge.attach(messageAddress);

    let from_chain_id = 1123;
    let from_id = 7;
    let from_sender = "0x9cc4669bb997c40579f89E08980B99218abaE3FE";
    let contract_address = '0x8Ac2C830532d7203a12C4C32C0BE4d3d15917534';
    let data = '0x1234';
    let signatures = ['0x74ab5ab0dccdb9ea403a181d4272102dc4ccd6f4fae4d2d3386099094d9a2c1315c183629af6237dc6499f02aa7b2dbed773f3bc9848ea090e6ae957fbee25031c'];
    let sendTx = await instance.send(from_chain_id, from_id, from_sender, contract_address, data, signatures);
    const sendTxReceipt = await sendTx.wait(1);
    console.log("sendTxReceipt:", sendTxReceipt.hash);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })