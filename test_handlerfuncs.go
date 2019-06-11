package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/peterxu30/cloudchain"
)

// Reset returns a handlerfunc that resets the backing CloudChain by deleting the CloudChain and reinitializing it.
func (s *TestServer) Reset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := DeleteTestCloudChain(r.Context())
		if err != nil {
			panic(err)
		}
		GetTestCloudChain(r.Context())
		fmt.Fprintln(w, "Test environment reset.")
	}
}

// AddFiftyBlocksTest starts 50 goroutines that add 1 block each to the cloudchain. The cloudchain is then iterated through to verify 50 blocks with the correct values were added.
func (s *TestServer) AddFiftyBlocksTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "AddFiftyBlocksTest started. \nThis test will add 50 blocks asynchronously to the cloudchain and then verify 50 blocks with the correct values have been added.\n")

		cc := GetTestCloudChain(r.Context())
		var numBlocks = 50
		var wg sync.WaitGroup
		wg.Add(numBlocks)

		expected := make(map[string]bool)
		for i := 0; i < numBlocks; i++ {
			msg := strconv.Itoa(i)
			go func(ctx context.Context, msg string) {
				defer wg.Done()
				recieved, errorChannel := cc.AddBlockExperimental(ctx, []byte(msg))

				select {
				case block := <-recieved:
					fmt.Fprintf(w, "Added block with data %v\n", string(block.Data))
				case err := <-errorChannel:
					fmt.Fprintf(w, "Encountered error %s\n", err.Error())
				}

			}(r.Context(), msg)

			expected[msg] = false
		}

		wg.Wait()
		fmt.Fprintln(w, "All blocks added.")

		iter, err := cc.Iterator()
		if err != nil {
			fmt.Fprintf(w, "Iterator could not be created: %s\n", err.Error())
			return
		}

		for i := 0; i < numBlocks; i++ {
			block, err := iter.Next(r.Context())
			if _, ok := err.(*cloudchain.StopIterationError); ok {
				fmt.Fprintf(w, "Reached end of CloudChain prematurely: %s\n", err.Error())
				break
			} else if err != nil {
				fmt.Fprintf(w, "Block could not be retrieved: %s\n", err.Error())
				continue
			}

			if block == nil {
				fmt.Fprintln(w, "Block is nil.")
				break
			}

			msg := string(block.Data)
			if _, ok := expected[msg]; !ok {
				fmt.Fprintf(w, "Found block with unexpected value: %s\n", msg)
			} else {
				expected[msg] = !expected[msg]
			}
		}

		missingValues := make([]string, 0)
		for k, v := range expected {
			if !v {
				missingValues = append(missingValues, k)
			}
		}

		if len(missingValues) > 0 {
			fmt.Fprintf(w, "Missing blocks with values %v", missingValues)
		} else {
			fmt.Fprintf(w, "Successfully added and verified 50 blocks.")
		}
	}
}

// SimultaneouslyAddAndReadFiftyBlocksTest simultaneously add 50 blocks asynchronously to the cloudchain and read the blockchain from whatever the head currently is to the end.
// It then synchronously verifies that 50 blocks were added.
func (s *TestServer) SimultaneouslyAddAndReadFiftyBlocksTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "SimultaneouslyAddAndReadFiftyBlocksTest started. \nThis test will simultaneously add 50 blocks asynchronously to the cloudchain and read the blockchain from whatever the head currently is to the end.\n")

		cc := GetTestCloudChain(r.Context())
		var wg sync.WaitGroup
		wg.Add(100)

		for i := 0; i < 100; i++ {
			if i%2 == 0 {
				// Add blocks on evens
				msg := strconv.Itoa(i)
				go func(ctx context.Context, cc *cloudchain.CloudChain, msg string) {
					defer wg.Done()
					recieved, errorChannel := cc.AddBlockExperimental(ctx, []byte(msg))

					select {
					case block := <-recieved:
						fmt.Fprintf(w, "Added block with data %v\n", string(block.Data))
					case err := <-errorChannel:
						fmt.Fprintf(w, "Encountered error %s\n", err.Error())
					}
				}(r.Context(), cc, msg)
			} else {
				// Read blocks on odds
				go func(ctx context.Context, cc *cloudchain.CloudChain) {
					defer wg.Done()

					sleepTime := rand.Intn(10)
					time.Sleep(sleepTime * time.Second)

					iter, err := cc.Iterator()
					if err != nil {
						fmt.Fprintf(w, "Error creating iterator %s\n", err.Error())
					}

					for i := 0; i < 50; i++ {
						_, err := iter.Next(ctx)
						if _, ok := err.(*cloudchain.StopIterationError); ok {
							fmt.Fprintf(w, "Reached the end. Read %v blocks.\n", i)
							return
						} else if err != nil {
							fmt.Fprintf(w, "Encountered unexpected error %s\n", err.Error())
							return
						}
					}
				}(r.Context(), cc)
			}
		}

		wg.Wait()
		fmt.Fprintln(w, "All blocks added.")

		iter, err := cc.Iterator()
		if err != nil {
			fmt.Fprintf(w, "Iterator could not be created: %s\n", err.Error())
			return
		}

		blocksRead := 0
		// Then read serially and verify 50 blocks added.
		for i := 0; i < 50; i++ {
			block, err := iter.Next(r.Context())
			if _, ok := err.(*cloudchain.StopIterationError); ok {
				fmt.Fprintf(w, "Reached end of CloudChain prematurely: %s\n", err.Error())
				break
			} else if err != nil {
				fmt.Fprintf(w, "Block could not be retrieved: %s\n", err.Error())
				continue
			}

			if block == nil {
				fmt.Fprintln(w, "Block is nil.")
				break
			}

			msg := string(block.Data)
			fmt.Fprintf(w, "Read block with message %s\n", msg)
			blocksRead++
		}

		fmt.Fprintf(w, "Read a total of %v blocks out of 50.", blocksRead)
	}
}
