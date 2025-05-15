package models

type CpuInfo struct {
	CpuUsage []float64
}

/*
	func (c *CpuInfo) CpuUsageString(cpuPercent []float64) {
		f := ""
		for _, i := range cpuPercent {
			f += fmt.Sprintf("%.1f", i) + ","
		}
		c.CpuUsage = f
	}
*/
