const {ethers, upgrades, network} = require("hardhat");

const messageAddress = "0x548B00A60316531cE781411a9108AA0B6300424A";

async function main() {
    /**
     * b2dev: yarn hardhat run scripts/message/send.js --network b2dev
     * b2: yarn hardhat run scripts/message/send.js --network b2
     */

    const [owner] = await ethers.getSigners()
    console.log("Owner Address:", owner.address);

    const B2MessageBridge = await ethers.getContractFactory("B2MessageBridge");
    const instance = await B2MessageBridge.attach(messageAddress);

    let from_chain_id = 1123;
    let from_id = 4;
    let from_sender =  owner.address;
    let contract_address = '0x67292d3b05C9fc8391372C23826236983C23b9a2';
    let data = '0x1234';
    let signatures = [];
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