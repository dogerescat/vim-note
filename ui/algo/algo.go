package algo

import (
	"sort"
)

type pair struct {
	a int
	b string
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func MatchString(str string, list []string) []string {
	siz := len(list)
	var data []pair
	for k := 0; k < siz; k++ {
		s := list[k]
		n := len(str)
		m := len(s)
		var dp [100][100]int
		dp[0][0] = 0
		for i := 1; i <= n; i++ {
			for j := 1; j <= m; j++ {
				dp[i][j] = 0
				if str[i-1] == s[j-1] {
					dp[i][j] = max(dp[i][j], dp[i-1][j-1]+1)
				}
				dp[i][j] = max(dp[i][j], dp[i-1][j])
				dp[i][j] = max(dp[i][j], dp[i][j-1])
			}
		}
		max_con := 0
		cnt := 100
		prev_con := false
		c := 0
		x := n
		y := m
		for x > 0 && y > 0 {
			if dp[x][y] == dp[x-1][y] {
				prev_con = false
				cnt = 0
				x--
			} else if dp[x][y] == dp[x][y-1] {
				prev_con = false
				cnt = 0
				y--
			} else {
				if prev_con {
					cnt++
					max_con = max(max_con, cnt)
				}
				prev_con = true
				x--
				y--
				c = y
			}
		}
		data = append(data, pair{(dp[n][m] + (100 - c) + 100*max_con), s})
	}
	sort.SliceStable(data, func(i, j int) bool { return data[i].a > data[j].a })
	var res []string
	for i := 0; i < siz; i++ {
		res = append(res, data[i].b)
	}
	return res
}
