package group

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestMembers(t *testing.T) {
	gid := "111111111111111"
	m := &Member{
		UnionId:   "1",
		NickName:  "1",
		ShowPhone: true,
	}

	start := time.Now()
	asyncLimit := make(chan bool, 20)
	var wg sync.WaitGroup
	for i := 1; i <= 200; i++ {
		asyncLimit <- true
		wg.Add(1)

		go func(i int) {
			defer func() {
				wg.Done()
				<-asyncLimit
			}()

			m.UnionId = fmt.Sprintf("%09d", i)
			err := AddMember(context.Background(), gid, m)
			if err != nil {
				t.Fatalf("%+v", err)
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("add 200 member, took %s\n", time.Now().Sub(start).String())

	// 查询成员列表
	start = time.Now()
	ms, err := ListDetailMembers(context.Background(), gid)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	fmt.Printf("list 200 member, took %s\n", time.Now().Sub(start).String())

	if len(ms) < 200 {
		t.Fatalf("list detail member faileds, get %d members", len(ms))
	}

	if ms[199].UnionId != fmt.Sprintf("%09d", 200) {
		t.Fatalf("get the 200 detail member faileds, get %v", ms[199])
	}
}
