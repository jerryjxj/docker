package saptail

import (
    "testing"
    //redis "github.wdf.sap.corp/sigma-anywhere/saptail/redis"
    "time"
    //redis "github.wdf.sap.corp/sigma-anywhere/saptail/redis"
    "github.com/docker/docker/daemon/logger"
)

const redisurl = "127.0.0.1"
const redisport = "6379"

func Test_Message_lazymode_false(t *testing.T) {
    tail, err := New(redisurl, redisport, 0)
    if (err != nil) {
        t.Error(err)
    }
    defer tail.Close()
    msg := logger.Message{
        Line:      []byte("safdfa"),
        Source:    "filename",
        Timestamp: time.Now(),
        Attrs: map[string]string{
            "pod":       "liuzheng",
            "landscape": "us",
            "namespace": "cy",
            "container": "liuzheng-container",
            "node":      "192.168.0.1",
        },
        //Stat_t: redis.Stat_t{
        //    Dev:  333,
        //    Ino:  222,
        //    Size: 333,
        //},
    }
    tail.Message(&msg, false)
}
func Test_Message_lazymode_true(t *testing.T) {
    tail, err := New(redisurl, redisport, 0)
    if (err != nil) {
        t.Error(err)
    }
    defer tail.Close()

    msg := logger.Message{
        Line:      []byte("safdfa"),
        Source:    "filename",
        Timestamp: time.Now(),
        //Stat_t: redis.Stat_t{
        //    Dev:  333,
        //    Ino:  222,
        //    Size: 333,
        //},
    }
    tail.Message(&msg, true)
}
