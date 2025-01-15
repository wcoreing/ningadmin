package common

import (
	"time"
)

func Task(jobs []func(time time.Time)) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		//从定时器中获取数据
		t := <-ticker.C
		for _, f := range jobs {
			go f(t)
		}
	}
	//阻塞
	for {
		time.Sleep(time.Second * 1)
	}

}
