package daemon

import (
	"errors"
	"fmt"
	"log"
	"matsya/src/config"
	"matsya/src/rpc"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

func processBlocks(counter uint64, done chan<- any) {
	fmt.Printf("counter :: %d \n", counter)
	time.Sleep(time.Second)
	// done <- "done"
}

func getCurrentBlock() (*types.Block, error) {
	client := rpc.GetNodeProvider()
	if client == nil {
		return nil, errors.New("COULD NOT GET NODE PROVIDER")
	}

	latestBlock := rpc.GetLatestBlock(client)

	return latestBlock, nil
}

func getTargetBlock() (uint64, error) {
	targetBlockUint64 := config.GetUint64("target")

	targetString := os.Getenv("TARGET")

	if targetString != "" {
		targetStringInt, err := strconv.Atoi(targetString)
		if err != nil {
			log.Fatal(err)
		}
		targetBlockUint64 = uint64(targetStringInt)
	}

	return targetBlockUint64, nil
}

func LongRunningProcess(ch chan<- string) {
	targetBlock, err := getTargetBlock()
	if err != nil {
		log.Fatal(err)
	}

	latestBlock, _ := getCurrentBlock()

	counter := latestBlock.NumberU64()
	fmt.Println("Latest block number :", counter)

	for {
		if counter < targetBlock {
			ch <- "Reached the target block"
			break
		}

		done := make(chan any)
		defer close(done)

		processBlocks(counter, done)
		// <-done

		counter--
	}
}
