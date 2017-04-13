package redis

import (
    "github.com/garyburd/redigo/redis"
    "github.com/pkg/errors"
    "github.com/docker/docker/daemon/logger"
    "log"
)

type Redis struct {
    Url    string
    Port   string
    DB     int
    Client redis.Conn
}

type RedisInterface interface {
    GetOffset(filename string, params ...int) (int64, error)
    SetOffset(filename string, Offset uint64, params ...int) (error)
    GetFileStat(filename string, params ...int) (Stat_t, error)
    SetFileStat(filename string, stat Stat_t, params ...int) (error)
    SendMessage(msg *logger.Message, lazymode bool) (error)
    HealthCheck() (string, error)
    Close() (error)
}
type Stat_t struct {
    Dev       uint64
    Ino       uint64
    Nlink     uint16
    Mode      uint16
    Uid       uint32
    Gid       uint32
    Rdev      uint64
    Size      int64
    Blksize   int64
    Blocks    int64
    Atimespec int64
    Mtimespec int64
    Ctimespec int64
}

//type LogAttributes map[string]string
//type Message struct {
//    Line      []byte
//    Source    string
//    Timestamp time.Time
//    Attrs     LogAttributes
//    Partial   bool
//    Offset    uint64
//    Stat_t    Stat_t
//}

func checkErr(err error, msg ...interface{}) {
    if err != nil {
        log.Printf("%v", err)
        log.Printf("%v", msg)
    }
}
func New(redisurl, redisport string, DB int) (RedisInterface, error) {
    client, err := redis.Dial("tcp", redisurl+":"+redisport)
    if err != nil {
        return nil, err
    }
    client.Do("SELECT", DB)
    var R RedisInterface
    R = &Redis{
        Client: client,
        Url:    redisurl,
        Port:   redisport,
        DB:     DB,
    }
    r, err := R.HealthCheck()
    if (r == "PONG") {
        return R, nil
    } else {
        return nil, err
    }
}

func (rdclient *Redis) GetOffset(filename string, params ...int) (int64, error) {
    var DB int
    if (len(params) == 0) {
        DB = 2
    } else {
        DB = params[0]
    }
    rdclient.Client.Do("SELECT", DB)

    reply, err := rdclient.Client.Do("GET", filename)
    checkErr(err)
    if (reply == nil) {
        return 0, errors.New("file does not exist")
    } else {
        re, err := redis.Int64(reply, nil)
        return re, err
    }
}
func (rdclient *Redis) SetOffset(filename string, Offset uint64, params ...int) (error) {
    var DB int
    if (len(params) == 0) {
        DB = 2
    } else {
        DB = params[0]
    }
    rdclient.Client.Do("SELECT", DB)
    rdclient.readmeOffset(DB)

    _, err := rdclient.Client.Do("SET", filename, Offset)
    return err

}

func (rdclient *Redis) GetFileStat(filename string, params ...int) (Stat_t, error) {
    var DB int
    if (len(params) == 0) {
        DB = 3
    } else {
        DB = params[0]
    }
    rdclient.Client.Do("SELECT", DB)

    reply, err := redis.Values(rdclient.Client.Do("HGETALL", filename))
    var stat Stat_t
    err = redis.ScanStruct(reply, &stat)
    checkErr(err)

    return stat, err
}
func (rdclient *Redis) SetFileStat(filename string, stat Stat_t, params ...int) (error) {
    var DB int

    if (len(params) == 0) {
        DB = 3
    } else {
        DB = params[0]
    }
    rdclient.Client.Do("SELECT", DB)
    rdclient.readmeFileStat(DB)

    _, err := rdclient.Client.Do("HMSET", redis.Args{filename}.AddFlat(&stat)...)
    return err
}

func (rdclient *Redis) HealthCheck() (string, error) {
    res, err := rdclient.Client.Do("ping")
    s, err := redis.String(res, err)
    return s, err
}

func (rdclient *Redis) Close() (error) {
    return rdclient.Client.Close()
}

func (rdclient *Redis) SendMessage(msg *logger.Message, lazymode bool) (error) {
    var err error
    //rdclient.Client.Do("SELECT", rdclient.DB)
    if lazymode {

    } else {
        _, err = rdclient.Client.Do("LPUSH", msg.Attrs["pod"], msg.Line)
        checkErr(err, msg)
        //rdclient.SetOffset(msg.Source, msg.Offset)
        //rdclient.SetFileStat(msg.Source, msg.Stat_t)
    }
    return err
}
func (rdclient *Redis) readmeFileStat(DB int) {
    rdclient.Client.Do("SELECT", DB)
    readme := "This db is for FileStat, you can use `HGETALL <KEY>`"
    reply, err := redis.String(rdclient.Client.Do("GET", "README"))
    if (err != nil) {
        rdclient.Client.Do("SET", "README", readme)
    } else if (reply != readme) {
        rdclient.Client.Do("SET", "README", readme)
    }
}
func (rdclient *Redis) readmeOffset(DB int) {
    rdclient.Client.Do("SELECT", DB)
    reply, err := redis.String(rdclient.Client.Do("GET", "README"))
    readme := "This db is for Offset, you can use `GET <KEY>`"
    if (err != nil) {
        rdclient.Client.Do("SET", "README", readme)
    } else if (reply != readme) {
        rdclient.Client.Do("SET", "README", readme)
    }
}
