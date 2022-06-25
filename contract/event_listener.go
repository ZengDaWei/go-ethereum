package contract

import (
	"context"
	"github.com/labstack/gommon/color"
	"log"
	"math/big"
	"sort"
	"strconv"
	"time"
)

type (
	SuccessCallBack func(*big.Int)
)

func LoopListenerEvent(syncInterval uint, reqBlockLimit uint, fromBlockNumber *uint, rpcEndpoint string, cb SuccessCallBack) {

	//时间间隔（秒为单位）
	duration := time.Second * time.Duration(syncInterval)

	ticker := time.NewTicker(duration)
	ctx := context.Background()

	client, err := GetClient(rpcEndpoint)
	if err != nil {
		panic(err)
	}
loop:
	for {
		select {
		case <-ticker.C:
			currentBlockNumber, err := client.BlockNumber(ctx)
			if err != nil {
				color.Yellow("连接断开,正在自动重试...")
				goto loop
			}

		handle:
			value := uint(currentBlockNumber) - *fromBlockNumber

			//有新区块未读，并达到了读取标准（即未读区块 > 数据库中设置的limit）
			if value >= reqBlockLimit {
				value = reqBlockLimit
			}

			from := *big.NewInt(int64(*fromBlockNumber) + 1)
			to := *big.NewInt(int64(*fromBlockNumber + value))

			if from.Int64() > int64(currentBlockNumber) {
				goto loop
			}

			log.Println("current: " + strconv.FormatUint(currentBlockNumber, 10) + ", from: " + from.Text(10) + ", to : " + to.Text(10))

			//获取事件
			eventLogs, err := GetEventLogs(&from, &to)
			if err != nil {
				HandleError(err)
				goto loop
			} else {
				//事件排序
				sortedLogs := logsSupportSort{logs: eventLogs}
				sort.Sort(&sortedLogs)
				//处理事件
				for _, eventLog := range sortedLogs.logs {
					err := HandleEvent(eventLog)
					if err != nil {
						HandleError(err)
						goto loop
					}
				}

				cb(&to)

				/**
				有新区块未读，并达到了读取标准（即未读区块 > 数据库中设置的limit）
				这里再次这样判断的原因是：如果未处理的区块过多，可以分批一次性处理完，不用每次处理之后都要再等待
				*/
				if (uint(currentBlockNumber) - reqBlockLimit) >= *fromBlockNumber {
					goto handle
				}
			}
		}
	}

}

func HandleError(err error) {
	color.Red("出现异常")
	panic(err)
}
