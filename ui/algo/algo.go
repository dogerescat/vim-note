package algo

import (
	"sort"

	"github.com/junegunn/fzf/src/algo"
	"github.com/junegunn/fzf/src/util"
)

type stringResult struct {
	name   string
	result algo.Result
}

func MatchString(str string, list []string) []string {
	pattern := []rune(str)
	var data []stringResult
	for _, s := range list {
		input := util.RunesToChars([]rune(s))
		slab := util.MakeSlab(100*1024, 2048)
		result, _ := algo.FuzzyMatchV2(false, true, true, &input, pattern, false, slab)
		data = append(data, stringResult{name: s, result: result})
	}
	sort.SliceStable(data, func(i, j int) bool {
		if data[i].result.Score > data[j].result.Score {
			return true
		} else if data[i].result.Score == data[j].result.Score {
			if data[i].result.Start > data[j].result.Start {
				return true
			} else if data[i].result.Start == data[j].result.Start {
				if data[i].result.End > data[j].result.End {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		} else {
			return false
		}
	})
	var res []string
	for _, sr := range data {
		res = append(res, sr.name)
	}
	return res
}
