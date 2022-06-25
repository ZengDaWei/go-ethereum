package contract

import (
	"github.com/barkimedes/go-deepcopy"
	"github.com/ethereum/go-ethereum/core/types"
	"sort"
)

type logsSupportSort struct {
	logs []types.Log
}

func (P *logsSupportSort) Len() int {
	return len(P.logs)
}
func (P *logsSupportSort) Less(i, j int) bool {
	if P.logs[i].BlockNumber < P.logs[j].BlockNumber || (P.logs[i].BlockNumber == P.logs[j].BlockNumber && P.logs[i].Index < P.logs[j].Index) {
		return true
	}
	return false
}
func (P *logsSupportSort) Swap(i, j int) {
	tmpI, _ := deepcopy.Anything(P.logs[i])
	tmpJ, _ := deepcopy.Anything(P.logs[j])
	P.logs[i] = tmpJ.(types.Log)
	P.logs[j] = tmpI.(types.Log)
}

func SortLogs(logsSlice []types.Log) {
	sortStruct := logsSupportSort{
		logs: logsSlice,
	}
	sort.Sort(&sortStruct)
}
