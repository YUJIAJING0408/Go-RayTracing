package utils

import (
	"testing"
	"time"
)

func TestProcessBar(t *testing.T) {
	total := int64(1000)
	bar := NewProgressBar(total, 50)

	for i := int64(0); i <= total; i++ {
		bar.Add(1)
		time.Sleep(50 * time.Millisecond) // 模拟任务处理
	}
}
