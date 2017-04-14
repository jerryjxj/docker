package saptail

import (
    "testing"
    "time"
    "github.com/docker/docker/daemon/logger"
)

const redisurl = "127.0.0.1"
const redisport = "6379"

func Test_Message_lazymode_false(t *testing.T) {
    tail, err := New(redisurl, redisport, 0)
    if (err != nil) {
        t.Error(err)
    }
    t.Log("newed")
    t.Log(tail)
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
    }
    tail.Message(&msg, true)
}
