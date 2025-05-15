package metric

import (
	"context"
	"main/internal/models"
	"math"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func Metrics() (models.Data, error) {
	var CpuInfo models.CpuInfo
	var VramInfo models.VramInfo
	var Data models.Data
	c, err := cpuMetrics()
	if err != nil {
		return Data, err
	}
	for i, v := range c {
		c[i] = math.Round(v*10) / 10
	}
	CpuInfo.CpuUsage = c
	vram, err := vramMetrics()
	if err != nil {
		return Data, err
	}
	vram.UsedPercent = math.Round(vram.UsedPercent*10) / 10
	VramInfo.VramConv(vram)
	Data = models.Data{Cpu: CpuInfo, Vram: VramInfo}
	return Data, err
}

func cpuMetrics() ([]float64, error) {
	return cpu.PercentWithContext(context.Background(), time.Second, true)
}

func vramMetrics() (*mem.VirtualMemoryStat, error) {
	return mem.VirtualMemory()
}
