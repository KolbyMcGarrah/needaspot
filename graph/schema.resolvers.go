package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"

	"github.com/KolbyMcGarrah/nas/graph/generated"
	"github.com/KolbyMcGarrah/nas/graph/model"
	"github.com/KolbyMcGarrah/nas/internal/auth"
	"github.com/KolbyMcGarrah/nas/internal/pkg/jwt"
	"github.com/KolbyMcGarrah/nas/internal/requests"
	"github.com/KolbyMcGarrah/nas/internal/users"
)

func (r *mutationResolver) CreateRequest(ctx context.Context, input model.NewRequest) (*model.Request, error) {
	//get user object from ctx and if user is not set then we respnosd with access denied.
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.Request{}, fmt.Errorf("access denied")
	}

	var request requests.Request

	request.Title = input.Title
	request.Location = input.Location
	request.Workout = input.Workout
	request.User = user

	requestId := request.Save()
	graqhqlUser := &model.User{
		ID:       user.ID,
		Username: user.Username,
		Age:      user.Age,
		Gender:   user.Gender,
		Level:    user.Level,
	}

	return &model.Request{ID: strconv.Itoa(requestId), Title: request.Title, Location: request.Location, Workout: request.Workout, User: graqhqlUser}, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	user.Age = input.Age
	user.Gender = input.Gender
	user.Level = input.Level
	user.Create()
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	correct := user.Authenticate()
	if !correct {
		//1
		return "", &users.InvalidLogin{}
	}
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	username, err := jwt.ParseToken(input.Token)
	if err != nil {
		return "", fmt.Errorf("access denied")
	}
	token, err := jwt.GenerateToken(username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) CreateComment(ctx context.Context, input model.NewComment) (*model.Comment, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.Comment{}, fmt.Errorf("access denied")
	}
	var comment model.Comment
	return &comment, nil
}

func (r *mutationResolver) CreateInterest(ctx context.Context, input model.NewInterest) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Requests(ctx context.Context) ([]*model.Request, error) {
	var resultRequests []*model.Request
	var dbRequests []requests.Request
	dbRequests = requests.GetAll()
	for _, request := range dbRequests {
		graphqlUser := &model.User{
			ID:       request.User.ID,
			Username: request.User.Username,
			Age:      request.User.Age,
			Gender:   request.User.Gender,
			Level:    request.User.Level,
		}
		resultRequests = append(resultRequests, &model.Request{ID: request.ID, Title: request.Title, Workout: request.Workout, Location: request.Location, User: graphqlUser})
	}
	return resultRequests, nil
}

func (r *queryResolver) Comments(ctx context.Context) ([]*model.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Interests(ctx context.Context) ([]*model.Interest, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
