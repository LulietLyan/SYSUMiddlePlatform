package control

import (
	"backend/response"
	"fmt"
	"net"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type SystemInfo struct {
	CPUCoreCount    int     `json:"cpu_core_count"`
	CPUUsagePercent float64 `json:"cpu_usage_percent"`
	TotalMemory     uint64  `json:"total_memory"`
	UsedMemory      uint64  `json:"used_memory"`
	TotalSwap       uint64  `json:"total_swap"`
	UsedSwap        uint64  `json:"used_swap"`
	TotalDisk       uint64  `json:"total_disk"`
	UsedDisk        uint64  `json:"used_disk"`
}

var CPUCoreCount int
var CPUUsagePercent []int
var TotalMemory uint64
var UsedMemory []uint64
var TotalSwap uint64
var UsedSwap []uint64
var TotalDisk uint64
var UsedDisk []uint64
var Date []string

func dailyTask() {
	maxLength := 30
	fmt.Println("定时任务执行时间:", time.Now())
	cpuPercentages, _ := cpu.Percent(time.Second, false)
	var sum float64
	for _, cpuPercentage := range cpuPercentages {
		sum += cpuPercentage
	}
	if len(CPUUsagePercent) >= maxLength {
		CPUUsagePercent = CPUUsagePercent[1:]
	}
	CPUUsagePercent = append(CPUUsagePercent, int(sum/float64(len(cpuPercentages))))

	// Memory info
	vmStat, _ := mem.VirtualMemory()
	if len(UsedMemory) >= maxLength {
		UsedMemory = UsedMemory[1:]
	}
	UsedMemory = append(UsedMemory, vmStat.Used/(1024*1024))

	// Swap memory info
	swapStat, _ := mem.SwapMemory()
	if len(UsedSwap) >= maxLength {
		UsedSwap = UsedSwap[1:]
	}
	UsedSwap = append(UsedSwap, swapStat.Used/(1024*1024))

	// Disk info
	diskStat, _ := disk.Usage("/")
	if len(UsedDisk) >= maxLength {
		UsedDisk = UsedDisk[1:]
	}
	UsedDisk = append(UsedDisk, diskStat.Used/(1024*1024*1024))

	if len(Date) >= maxLength {
		Date = Date[1:]
	}
	Date = append(Date, time.Now().Format("2006-01-02 15:04"))
}

func InitSysInfo() error {
	var err error
	//cpu数量
	CPUCoreCount, err = cpu.Counts(true)
	if err != nil {
		return err
	}
	//总内存
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	TotalMemory = vmStat.Total / (1024 * 1024)
	//交换区
	swapStat, err := mem.SwapMemory()
	if err != nil {
		return err
	}
	TotalSwap = swapStat.Total / (1024 * 1024)
	//磁盘
	diskStat, err := disk.Usage("/")
	if err != nil {
		return err
	}
	TotalDisk = diskStat.Total / (1024 * 1024 * 1024)
	// 创建一个新的 cron 实例
	c := cron.New(cron.WithSeconds())

	// 添加一个每天执行的任务
	// _, err = c.AddFunc("0 0 0 * * *", dailyTask) // 每天零点执行
	_, err = c.AddFunc("0 0 * * * *", dailyTask) // 每小时运行
	if err != nil {
		panic(err)
	}

	// 启动 cron 调度器
	c.Start()

	// 确保 cron 调度器在应用退出时停止
	// defer c.Stop()
	return nil
}

func GetSystemInfo(c *gin.Context) {
	type Sysinfo struct {
		Info    string `json:"info"`
		IP      string `json:"IP"`
		Runtime int    `json:"runtime"`
	}
	// 获取系统启动时间
	info, _ := host.BootTime()

	uptime := time.Since(time.Unix(int64(info), 0)).Hours()
	sysinfo := Sysinfo{Info: runtime.GOOS, Runtime: int(uptime)}
	ipAddr, err := getLocalIP()
	if err == nil {
		sysinfo.IP = ipAddr
	}

	cpuPercentages, _ := cpu.Percent(time.Second, false)
	var sum float64
	for _, cpuPercentage := range cpuPercentages {
		sum += cpuPercentage
	}

	// Memory info
	vmStat, _ := mem.VirtualMemory()

	// Swap memory info
	swapStat, _ := mem.SwapMemory()

	// Disk info
	diskStat, _ := disk.Usage("/")

	response.Success(c, gin.H{"sysinfo": sysinfo, "cpuinfo": fmt.Sprintf("%d", CPUCoreCount) + "核心",
		"curcpu": int(sum / float64(len(cpuPercentages))), "curmem": vmStat.Used / (1024 * 1024),
		"curshare": swapStat.Used / (1024 * 1024), "curdisk": diskStat.Used / (1024 * 1024 * 1024), "maxmem": TotalMemory, "maxshare": TotalSwap, "maxdisk": TotalDisk,
		"dateX": Date, "cpuY": CPUUsagePercent, "memY": UsedMemory, "shareY": UsedSwap, "diskY": UsedDisk}, "")
}

// 获取本机IP地址
func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("未找到有效的IP地址")
}

// 获取系统运行时长
func getUptime() string {
	uptime := time.Since(startTime())
	return uptime.String()
}

// 获取系统启动时间
func startTime() time.Time {
	info, err := os.Stat("/proc")
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}
