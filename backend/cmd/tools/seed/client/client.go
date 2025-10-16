package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/hay-kot/hookfeed/backend/internal/data/dtos"
	"github.com/hay-kot/hookfeed/backend/pkgs/httpclient"
)

type Client struct {
	http *httpclient.Client

	Auth      Auth
	User      User
	Lists     Lists
	ListItems ListItems
}

func New(c *http.Client, baseurl string) Client {
	client := httpclient.New(c, baseurl)
	client.Use(httpclient.MwJSON)
	return Client{
		http:      client,
		Auth:      Auth{http: client},
		User:      User{http: client},
		Lists:     Lists{http: client},
		ListItems: ListItems{http: client},
	}
}

type Auth struct {
	http *httpclient.Client
}

func (a *Auth) Login(ctx context.Context, data dtos.UserAuthenticate) (dtos.UserSession, error) {
	const path = "/api/v1/users/login"
	sess, err := post[dtos.UserSession](a.http, ctx, path, data)
	if err != nil {
		return dtos.UserSession{}, err
	}

	a.http.SetBearer(sess.Token)
	return sess, nil
}

func (a *Auth) Register(ctx context.Context, data dtos.UserRegister) (dtos.User, error) {
	const path = "/api/v1/users/register"
	return post[dtos.User](a.http, ctx, path, data)
}

type User struct {
	http *httpclient.Client
}

func (u *User) Self(ctx context.Context) (dtos.User, error) {
	const path = "/api/v1/users/self"
	return get[dtos.User](u.http, ctx, path)
}

func (u *User) SelfUpdate(ctx context.Context, body dtos.UserUpdate) (dtos.User, error) {
	const path = "/api/v1/users/self"
	return patch[dtos.User](u.http, ctx, path, body)
}

type Lists struct {
	http *httpclient.Client
}

// GetAll retrieves all packing lists with pagination and optional query parameters
func (l *Lists) GetAll(ctx context.Context, query dtos.PackingListQuery) (dtos.PaginationResponse[dtos.PackingList], error) {
	path := "/api/v1/packing-lists"

	// Build query parameters
	var params []string
	if query.OrderBy != "" {
		params = append(params, "orderBy="+query.OrderBy)
	}
	if query.Skip > 0 {
		params = append(params, fmt.Sprintf("skip=%d", query.Skip))
	}
	if query.Limit > 0 {
		params = append(params, fmt.Sprintf("limit=%d", query.Limit))
	}

	// Add query string to path if we have parameters
	if len(params) > 0 {
		path += "?" + strings.Join(params, "&")
	}

	return get[dtos.PaginationResponse[dtos.PackingList]](l.http, ctx, path)
}

// Get retrieves a specific packing list by ID
func (l *Lists) Get(ctx context.Context, id string) (dtos.PackingList, error) {
	path := "/api/v1/packing-lists/" + id
	return get[dtos.PackingList](l.http, ctx, path)
}

// Create creates a new packing list
func (l *Lists) Create(ctx context.Context, data dtos.PackingListCreate) (dtos.PackingList, error) {
	const path = "/api/v1/packing-lists"
	return post[dtos.PackingList](l.http, ctx, path, data)
}

// Update updates an existing packing list
func (l *Lists) Update(ctx context.Context, id uuid.UUID, data dtos.PackingListUpdate) (dtos.PackingList, error) {
	path := "/api/v1/packing-lists/" + id.String()
	return patch[dtos.PackingList](l.http, ctx, path, data)
}

// Delete deletes a packing list by ID
func (l *Lists) Delete(ctx context.Context, id string) error {
	path := "/api/v1/packing-lists/" + id
	return delete(l.http, ctx, path)
}

type ListItems struct {
	http *httpclient.Client
}

// Create creates a new packing list item
func (li *ListItems) Create(ctx context.Context, listID uuid.UUID, data dtos.PackingListItemCreate) (dtos.PackingListItem, error) {
	path := "/api/v1/packing-lists/" + listID.String() + "/items"
	return post[dtos.PackingListItem](li.http, ctx, path, data)
}

// Update updates an existing packing list item
func (li *ListItems) Update(ctx context.Context, listID uuid.UUID, itemID uuid.UUID, data dtos.PackingListItemUpdate) (dtos.PackingListItem, error) {
	path := "/api/v1/packing-lists/" + listID.String() + "/items/" + itemID.String()
	return put[dtos.PackingListItem](li.http, ctx, path, data)
}

// Delete deletes a packing list item
func (li *ListItems) Delete(ctx context.Context, listID uuid.UUID, itemID uuid.UUID) error {
	path := "/api/v1/packing-lists/" + listID.String() + "/items/" + itemID.String()
	return delete(li.http, ctx, path)
}
