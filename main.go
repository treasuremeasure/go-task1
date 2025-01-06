package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func main() {

	//send get request to the host on srv.msk01.gigacorp.local
	//and print the response

	errorCount := 0

	for {

		if errorCount >= 3 {
			fmt.Println("Unable to fetch server statistic")
			errorCount = 0
		}

		resp, err := http.Get("http://srv.msk01.gigacorp.local")

		//check if http request was successful
		if err != nil {
			errorCount++
			continue
		}

		//check if http responce status code is ok
		if resp.StatusCode == 200 {

			//check if http responce contents are ok
			contentType := resp.Header.Get("Content-Type")

			if contentType != "text/plain; charset=utf-8" {
				fmt.Println("bad content type, received:", contentType)
				errorCount++
				continue
			}

		} else {
			//http responce status code is not ok (not 200)
			fmt.Println("http responce status not 200")
			errorCount++
			continue
		}

		//read responce body
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Println("unable to read response body")
			errorCount++
			continue
		}

		//split responce body by comma
		values := strings.Split(string(body), ",")

		//check if values are ok
		if len(values) != 7 {
			fmt.Println("wrong amount of values")
			errorCount++
			continue
		}

		averageLoad, err := strconv.Atoi(values[0])
		if err != nil {
			continue
		}

		maxRAM, err := strconv.Atoi(values[1])
		if err != nil {
			continue
		}

		usedRAM, err := strconv.Atoi(values[2])
		if err != nil {
			continue
		}

		diskSpaceBytes, err := strconv.Atoi(values[3])
		if err != nil {
			continue
		}

		usedDiskSpaceBytes, err := strconv.Atoi(values[4])
		if err != nil {
			continue
		}

		netThroughputBytesPerS, err := strconv.Atoi(values[5])
		if err != nil {
			continue
		}

		netLoadBytesPerS, err := strconv.Atoi(values[6])
		if err != nil {
			continue
		}

		if averageLoad > 30 {
			fmt.Println("Load Average is too high:", averageLoad)
		}

		if float64(usedRAM)/float64(maxRAM) > 0.8 {
			//ram_usage_percent_str := strconv.FormatFloat((float64(used_ram)/float64(max_ram))*100, 'f', 2, 64)
			ramUsagePercentStr := strconv.FormatInt(int64(math.Round((float64(usedRAM)/float64(maxRAM))*100)), 10) + "%"
			fmt.Println("Memory usage too high:", ramUsagePercentStr)
		}

		if float64(usedDiskSpaceBytes)/float64(diskSpaceBytes) > 0.9 {
			diskUsagePercentStr := strconv.FormatInt(int64((diskSpaceBytes-usedDiskSpaceBytes)/1048576), 10)
			fmt.Println("Free disk space is too low:", diskUsagePercentStr, "Mb left")
		}

		if float64(netLoadBytesPerS)/float64(netThroughputBytesPerS) > 0.9 {
			netAvailableMbitsStr := strconv.FormatInt(int64((netThroughputBytesPerS-netLoadBytesPerS)/1000000), 10)
			fmt.Println("Network bandwidth usage high:", netAvailableMbitsStr, "Mbit/s available")
		}

	}
}
