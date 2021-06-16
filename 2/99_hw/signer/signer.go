package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func ExecutePipeline(jobs ...job) {
	channelsPool := make([]chan interface{}, 0)
	for i := 0; i < len(jobs)+1; i++ {
		channelsPool = append(channelsPool, make(chan interface{}))
	}

	wg := &sync.WaitGroup{}

	for i, jobFunc := range jobs {
		wg.Add(1)
		go func(in, out chan interface{}, jobFunc job, wg *sync.WaitGroup) {
			jobFunc(in, out)
			defer wg.Done()
			defer close(out)
		}(channelsPool[i], channelsPool[i+1], jobFunc, wg)
	}

	wg.Wait()
}

/*
 Calculate crc32() hash and crc32(md5()) from input data
 and pass it further to MultiHash
*/
func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for rawData := range in {
		data, ok := rawData.(int)
		if !ok {
			fmt.Print("cant convert data to int in SingleHash function")
		}
		stringData := strconv.Itoa(data)

		wg.Add(1)
		md5Hash := DataSignerMd5(stringData)
		go func(wg *sync.WaitGroup, stringData, md5Hash string, out chan interface{}) {
			calcSingleHash(stringData, md5Hash, out)
			defer wg.Done()
		}(wg, stringData, md5Hash, out)
	}
	wg.Wait()
}

/*
 Calculate two separate crc32 hashes and concatenate
 them with '~' sign. The result is SingleHash value
*/
func calcSingleHash(data, md5Hash string, out chan interface{}) {
	wg := &sync.WaitGroup{}

	crcChannel := make(chan string)
	crcChannelWithMd5 := make(chan string)

	wg.Add(2)
	go func(data string, channel chan string, wg *sync.WaitGroup) {
		result := DataSignerCrc32(data)
		channel <- result
		defer wg.Done()
		defer close(channel)
	}(data, crcChannel, wg)

	go func(data string, channel chan string, wg *sync.WaitGroup) {
		result := DataSignerCrc32(data)
		channel <- result
		defer wg.Done()
		defer close(channel)
	}(md5Hash, crcChannelWithMd5, wg)

	result := <-crcChannel + "~" + <-crcChannelWithMd5

	out <- result
	wg.Wait()
}

/*
 Calculate MultiHash and pass it further to CombineResults
*/
func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for rawData := range in {
		singleHashResult, ok := rawData.(string)
		if !ok {
			fmt.Print("cant convert data to string in MultiHash function")
		}

		wg.Add(1)

		go func(wg *sync.WaitGroup, data string, out chan interface{}) {
			calcMultiHash(data, out)
			defer wg.Done()
		}(wg, singleHashResult, out)

	}
	wg.Wait()
}

/*
 Calculate six crc32 hashes (index + value from SingleHash) and concatenate them
*/
func calcMultiHash(data string, out chan interface{}) {
	wg := &sync.WaitGroup{}

	hashes := make([]string, 6)
	wg.Add(6)
	for i := 0; i < 6; i++ {
		th := strconv.Itoa(i)
		go func(data string, wgForClosingChannel *sync.WaitGroup, index int, hashes []string) {
			hash := DataSignerCrc32(data)
			hashes[index] = hash
			defer wgForClosingChannel.Done()
		}(th+data, wg, i, hashes)
	}
	wg.Wait()

	result := strings.Join(hashes, "")
	out <- result
}

/*
 Sort all of the results and concatenate them with "_" sign
*/
func CombineResults(in, out chan interface{}) {
	combinedResult := make([]string, 0)

	for rawData := range in {
		data, ok := rawData.(string)
		if !ok {
			fmt.Print("cant convert data to string in CombineResult function")
		}
		combinedResult = append(combinedResult, data)
	}

	sort.Strings(combinedResult)
	result := strings.Join(combinedResult, "_")

	out <- result
}
