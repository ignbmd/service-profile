package lib

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/lib/aggregates"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func CreateWalletHistoryPremiumPackage(c *request.CreateWalletHistoryPremium) (*mongo.InsertOneResult, error) {
	opts := options.Update().SetUpsert(true)
	whCol := db.Mongodb.Collection("wallet_histories")
	wCol := db.Mongodb.Collection("wallets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": c.SmartbtwID, "type": models.BONUS, "delete_at": nil}
	wModels := models.Wallet{}
	err := wCol.FindOne(ctx, filter).Decode(&wModels)
	if err != nil {
		return nil, fmt.Errorf("data not found")
	}

	point := (0.2 / 100) * c.Price

	payload := models.WalletHistory{
		WalletID:    wModels.ID,
		SmartbtwID:  c.SmartbtwID,
		Point:       float32(math.Ceil(float64(point))),
		Description: c.Description,
		Status:      string(models.IN),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}
	res, err1 := whCol.InsertOne(ctx, payload)
	if err1 != nil {
		return nil, err1
	}
	payload1 := models.Wallet{
		SmartbtwID: wModels.SmartbtwID,
		Point:      wModels.Point + float32(math.Ceil(float64(point))),
		Type:       wModels.Type,
		CreatedAt:  wModels.CreatedAt,
		UpdatedAt:  time.Now(),
		DeletedAt:  nil,
	}
	update := bson.M{"$set": payload1}
	_, err2 := wCol.UpdateByID(ctx, wModels.ID, update, opts)
	if err2 != nil {
		return nil, err2
	}

	return res, nil
}

func CreateWalletHistoryUKA(c *request.CreateWalletHistory) (*mongo.InsertOneResult, error) {
	opts := options.Update().SetUpsert(true)
	whCol := db.Mongodb.Collection("wallet_histories")
	wCol := db.Mongodb.Collection("wallets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": c.SmartbtwID, "type": models.BONUS, "delete_at": nil}
	wModels := models.Wallet{}
	err := wCol.FindOne(ctx, filter).Decode(&wModels)
	if err != nil {
		return nil, fmt.Errorf("data not found")
	}

	payload := models.WalletHistory{
		WalletID:    wModels.ID,
		SmartbtwID:  c.SmartbtwID,
		Point:       10,
		Description: c.Description,
		Status:      string(models.IN),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}
	res, err1 := whCol.InsertOne(ctx, payload)
	if err1 != nil {
		return nil, err1
	}
	payload1 := models.Wallet{
		SmartbtwID: wModels.SmartbtwID,
		Point:      wModels.Point + 10,
		Type:       wModels.Type,
		CreatedAt:  wModels.CreatedAt,
		UpdatedAt:  time.Now(),
		DeletedAt:  nil,
	}
	update := bson.M{"$set": payload1}
	_, err2 := wCol.UpdateByID(ctx, wModels.ID, update, opts)
	if err2 != nil {
		return nil, err2
	}

	return res, nil
}

func CreateWalletHistoryInvitePeople(c *request.CreateWalletHistory) (*mongo.InsertOneResult, error) {
	opts := options.Update().SetUpsert(true)
	whCol := db.Mongodb.Collection("wallet_histories")
	wCol := db.Mongodb.Collection("wallets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": c.SmartbtwID, "type": models.BONUS, "delete_at": nil}
	wModels := models.Wallet{}
	err := wCol.FindOne(ctx, filter).Decode(&wModels)
	if err != nil {
		return nil, fmt.Errorf("data not found")
	}

	payload := models.WalletHistory{
		WalletID:    wModels.ID,
		SmartbtwID:  c.SmartbtwID,
		Point:       100,
		Description: c.Description,
		Status:      string(models.IN),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}
	res, err1 := whCol.InsertOne(ctx, payload)
	if err1 != nil {
		return nil, err1
	}
	payload1 := models.Wallet{
		SmartbtwID: wModels.SmartbtwID,
		Point:      wModels.Point + 100,
		Type:       wModels.Type,
		CreatedAt:  wModels.CreatedAt,
		UpdatedAt:  time.Now(),
		DeletedAt:  nil,
	}
	update := bson.M{"$set": payload1}
	_, err2 := wCol.UpdateByID(ctx, wModels.ID, update, opts)
	if err2 != nil {
		return nil, err2
	}

	return res, nil
}

func ReceiveNewWalletPoint(c *request.ReceiveWallet) (*mongo.InsertOneResult, error) {
	opts := options.Update().SetUpsert(true)
	whCol := db.Mongodb.Collection("wallet_histories")
	wCol := db.Mongodb.Collection("wallets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": c.SmartbtwID, "type": c.Type, "delete_at": nil}
	wModels := models.Wallet{}
	err := wCol.FindOne(ctx, filter).Decode(&wModels)
	if err != nil {
		// return nil, fmt.Errorf("data not found")
		resCr, errCrWallet := CreateWallet(&request.CreateWallet{
			SmartbtwID: c.SmartbtwID,
			Point:      c.Point,
			Type:       c.Type,
		})
		if errCrWallet != nil {
			return nil, errCrWallet
		} else {
			return resCr, nil
		}
	}

	trType := "GENERAL"
	if c.TransactionType != "" {
		trType = c.TransactionType
	}

	payload := models.WalletHistory{
		WalletID:    wModels.ID,
		AmountPay:   c.AmountPay,
		SmartbtwID:  c.SmartbtwID,
		Point:       c.Point,
		Description: c.Description,
		Status:      string(models.IN),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Type:        trType,
		DeletedAt:   nil,
	}
	res, err1 := whCol.InsertOne(ctx, payload)
	if err1 != nil {
		return nil, err1
	}
	payload1 := models.Wallet{
		SmartbtwID: wModels.SmartbtwID,
		Point:      wModels.Point + c.Point,
		Type:       wModels.Type,
		CreatedAt:  wModels.CreatedAt,
		UpdatedAt:  time.Now(),
		DeletedAt:  nil,
	}
	update := bson.M{"$set": payload1}
	_, err2 := wCol.UpdateByID(ctx, wModels.ID, update, opts)
	if err2 != nil {
		return nil, err2
	}

	return res, nil
}

func CheckCoinDrillo(smartbtw_id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	wCol := db.Mongodb.Collection("wallets")
	filter := bson.M{"smartbtw_id": smartbtw_id, "type": models.BONUS, "deleted_at": nil}
	filterWd := bson.M{"smartbtw_id": smartbtw_id, "type": models.DEFAULT, "deleted_at": nil}
	wModels := models.Wallet{}
	wdModels := models.Wallet{}
	err := wCol.FindOne(ctx, filter).Decode(&wModels)
	if err != nil {
		return fmt.Errorf("data not found")
	}
	err = wCol.FindOne(ctx, filterWd).Decode(&wdModels)
	if err != nil {
		return fmt.Errorf("data not found")
	}
	value := ((wModels.Point + wdModels.Point) / 100) - 1
	if value == -1 {
		return fmt.Errorf("BTW coin anda tidak cukup, silahkan top up")
	}
	return nil
}

func CoinCuttingMasaAI(smartbtw_id int) error {
	opts := options.Update().SetUpsert(true)
	whCol := db.Mongodb.Collection("wallet_histories")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	wCol := db.Mongodb.Collection("wallets")
	filter := bson.M{"smartbtw_id": smartbtw_id, "type": models.BONUS, "deleted_at": nil}
	filterWd := bson.M{"smartbtw_id": smartbtw_id, "type": models.DEFAULT, "deleted_at": nil}
	wModels := models.Wallet{}
	wdModels := models.Wallet{}
	err := wCol.FindOne(ctx, filter).Decode(&wModels)
	if err != nil {
		return fmt.Errorf("data not found")
	}
	err = wCol.FindOne(ctx, filterWd).Decode(&wdModels)
	if err != nil {
		return fmt.Errorf("data not found")
	}
	valueBonus := (wModels.Point / 100)
	trType := "GENERAL"
	if valueBonus >= 1 {
		payloadBonus := models.WalletHistory{
			WalletID:    wModels.ID,
			SmartbtwID:  smartbtw_id,
			Point:       100,
			Description: "penggunaan drillo chat",
			Status:      string(models.OUT),
			Type:        trType,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		}
		_, errBns := whCol.InsertOne(ctx, payloadBonus)
		if errBns != nil {
			return errBns
		}
		payloadBnsWallet := models.Wallet{
			SmartbtwID: smartbtw_id,
			Point:      wModels.Point - 100,
			Type:       string(models.BONUS),
			CreatedAt:  wModels.CreatedAt,
			UpdatedAt:  time.Now(),
			DeletedAt:  nil,
		}
		fil := models.UpdateFilter{
			ID:   wModels.ID,
			Type: string(models.BONUS),
		}
		updateBnsW := bson.M{"$set": payloadBnsWallet}
		_, errBnsW := wCol.UpdateOne(ctx, fil, updateBnsW, opts)
		if errBnsW != nil {
			return errBnsW
		}
	} else {
		payloadBonus := models.WalletHistory{
			WalletID:    wModels.ID,
			SmartbtwID:  smartbtw_id,
			Point:       100,
			Description: "penggunaan drillo chat",
			Status:      string(models.OUT),
			Type:        trType,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		}
		_, errBns := whCol.InsertOne(ctx, payloadBonus)
		if errBns != nil {
			return errBns
		}
		payloadBnsWallet := models.Wallet{
			SmartbtwID: smartbtw_id,
			Point:      wdModels.Point - 100,
			Type:       string(models.DEFAULT),
			CreatedAt:  wModels.CreatedAt,
			UpdatedAt:  time.Now(),
			DeletedAt:  nil,
		}
		fil := models.UpdateFilter{
			ID:   wdModels.ID,
			Type: string(models.DEFAULT),
		}
		updateBnsW := bson.M{"$set": payloadBnsWallet}
		_, errBnsW := wCol.UpdateOne(ctx, fil, updateBnsW, opts)
		if errBnsW != nil {
			return errBnsW
		}
	}
	return nil
}

