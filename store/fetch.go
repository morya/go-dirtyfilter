package store

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
)

type FetchConfig struct {
	Remote string
}

type FetchResp struct {
	Code int
	Data []string
}

type FetchStore struct {
	version uint64
	cfg     FetchConfig

	lock sync.Mutex

	tmp     *sync.Map
	persist *sync.Map
}

func fetch(remote string) (*FetchResp, error) {
	resp, err := http.Get(remote)
	if err != nil {
		return nil, fmt.Errorf("%w fetch failed", err)
	}
	defer resp.Body.Close()

	var f = &FetchResp{}
	dec := json.NewDecoder(resp.Body)

	if err := dec.Decode(f); err != nil {
		return nil, errors.Wrap(err, "decode json failed")
	}

	return f, nil
}

// NewFetchStore 创建敏感词内存存储
func NewFetchStore(config FetchConfig) (*FetchStore, error) {
	store := &FetchStore{
		cfg: config,

		tmp:     &sync.Map{},
		persist: &sync.Map{},
	}

	if err := store.Update(); err != nil {
		return nil, err
	}

	return store, nil
}

// Write Write
func (ms *FetchStore) Write(words ...string) error {
	return errors.New("later")
}

func (ms *FetchStore) Update() error {
	fetchResp, err := fetch(ms.cfg.Remote)
	if err != nil {
		return err
	}

	if fetchResp.Code == 0 && len(fetchResp.Data) > 0 {
		ms.tmp = &sync.Map{}

		for _, d := range fetchResp.Data {
			ms.tmp.Store(d, 1)
		}

		ms.persist = ms.tmp
		atomic.AddUint64(&ms.version, 1)
	}

	return nil
}

// Read Read
func (ms *FetchStore) Read() <-chan string {
	chResult := make(chan string)
	go func() {
		ms.persist.Range(func(k, v interface{}) bool {
			chResult <- v.(string)
			return true
		})
	}()
	return chResult
}

// ReadAll ReadAll
func (ms *FetchStore) ReadAll() ([]string, error) {
	result := make([]string, 8)
	ms.persist.Range(func(k, v interface{}) bool {
		result = append(result, v.(string))
		return true
	})
	return result, nil
}

// Remove Remove
func (ms *FetchStore) Remove(words ...string) error {
	if len(words) == 0 {
		return nil
	}
	for i, l := 0, len(words); i < l; i++ {
		ms.persist.Delete(words[i])
	}
	atomic.AddUint64(&ms.version, 1)
	return nil
}

// Version Version
func (ms *FetchStore) Version() uint64 {
	return ms.version
}
