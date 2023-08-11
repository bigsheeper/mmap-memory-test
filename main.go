package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/containerd/cgroups"
)

const (
	filePath  = "/tmp/data"
	batchSize = 500 * 1024 * 1024
	fileSize  = 5 * 1024 * 1024 * 1024
)

func PrintFullMemInfo() {
	memInfoFile, err := os.Open("/sys/fs/cgroup/memory/docker/2b95e9f9042b7cafddfc9b76f0d55251a5c4458341ba010ff8810662618d4844/memory.stat")
	if err != nil {
		panic(err)
	}
	defer memInfoFile.Close()

	buffer := make([]byte, 1024)
	_, err = memInfoFile.Read(buffer)
	if err != nil {
		panic(err)
	}
	// TODO
}

func PrintMeminfo(msg string) {
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
	// ref: <https://github.com/docker/cli/blob/e57b5f78de635e6e2b688686d10b830c4747c4dc/cli/command/container/stats_helpers.go#L239>
	usage := stats.Memory.Usage.Usage
	usage2 := usage - stats.Memory.TotalInactiveFile - stats.Memory.TotalActiveFile
	fmt.Printf(">>>>>>>>>>>>>>>> %s\n", msg)
	fmt.Printf("total:%.2fMB\n", float64(stats.Memory.Usage.Limit)/1024/1024)
	fmt.Printf("used:%.2fMB\n", float64(usage)/1024/1024)
	fmt.Printf("inactive(file):%.2fMB\n", float64(stats.Memory.TotalInactiveFile)/1024/1024)
	fmt.Printf("active(file):%.2fMB\n", float64(stats.Memory.TotalActiveFile)/1024/1024)
	fmt.Printf("used-active-inactive:%.2fMB\n", float64(usage2)/1024/1024)
	//fmt.Printf("inactive(anon):%.2fMB\n", float64(stats.Memory.TotalInactiveAnon)/1024/1024)
	//fmt.Printf("active(anon):%.2fMB\n", float64(stats.Memory.TotalActiveAnon)/1024/1024)
	//fmt.Printf("cached:%.2fMB\n", float64(stats.Memory.TotalCache)/1024/1024)
	//fmt.Printf("mapped:%.2fMB\n", float64(stats.Memory.TotalMappedFile)/1024/1024)
	fmt.Println("")
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
		fmt.Printf("wrote %d bytes to file!\n", batchSize)
	}
}

func Mmap() *ReaderAt {
	readAt, err := Open(path.Clean(filePath))
	if err != nil {
		panic(err)
	}
	return readAt
}

func ReadOnce(r *ReaderAt) []byte {
	var dummy = make([]byte, 1024*1024)
	//fmt.Printf("read total %d bytes\n", r.Len())
	for i := 0; i < r.Len(); i += 4 << 10 {
		dummy[i%1024] = r.At(i)
	}
	return dummy
}

func RemoveFile() {
	err := os.Remove(filePath)
	if err != nil {
		panic(err)
	}
}

func Malloc(size int) []byte {
	fmt.Printf("malloc %.2fGB...\n", float64(size)/1024/1024/1024)
	b := make([]byte, size)
	for i := 0; i < len(b); i++ {
		b[i] = byte(i % 10)
	}
	return b
}

func ReadBench(r *ReaderAt) {
	start := time.Now()
	times := 10000
	for i := 0; i < times; i++ {
		ReadOnce(r)
	}
	msg := fmt.Sprintf("ReadBench: run %d times, time taken: %s", times, time.Since(start))
	PrintMeminfo(msg)
}

func main() {
	PrintMeminfo("")
	time.Sleep(1 * time.Second)

	//WriteFile()
	//PrintMeminfo("WriteFile")
	//time.Sleep(1 * time.Second)

	PrintMeminfo("")
	time.Sleep(1 * time.Second)

	r := Mmap()
	PrintMeminfo("mmaped")
	time.Sleep(1 * time.Second)
	
	for i := 0; i < 2; i++ {
		_ = ReadOnce(r)
		PrintMeminfo("read mmap file")
		time.Sleep(1 * time.Second)
	}
	//ReadBench(r)

	r.Munmap()
	PrintMeminfo("munmap")
	time.Sleep(1 * time.Second)

	Malloc(4 * 1024 * 1024 * 1024)
	PrintMeminfo("")
	time.Sleep(1 * time.Second)

	Malloc(3891 * 1024 * 1024) // 3.8GB
	PrintMeminfo("")
	time.Sleep(1 * time.Second)

	//RemoveFile()
}
