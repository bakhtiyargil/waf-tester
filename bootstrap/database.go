package bootstrap

import (
	"context"
	"fmt"
	"net/url"
	"time"
	"waf-tester/mongo"
)

const (
	DbScheme = "mongodb"
)

func InitMongoDatabase() mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbHost := App.Config.Database.Mongo.Host
	dbPort := App.Config.Database.Mongo.Port
	dbUser := App.Config.Database.Mongo.Username
	dbPass := App.Config.Database.Mongo.Password
	dbName := App.Config.Database.Mongo.Name

	mongodbURI := buildURI(DbScheme, dbHost, dbPort, dbName, dbUser, dbPass)
	client, err := mongo.NewClient(mongodbURI)
	if err != nil {
		App.Logger.FatalF("mongodb init error: %v", err)
	}

	err = client.Connect(ctx)
	if err != nil {
		App.Logger.FatalF("mongodb connect error: %v", err)
	}

	err = client.Ping(ctx)
	if err != nil {
		App.Logger.FatalF("mongodb ping error: %v", err)
	}

	return client
}

func buildURI(scheme, host, port, name, user, pass string) string {
	u := &url.URL{
		Scheme: scheme,
		User:   url.UserPassword(user, pass),
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   "/" + name,
	}
	query := url.Values{}
	query.Set("authSource", "admin")
	u.RawQuery = query.Encode()
	return u.String()
}
