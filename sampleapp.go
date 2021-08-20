package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	pb "git.local/go-app/models"
	"git.local/go-app/repository"
	"git.local/go-app/services"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type Config struct {
	HTTPPort string `default:":8443"`
	GrpcPort string `default:":50051"`
}

type Sampleapp struct {
	*sync.Mutex
	Cf        Config
	stopchan  chan struct{} // signal to stop scheduling
	isrunning bool
	pb.UnimplementedSampleAPIServer
	ordermgr *Ordermgr
}

func NewSampleapp() *Sampleapp {
	var cf Config
	envconfig.MustProcess("sampleapp", &cf)
	app := &Sampleapp{
		Mutex:    &sync.Mutex{},
		Cf:       cf,
		stopchan: make(chan struct{}),
		ordermgr: NewOrdermgr(repository.NewDB()),
	}
	return app
}

func (app *Sampleapp) ServeHTTP() {
	fmt.Println("ServeHTTP at localhost" + app.Cf.HTTPPort)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong from sampleapp",
		})
	})
	r.GET("/songs", func(c *gin.Context) {
		db := repository.NewPostgreSQLDB()
		query := ""
		if c.Query("query") != "" {
			query = c.Query("query")
		}
		resp, err := db.ListSong(query)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp)
		jsonPessoal, errr := json.Marshal(resp)
		if errr != nil {
			log.Fatal(errr)
		}
		c.JSON(200, gin.H{
			"code":  200,
			"error": "0",
			"data":  string(jsonPessoal),
		})
	})
	r.GET("/album", func(c *gin.Context) {
		db := repository.NewPostgreSQLDB()
		resp, err := db.ListAlbum()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp)
		jsonPessoal, errr := json.Marshal(resp)
		if errr != nil {
			log.Fatal(errr)
		}
		c.JSON(200, gin.H{
			"code":  200,
			"error": "0",
			"data":  string(jsonPessoal),
		})
	})
	r.GET("/album/:id", func(c *gin.Context) {
		db := repository.NewPostgreSQLDB()
		id := c.Param("id")
		_id, _ := strconv.Atoi(id)
		resp, songs, err := db.ReadAlbum(_id)
		if err != nil {
			log.Fatal(err)
		}
		_, errr := json.Marshal(resp)
		jsonPessoal2, errr := json.Marshal(songs)
		if errr != nil {
			log.Fatal(errr)
		}
		c.JSON(200, gin.H{
			"code":  200,
			"error": "0",
			"data":  string(jsonPessoal2),
		})
	})
	r.GET("/category/:id", func(c *gin.Context) {
		db := repository.NewPostgreSQLDB()
		id := c.Param("id")
		_id, _ := strconv.Atoi(id)
		resp, songs, err := db.ReadCategory(_id)
		fmt.Println(resp, songs, "test detail")
		if err != nil {
			log.Fatal(err)
		}
		_, errr := json.Marshal(resp)
		jsonPessoal2, errr := json.Marshal(songs)
		if errr != nil {
			log.Fatal(errr)
		}
		c.JSON(200, gin.H{
			"code":  200,
			"error": "0",
			"data":  string(jsonPessoal2),
		})
	})
	r.POST("/upload", services.Upload)
	r.StaticFS("/file", http.Dir("public"))
	r.POST("/category", func(c *gin.Context) {
		db := repository.NewPostgreSQLDB()
		resp, err := db.ListCategory()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp)
		jsonPessoal, errr := json.Marshal(resp)
		if errr != nil {
			log.Fatal(errr)
		}
		c.JSON(200, gin.H{
			"code":  200,
			"error": "0",
			"data":  string(jsonPessoal),
		})
	})

	r.POST("/user", func(c *gin.Context) {
		data := pb.Account{
			Link: "aaa",
			Pic:  "https://image.shutterstock.com/image-photo/miniature-greenhouse-concept-alone-mini-260nw-1176115702.jpg",
			Type: "11",
			Name: "Nhạc nhẹ",
		}

		data1 := pb.Account{
			Link: "aaa",
			Pic:  "https://image.shutterstock.com/image-photo/miniature-greenhouse-concept-alone-mini-260nw-1176115702.jpg",
			Type: "111111111111 ffff",
			Name: "Tuyển tập nhạc của trần điệp",
		}
		resp := [5]pb.Account{data, data1, data, data, data}
		fmt.Println(resp)
		jsonPessoal, errr := json.Marshal(resp)
		if errr != nil {
			log.Fatal(errr)
		}
		c.JSON(200, gin.H{
			"code":  200,
			"error": "0",
			"data":  string(jsonPessoal),
		})
	})
	x509, errTls := tls.LoadX509KeyPair(os.Getenv("SSLCRT"), os.Getenv("SSLKEY"))
	if errTls != nil {
		fmt.Println(errTls)
	}
	fmt.Println(os.Getenv("SSLCRT"), "SSLKEY")
	var server *http.Server
	server = &http.Server{
		Addr:    app.Cf.HTTPPort,
		Handler: r,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{x509},
		},
	}
	server.ListenAndServeTLS("", "")
}

func (app *Sampleapp) BatchAsync() chan struct{} {
	fmt.Println("BATCHASYNC-STARTTT")
	app.Lock()
	defer app.Unlock()
	if app.isrunning {
		return app.stopchan
	}
	app.isrunning = true
	ticker := time.NewTicker(5 * time.Second)
	var tickerCount int64
	go func() {
		for {
			select {
			case <-ticker.C:
				atomic.AddInt64(&tickerCount, 1)
				go func() { fmt.Println("TODO", tickerCount) }()
			case <-app.stopchan:
				ticker.Stop()
				app.isrunning = false
				fmt.Println("BATCHASYNC-ENDDD")
				return
			}
		}
	}()

	return app.stopchan
}

func (app *Sampleapp) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (app *Sampleapp) ServeGrpc() {
	fmt.Println("ServeGrpc at localhost" + app.Cf.GrpcPort)
	lis, err := net.Listen("tcp", app.Cf.GrpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSampleAPIServer(s, app)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (app *Sampleapp) CreateOrder(ctx context.Context, in *pb.Order) (*pb.Order, error) {
	// TODO perm
	out, err := app.ordermgr.CreateOrder(in)
	if err != nil {
		// TODO errconv
		return nil, err
	}
	return out, err
}

func (app *Sampleapp) DeleteOrder(context.Context, *pb.Id) (*pb.Empty, error) {
	// TODO
	return nil, nil
}

func (app *Sampleapp) ListOrders(context.Context, *pb.Id) (*pb.Orders, error) {
	// TODO
	return nil, nil
}

func (app *Sampleapp) UpdateOrder(context.Context, *pb.Order) (*pb.Order, error) {
	// TODO
	return nil, nil
}

func (app *Sampleapp) ReadOrder(ctx context.Context, in *pb.Id) (*pb.Order, error) {
	// TODO perm
	out, err := app.ordermgr.ReadOrder(in.GetId())
	if err != nil {
		// TODO errconv
		return nil, err
	}
	return out, err
}
