// Package cache 缓存工具类，可以缓存各种类型包括 struct 对象
package cache

import (
	"fmt"
	cache2 "github.com/astaxie/beego/cache"
	"strconv"
	"sync"
	"time"
)

type CacheService struct {
	Store cache2.Cache
}

var once sync.Once
var Cache *CacheService

func CacheStore() *CacheService {
	once.Do(func() {
		bm, err := cache2.NewCache("file", `{"CachePath":"./temp/cache","FileSuffix":".cache","DirectoryLevel":"2","EmbedExpiry":"120"}`)
		if err != nil {
			print(err)
		}
		Cache = &CacheService{
			Store: bm,
		}
	})
	return Cache
}


func Set(key string, obj string, expireTime int) error {
	return CacheStore().Store.Put(key, obj, time.Duration(expireTime)*time.Second)
}



func Del(key string) error {
	return CacheStore().Store.Delete(key)
}




func Get(key string) string {
	return CacheStore().Store.Get(key).(string)
}



func IsExist(key string) bool{
	b := CacheStore().Store.IsExist(key)
	return b && len(fmt.Sprint(CacheStore().Store.Get(key)))>0
}

func SetInt(key string, obj int, expireTime int) error {
	err :=CacheStore().Store.Put(key, obj, time.Duration(expireTime)*time.Second)
	return err
}

func GetInt(key string) int {
	a:=CacheStore().Store.Get(key)
	b := fmt.Sprint(a)
	i,err := strconv.Atoi(b)
	if err!=nil {
		return 0
	}
	return i
}