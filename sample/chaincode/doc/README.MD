# 梧桐链智能合约开发及使用

## 目录
### [梧桐链智能合约开发--shim.ChaincodeStubInterface用法](#梧桐链智能合约开发--shim.ChaincodeStubInterface用法)
* [开发介绍](#1.开发介绍)
* [智能合约模板](#2.智能合约模板)
* [shim.ChaincodeStubInterface接口详解](#3.shim.ChaincodeStubInterface接口详解)
    * [获得调用的参数](#3.1获得调用的参数)
    * [对worldstate数据进行增删改查](#3.2对worldstate数据进行增删改查)
    * [合约调合约](#3.3合约调合约)
### [梧桐链智能合约使用--通过sdk使合约生效](#梧桐链智能合约使用--通过sdk使合约生效)
* [使用介绍](#1.使用介绍)
* [合约安装](#2.合约安装)
* [合约调用](#3.合约调用)
* [合约销毁](#4.合约销毁)

## 梧桐链智能合约开发--shim.ChaincodeStubInterface用法

## 1.开发介绍
梧桐链的应用开发一共分为两个部分，智能合约开发和基于SDK的application开发。目前梧桐链智能合约暂时只支持go语言的版本，所以本篇文档主要介绍go版本的智能合约开发。

## 2.智能合约模板
智能合约的go代码需要定义一个智能合约结构体，然后在这个结构体上定义invoke和init两个函数。init函数是智能合约第一次启动的时候调用的，会且只会执行一次。invoke是sdk掉用合约时的入口，所有的调用操作都需要通过invoke函数。最后还需要定义一个main函数，作为合约启动的入口。以下是一个智能合约模板：
```go
package main

import (
	"fmt"
	"github.com/tjfoc/tjfoc/core/chaincode/shim"
	pb "github.com/tjfoc/tjfoc/protos/chaincode"
)

//合约结构体
type SmartContract struct {
}

//Init 实现智能合约接口方法
func (c *SmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success([]byte("init successful"))
}

//Invoke 实现智能合约接口方法
func (c *SmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	method, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)
   if function == "test1" {//自定义函数名称
      return c.test1(stub, args)//定义调用的函数
   }
   return shim.Error("Invalid method")
}

func (c *SmartContract) test1(stub shim.ChaincodeStubInterface, args []string) pb.Response{
   return shim.Success([]byte("Called test1"))
}

//程序入口
func main() {
	//注册自定义合约
	err := shim.Start(new(SmartContract)
	if err != nil {
		log.Printf("Error starting my smart contract : %s", err)
	}
}
```
## 3.shim.ChaincodeStubInterface接口详解
通过以上示例不难发现，在invoke和init的时候都会传入shim.ChaincodeStubInterface类型的参数，该参数提供了很多方便我们实现智能合约业务逻辑的方法。
### 3.1获得调用的参数
在模板的invoke中，我们可以看到是通过传入的参数来决定具体调用哪个方法的。我们在示例中使用了GetFunctionAndParameters来获取参数，除了该方法，还提供了其他的几个方法，大同小异，可以根据自己的开发需求和喜好进行调用：
* GetArgs() [][]byte 数组第一个字是Function，剩下的都是Parameter
* GetStringArgs() []string 数组第一个字是Function，剩下的都是Parameter
* GetFunctionAndParameters() (string, []string) 将数组分成两部分返回，第一部分是Function，第二部分是Parameter
### 3.2对worldstate数据进行增删改查
整个智能合约部分最核心的就是对worldstate数据进行增删改查，梧桐链上提供了以下结构来实现该操作：
#### 增加数据PutState(key string, value []byte) error
因为梧桐链的worldstate来说，修改和增加其实是同一个操作。因为worldstate数据库采用的是levelDB这样一个key-value数据库，如果我们指定的key在数据库中已经存在，就是修改操作，如果不存在，就是添加操作。对于实际的系统来说，我们的key可能是根据病例号，或者系统自动分配的自增实体+ID,而value则是一个对象进行JSON序列化之后的byte数组。举个例子，我们定义一个Resident的Struct,然后插入一个居民数据：
```go
type Resident struct{
    Id string
    Name string
    Sex string
}

func (c *SmartContract) testPutState(stub shim.ChaincodeStubInterface, args []string){
    resident1:=Resident{"111111111111111111","Bob","male"}
    residentJsonBytes,err:=json.Marshal(resident1)//Json序列化
    if err!=nil{
        return shim.Error(err.Error())
    }
    err = stub.PutState("Resident1",residentJsonBytes)
    if err!=nil{
        return shim.Error(err.Error())
    }
    return shim.Success([]byte{"Add a Resident !"})
}
```
#### 删除数据DelState(key string) error
这个方法是根据key去删除对应的数据，如果没有对应的key,删除失败：
```go
func (c *SmartContract) testDelState(stub shim.ChaincodeStubInterface, args []string){
err = stub.DelState("Resident1")
if err!=nil{
    return shim.Error(err.Error())
}
return shim.Success([]byte{"Delete a Resident !"})
}
```
#### 批量删除数据DelStaten(keyn []string) error
这个是一个可以根据多个key去删除数据的方法：
```go
func (c *SmartContract) testDelStaten(stub shim.ChaincodeStubInterface, args []string){
    err = stub.DelStaten([]string{"Resident1","Resident1"})
    if err!=nil{
        return shim.Error(err.Error())
    }
    return shim.Success([]byte{"Delete Residents !"})
}
```
#### 查询数据GetState(key string) ([]byte, error)
查询数据需要注意的是我们需要将得到的数据从json转换回来，比如我们想获取一个人的名字：
```go
func (c *SmartContract) testGetState(stub shim.ChaincodeStubInterface, args []string){
    residentJsonBytes,err := stub.GetState("Resident1")
    if err!=nil{
        return shim.Error(err.Error())
    }
    var resident Resident
    err = json.Unmarshal(residentJsonBytes,&resident)//反序列化
    fmt.Println("Read Resident from DB, name:"+resident.Name)
    return shim.Success(nil)
}
```
#### 批量查询数据GetStaten(keyn []string) (map[string]string, error)
批量查询数据是根据多个key去获取想要的数据。
#### 根据前缀查询数据GetStateByPrefix(key string) (map[string]string, error)
我们还可以根据前缀去查询你想要的数据，一般就是将数据分好类别，比如将男居民和女居民分开：男居民使用Male_Resident作为key存储，女居民使用Female_Resident作为key存储，则获取所有男居民可以这样写：
```go
 stub.GetStateByPrefix("Male")
```
获取所有女居民：
```go
 stub.GetStateByPrefix("Female")
```
### 3.3合约调合约
#### InvokeChaincode(chaincodeName string, chaincodeVersion string, args [][]byte) pb.Response
根据另一个合约的合约名和合约版本可以直接将参数从这个合约传入另一个合约，从而实现调用
## 梧桐链智能合约使用--通过sdk使合约生效
## 1.使用介绍
合约开发成功之后，下一步就是如何使用这个合约。在梧桐链中，对合约的操作都是通过向节点发送交易实现的。为了方便用户操作，梧桐链提供了的sdk服务中就有很方便的对合约进行安装、调用和销毁的接口。  
*注：如果要升级合约的话，就进行一个新的合约安装操作，再进行一个旧合约的销毁操作*
## 2.合约安装
先需要启动一个sdk服务，将合约的安装命令封装成一个post请求发送到你的sdk服务器，url为：
```
http://sdk ip/contract/install
```
请求内容是一个JSON字符串，为：
```
{
    "Name": "合约名",
    "Version": "合约版本号",
    "File": "base64('合约代码文件内容')"
}
```
发送成功之后将会收到：
```
{
    "State": 200,
    "Data": "安装合约的交易的hash",
    "Message": "success"
}
```
## 3.合约调用
合约安装成功之后，就可以对合约进行调用了，调用方法还是封装一个post请求到你的sdk服务器，url为：
```
http://sdk ip/createnewtransaction
```
请求内容是一个JSON字符串，为：
```
{
    "Name": "合约名",
    "Version": "合约版本",
    "Method": "调用合约中的方法名",
    "Args":  ["参数1", "参数2", "参数3", "参数4"]
}
```
发送成功之后将会收到
```
{
    "State": 200,
    "Data": "调用合约的交易的hash",
    "Message": "success"
}
```
## 4.合约销毁
如果想删除一个合约，可以向你的sdk服务器发送一个删除合约的请求，url为：
```
http://sdk ip/contract/destroy
```
请求内容是一个JSON字符串，为：
```
{
    "Name": "合约名",
    "Version": "合约版本"
}
```
发送成功之后将会收到
```
{
    "State": 200,
    "Data": "删除合约的交易的hash",
    "Message": "success"
}
```
