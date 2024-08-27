const {ethers, network} = require("hardhat");

let messageAddress;

async function main() {
    /**
     * b2dev: yarn hardhat run scripts/message/grant_role.js --network b2dev
     * as: yarn hardhat run scripts/message/grant_role.js --network as
     */

    const [owner] = await ethers.getSigners()
    console.log(owner.address)
    // return

    if (network.name == 'b2dev') {
        messageAddress = "0x5c2646996eEe3ECf865BEfA2De24e5BbE1C552Ba";
    } else if (network.name == 'as') {
        messageAddress = "0x2A82058E46151E337Baba56620133FC39BD5B71F";
    }

    // bridge
    const bridge = await ethers.getContractFactory("B2MessageBridge");
    const instance = await bridge.attach(messageAddress);

    let grantRole = "0x7e108D9Da2d079cEe26177159c47a2f71EF465B6";
    // const grantRoleTx = await instance.grantRole(await instance.VALIDATOR_ROLE(), grantRole);
    // const res = await grantRoleTx.wait();
    // console.log("status:", res.status);
    //
    hasRole = await instance.hasRole(await instance.VALIDATOR_ROLE(), grantRole);
    console.log("hasRole:", hasRole)


    // const revokeRoleTx = await instance.revokeRole(await instance.ADMIN_ROLE(), grantRole);
    // const res = await revokeRoleTx.wait();
    // console.log("status:", res.status);

    // let hasRole = await instance.hasRole(await instance.UPGRADE_ROLE(), grantRole);
    // console.log("hasRole:", hasRole)
    // if (hasRole) {
    //     const grantRoleTx = await instance.renounceRole(await instance.UPGRADE_ROLE(), grantRole);
    //     const res = await grantRoleTx.wait();
    //     console.log("status:", res.status);
    // }

}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })