package redis

import (
    "testing"
    "time"
    "github.com/docker/docker/daemon/logger"

    //redis "github.wdf.sap.corp/sigma-anywhere/saptail/redis"
)

const redisurl = "127.0.0.1"
const redisport = "6379"

func Test_NewConnection_1(t *testing.T) {
    client, err := New(redisurl, redisport, 0)
    if (err != nil) {
        t.Error(err)
    }
    res, err := client.HealthCheck()
    if (err != nil) {
        t.Error(err)
    }
    if (res == "PONG") {
        t.Log(res)
    } else {
        t.Fatalf("the result is `%v`, we need `PONG`", res)
    }
    err = client.Close()
    if (err != nil) {
        t.Error(err)
    }
}

func Test_SetOffset_1(t *testing.T) {
    client, err := New(redisurl, redisport, 0)
    if (err != nil) {
        t.Error(err)
    }
    defer client.Close()

    client.SetOffset("sss", 3)
}
func Test_SetOffset_2(t *testing.T) {
    client, err := New(redisurl, redisport, 0)
    if (err != nil) {
        t.Error(err)
    }
    defer client.Close()

    msg := logger.Message{
        Line:      []byte("safdfa"),
        Source:    "filename",
        Timestamp: time.Now(),
    }
    t.Log(9)
    err = client.SetOffset(msg.Source, 9)
    if (err != nil) {
        t.Error(err)
    }
}
func Test_GetOffset_2(t *testing.T) {
    client, err := New(redisurl, redisport, 0)
    if (err != nil) {
        t.Error(err)
    }
    defer client.Close()

    msg := logger.Message{
        Line:      []byte("safdfa"),
        Source:    "filename",
        Timestamp: time.Now(),
    }
    t.Log(9)
    offset, err := client.GetOffset(msg.Source)
    if (err != nil) {
        t.Error(err)
    }
    if (offset == 9) {
        t.Log("Good!")
    } else {
        t.Errorf("did not match, we need:\n %v\n but got \n %v", 9, offset)

    }
}
func Test_SetFileStat_1(t *testing.T) {
    client, err := New(redisurl, redisport, 0)
    if (err != nil) {
        t.Error(err)
    }
    defer client.Close()

    //msg := logger.Message{
    //    Line:      []byte("safdfa"),
    //    Source:    "filename",
    //    Timestamp: time.Now(),
    //    //Partial:   false,
    //    //Stat_t: Stat_t{
    //    //    Dev:  333,
    //    //    Ino:  222,
    //    //    Size: 333,
    //    //},
    //}
    //client.SetFileStat(msg.Source, msg.Stat_t)
}
func Test_GetFileStat_1(t *testing.T) {
    client, err := New(redisurl, redisport, 0)
    if (err != nil) {
        t.Error(err)
    }
    defer client.Close()

    //msg := logger.Message{
    //    Line:      []byte("safdfa"),
    //    Source:    "filename",
    //    Timestamp: time.Now(),
    //    //Partial:   false,
    //    //Stat_t: Stat_t{
    //    //    Dev:  333,
    //    //    Ino:  222,
    //    //    Size: 333,
    //    //},
    //}
    //res, _ := client.GetFileStat(msg.Source)
    //if (msg.Stat_t == res) {
    //    t.Log("OK")
    //} else {
    //    t.Errorf("did not match, we need:\n %v\n but got \n %v", msg.Stat_t, res)
    //}
}
