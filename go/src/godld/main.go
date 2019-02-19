// parallel download, with traditional waitgroup
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

type Target struct {
	Filename        string // Filename can be derived from Addr
	Addr            string
	Length          int    // total length of download target
	Status          []bool // downloading status. File at addr is divided into chunks, and the downloaded size is recorded in Status
	NumOfDownloader int
}

func (t *Target) GetLength() error {
	res, err := http.Head(t.Addr)
	if err != nil {
		fmt.Printf("%v when Head %v\n", err, t.Addr)
		return err
	}
	headers := res.Header
	fmt.Println(headers)
	l, err := strconv.Atoi(headers["Content-Length"][0])
	if err != nil {
		fmt.Printf("%v when getting length of %v\n", err, t.Addr)
		return err
	}
	t.Length = l
	return nil
}

func (t *Target) AcceptRangeRequestP() (bool, error) {
	res, err := http.Head(t.Addr)
	if err != nil {
		fmt.Printf("%v when Head %v\n", err, t.Addr)
		return false, err
	}
	headers := res.Header
	fmt.Println(headers)
	if headers["Accept-Ranges"] == nil {
		return false, nil
	}
	ar, err := strconv.Atoi(headers["Accept-Ranges"][0])
	if err != nil {
		fmt.Printf("%v when check Accept-Ranges header %v\n", err, t.Addr)
		return false, err
	}
	fmt.Printf("Accept-Ranges: %v\n", ar)
	if string(ar) == "bytes" {
		return true, nil
	} else {
		return false, nil
	}
}

// The partial file will be saved in t.Filename.idx.
func (t *Target) Download(start, end, idx int) error {
	fmt.Printf("start: %v, stop: %v\n, index: %v", start, end, idx)
	client := &http.Client{}
	req, err := http.NewRequest("GET", t.Addr, nil)
	if err != nil {
		fmt.Printf("%v when creating http request for %v\n", err, t.Addr)
		return err
	}
	if idx != -1 {
		req.Header.Add("Range", makeRangeHeader(start, end))
	}
	fmt.Println(req)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if idx == -1 {
		err = ioutil.WriteFile(t.Filename, r, 0x777)
	} else {
		err = ioutil.WriteFile(t.Filename+"."+strconv.Itoa(idx), r, 0x777)
	}
	if err != nil {
		fmt.Println(err)
		return err
	}
	if idx == -1 {
		t.Status[0] = true
	} else {
		t.Status[idx] = true
		//t.Filename = "a"
		wg.Done()
	}
	return nil
}

// helper functions for dispatcher
func makeRangeHeader(start, end int) string { // end is not inclusive
	return "bytes=" + strconv.Itoa(start) + "-" + strconv.Itoa(end-1)
}

func makeRange(length, n int) ([]int, []int) {
	var start []int = make([]int, n)
	var stop []int = make([]int, n)

	if length%n == 0 {
		step := length / n
		for i := 0; i < n; i++ {
			start[i] = i * step
			if i == n-1 {
				stop[i] = length
			} else {
				stop[i] = (i+1)*step - 1
			}
		}
		return start, stop
	} else {
		step := length / (n - 1)
		for i := 0; i < n-1; i++ {
			start[i] = i * step
			stop[i] = (i+1)*step - 1
		}
		start[n-1] = (n - 1) * step
		stop[n-1] = length
		return start, stop
	}
}

// Dispatch task to downloader and wait for their completion.
func (t *Target) Dispatch() error {
	supportRange, err := t.AcceptRangeRequestP()
	if supportRange == false || err != nil {
		fmt.Println("server does not support range request. switch to 1 downloader.")
		t.NumOfDownloader = 1
	}
	t.Status = make([]bool, t.NumOfDownloader)
	if t.NumOfDownloader == 1 {
		t.Download(0, 0, -1)
		return nil
	}
	start, stop := makeRange(t.Length, t.NumOfDownloader)
	// start and stop are of the same length
	for i := 0; i < len(start); i++ {
		wg.Add(1)
		go t.Download(start[i], stop[i], i)
	}
	wg.Wait()
	fmt.Printf("status: %v\n", t.Status)
	// Assemble the file
	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %v filename url\n", os.Args[0])
		os.Exit(-1)
	}
	var t *Target = &Target{os.Args[1],
		os.Args[2],
		//"http://www.cs.jhu.edu/~mdredze/datasets/sentiment/unprocessed.tar.gz",
		//"http://randd.defra.gov.uk/Document.aspx?Document=10043_R66141HouseholdElectricitySurveyFinalReportissue4.pdf",
		0,
		nil,
		10} // 10 downloader in total
	err := t.GetLength()
	fmt.Println(err)
	fmt.Println(t)
	fmt.Println(t.AcceptRangeRequestP())
	fmt.Println(makeRange(t.Length, 10))
	//err = t.Download(0, 158, 0)
	err = t.Dispatch()
}
