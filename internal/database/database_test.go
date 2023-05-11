package database

import (
	"fmt"
	"testing"
	"time"
)

// func TestPut(t *testing.T) {
// 	fi := FileInfo{
// 		Address: "0x123456",
// 		SType:   1,
// 		Name:    "example.jpg",
// 		Mid:     "abc123",
// 		Size:    1024,
// 		OnChain: false,
// 	}

// 	flag, err := Put(fi)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if !flag {
// 		t.Error("put failed")
// 	}

// 	re, err := GetNotOnChain()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if len(re) != 1 {
// 		t.Error("not right")
// 	}

// 	f, err := Get(fi.Address, fi.Mid, 0)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	t.Error(f)

// 	l, err := List(fi.Address, fi.SType)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	fmt.Println(l)

// 	r, err := Delete(fi.Address, fi.Mid, fi.SType)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if !r {
// 		t.Error("put failed")
// 	}
// }

func TestWriteCheck(t *testing.T) {
	fi := FileInfo{
		Address:    "0x123456",
		SType:      1,
		Name:       "example.jpg",
		Mid:        "abc123",
		Size:       1024,
		ModTime:    time.Now(),
		UserDefine: "",
	}

	wc := NewWriteCheck()
	res, err := wc.Write(fi)
	if err != nil || !res {
		t.Error(err)
	}
	fi.Mid = "abc124"
	fi.Name = "111"
	res, err = wc.Write(fi)
	if err != nil || !res {
		t.Error(err)
	}
	go func() {
		fmt.Println("read")
		err := wc.Read()
		if err != nil {
			t.Error(err)
		}
	}()

	time.Sleep(time.Second * 10)
	l, err := List("0x123456", 1)
	for _, fi := range l {
		fmt.Println(fi)
	}
}
