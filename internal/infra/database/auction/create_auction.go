package auction

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/brunoofgod/goexpert-lesson-6/configuration/logger"
	"github.com/brunoofgod/goexpert-lesson-6/internal/entity/auction_entity"
	"github.com/brunoofgod/goexpert-lesson-6/internal/internal_error"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	StartAuctionExpirationWatcher(ctx, ar, auctionEntity)

	return nil
}

func StartAuctionExpirationWatcher(ctx context.Context, ar *AuctionRepository, auction *auction_entity.Auction) {
	durationStr := os.Getenv("AUCTION_DURATION_SECONDS")
	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		logger.Error("Erro ao converter AUCTION_DURATION_SECONDS: %v", err)
		return
	}

	go func() {
		time.Sleep(time.Duration(duration) * time.Second)
		if err := closeAuction(ctx, ar, auction.Id); err != nil {
			logger.Error("Erro ao fechar leilão ID: %v", err)
		}
	}()
}

func closeAuction(ctx context.Context, ar *AuctionRepository, auctionID string) error {
	auction, err := ar.FindAuctionById(ctx, auctionID)
	if err != nil {
		return err
	}

	auction.Status = auction_entity.Completed

	// Constrói o filtro para encontrar o leilão pelo ID
	filter := bson.M{"_id": auctionID}

	// Define a atualização do status do leilão
	update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}

	// Atualiza o leilão no MongoDB
	_, errUpdate := ar.Collection.UpdateOne(ctx, filter, update)
	if errUpdate != nil {
		return errUpdate // Retorna diretamente o erro do MongoDB
	}

	return nil
}
