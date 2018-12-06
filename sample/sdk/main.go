/*
1.基础链存证交易发起示例
2.基础链utxo交易的，发行token和转账
*/

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (

	sdkUrl = "http://your ip:8888/" //需要更换成实际的物理机ip
	//存证交易的请求body
	store = `{
	"Data": "3ddd"
}`
	//用户1的私钥的base64编码
	privkey1 = "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tDQpNSUdUQWdFQU1CTUdCeXFHU000OUFnRUdDQ3FCSE05VkFZSXRCSGt3ZHdJQkFRUWdTYTh5cjlzMGphZ3NDZXIzDQp1R3dVeVpneEoybUNSRlRWejFQU3R3MjA4TzJnQ2dZSUtvRWN6MVVCZ2kyaFJBTkNBQVNmVStIZWxLSGV4TTdRDQpKd2c5NDZIWCt6OUpsTE9CVWlrZkRXOW9pVmNNWnJsRGZrMndwQ29pZ25IbGVyWEM3Z1BhbnNFemZuajhDdEc3DQpqM05mR0ZwRg0KLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLQ=="
	//用户2的私钥的base64编码
	privkey2 = "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tDQpNSUdUQWdFQU1CTUdCeXFHU000OUFnRUdDQ3FCSE05VkFZSXRCSGt3ZHdJQkFRUWdVOG01MVYzRHRwSDNtTTdVcHhCS1E2eTNCVWZEM2IxOEJlNWsrVGxncjlPZ0NnWUlLb0VjejFVQmdpMmhSQU5DQUFRSDdzdTNuMlhhb1Zjbg0KSHU0dE9WR05GdkVkYmFTNUZFczljNC9JVEJuV0ZneHBpaDNOVmJTejBjc01RUlQ4RTNnOHN3cjc1ZnNxek5uWQ0KdVF3cUNqb1INCi0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0="
	//token发行的请求body
	body = `{
	"Prikey":"LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tDQpNSUdUQWdFQU1CTUdCeXFHU000OUFnRUdDQ3FCSE05VkFZSXRCSGt3ZHdJQkFRUWdTYTh5cjlzMGphZ3NDZXIzDQp1R3dVeVpneEoybUNSRlRWejFQU3R3MjA4TzJnQ2dZSUtvRWN6MVVCZ2kyaFJBTkNBQVNmVStIZWxLSGV4TTdRDQpKd2c5NDZIWCt6OUpsTE9CVWlrZkRXOW9pVmNNWnJsRGZrMndwQ29pZ25IbGVyWEM3Z1BhbnNFemZuajhDdEc3DQpqM05mR0ZwRg0KLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLQ==",
	"TokenType": "JD011155",
	"Amount": "20000",
	"Data": "发行token"
}`
	//转账的请求body
	tbody = `{
	"PriKey": "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tDQpNSUdUQWdFQU1CTUdCeXFHU000OUFnRUdDQ3FCSE05VkFZSXRCSGt3ZHdJQkFRUWdTYTh5cjlzMGphZ3NDZXIzDQp1R3dVeVpneEoybUNSRlRWejFQU3R3MjA4TzJnQ2dZSUtvRWN6MVVCZ2kyaFJBTkNBQVNmVStIZWxLSGV4TTdRDQpKd2c5NDZIWCt6OUpsTE9CVWlrZkRXOW9pVmNNWnJsRGZrMndwQ29pZ25IbGVyWEM3Z1BhbnNFemZuajhDdEc3DQpqM05mR0ZwRg0KLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLQ==",
	"Data":"转账100",
	"Token":[ {"TokenType":"JD011155","ToAddr":"36FDtMXa6bDkD1URuALKvyZ5Rczrdfi74xUwpmsNn9mxpq3TEC","Amount":"100"}]
	}`
)

//发送psot请求
func post(s, body, log string) {
	client := &http.Client{}
	res, err := http.NewRequest("POST", sdkUrl+s, strings.NewReader(body))
	if err != nil {
		fmt.Println(err)
	}

	res.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(res)
	if err != nil {
		fmt.Println(err)
	}
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(log, string(b))
}

//get请求
func get(s, log string) {
	client := &http.Client{}
	resp, err := client.Get(sdkUrl + s)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(log, string(b))
}

func txstore() {
	post("store", store, "存证:")
}

func txissue() {
	post("utxo/issuetoken", body, "发行token")
	//为了确保发行token已上链,此处休眠2秒
	time.Sleep(time.Second * 2)
	get("utxo/balance?key="+privkey1, "查询privkey1余额:")
}

func transfer() {
	post("utxo/transfer", tbody, "转账:")
	time.Sleep(time.Second * 6)
	get("utxo/balance?key="+privkey2, "查询privkey2余额:")
}

func main() {
	//存证交易
	txstore()
	//发行token
	txissue()
	//转账
	transfer()
}

//运行结果如下:
/*
存证: {"State":200,"Message":"success","Data":{"Figure":"XzRfF5rKgfHCpxFsv1ncMygD2HUHcehQpw/rsp6qAJo=","OK":true}}
发行token {"State":400,"Message":"Token type JD011155 has been issued","Data":null}
查询privkey1余额: {"State":200,"Message":"success","Data":{"Detail":{"JD011155":"19700"},"Total":"19700"}}
转账: {"State":200,"Message":"success","Data":"xuQCsWmuK9U7810xivDGbAj3BKxRZ13uvztuLWZABks="}
查询privkey2余额: {"State":200,"Message":"success","Data":{"Detail":{"JD011155":"300"},"Total":"300"}}
*/

/*
其中发行token的返回码为400，因为以上结果为执行了3次该程序，由于相同token不能重复发行，所以结果返回异常。
*/
