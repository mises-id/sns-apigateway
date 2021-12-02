package factories

import (
	"context"

	"github.com/bluele/factory-go/factory"
	"github.com/google/uuid"
	"github.com/mises-id/socialsvc/app/models"
	"github.com/mises-id/socialsvc/app/models/enum"
	"github.com/mises-id/socialsvc/lib/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserFactory = factory.NewFactory(
	&models.User{},
).SeqInt("UID", func(n int) (interface{}, error) {
	return uint64(n), nil
}).Attr("Username", func(args factory.Args) (interface{}, error) {
	return "", nil
}).Attr("Gender", func(args factory.Args) (interface{}, error) {
	return enum.GenderFemale, nil
}).Attr("Misesid", func(args factory.Args) (interface{}, error) {
	return uuid.New().String(), nil
}).Attr("Mobile", func(args factory.Args) (interface{}, error) {
	return "", nil
}).Attr("Email", func(args factory.Args) (interface{}, error) {
	return "", nil
}).Attr("Address", func(args factory.Args) (interface{}, error) {
	return "", nil
}).Attr("AvatarID", func(args factory.Args) (interface{}, error) {
	return uint64(0), nil
}).OnCreate(func(args factory.Args) error {
	user := args.Instance().(*models.User)
	_, err := db.DB().Collection("users").InsertOne(context.Background(), user)
	db.DB().Collection("counters").FindOneAndUpdate(context.Background(), bson.M{"_id": "userid"},
		bson.M{"$set": bson.M{"seq": user.UID}},
		options.FindOneAndUpdate().SetUpsert(true),
	)
	println("create user ", user.UID, user.Misesid)
	return err
})

func InitUsers(args ...*models.User) {
	for _, arg := range args {
		UserFactory.MustCreateWithOption(map[string]interface{}{
			"UID":      arg.UID,
			"Username": arg.Username,
			"Misesid":  arg.Misesid,
			"Gender":   arg.Gender,
			"Mobile":   arg.Mobile,
			"Email":    arg.Email,
			"Address":  arg.Address,
			"AvatarID": arg.AvatarID,
		})
	}
}
