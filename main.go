package main

import (
	"errors"
	"fmt"
	"github.com/containerd/cgroups"
	"golang.org/x/exp/rand"
	_ "net/http/pprof"
	"os"
	"path"
	"time"
)

const (
	filePath  = "/tmp/data"
	batchSize = 500 * 1024 * 1024
	fileSize  = 5 * 1024 * 1024 * 1024
)

func printMemInfo() {
	memInfoFile, err := os.Open("/sys/fs/cgroup/memory/docker/2b95e9f9042b7cafddfc9b76f0d55251a5c4458341ba010ff8810662618d4844/memory.stat")
	if err != nil {
		fmt.Printf("Failed to open /proc/meminfo: %v\n", err)
		return
	}
	defer memInfoFile.Close()

	buffer := make([]byte, 1024)
	_, err = memInfoFile.Read(buffer)
	if err != nil {
		fmt.Printf("Failed to read /proc/meminfo: %v\n", err)
		return
	}
}

func PrintActiveInactive() {
	control, err := cgroups.Load(cgroups.V1, cgroups.RootPath)
	if err != nil {
		panic(err)
	}
	stats, err := control.Stat(cgroups.IgnoreNotExist)
	if err != nil {
		panic(err)
	}
	if stats.Memory == nil || stats.Memory.Usage == nil {
		panic(errors.New("cannot find memory usage info from cGroups"))
	}
	// 1. usage
	// ref: <https://github.com/docker/cli/blob/e57b5f78de635e6e2b688686d10b830c4747c4dc/cli/command/container/stats_helpers.go#L239>
	//usage := stats.Memory.Usage.Usage
	//if inactiveFile < usage {
	//	usage = usage - stats.Memory.TotalInactiveFile
	//}
	//fmt.Printf("total:%.2fMB\n", float64(stats.Memory.Usage.Limit)/1024/1024)
	//fmt.Printf("used:%.2fMB\n", float64(usage)/1024/1024)
	fmt.Printf("%.2f ", float64(stats.Memory.TotalInactiveFile)/1024/1024)
	fmt.Printf("%.2f ", float64(stats.Memory.TotalActiveFile)/1024/1024)
}

func PrintMeminfo() {
	control, err := cgroups.Load(cgroups.V1, cgroups.RootPath)
	if err != nil {
		panic(err)
	}
	stats, err := control.Stat(cgroups.IgnoreNotExist)
	if err != nil {
		panic(err)
	}
	if stats.Memory == nil || stats.Memory.Usage == nil {
		panic(errors.New("cannot find memory usage info from cGroups"))
	}
	// 1. usage
	// ref: <https://github.com/docker/cli/blob/e57b5f78de635e6e2b688686d10b830c4747c4dc/cli/command/container/stats_helpers.go#L239>
	usage := stats.Memory.Usage.Usage
	if stats.Memory.TotalInactiveFile < usage {
		usage = usage - stats.Memory.TotalInactiveFile - stats.Memory.TotalActiveFile
	}
	fmt.Printf("----------------\n")
	fmt.Printf("total:%.2fMB\n", float64(stats.Memory.Usage.Limit)/1024/1024)
	fmt.Printf("used:%.2fMB\n", float64(usage)/1024/1024)
	fmt.Printf("inactive(file):%.2fMB\n", float64(stats.Memory.TotalInactiveFile)/1024/1024)
	fmt.Printf("active(file):%.2fMB\n", float64(stats.Memory.TotalActiveFile)/1024/1024)
	//fmt.Printf("inactive(anon):%.2fMB\n", float64(stats.Memory.TotalInactiveAnon)/1024/1024)
	//fmt.Printf("active(anon):%.2fMB\n", float64(stats.Memory.TotalActiveAnon)/1024/1024)
	//fmt.Printf("cached:%.2fMB\n", float64(stats.Memory.TotalCache)/1024/1024)
	fmt.Printf("mapped:%.2fMB\n", float64(stats.Memory.TotalMappedFile)/1024/1024)
	fmt.Printf("----------------\n")
}

func WriteFile() {
	f, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}()
	batches := fileSize / batchSize
	buf := make([]byte, batchSize)
	for j := 0; j < batchSize; j++ {
		buf[j] = byte(11)
	}
	for i := 0; i < batches; i++ {
		_, err := f.Write(buf)
		if err != nil {
			panic(err)
		}
		//fmt.Printf("wrote %d bytes to file!\n", n2)
	}

	//// then we can call rand.Read.
	//_, err = rand.Read(buf)
	//if err != nil {
	//	panic(err)
	//}
}

