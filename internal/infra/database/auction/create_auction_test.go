package auction

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/brunoofgod/goexpert-lesson-6/internal/entity/auction_entity"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupMongoContainer(t *testing.T) (context.Context, *mongo.Client, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp"),
	}

	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)

	// Obtém a URL do MongoDB no container
	mongoHost, err := mongoC.Host(ctx)
	assert.NoError(t, err)

	mongoPort, err := mongoC.MappedPort(ctx, "27017")
	assert.NoError(t, err)

	uri := "mongodb://" + mongoHost + ":" + mongoPort.Port()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	assert.NoError(t, err)

	// Função para encerrar o container
	cleanup := func() {
		client.Disconnect(ctx)
		mongoC.Terminate(ctx)
	}

	return ctx, client, cleanup
}

func TestAuctionAutoClose(t *testing.T) {
	ctx, client, cleanup := setupMongoContainer(t)
	defer cleanup()

	db := client.Database("testdb")
	collection := db.Collection("auctions")
	collection.Drop(ctx) // Limpa a coleção antes do teste

	// Define a duração do leilão como 2 segundos
	os.Setenv("AUCTION_DURATION_SECONDS", "2")

	auction := &auction_entity.Auction{
		Id:     "123",
		Status: auction_entity.Active,
	}

	auctionRepository := NewAuctionRepository(db)
	err := auctionRepository.CreateAuction(context.Background(), auction)
	assert.Nil(t, err)

	// Aguarda um pouco mais que a duração definida
	time.Sleep(3 * time.Second)

	// Verifica se o leilão foi fechado automaticamente
	updatedAuction, errFind := auctionRepository.FindAuctionById(ctx, "123")
	assert.Nil(t, errFind)

	assert.NotNil(t, updatedAuction)
	assert.Equal(t, auction_entity.Completed, updatedAuction.Status)
}

/*
import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/brunoofgod/goexpert-lesson-6/internal/entity/auction_entity"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestAuctionAutoClose(t *testing.T) {
	// Configuração do MongoDB em memória
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	db := mt.Client.Database("testdb")
	collection := db.Collection("auctions")
	collection.Drop(context.Background()) // Limpa a coleção antes do teste

	// Define a duração do leilão como 2 segundos
	os.Setenv("AUCTION_DURATION", "2")

	auction := &auction_entity.Auction{
		Id:     "123",
		Status: auction_entity.Active,
	}

	auctionRepository := NewAuctionRepository(db)
	err := auctionRepository.CreateAuction(context.Background(), auction)
	assert.NoError(t, err)

	// Aguarda um pouco mais que a duração definida
	time.Sleep(3 * time.Second)

	// Verifica se o leilão foi fechado automaticamente
	var updatedAuction auction_entity.Auction
	mongoSingleResult := collection.FindOne(context.Background(), bson.M{"id": "123"}).Decode(&updatedAuction)

	assert.Empty(t, mongoSingleResult.Error())
	assert.Equal(t, "closed", updatedAuction.Status)
}
*/
