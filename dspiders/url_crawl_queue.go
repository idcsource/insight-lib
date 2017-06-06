// Copyright 2016-2017
// CoderG the 2016 project
// Insight 0+0 [ 洞悉 0+0 ]
// InDimensions Construct Source [ 忆黛蒙逝·建造源 ] -> idcsource@gmail.com
// Stephen Fire Meditation Qin [ 火志溟 ] -> firemeditation@gmail.com
// Use of this source code is governed by GNU LGPL v3 license

package dspiders

import (
	"fmt"
)

// The url crawl queue
type UrlCrawlQueue struct {
	urlchan chan UrlBasic
	list    map[string]bool
	count   uint
}

// Initialize the url crawl queue, the chan's length is const URL_CRAWL_QUEUE_CAP
func NewUrlCrawlQueue() (u *UrlCrawlQueue) {
	u = &UrlCrawlQueue{
		urlchan: make(chan UrlBasic, URL_CRAWL_QUEUE_CAP),
		list:    make(map[string]bool),
		count:   0,
	}
	return
}

// Add a url basic information to the url crawl queue
func (u *UrlCrawlQueue) Add(ub UrlBasic) (err error) {
	if u.count == URL_CRAWL_QUEUE_CAP {
		err = fmt.Errorf("The queue is full.")
		return
	}
	if _, find := u.list[ub.Url]; find == true {
		return
	}
	u.urlchan <- ub
	u.count++
	return
}

// Get one url basic information from the url crawl queue
func (u *UrlCrawlQueue) Get() (ub UrlBasic, err error) {
	if u.count == 0 {
		err = fmt.Errorf("The queue is empty.")
		return
	}
	ub = <-u.urlchan
	delete(u.list, ub.Url)
	u.count--
	return
}

// Get the queue's length
func (u *UrlCrawlQueue) Count() (count uint) {
	return u.count
}

// List all url in the queue
func (u *UrlCrawlQueue) List() (list []string) {
	listlen := len(u.list)
	list = make([]string, listlen)
	i := 0
	for url := range u.list {
		list[i] = url
		i++
	}
	return
}