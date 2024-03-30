package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/ssz"

	_ "net/http/pprof"
)

func main() {
	body, _ := ioutil.ReadFile("state2.ssz")

	test1(body)
	test2(body)

	f, _ := os.Create("mem.pprof")
	pprof.WriteHeapProfile(f)
	f.Close()
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func test1(body []byte) {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		fmt.Printf("unmarshal time: %v\n", elapsed)
	}()

	fmt.Printf("\nfastssz / go-eth2-client\n")

	t := new(deneb.BeaconState)
	err := t.UnmarshalSSZ(body)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	runtime.GC()
	printMemUsage()

	//root, _ := t.HashTreeRoot()
	//fmt.Printf("state root: 0x%x\n", root)
	fmt.Printf("gvr: 0x%x\n", t.GenesisValidatorsRoot)

	_, err = t.MarshalSSZ()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
}

func test2(body []byte) {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		fmt.Printf("unmarshal time: %v\n", elapsed)
	}()

	fmt.Printf("\npk's dynamic ssz\n")

	t := new(deneb.BeaconState)
	err := ssz.UnmarshalSSZ(t, body)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	runtime.GC()
	printMemUsage()

	//root, _ := t.HashTreeRoot()
	//fmt.Printf("state root: 0x%x\n", root)
	fmt.Printf("gvr: 0x%x\n", t.GenesisValidatorsRoot)

	_, err = ssz.MarshalSSZ(t, nil)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	/*
		if len(buf) != len(body) {
			fmt.Printf("size mismatch: %v / %v\n", len(buf), len(body))
		}
		for i := 0; i < len(body); i++ {
			if body[i] != buf[i] {
				fmt.Printf("ssz mismatch: %v : %v / %v\n", i, buf[i], body[i])
				break
			}
		}
	*/
}
