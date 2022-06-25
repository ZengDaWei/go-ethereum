# Go ethereum

基于 [go-ethereum](github.com/ethereum/go-ethereum) 实现的监听 evm 链工具包

## 使用

安装
```shell
$ go get github.com/ZengDaWei/go-ethereum
```

案例：

```go
package main

import (
	"fmt"
	"math/big"

	"github.com/ZengDaWei/go-ethereum/contract"
	"github.com/ZengDaWei/go-ethereum/route"
	"github.com/ethereum/go-ethereum/core/types"
)

func main() {
	route.Add("0x2d184680AAA47bcAF747E91070d01D56CB4982d5", "Test(uint256)", Hello)

	var from uint = 0

	contract.Run(3, 10, &from, "http://127.0.0.1:8545", func(i *big.Int) {
		from += 3
	}, func(err error) {
		panic(err)
	})
}

func Hello(log types.Log) error {
	fmt.Println("test")
	return nil
}

```

- 可自定义轮训时间间隔
- 可自定义获取区块数量
- 可自定义事件处理回调
- 可自定义事件处理异常回调

## 注册事件

```go
package main

import (
	"fmt"
	"github.com/ZengDaWei/go-ethereum/route"
	"github.com/ethereum/go-ethereum/core/types"
)

func main() {
  // 合约地址，事件签名，事件对应函数
	route.Add("0x2d184680AAA47bcAF747E91070d01D56CB4982d5", "Test(uint256)", Hello)

}

func Hello(log types.Log) error {
	fmt.Println("test")
	return nil
}
```

## 运行监听事件

```go
package main

import (
	"fmt"
	"math/big"

	"github.com/ZengDaWei/go-ethereum/contract"
	"github.com/ZengDaWei/go-ethereum/route"
	"github.com/ethereum/go-ethereum/core/types"
)

func main() {
	route.Add("0x2d184680AAA47bcAF747E91070d01D56CB4982d5", "Test(uint256)", TestEventHandler)

	var syncInterval uint = 3
	var blockLimit uint = 10
	var rpcEndpoint = "http://127.0.0.1:8545"
	var fromBlockNumber uint = 0

	contract.Run(syncInterval, blockLimit, &fromBlockNumber, rpcEndpoint, func(i *big.Int) {
		fromBlockNumber += 3
	}, func(err error) {
		panic(err)
	})
}

func TestEventHandler(log types.Log) error {
	fmt.Println("test")
	return nil
}

```