package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	txId := "fee63a7c0cf184824ba4fb5ea0be625ac673ab0e97ecf5e386b6805ac6ab869e"
	//fromAddress := "bc1q2uduc360s3s3jww4nwnzn3yhad7ry0ldjkfr79"
	//toAddress := "0xeF9B26046a2392C956320200eE0818543aA96aB7"
	//amount := decimal.New(1212, 0)
	//data := message.EncodeSendData(txId, fromAddress, toAddress, amount)
	//fmt.Println(hexutil.Encode(data))

	//num := big.NewInt(0)
	//num.SetBytes(hexutil.MustDecode(txId))
	//fmt.Println(num.Text(16))
	fmt.Println(common.HexToHash(txId))

	// 0xfee63a7c0cf184824ba4fb5ea0be625ac673ab0e97ecf5e386b6805ac6ab869e0000000000000000000000000000000000000000000000000000000000000080000000000000000000000000ef9b26046a2392c956320200ee0818543aa96ab700000000000000000000000000000000000000000000000000000000000004bc000000000000000000000000000000000000000000000000000000000000002a626331713275647563333630733373336a7777346e776e7a6e3379686164377279306c646a6b6672373900000000000000000000000000000000000000000000
	// 0xfee63a7c0cf184824ba4fb5ea0be625ac673ab0e97ecf5e386b6805ac6ab869e0000000000000000000000000000000000000000000000000000000000000080000000000000000000000000ef9b26046a2392c956320200ee0818543aa96ab700000000000000000000000000000000000000000000000000000000000004bc000000000000000000000000000000000000000000000000000000000000002a626331713275647563333630733373336a7777346e776e7a6e3379686164377279306c646a6b6672373900000000000000000000000000000000000000000000
	//bclient, err := rpcclient.New(&rpcclient.ConnConfig{
	//	Host:         "129.226.198.246:8332",
	//	Host: "indulgent-floral-bird.btc.quiknode.pro/1df147114bdfd6f119d1a5eef706c91b3710ca6e",
	//	//Host:         "divine-patient-meadow.btc-testnet.quiknode.pro/36ade4989e57ddd507d9790dec479606cbe8c0c6",
	//	User:         "bnan",
	//	Pass:         "bnan2021",
	//	HTTPPostMode: true,  // Bitcoin core only supports HTTP POST mode
	//	DisableTLS:   false, // Bitcoin core does not provide TLS by default
	//}, nil)
	//if err != nil {
	//	fmt.Println("new: ", err.Error())
	//	return
	//}
	//ChainInfo, err := bclient.GetBlockCount()
	//if err != nil {
	//	fmt.Println("get: ", err.Error())
	//	return
	//}
	//fmt.Println("ChainInfo: ", ChainInfo)
}
