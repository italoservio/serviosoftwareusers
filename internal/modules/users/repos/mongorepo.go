package repos

import (
	"context"
	"errors"
	"github.com/italoservio/serviosoftwareusers/internal/modules/users/models"
	"github.com/italoservio/serviosoftwareusers/pkg/db"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

type MongoUsersRepo struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewMongoUsersRepo(d *db.DB) *MongoUsersRepo {
	database := d.Client.Database("users")
	collection := database.Collection("users")

	repo := MongoUsersRepo{db: database, coll: collection}
	repo.createIndices()

	return &repo
}

func (r *MongoUsersRepo) createIndices() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	if err != nil {
		panic(err)
	}
}

func (r *MongoUsersRepo) GetByID(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var filter = bson.M{"_id": objectID, "deletedAt": nil}
	err = r.coll.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			return nil, nil
		}

		return nil, err
	}

	user.Password = ""
	return &user, nil
}

func (r *MongoUsersRepo) GetByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User

	var filter = bson.M{"email": email, "deletedAt": bson.M{"$eq": nil}}
	err := r.coll.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *MongoUsersRepo) Create(user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.DeletedAt = nil

	inserted, err := r.coll.InsertOne(ctx, user)

	if err != nil {
		return nil, err
	}

	user.ID = inserted.InsertedID.(bson.ObjectID)
	user.Password = ""
	return user, nil
}

func (r *MongoUsersRepo) Update(user *models.User) (*models.User, error) {
	updateDoc := bson.M{}

	if user.FirstName != "" {
		updateDoc["firstName"] = user.FirstName
	}

	if user.LastName != "" {
		updateDoc["lastName"] = user.LastName
	}

	if user.FullName != "" {
		updateDoc["fullName"] = user.FirstName + " " + user.LastName
	}

	if user.Email != "" {
		updateDoc["email"] = user.Email
	}

	if user.Password != "" {
		updateDoc["password"] = user.Password
	}

	if user.Roles != nil && len(user.Roles) > 0 {
		updateDoc["roles"] = user.Roles
	}

	if len(updateDoc) == 0 {
		return nil, errors.New("no fields to update")
	}

	updateDoc["updatedAt"] = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": user.ID}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err := r.coll.FindOneAndUpdate(ctx, filter, bson.M{"$set": updateDoc}, opts).Decode(&user)
	return user, err
}

func (r *MongoUsersRepo) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *MongoUsersRepo) List(input ListInput) (*ListOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	match := bson.M{"deletedAt": bson.M{"$eq": nil}}

	if input.Email != nil && len(*input.Email) > 0 {
		match["email"] = bson.M{"$in": *input.Email}
	}

	if input.FullName != nil && len(*input.FullName) > 0 {
		match["fullName"] = bson.M{"$in": *input.FullName}
	}

	if input.Roles != nil && len(*input.Roles) > 0 {
		match["roles"] = bson.M{"$all": *input.Roles}
	}

	sortOrder := 1
	if input.Order == "desc" {
		sortOrder = -1
	}

	cursor, err := r.coll.Aggregate(ctx, mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$facet", Value: bson.M{
			"total": []bson.M{{"$count": "count"}},
			"items": []bson.M{
				{"$sort": bson.M{input.SortBy: sortOrder}},
				{"$skip": input.Offset},
				{"$limit": input.Limit},
			}},
		}},
	})
	defer cursor.Close(ctx)

	if err != nil {
		return nil, err
	}

	var output ListOutput
	if err = cursor.All(ctx, &output); err != nil {
		return nil, err
	}

	return &output, nil
}