func Mmap() *ReaderAt {
	readAt, err := Open(path.Clean(filePath))
	if err != nil {
		panic(err)
	}
	return readAt
}

func ReadMmapFile(duration time.Duration, r *ReaderAt) []byte {
	var p []byte
	timer := time.After(duration)
	for {
		select {
		case <-timer:
			return p
		default:
			// ReadAt
			//p = make([]byte, fileSize)
			//_, err := r.ReadAt(p, 0)
			//if err != nil {
			//	panic(err)
			//}
			// At
			//for i := 0; i < fileSize; i++ {
			//	_ = r.At(i)
			//}
			_ = r.At(rand.Intn(fileSize))
		}
	}
}

func ReadOnce(r *ReaderAt) {
	var dummy float64
	for i := 0; i < r.Len(); i++ {
		dummy += float64(r.At(i))
	}
	//r.data = nil
	//for i := 0; i < len(r.data); i++ {
	//	_ = r.data[i]
	//}
}

func Print(duration time.Duration) {
	timer := time.After(duration)
	ticker := time.NewTicker(200 * time.Millisecond)
	for {
		select {
		case <-timer:
			return
		case <-ticker.C:
			PrintMeminfo()
		}
	}
}

func RemoveFile() {
	err := os.Remove(filePath)
	if err != nil {
		panic(err)
	}
}

func Malloc(size int) []byte {
	b := make([]byte, size)
	return b
}

func main() {
	//go func() {
	//	log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	//}()

	//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>> before start")
	//PrintMeminfo()
	//
	//WriteFile()
	//time.Sleep(1 * time.Second)
	//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>> after WriteFile")
	//PrintMeminfo()

	//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>> before mmap")
	//PrintMeminfo()
	//
	//r := Mmap()
	//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>> after mmap")
	//PrintMeminfo()
	//time.Sleep(1 * time.Second)
	//
	//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>> read mmap file")
	//// read test 1
	////go ReadMmapFile(60*time.Second, r)
	////go Print(60 * time.Second)
	////time.Sleep(60 * time.Second)
	//// read test 2
	//ReadOnce(r)
	//ReadOnce(r)
	//PrintMeminfo()
	//time.Sleep(1 * time.Second)

	//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>> after munmap")
	//err := r.Close()
	//if err != nil {
	//	panic(err)
	//}
	////RemoveFile()
	//for i := 0; i < 1; i++ {
	//	time.Sleep(1 * time.Second)
	//	PrintMeminfo()
	//}

	//debug.FreeOSMemory()

	//time.Sleep(1 * time.Second)
	//PrintMeminfo()

	fmt.Println("malloc b......")
	b := Malloc(4 * 1024 * 1024 * 1024)
	for i := 0; i < 1; i++ {
		time.Sleep(1 * time.Second)
		PrintMeminfo()
	}

	fmt.Println("malloc c......")
	c := Malloc(3584 * 1024 * 1024) // 3.5GB
	for i := 0; i < 1; i++ {
		time.Sleep(1 * time.Second)
		PrintMeminfo()
	}

	//fmt.Println("malloc d......")
	//d := Malloc(3 * 1024 * 1024 * 1024)
	//for i := 0; i < 5; i++ {
	//	time.Sleep(1 * time.Second)
	//	PrintMeminfo()
	//}

	//var batches = len(b) / batchSize
	//var memSize = 0
	//for i := 0; i < batches; i++ {
	//	for j := 0; j < batchSize; j++ {
	//		index := i*batchSize + j
	//		b[index] = byte(index % 10)
	//	}
	//	memSize += batchSize
	//	PrintActiveInactive()
	//	fmt.Printf("%d\n", memSize)
	//}

	fmt.Println("write b...")
	for i := 0; i < len(b); i++ {
		b[i] = byte(i % 10)
	}
	for i := 0; i < 1; i++ {
		time.Sleep(1 * time.Second)
		PrintMeminfo()
	}

	fmt.Println("write c...")
	for i := 0; i < len(c); i++ {
		c[i] = byte(i % 10)
	}
	for i := 0; i < 1; i++ {
		time.Sleep(1 * time.Second)
		PrintMeminfo()
	}

	//for i := 0; i < len(d); i++ {
	//	d[i] = byte(i % 10)
	//}
	//for i := 0; i < 5; i++ {
	//	time.Sleep(1 * time.Second)
	//	PrintMeminfo()
	//}

	//RemoveFile()
}
