package factories

import (
	"context"
	"time"

	"github.com/bluele/factory-go/factory"
	"github.com/mises-id/socialsvc/app/models"
	"github.com/mises-id/socialsvc/lib/db"
)

var FollowFactory = factory.NewFactory(
	&models.Follow{},
).Attr("ID", func(args factory.Args) (interface{}, error) {
	return nil, nil
}).Attr("FromUID", func(args factory.Args) (interface{}, error) {
	return uint64(0), nil
}).Attr("ToUID", func(args factory.Args) (interface{}, error) {
	return uint64(0), nil
}).Attr("IsFriend", func(args factory.Args) (interface{}, error) {
	return false, nil
}).Attr("CreatedAt", func(args factory.Args) (interface{}, error) {
	return time.Now(), nil
}).Attr("UpdatedAt", func(args factory.Args) (interface{}, error) {
	return time.Now(), nil
}).OnCreate(func(args factory.Args) error {
	_, err := db.DB().Collection("follows").InsertOne(context.Background(), args.Instance())
	return err
})
