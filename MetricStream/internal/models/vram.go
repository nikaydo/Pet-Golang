package models

import (
	"github.com/shirou/gopsutil/mem"
)

type VramInfo struct {
	VramUsage float64
	VramTotal uint64
	VramUsed  uint64
}

func (c *VramInfo) VramConv(vram *mem.VirtualMemoryStat) {
	c.VramUsage = vram.UsedPercent
	c.VramTotal = vram.Total
	c.VramUsed = vram.Used
}
