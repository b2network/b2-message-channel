const {ethers, upgrades, network} = require("hardhat");

let messageAddress;
// let businessAddress;
// let to_chain_id;
// let data;

async function main() {
    /**
     * b2dev: yarn hardhat run scripts/message/call.js --network b2dev
     * as: yarn hardhat run scripts/message/call.js --network as
     */

    const [owner] = await ethers.getSigners()
    console.log("Owner Address:", owner.address);

    if (network.name == 'b2dev') {
        messageAddress = "0x5c2646996eEe3ECf865BEfA2De24e5BbE1C552Ba";
        // businessAddress = "0x8Ac2C830532d7203a12C4C32C0BE4d3d15917534";
        // to_chain_id = 421614;
        // data = '0x1234567890';
    } else if (network.name == 'as') {
        messageAddress = "0x2A82058E46151E337Baba56620133FC39BD5B71F";
        // businessAddress = "0x91171cf194a4B66Bd459Ada038397c7e890FB9D4";
        // to_chain_id = 1123;
        // data = '0x1234';
    }

    const B2MessageBridge = await ethers.getContractFactory("B2MessageBridge");
    const instance = await B2MessageBridge.attach(messageAddress);

    // let tx1 = await instance.setWeight(1);
    // await tx1.wait(1);

    let tx = await instance.setWeight(1);
    await tx.wait(1);
    console.log(await instance.weight());
    let from_chain_id = "0";
    let from_id = "68130661828825923178236639338972550787299552799811069952048480569636815930548";
    let from_sender = "0x0000000000000000000000000000000000000000";
    let to_chain_id = "1123";
    let contract_address = "0xAaF7f27EA526B29cF4dA3b957Ed0C388070FcCE5";
    let data = "0x96a0968b0f115755b4f63794cef60c5ce35a8cc9633f2ecbdb4b30ea3f0084b40000000000000000000000000000000000000000000000000000000000000080000000000000000000000000de39679d030c63ae9c5f39d1ba500425e4d6ebe50000000000000000000000000000000000000000000000000000000000989680000000000000000000000000000000000000000000000000000000000000002a746231717a3464337577326330326861736d63733672796c653761326a617678666a6e6a667a7064707a00000000000000000000000000000000000000000000";

    // 0x9cc4669bb997c40579f89E08980B99218abaE3FE
    // 0xfd1d24ee09b1263bcbd62badbd7bb9c295a05eea85126274f36abf20c9b57499
    // 0x940d61022ba8c7cf9784b6e58fb9da6e143e630074e356c2a236b4280b7d688b
    // 0x0a048b9daf2b41df9ee6ddaa874c2099636a720f1a539670d8420da9dca410e0
    let hash = await instance.SendHash(from_chain_id, from_id, from_sender, to_chain_id, contract_address, data);
    console.log("hash:", hash);

    let signature = "0xa3b185d9dc2a27a4b6af23c03bc4e4a2887a6e769529aa9b967c37d0414a12523b214d3d73dae05acf6122967ef2a00263ccf4ad5a0710470c4639976b3178571b";
    let v = await instance.verify(from_chain_id, from_id,
        from_sender, to_chain_id, contract_address, data, signature);
    console.log("verify:", v)

    // let callTx = await instance.call(to_chain_id, businessAddress, data);
    // const callTxReceipt = await callTx.wait(1);
    // console.log("callTxReceipt:", callTxReceipt.hash);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })