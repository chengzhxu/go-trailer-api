package testing

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"go-trailer-api/pkg/util"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	c := make(chan int)
	go func() {
		fmt.Println("ready to send in g1")
		c <- 1
		fmt.Println("send 1 to chan")
		fmt.Println("goroutine start sleep 1 second")
		time.Sleep(time.Second)
		fmt.Println("goroutine end sleep")
		c <- 2
		fmt.Println("send 2 to chan")
	}()

	fmt.Println("main thread start sleep 1 second")
	time.Sleep(time.Second)
	fmt.Println("main thread and sleep")
	i := <-c
	fmt.Printf("receive %d\n", i)
	i = <-c
	fmt.Printf("receive %d\n", i)
	time.Sleep(time.Second)
}

func TestDefer(t *testing.T) {
	defer func() { fmt.Println("1") }()
	defer func() { fmt.Println("2") }()
	defer func() { fmt.Println("3") }()

	panic("catch")
}

type student struct {
	Name string
	Age  int
}

func TestBird(t *testing.T) {
	m := make(map[int]*student)
	stus := []student{
		{Name: "Zhao", Age: 12},
		{Name: "Qian", Age: 18},
		{Name: "Sun", Age: 30},
	}

	for _, stu := range stus {
		m[stu.Age] = &stu
		fmt.Println(stu.Age)
	}

	for _, v := range m {
		//fmt.Printf("%s", k)
		fmt.Println(v.Age)
	}
}

var client *elastic.Client
var host = "http://127.0.0.1:9200/"

type Employee struct {
	Level      string `json:"level"`
	Timestamp  string `json:"timestamp"`
	Caller     int    `json:"caller"`
	Msg        string `json:"msg"`
	Data       string `json:"data"`
	Stacktrace string `json:"stacktrace"`
}

func TestEs(t *testing.T) {
	errorlog := log.New(os.Stdout, "APP", log.LstdFlags)
	var err error
	client, err = elastic.NewClient(elastic.SetErrorLog(errorlog), elastic.SetURL(host))
	if err != nil {
		panic(err)
	}
	info, code, err := client.Ping(host).Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	esversion, err := client.ElasticsearchVersion(host)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)

	contents := util.GetLogContents()
	for _, line := range contents {

		fmt.Println(line)
		_, err := client.Index().
			Index("kevin").
			Type("go-bird").
			Id(strconv.FormatInt(time.Now().UnixNano()/1e6, 10)).
			BodyJson(line).
			Do(context.Background())
		if err != nil {
			fmt.Println(err)
		}

		//time.Sleep(time.Second)

		//fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put1.Id, put1.Index, put1.Type)
	}

	//e1 := `{"level":"error","timestamp":"2022-03-10T17:59:02+0800","caller":"logging/log.go:53","msg":"","data":"Error 1062: Duplicate entry 'kevin' for key 'inx_username'","stacktrace":"go-trailer-api/pkg/logging.Error\n\t/Users/kevin/Documents/code/GoProjects/go-trailer-api/pkg/logging/log.go:53\ngo-trailer-api/pkg/model/bird/userModel.AddUser\n\t/Users/kevin/Documents/code/GoProjects/go-trailer-api/pkg/model/bird/userModel/user.go:42\ngo-trailer-api/pkg/service/bird/userService.User.Add\n\t/Users/kevin/Documents/code/GoProjects/go-trailer-api/pkg/service/bird/userService/user.go:28\ngo-trailer-api/routers/bird/user.AddUser\n\t/Users/kevin/Documents/code/GoProjects/go-trailer-api/routers/bird/user/user.go:48\ngithub.com/gin-gonic/gin.(*Context).Next\n\t/Users/kevin/go/pkg/mod/github.com/gin-gonic/gin@v1.6.3/context.go:161\ngithub.com/gin-gonic/gin.RecoveryWithWriter.func1\n\t/Users/kevin/go/pkg/mod/github.com/gin-gonic/gin@v1.6.3/recovery.go:83\ngithub.com/gin-gonic/gin.(*Context).Next\n\t/Users/kevin/go/pkg/mod/github.com/gin-gonic/gin@v1.6.3/context.go:161\ngithub.com/gin-gonic/gin.LoggerWithConfig.func1\n\t/Users/kevin/go/pkg/mod/github.com/gin-gonic/gin@v1.6.3/logger.go:241\ngithub.com/gin-gonic/gin.(*Context).Next\n\t/Users/kevin/go/pkg/mod/github.com/gin-gonic/gin@v1.6.3/context.go:161\ngithub.com/gin-gonic/gin.(*Engine).handleHTTPRequest\n\t/Users/kevin/go/pkg/mod/github.com/gin-gonic/gin@v1.6.3/gin.go:409\ngithub.com/gin-gonic/gin.(*Engine).ServeHTTP\n\t/Users/kevin/go/pkg/mod/github.com/gin-gonic/gin@v1.6.3/gin.go:367\nnet/http.serverHandler.ServeHTTP\n\t/usr/local/opt/go/libexec/src/net/http/server.go:2887\nnet/http.(*conn).serve\n\t/usr/local/opt/go/libexec/src/net/http/server.go:1952"}`
	//put1, err := client.Index().
	//	Index("kevin").
	//	Type("go-bird").
	//	Id(strconv.Itoa(util.GetNowTimeStamp())).
	//	BodyJson(e1).
	//	Do(context.Background())
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put1.Id, put1.Index, put1.Type)
}

func TestWriteEs(t *testing.T) {
	//contents := util.GetLogContents()

	//yesTime := time.Now().AddDate(0, 0, -1)
	//yesDate := yesTime.Format("2006-01-02")

	//for _, line := range contents {
	//
	//}

}
