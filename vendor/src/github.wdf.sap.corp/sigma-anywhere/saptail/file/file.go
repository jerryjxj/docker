package file

import (
    "os"

    "bufio"
    "io"
    "reflect"
    //"encoding/json"
    "time"
    //"strings"
    //"github.com/elastic/beats/libbeat/outputs/mode/modetest"
    redis "github.wdf.sap.corp/sigma-anywhere/saptail/redis"
    "github.wdf.sap.corp/sigma-anywhere/saptail"
    "fmt"
    "errors"
)

type MapStr map[string]interface{}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

type Stat_t redis.Stat_t

type fileReader struct {
    file   string
    offset int64
}
type FStat struct {
    filename string
    stat     Stat_t
    offset   int64
}

func (f *fileReader) Read(p []byte) (n int, err error) {
    reader, err := os.Open(f.file)
    defer reader.Close()
    if err != nil {
        return 0, err
    }
    reader.Seek(f.offset, 0)

    n, err = reader.Read(p)

    if err == io.EOF {
        time.Sleep(1 * time.Second)
    }
    f.offset += int64(n)
    return n, err
}

func GetFileStat(filename string) Stat_t {
    //logp.Debug("saptail", "getFileStat %v", filename)
    finfo, _ := os.Stat(filename)
    ffield := reflect.ValueOf(finfo.Sys()).Elem()
    return Stat_t{
        Dev:ffield.FieldByName("Dev").Uint(),
        Mode:uint16(ffield.FieldByName("Mode").Uint()),
        Nlink:uint16(ffield.FieldByName("Nlink").Uint()),
        Ino:uint64(ffield.FieldByName("Ino").Uint()),
        Uid:uint32(ffield.FieldByName("Uid").Uint()),
        Gid:uint32(ffield.FieldByName("Gid").Uint()),
        Rdev:ffield.FieldByName("Rdev").Uint(),
        //Pad_cgo_0: ffield.FieldByName("Pad_cgo_0").Interface(),
        Atimespec:ffield.FieldByName("Atim").Field(0).Int(),
        Mtimespec:ffield.FieldByName("Mtim").Field(0).Int(),
        Ctimespec:ffield.FieldByName("Ctim").Field(0).Int(),
        Size:ffield.FieldByName("Size").Int(),
        Blocks:ffield.FieldByName("Blocks").Int(),
        Blksize:ffield.FieldByName("Blksize").Int(),
        //Flags:uint32(ffield.FieldByName("Flags").Uint()),
    }
}

func MonitorFile(rdclient *redis.Redis, filename string, zip bool) {
    go monitorFile(rdclient, filename, zip)
}
func monitorFile(rdclient *redis.Redis, filename string, zip bool) {
    var err error
    offset, err := rdclient.GetOffset(filename)
    checkErr(err)
    var lines string
    file := &fileReader{filename, offset}
    br := bufio.NewReader(file)
    for {
        log, _, err := br.ReadLine()
        //logp.Debug("readline", "%v", string(log))
        if err == io.EOF {
            break
        }

        if err != nil {
            //logp.Err("err: %v", err)
            continue
        }
        if len(lines) == 0 {
            lines = string(log)
        } else {
            lines = lines + "\n" + string(log)
        }
    }
    var message interface{}

    if zip {
        message = saptail.DoZlibCompress([]byte(lines))
    } else {
        message = lines
    }
    e := MapStr{
        "@timestamp":       time.Now().UTC(),
        "source":        file.file,
        "message":        message,
        "offset":         offset,
        "zip":        zip,
        //"type":        m.Type,
        //"version":        m.version,
        //"namespace":        m.namespace,
        //"landscape":    m.landscape,
        //"servicename":        m.servicename,
        //"bitid":        m.bitid,
        //"starttime":    m.starttime,
    }
    rdclient.SetOffset(filename, file.offset)
    fmt.Println(e)
}
