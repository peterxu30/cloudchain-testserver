package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"

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

//TODO: Make it a variable number of blocks to add
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

//TODO: Test to simultaneously add blocks and read from current head to end
