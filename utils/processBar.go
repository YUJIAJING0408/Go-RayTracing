package utils

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type ProgressBar struct {
	total     int64     // 总任务量
	current   int64     // 当前进度
	width     int       // 进度条宽度(字符数)
	startTime time.Time // 开始时间
	sync.Mutex
}

func NewProgressBar(total int64, width int) *ProgressBar {
	return &ProgressBar{
		total: total,
		width: width,
	}
}

func (p *ProgressBar) Add(n int64) {
	if p.current == 0 {
		p.startTime = time.Now()
	}
	p.Lock()
	p.current += n
	p.Unlock()
	if p.current >= p.total {
		p.current = p.total
	}
	p.Render()
}

func (p *ProgressBar) Render() {
	// 计算完成百分比
	percent := float64(p.current) / float64(p.total)

	// 计算已完成的条带长度
	completedWidth := int(percent * float64(p.width))
	// 创建进度条字符串
	bar := strings.Repeat("=", completedWidth)
	if completedWidth < p.width {
		bar += ">" + strings.Repeat(" ", p.width-completedWidth-1)
	} else {
		bar = strings.Repeat("=", p.width)
	}
	// 计算耗时
	elapsed := time.Since(p.startTime).Round(time.Second)

	// 计算剩余时间
	var remaining string
	if p.current == 0 {
		remaining = "?" // 初始状态显示未知
	} else if p.current >= p.total {
		remaining = "0s" // 完成状态显示0秒
	} else {
		// 准确计算剩余时间
		remaining = time.Duration(
			float64(elapsed) * float64(p.total-p.current) / float64(p.current),
		).Round(time.Second).String()
	}
	// 输出进度条
	fmt.Printf("\r[%s] %6.2f%% (%d/%d) %v/%s",
		bar,
		percent*100,
		p.current,
		p.total,
		elapsed,
		remaining)

	// 完成后换行
	if p.current >= p.total {
		fmt.Println()
	}
}
