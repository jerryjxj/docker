package saptail

import (
    "bytes"
    "compress/zlib"
    "io"
    "log"
    redis "github.wdf.sap.corp/sigma-anywhere/saptail/redis"
    "github.com/docker/docker/daemon/logger"
)

// MetricSet type defines all fields of the MetricSet
// As a minimum it must inherit the mb.BaseMetricSet fields, but can be extended with
// additional entries. These variables can be used to persist data or configuration between
// multiple fetch calls.
//type MetricSet struct {
//    Paths       []string
//    Client      *redis.Conn
//    //listfile    []string
//    version     string
//    landscape   string
//    namespace   string
//    servicename string
//    Type        string
//    //RPCserver   string
//    Zip         bool
//}

type Saptail struct {
    Redis redis.RedisInterface
}
type SaptailInterface interface {
    Message(msg *logger.Message, lazymode bool) (error)
    Close() (error)
}

func LoadConfig() {

}
func checkErr(err error) {
    if err != nil {
        log.Printf("%v", err)
    }
}
func DoZlibCompress(src []byte) []byte {
    var in bytes.Buffer
    w := zlib.NewWriter(&in)
    w.Write(src)
    w.Close()
    return in.Bytes()
}
func DoZlibUnCompress(compressSrc []byte) []byte {
    b := bytes.NewReader(compressSrc)
    var out bytes.Buffer
    r, _ := zlib.NewReader(b)
    io.Copy(&out, r)
    return out.Bytes()
}

func New(redisurl, redisport string, DB int) (SaptailInterface, error) {
    rd, err := redis.New(redisurl, redisport, DB)
    return &Saptail{Redis: rd}, err
}

func (st *Saptail) Message(msg *logger.Message, lazymode bool) (error) {

    var err error
    if (lazymode) {
        // lazy to send msg, buffer 10M, per second send once
        err = st.Redis.SendMessage(msg, true)
    } else {
        err = st.Redis.SendMessage(msg, false)
    }
    checkErr(err)
    return err
}
func (st *Saptail) Close() (error) {
    return st.Redis.Close()
}
