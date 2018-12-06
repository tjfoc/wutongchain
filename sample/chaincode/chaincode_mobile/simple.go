/*
=========示例背景=========
手机供应商节点保存手机货源信息，手机卖家节点同时会同步到数据。
整个手机的售卖和进货流程在整个联盟链内保存。卖家和供应商随时可
查。
*/
//本示例模拟了简单的手机商店中的手机管理
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/tjfoc/tjfoc/core/chaincode/shim"
	pb "github.com/tjfoc/tjfoc/protos/chaincode"
)

//合约方法处理器
type handler func(stub shim.ChaincodeStubInterface, args []string) pb.Response

//合约结构体
type SmartContract struct {
	handlerMap map[string]handler
}

func newSmartContract() *SmartContract {
	cc := &SmartContract{}
	cc.handlerMap = map[string]handler{
		"initMobile":        cc.initMobile,
		"createMobile":      cc.createMobile,
		"getAllMobile":      cc.getAllMobile,
		"queryMobile":       cc.queryMobile,
		"changeMobileCount": cc.changeMobileCount,
		"deleteMobile":      cc.deleteMobile,
	}
	return cc
}

//手机结构体
type Mobile struct {
	Brand string  `json:"brand"`
	Model string  `json:"model"`
	Price float32 `json:"price"`
	Color string  `json:"color"`
	Count int     `json:"count"`
}

//初始化手机信息,将手机的信息添加到链上
//由供应商操作这个方法，用于系统开始运行时，初始化库存手机信息
func (c *SmartContract) initMobile(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	mobiles := []Mobile{
		Mobile{Brand: "HUAWEI", Model: "Mate 20", Price: 5000.00, Color: "black", Count: 100},
		Mobile{Brand: "Apple", Model: "iphoneXS", Price: 10000.00, Color: "red", Count: 100},
		Mobile{Brand: "OPPO", Model: "R11", Price: 3000.00, Color: "bule", Count: 100},
		Mobile{Brand: "VIVO", Model: "X6", Price: 2000.00, Color: "white", Count: 100},
		Mobile{Brand: "samsung", Model: "Galaxy A9s", Price: 7000.00, Color: "gray", Count: 100},
	}
	i := 0
	for i < len(mobiles) {
		fmt.Println("i is ", i)
		mobileAsBytes, _ := json.Marshal(mobiles[i])
		stub.PutState("Mobile"+strconv.Itoa(i), mobileAsBytes)
		fmt.Println("Added", mobiles[i])
		i = i + 1
	}
	return shim.Success(nil)
}

//创建新手机，添加新的手机信息到链上
//由供应商操作这个方法，用于手机进货，共享手机信息
func (c *SmartContract) createMobile(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	price, _ := strconv.ParseFloat(args[2], 32)
	count, _ := strconv.Atoi(args[4])
	mobile := Mobile{Brand: args[0], Model: args[1], Price: float32(price), Color: args[3], Count: count}

	mobileAsBytes, _ := json.Marshal(mobile)
	stub.PutState(args[5], mobileAsBytes)

	return shim.Success(nil)
}

//查询手机信息，根据手机的唯一ID--Mobile(?)查询指定手机的信息
//由商家和卖家操作，用于查询手机库存量
func (c *SmartContract) queryMobile(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	mobileAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(fmt.Sprintf("Invalid key:[%s]", args[0]))
	}
	return shim.Success(mobileAsBytes)
}

//获取所有手机,根据手机的唯一ID--Mobile(?)的前缀，查询所有手机手机的信息
//由商家和卖家操作，用于查询手机库存量
func (c *SmartContract) getAllMobile(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	mobiles, _ := stub.GetStateByPrefix("Mobile")
	mobilesBytes, _ := json.Marshal(mobiles)
	return shim.Success(mobilesBytes)
}

//修改手机数量,根据手机的唯一ID--Mobile(?)，修改指定手机的数量
//由商家操作，用于进货时增加手机数量，同时卖家节点上也会增加手机数量
//由卖家操作，用于售卖手机之后数量更新，同时商家根据卖家售卖信息发货
func (c *SmartContract) changeMobileCount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	count, _ := strconv.Atoi(args[0])
	mobileAsBytes, _ := stub.GetState(args[1])
	mobile := Mobile{}
	json.Unmarshal(mobileAsBytes, &mobile)
	mobile.Count = count

	newMobileAsBytes, _ := json.Marshal(mobile)
	stub.PutState(args[1], newMobileAsBytes)
	fmt.Println("Update", args[1])

	return shim.Success(nil)
}

//删除手机，根据手机的唯一ID--Mobile(?)，删除指定手机的信息
//注：只是移除worldsate中的手机信息，链上依然保存
//由商家操作，用于某手机下架或者某些原因不再售卖，同时卖家节点上也会失去这个手机的信息
func (c *SmartContract) deleteMobile(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	stub.DelState(args[0])
	return shim.Success(nil)
}

func (c *SmartContract) getHandler(method string) (handler, bool) {
	h, ok := c.handlerMap[method]
	return h, ok
}

//Init 实现智能合约接口方法
func (c *SmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success([]byte("init successful"))
}

//Invoke 实现智能合约接口方法
func (c *SmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	method, args := stub.GetFunctionAndParameters()
	if h, ok := c.getHandler(method); ok {
		return h(stub, args)
	}
	return shim.Error("Invalid method")
}

//程序入口
func main() {
	//注册自定义合约
	err := shim.Start(newSmartContract())
	if err != nil {
		log.Printf("Error starting my SmartContract : %s", err)
	}
}
