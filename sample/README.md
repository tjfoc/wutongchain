存证和转账的sample在当前目录下的sdk目录下。

1. store = `{
	"Data": "3ddd"
}`为存证请求的body,存证的具体内容为“3ddd”

2. privkey1 = "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tDQpNSUdUQWdFQU1CTUdCeXFHU000OUFnRUdDQ3FCSE05VkFZSXRCSGt3ZHdJQkFRUWdTYTh5cjlzMGphZ3NDZXIzDQp1R3dVeVpneEoybUNSRlRWejFQU3R3MjA4TzJnQ2dZSUtvRWN6MVVCZ2kyaFJBTkNBQVNmVStIZWxLSGV4TTdRDQpKd2c5NDZIWCt6OUpsTE9CVWlrZkRXOW9pVmNNWnJsRGZrMndwQ29pZ25IbGVyWEM3Z1BhbnNFemZuajhDdEc3DQpqM05mR0ZwRg0KLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLQ=="

privkey1为用户的私钥的base64编码。

3. 	body = `{
	"Prikey":"LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tDQpNSUdUQWdFQU1CTUdCeXFHU000OUFnRUdDQ3FCSE05VkFZSXRCSGt3ZHdJQkFRUWdTYTh5cjlzMGphZ3NDZXIzDQp1R3dVeVpneEoybUNSRlRWejFQU3R3MjA4TzJnQ2dZSUtvRWN6MVVCZ2kyaFJBTkNBQVNmVStIZWxLSGV4TTdRDQpKd2c5NDZIWCt6OUpsTE9CVWlrZkRXOW9pVmNNWnJsRGZrMndwQ29pZ25IbGVyWEM3Z1BhbnNFemZuajhDdEc3DQpqM05mR0ZwRg0KLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLQ==",
	"TokenType": "JD011155",
	"Amount": "20000",
	"Data": "发行token"
}`
body:token发行的请求body。
prikey:参考2，为用户的私钥的base64编码。
TokenType: token的类型，不可重复发行。
Amount:发行token的数值。
Data:本次发行token的描述。

4. tbody = `{
	"PriKey": "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tDQpNSUdUQWdFQU1CTUdCeXFHU000OUFnRUdDQ3FCSE05VkFZSXRCSGt3ZHdJQkFRUWdTYTh5cjlzMGphZ3NDZXIzDQp1R3dVeVpneEoybUNSRlRWejFQU3R3MjA4TzJnQ2dZSUtvRWN6MVVCZ2kyaFJBTkNBQVNmVStIZWxLSGV4TTdRDQpKd2c5NDZIWCt6OUpsTE9CVWlrZkRXOW9pVmNNWnJsRGZrMndwQ29pZ25IbGVyWEM3Z1BhbnNFemZuajhDdEc3DQpqM05mR0ZwRg0KLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLQ==",
	"Data":"转账100",
	"Token":[ {"TokenType":"JD011155","ToAddr":"36FDtMXa6bDkD1URuALKvyZ5Rczrdfi74xUwpmsNn9mxpq3TEC","Amount":"100"}]
	}`

tbody:转账请求的body
prikey:参考2
Data:转账描述

Token:转账的具体信息，可能同时向多账户转账，[ {"TokenType":"JD011155","ToAddr":"36FDtMXa6bDkD1URuALKvyZ5Rczrdfi74xUwpmsNn9mxpq3TEC","Amount":"100"}，{""}],放入[]中，用{}包含，用逗号隔开，参考如上。

token中的TokenType:转出的token类型
token中的ToAddr:转入的地址(通过http://sdk:port/utxo/genaddrss 生成)
token中的Amount:转出token的数值

5. func main() {
	//存证交易
	txstore()
	//发行token
	txissue()
	//转账
	transfer()
}

main函数中为3中交易，分别为:存证，token发行，token转账。
