package worker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/pandeptwidyaop/golog"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/mockstruct"
	"smartbtw.com/services/profile/models"
)

func HandleWorkerDelivery() {

	deliveryList := []mockstruct.GenerateProgressReportMessage{}
	cmds, err := db.NewRedisCluster().TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		pipe.LRange(context.Background(), db.RAPORT_REDIS_QUEUE_BUILDY_KEY, 0, 0)
		pipe.LTrim(context.Background(), db.RAPORT_REDIS_QUEUE_BUILDY_KEY, 1, -1)
		return nil
	})

	if err != nil {
		logrus.Error(err)
	}

	resultStr := cmds[0].String()
	startIndex := strings.Index(resultStr, "[")
	endIndex := strings.LastIndex(resultStr, "]")
	if startIndex != -1 && endIndex != -1 {
		jsonArrayStr := strings.ReplaceAll(resultStr[startIndex:endIndex+1], "} {", "},{")
		if err := sonic.Unmarshal([]byte(jsonArrayStr), &deliveryList); err != nil {
			fmt.Println("Unmarshal Error")
			fmt.Println(err.Error())
			// return err
		}
	}

	fmt.Println("Total to Delivery :", len(deliveryList))

	if len(deliveryList) > 0 {

		g, _ := errgroup.WithContext(context.Background())

		for _, k := range deliveryList {
			ns := k
			g.Go(func() error {
				err := lib.BuildProgressRaport(ns)
				if err != nil {
					lib.StoreFailed(ns)
				}
				return err
			})
		}

		go func() {
			err := g.Wait()
			if err != nil {
				golog.Slack.Error(fmt.Sprintf("One build progress raport encountered an error: %s", err.Error()), err)
				return
			}
		}()

		if err := g.Wait(); err != nil {
			return
		}
	}

	time.Sleep(1 * time.Second)
}

func HandleWorkerRaportDelivery() {

	deliveryList := []models.GetResultRaportBody{}
	cmds, err := db.NewRedisCluster().TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		pipe.LRange(context.Background(), db.RAPORT_RESULT_REDIS_QUEUE_BUILDY_KEY, 0, 0)
		pipe.LTrim(context.Background(), db.RAPORT_RESULT_REDIS_QUEUE_BUILDY_KEY, 1, -1)
		return nil
	})

	if err != nil {
		logrus.Error(err)
	}

	resultStr := cmds[0].String()
	startIndex := strings.Index(resultStr, "[")
	endIndex := strings.LastIndex(resultStr, "]")
	if startIndex != -1 && endIndex != -1 {
		jsonArrayStr := strings.ReplaceAll(resultStr[startIndex:endIndex+1], "} {", "},{")
		if err := sonic.Unmarshal([]byte(jsonArrayStr), &deliveryList); err != nil {
			fmt.Println("Unmarshal Error")
			fmt.Println(err.Error())
			// return err
		}
	}

	fmt.Println("Total to Delivery :", len(deliveryList))

	if len(deliveryList) > 0 {

		g, _ := errgroup.WithContext(context.Background())

		for _, k := range deliveryList {
			ns := k
			g.Go(func() error {
				err := lib.BuildRaportByTaskID(uint(ns.SmartbtwID), ns.Program, ns.TaskID)
				if err != nil {
					lib.StoreFailedRaport(ns)
				}
				return err
			})
		}

		go func() {
			err := g.Wait()
			if err != nil {
				golog.Slack.Error(fmt.Sprintf("One build raport encountered an error: %s", err.Error()), err)
				return
			}
		}()

		if err := g.Wait(); err != nil {
			return
		}
	}

	time.Sleep(1 * time.Second)
}
