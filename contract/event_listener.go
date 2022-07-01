package contract

import (
	"context"
	"github.com/labstack/gommon/color"
	"math/big"
	"sort"
	"time"
)

type (
	BeginCallBack   func(ctx context.Context) context.Context
	SuccessCallBack func(ctx context.Context, value *big.Int) context.Context
	ErrorCallBack   func(ctx context.Context, err error) context.Context
)

func Run(ctx context.Context, syncInterval uint, reqBlockLimit uint, fromBlockNumber *uint, rpcEndpoint string, beginCb BeginCallBack, cb SuccessCallBack, erCb ErrorCallBack) {

	//时间间隔（秒为单位）
	duration := time.Second * time.Duration(syncInterval)

	ticker := time.NewTicker(duration)

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
				color.Yellow("Disconnected, retrying automatically...")
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

			//获取事件
			eventLogs, err := GetEventLogs(&from, &to)
			if err != nil {
				ctx = erCb(ctx, err)
				goto loop
			} else {
				//事件排序
				sortedLogs := logsSupportSort{logs: eventLogs}
				sort.Sort(&sortedLogs)
				//处理事件
				if sortedLogs.Len() > 0 {
					ctx = beginCb(ctx)
					for _, eventLog := range sortedLogs.logs {
						err := HandleEvent(eventLog)
						if err != nil {
							ctx = erCb(ctx, err)
							goto loop
						}
					}
				}

				ctx = cb(ctx, &to)

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
