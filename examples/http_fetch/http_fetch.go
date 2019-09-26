package main

import (
	"fmt"

	filter "github.com/morya/go-dirtyfilter"
	"github.com/morya/go-dirtyfilter/store"
)

var (
	filterText = `我是需要过滤的内容，内容为：**文*@@件**名，需要过滤。。。`
)

func main() {
	fs, err := store.NewFetchStore(store.FetchConfig{
		Remote: "http://test.domain.com/garbage_words",
	})
	if err != nil {
		panic(err)
	}
	filterManage := filter.NewDirtyManager(fs)
	result, err := filterManage.Filter().Filter(filterText, '*', '@')
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
