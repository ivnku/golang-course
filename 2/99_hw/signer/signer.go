package main

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
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
	defer wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for rawData := range in {
		data, ok := rawData.(int)
		if !ok {
			fmt.Print("cant convert data to int in SingleHash function")
		}
		stringData := strconv.Itoa(data)

		crcChannel := make(chan string)
		crcChannelWithMd5 := make(chan string)

		result := calcSingleHash(stringData, crcChannel, crcChannelWithMd5, wg)

		out <- result
	}

}

/*
 Calculate two separate crc32 hashes and concatenate
 them with '~' sign. The result is SingleHash value
*/
func calcSingleHash(data string, crcChannel, crcChannelWithMd5 chan string, wg *sync.WaitGroup) string {
	md5Hash := DataSignerMd5(data)

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

	leftSide := <-crcChannel
	rightSide := <-crcChannelWithMd5
	result := leftSide + "~" + rightSide

	wg.Wait()

	return result
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	wgForClosingChannel := &sync.WaitGroup{}

	for rawData := range in {
		singleHashResult, ok := rawData.(string)
		if !ok {
			fmt.Print("cant convert data to string in MultiHash function")
		}

		channel := make(chan struct{index int; hash string})
		wg.Add(6)
		wgForClosingChannel.Add(6)
		for i := 0; i < 6; i++ {
			th := strconv.Itoa(i)
			go func(data string, channel chan struct{index int; hash string}, wg *sync.WaitGroup, wgForClosingChannel *sync.WaitGroup, index int) {
				result := DataSignerCrc32(data)
				resultStruct := struct {
					index int
					hash string
				}{
					index: index,
					hash: result,
				}
				channel <- resultStruct
				defer wg.Done()
				defer wgForClosingChannel.Done()
			}(th+singleHashResult, channel, wg, wgForClosingChannel, i)
		}
		go func(wgForClosingChannel *sync.WaitGroup, channel chan struct{index int; hash string}) {
			wgForClosingChannel.Wait()
			close(channel)
		}(wgForClosingChannel, channel)

		resultMap := make(map[int]string)
		for thHash := range channel {
			resultMap[thHash.index] = thHash.hash
		}

		result := ""
		for i := 0; i < 6; i++ {
			result += resultMap[i]
		}
		out <- result

		wg.Wait()
	}
}

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

func main() {
	testExpected := "1173136728138862632818075107442090076184424490584241521304_1696913515191343735512658979631549563179965036907783101867_27225454331033649287118297354036464389062965355426795162684_29568666068035183841425683795340791879727309630931025356555_3994492081516972096677631278379039212655368881548151736_4958044192186797981418233587017209679042592862002427381542_4958044192186797981418233587017209679042592862002427381542"
	testResult := "NOT_SET"

	var (
		DataSignerSalt         string = "" // на сервере будет другое значение
		OverheatLockCounter    uint32
		OverheatUnlockCounter  uint32
		DataSignerMd5Counter   uint32
		DataSignerCrc32Counter uint32
	)
	OverheatLock = func() {
		atomic.AddUint32(&OverheatLockCounter, 1)
		for {
			if swapped := atomic.CompareAndSwapUint32(&dataSignerOverheat, 0, 1); !swapped {
				fmt.Println("OverheatLock happend")
				time.Sleep(time.Second)
			} else {
				break
			}
		}
	}
	OverheatUnlock = func() {
		atomic.AddUint32(&OverheatUnlockCounter, 1)
		for {
			if swapped := atomic.CompareAndSwapUint32(&dataSignerOverheat, 1, 0); !swapped {
				fmt.Println("OverheatUnlock happend")
				time.Sleep(time.Second)
			} else {
				break
			}
		}
	}
	DataSignerMd5 = func(data string) string {
		atomic.AddUint32(&DataSignerMd5Counter, 1)
		OverheatLock()
		defer OverheatUnlock()
		data += DataSignerSalt
		dataHash := fmt.Sprintf("%x", md5.Sum([]byte(data)))
		time.Sleep(10 * time.Millisecond)
		return dataHash
	}
	DataSignerCrc32 = func(data string) string {
		atomic.AddUint32(&DataSignerCrc32Counter, 1)
		data += DataSignerSalt
		crcH := crc32.ChecksumIEEE([]byte(data))
		dataHash := strconv.FormatUint(uint64(crcH), 10)
		time.Sleep(time.Second)
		return dataHash
	}

	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	//inputData := []int{0, 1}

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			dataRaw := <-in
			data, ok := dataRaw.(string)
			if !ok {
				fmt.Println("cant convert result data to string")
			}
			testResult = data
		}),
	}

	start := time.Now()

	ExecutePipeline(hashSignJobs...)

	end := time.Since(start)

	expectedTime := 3 * time.Second

	if testExpected != testResult {
		fmt.Printf("results not match\nGot: %v\nExpected: %v", testResult, testExpected)
		fmt.Println()
	}

	if end > expectedTime {
		fmt.Printf("execition too long\nGot: %s\nExpected: <%s", end, time.Second*3)
		fmt.Println()
	}

	// 8 потому что 2 в SingleHash и 6 в MultiHash
	if int(OverheatLockCounter) != len(inputData) ||
		int(OverheatUnlockCounter) != len(inputData) ||
		int(DataSignerMd5Counter) != len(inputData) ||
		int(DataSignerCrc32Counter) != len(inputData)*8 {
		fmt.Println("not enough hash-func calls")
	}
}