func GetStudentWalletHistory(SmartBTWID int, walletType *string) ([]bson.M, error) {
	var wallet []bson.M
	collection := db.Mongodb.Collection("wallets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pipel := aggregates.GetStudentWalletHistoryByWalletType(SmartBTWID, walletType)
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := collection.Aggregate(ctx, pipel, opts...)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &wallet)
	if err != nil {
		return nil, err
	}

	if wallet != nil {
		return wallet, nil
	}

	return []bson.M{}, nil
}

func ChargeStudentWallet(c *request.ChargeWallet) error {
	opts := options.Update().SetUpsert(true)
	whCol := db.Mongodb.Collection("wallet_histories")
	wCol := db.Mongodb.Collection("wallets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var walletData []models.Wallet

	pipel := aggregates.GetStudentWalletBalance(c.SmartbtwID, false)
	optsW := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := wCol.Aggregate(ctx, pipel, optsW...)
	if err != nil {
		return err
	}

	err = cursor.All(ctx, &walletData)
	if err != nil {
		return err
	}

	if walletData == nil {
		return errors.New("^user doesnt have wallet")
	}

	bonusDeduction := float32(0)
	walletDeduction := float32(0)
	isUseBonus := false
	bnsW := models.Wallet{}
	dftW := models.Wallet{}

	for _, wlt := range walletData {
		if wlt.Type == string(models.BONUS) {
			bnsW = wlt
		} else if wlt.Type == string(models.DEFAULT) {
			dftW = wlt
		}
	}

	if bnsW.Point > 0 {
		if bnsW.Point >= c.Point {
			bonusDeduction = c.Point
		} else if bnsW.Point < c.Point {
			bonusDeduction = bnsW.Point
		}
		isUseBonus = true
	}

	if dftW.Point >= 0 {
		if isUseBonus {
			if dftW.Point >= (c.Point - bonusDeduction) {
				walletDeduction = c.Point - bonusDeduction
			} else if dftW.Point < (c.Point - bonusDeduction) {
				return errors.New("^insufficient wallet balance")
			}
		} else {
			if dftW.Point >= (c.Point) {
				walletDeduction = c.Point
			} else if dftW.Point < (c.Point) {
				return errors.New("^insufficient wallet balance")
			}
		}
	} else {
		return errors.New("^insufficient wallet balance")
	}
	trType := "GENERAL"
	if c.TransactionType != "" {
		trType = c.TransactionType
	}
	if (!isUseBonus) || (isUseBonus && walletDeduction > 0) {
		payload := models.WalletHistory{
			WalletID:    dftW.ID,
			SmartbtwID:  c.SmartbtwID,
			AmountPay:   0,
			Point:       walletDeduction,
			Description: c.Description,
			Status:      string(models.OUT),
			CreatedAt:   time.Now(),
			Type:        trType,
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		}
		_, err1 := whCol.InsertOne(ctx, payload)
		if err1 != nil {
			return err1
		}
	}
	payload1 := models.Wallet{
		SmartbtwID: c.SmartbtwID,
		Point:      dftW.Point - walletDeduction,
		Type:       dftW.Type,
		CreatedAt:  dftW.CreatedAt,
		UpdatedAt:  time.Now(),
		DeletedAt:  nil,
	}
	update := bson.M{"$set": payload1}
	_, err2 := wCol.UpdateByID(ctx, dftW.ID, update, opts)
	if err2 != nil {
		return err2
	}

	if isUseBonus {

		payloadBonus := models.WalletHistory{
			WalletID:    bnsW.ID,
			SmartbtwID:  c.SmartbtwID,
			Point:       bonusDeduction,
			Description: c.Description,
			Status:      string(models.OUT),
			Type:        trType,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		}
		_, errBns := whCol.InsertOne(ctx, payloadBonus)
		if errBns != nil {
			return errBns
		}
		payloadBnsWallet := models.Wallet{
			SmartbtwID: c.SmartbtwID,
			Point:      bnsW.Point - bonusDeduction,
			Type:       bnsW.Type,
			CreatedAt:  bnsW.CreatedAt,
			UpdatedAt:  time.Now(),
			DeletedAt:  nil,
		}
		updateBnsW := bson.M{"$set": payloadBnsWallet}
		_, errBnsW := wCol.UpdateByID(ctx, bnsW.ID, updateBnsW, opts)
		if errBnsW != nil {
			return errBnsW
		}
	}
	return nil
	// return res, nil
}
