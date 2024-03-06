package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/asaskevich/govalidator"
	pb "github.com/btwedutech/grpc/service/profile"
	"github.com/joho/godotenv"
	"github.com/pandeptwidyaop/golog"
	"github.com/pandeptwidyaop/gorabbit"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"smartbtw.com/services/profile/config"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/grpchandler"
	"smartbtw.com/services/profile/helpers"
	"smartbtw.com/services/profile/lib/wghttp"
	"smartbtw.com/services/profile/lib/worker"
	"smartbtw.com/services/profile/listener"
	"smartbtw.com/services/profile/server"
)

func main() {
	InitEnv()
	InitSlackLog()
	InitDB()
	InitGoValidator()
	listenPort := ":4000"

	useArgs := false
	if len(os.Args) > 1 {
		useArgs = true
	}

	if !useArgs {
		wg := new(sync.WaitGroup)
		wg.Add(2)

		app := server.SetupFiber()

		go func() {
			log.Fatal(app.Listen(listenPort))
			wg.Done()
		}()

		go func() {
			InitRabbitMQ()
			Listen()
			wg.Done()
		}()

		golog.Slack.Info("Service Up & Running")

		wg.Wait()
	} else {
		// Run server on port 4000
		appName := os.Getenv("APP_NAME")
		switch strings.ToLower(os.Args[1]) {
		case "grpc-only":

			InitRabbitMQ()

			err := db.Broker.Connect()
			if err != nil {
				panic(err)
			}

			log.Println("intializing gRPC server")
			lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config.GetGrpcServerHost(), config.GetGrpcServerPort()))
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
			log.Printf("configuring gRPC on host %s and port %s \n", config.GetGrpcServerHost(), config.GetGrpcServerPort())
			log.Println("initializing TLS certificate")
			cred, err := credentials.NewServerTLSFromFile(config.GetGrpcServerCertificatePath(), config.GetGrpcServerKeyPath())

			if err != nil {
				log.Fatalf("failed to load certificate: %v", err)
			}
			// Initialize the wait group
			log.Println("initializing wait group")
			wghttp.NewHttpWg()

			// Create a gRPC server object
			log.Println("initializing gRPC server with TLS")
			s := grpc.NewServer(grpc.Creds(cred))

			// Register the service or delivery
			log.Println("registering gRPC service")
			pb.RegisterProfileServer(s, &grpchandler.ProfileDelivery{})

			// Register reflection service on gRPC server.
			reflection.Register(s)
			log.Println("registering reflection service on gRPC server")

			// Start the server
			go func() {
				log.Println("gRPC server is running")
				err := s.Serve(lis)
				if err != nil {
					log.Fatalf("failed to serve: %v", err)
				}
			}()

			// Graceful shutdown
			interrupt := make(chan os.Signal, 1)
			signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
			<-interrupt

			wghttp.HttpWG.Wait()

			log.Println("gRPC server is shutting down")
			s.GracefulStop()
		case "http-only":
			InitRabbitMQ()

			err := db.Broker.Connect()
			if err != nil {
				panic(err)
			}

			// wg := new(sync.WaitGroup)
			// wg.Add(1)
			app := server.SetupFiber()

			go func() {
				golog.Slack.Info(fmt.Sprintf("%s: HTTP-Only Service Started (Ver: %s)", appName, config.ServiceVersion))
				log.Fatal(app.Listen(listenPort))
			}()

			interrupt := make(chan os.Signal, 1)
			signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
			<-interrupt

			wghttp.HttpWG.Wait()

			golog.Slack.Info(fmt.Sprintf("%s: HTTP-Only Service Stopped (Ver: %s)", appName, config.ServiceVersion))
			app.Shutdown()

			// go func() {
			// golog.Slack.Info(fmt.Sprintf("%s: HTTP-Only Service Started", appName))
			// log.Fatal(app.Listen(listenPort))
			// wg.Done()
			// }()

			// wg.Wait()
		case "consume-only":
			InitRabbitMQ()
			golog.Slack.Info(fmt.Sprintf("%s: Consume-Only Service Started (Ver: %s)", appName, config.ServiceVersion))
			Listen()

		case "worker-process":
			InitRabbitMQ()
			err := db.Broker.Connect()
			if err != nil {
				panic(err)
			}

			golog.Slack.Info(fmt.Sprintf("%s: Worker Started", appName))
			log.Println("Start Queue Process worker")
			exit := make(chan bool)

			go startQueueWorker(exit)

			interrupt := make(chan os.Signal, 1)
			signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
			<-interrupt
			log.Println("Shutting down worker")
			exit <- true
			log.Println("Worker is shutted down")

		case "worker-raport-process":
			InitRabbitMQ()
			err := db.Broker.Connect()
			if err != nil {
				panic(err)
			}

			golog.Slack.Info(fmt.Sprintf("%s: Worker Build Raport Started", appName))
			log.Println("Start Queue Process worker build raport")
			exit := make(chan bool)

			go startQueueWorkerRaport(exit)

			interrupt := make(chan os.Signal, 1)
			signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
			<-interrupt
			log.Println("Shutting down worker raport")
			exit <- true
			log.Println("Worker build raport is shutted down")

		default:
			panic("argument invalid")

		}
	}
}

func InitEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not detected using global variable")
	}
}

func InitSlackLog() {
	golog.New()
}

func InitDB() {
	InitMongoDB()
	db.NewElastic()
	db.NewFirebaseData()
	db.NewRedisCluster()
}

func InitMongoDB() {
	connection := os.Getenv("MONGODB_CONNECTION")
	database := os.Getenv("MONGODB_DATABASE")
	if connection == "" {
		e := errors.New("undefined MONGODB_CONNECTION")
		golog.Slack.Error("Undefined MONGODB_CONNECTION", e)
		log.Fatal(e)
	}
	if database == "" {
		e := errors.New("undefined MONGODB_DATABASE")
		golog.Slack.Error("Undefined MONGODB_DATABASE", e)
		log.Fatal(e)
	}
	err := db.Connect(connection, database)

	if err != nil {
		golog.Slack.Error("Error when initialize mongodb", err)
		log.Fatal(err)
	}
}

func InitRabbitMQ() {
	var err error
	conn := os.Getenv("RABBITMQ_CONNECTION")
	app := os.Getenv("APP_NAME")
	if conn == "" {
		golog.Slack.Error("Rabbit MQ URL Connection not set", nil)
		log.Panic("no varible found for RABBITMQ_CONNECTION")
	}

	if app == "" {
		log.Panic("app name not initialize")
	}

	db.Broker, err = gorabbit.New(conn, app, "GLOBAL_X")

	if err != nil {
		golog.Slack.Error("Error when create new connection to Rabbit MQ server", err)
		log.Panic(err)
	}
}

func Listen() {
	forever := make(chan bool)
	err := db.Broker.Connect()

	if err != nil {
		panic(err)
	}

	err = db.Broker.Bind([]string{
		"user.import",
		"user.created",
		"user.updated",
		"user.deleted",
		"user.create-elastic",
		"branch.created",
		"branch.updated",
		"wallet.created",
		"wallet.received",
		"wallet.cutting-masa-ai",
		"student.target.created",
		"history-ptk.created",
		"history-ptn.created",
		"history-assessment.created",
		"history-cpns.created",
		"history-cpns.time-consume.update",
		"user.data.updated",
		"school.updated",
		"study.program.updated",
		"target.score.updated",
		"wallet-history-invite.created",
		"wallet-history-premium.created",
		"student.target.updated",
		"profile.syncResult",
		"class-member.created",
		"class-member.updated",
		"class-member.deleted",
		"class-member.switch",
		"classroom.created",
		"classroom.updated",
		"user.upsert-profile-elastic",
		"user.upsert-compmap-elastic",
		"user.binsus.sync",
		"user.binsus.final-sync",
		"result-raport.generated",
		"progress-result-raport.generated",
		"raport-ptk.build",
		"raport-ptn.build",
		"raport-cpns.build",
		"result-raport.build-bulk.request",
		"progress-result-raport.build.queue",
	})

	if err != nil {
		log.Panic(err)
	}

	ds, err := db.Broker.Consume()
	if err != nil {
		log.Panic(err)
	}
	log.Println("Waiting for messages")
	for q, d := range ds {
		go db.Broker.HandleConsumedDeliveries(q, d, handleConsume)
	}
	<-forever
}

func InitGoValidator() {
	govalidator.SetNilPtrAllowedByRequired(true)
}

func handleConsume(mq gorabbit.RabbitMQ, q string, deliveries <-chan amqp.Delivery) {
	for d := range deliveries {
		fmt.Println("Received a message: ", string(d.RoutingKey))
		listener.ListenStudentBinding(&d)
		listener.ListenBranchBinding(&d)
		listener.ListenWalletBinding(&d)
		listener.ListenWalletHistoryBinding(&d)
		listener.ListenStudentTargetBinding(&d)
		listener.ListenScoreHistoryBinding(&d)
		listener.ListenClassMemberBinding(&d)
		listener.ListenClassroomBinding(&d)
		listener.ListenResultRaportBinding(&d)
		listener.ListenProgressResultRaportBinding(&d)
	}
}

func startQueueWorker(exit chan bool) {
	limiter := helpers.NewLimiter(1, helpers.InSecond(db.GetBuildRaportWorkerProcessInSecond()))
	for {
		select {
		case <-exit:
			return
		default:
			allow, err := limiter.IsAllowed("profile:rate-limit:progress-raport-build")
			if err != nil {
				log.Println(err)
				continue
			}

			if allow {
				log.Println("Looking for item to deliver...")
				worker.HandleWorkerDelivery()
			}

			time.Sleep(1 * time.Second)
		}

	}
}

func startQueueWorkerRaport(exit chan bool) {
	limiter := helpers.NewLimiter(1, helpers.InSecond(db.GetBuildRaportResultWorkerProcessInSecond()))
	for {
		select {
		case <-exit:
			return
		default:
			allow, err := limiter.IsAllowed("profile:rate-limit:raport-build")
			if err != nil {
				log.Println(err)
				continue
			}

			if allow {
				log.Println("Looking for item to deliver...")
				worker.HandleWorkerRaportDelivery()
			}

			time.Sleep(1 * time.Second)
		}

	}
}
