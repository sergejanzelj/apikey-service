package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vibeitco/apikey-service/model"
	"github.com/vibeitco/go-utils/auth"
	"github.com/vibeitco/go-utils/config"
	"github.com/vibeitco/go-utils/locker"
	"github.com/vibeitco/go-utils/storage"
	"github.com/vibeitco/go-utils/validation"
	"github.com/vibeitco/service-definitions/go/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	timeout = time.Second
)

type handler struct {
	auth.Enforcer
	config config.Core
	dao    storage.Store
	locker locker.Locker
}

func NewHandler(conf config.Core,
	dao storage.Store,
	locker locker.Locker) (model.ApiKeyServiceServer, error) {
	enforcer, err := auth.NewDefaultEnforcer()
	if err != nil {
		return nil, err
	}
	h := &handler{
		Enforcer: enforcer,
		config:   conf,
		dao:      dao,
		locker:   locker,
	}
	return h, nil
}

func (h *handler) ValidateAPIKey(ctx context.Context, req *model.ValidateApiKeyRequest) (*model.ValidateApiKeyResponse, error) {

	result, err := h.GetApiKeys(ctx, &model.GetApiKeysRequest{
		Key: req.GetApiKey(),
	})

	if req.ApiKey == "" {
		fmt.Println("API Key parameter is empty")
		return &model.ValidateApiKeyResponse{Valid: false, UserId: ""}, nil
	}

	userId := ""
	valid := false
	if err != nil {
		fmt.Println("Error reading from database")
		return nil, nil
	}

	if result.Total > 0 {
		userId = result.Items[0].UserId
		valid = true
	}

	return &model.ValidateApiKeyResponse{
		Valid:  valid,
		UserId: userId,
	}, err

}

func (h *handler) Status(ctx context.Context, req *common.GetServiceStatusRequest) (*common.ServiceStatus, error) {
	return &common.ServiceStatus{
		Service: h.config.Service,
		Env:     h.config.Env,
		Version: h.config.Version,
	}, nil
}

func NewUUID() string {
	uid := uuid.New()
	return uid.String()
}

func (h *handler) NewApiKey(ctx context.Context, req *model.NewApiKeyRequest) (*model.ApiKey, error) {
	if err := validation.Validate(req); err != nil {
		return nil, err
	}

	newApiKey := &apikey{
		ApiKey: model.ApiKey{
			Key:       NewUUID(),
			ValidFrom: 0,
			ValidTo:   1,
		},
	}

	err := h.dao.Insert(newApiKey)

	if err != nil {
		return nil, status.Error(codes.Internal, "error while inserting")
	}

	return &newApiKey.ApiKey, nil
}

/*

NewApiKey(ctx context.Context, in *NewApiKeyRequest, opts ...grpc.CallOption) (*ApiKey, error)
	UpdateApiKey(ctx context.Context, in *UpdateApiKeyRequest, opts ...grpc.CallOption) (*ApiKey, error)
	GetApiKey(ctx context.Context, in *GetApiKeyRequest, opts ...grpc.CallOption) (*ApiKey, error)
	GetApiKeys(ctx context.Context, in *GetApiKeysRequest, opts ...grpc.CallOption) (*ApiKeyList, error)
	ValidateAPIKey(ctx context.Context, in *ValidateApiKeyRequest, opts ...grpc.CallOption) (*ValidateApiKeyResponse, error)
	Status(ctx context.Context, in *common.GetServiceStatusRequest, opts ...grpc.CallOption) (*common.ServiceStatus, error)

*/
func (h *handler) UpdateApiKey(ctx context.Context, req *model.UpdateApiKeyRequest) (*model.ApiKey, error) {
	fmt.Println("Called update api")
	return nil, nil
}
func (h *handler) GetApiKey(ctx context.Context, req *model.GetApiKeyRequest) (*model.ApiKey, error) {
	fmt.Println("Called get single key api")

	/*if err := validation.Validate(req); err != nil {
		return nil, err
	}
	// order wrapper
	o := &apikey{
		ApiKey: model.ApiKey{
			Key: req.GetKey(),
		},
	}

	return &o.ApiKey, h.dao.One(o)


	if err := validation.Validate(req); err != nil {
		return nil, err
	}*/

	return nil, nil
}

func getCriteriaForListRequest(req *model.GetApiKeysRequest) storage.Criteria {
	var criterias []storage.Value

	if x := req.GetKey(); len(x) != 0 {
		values := []storage.Value{req.GetKey()}
		c := storage.Criteria{
			Operator: storage.OperatorIn,
			Field:    "key",
			Values:   values,
		}

		criterias = append(criterias, c)
	}

	if x := req.GetUserId(); len(x) != 0 {
		c := storage.Criteria{
			Operator: storage.OperatorEq,
			Field:    "userid",
			Values:   []storage.Value{x}}

		criterias = append(criterias, c)
	}

	if x := req.GetId(); len(x) != 0 {
		c := storage.Criteria{
			Operator: storage.OperatorEq,
			Field:    "_id",
			Values:   []storage.Value{x}}

		criterias = append(criterias, c)
	}

	switch len(criterias) {
	case 0:
		return storage.Criteria{}
	case 1:
		return criterias[0].(storage.Criteria)
	default:
		return storage.Criteria{Operator: storage.OperatorAnd, Values: criterias}
	}
}

func (h *handler) GetApiKeys(ctx context.Context, req *model.GetApiKeysRequest) (*model.ApiKeyList, error) {
	fmt.Println("Called get api keys api")

	slice := apikeys{}

	n, err := h.dao.List(&slice, storage.ListOpt{
		Limit:    25,
		Page:     0,
		Sort:     storage.SortCreatedDesc,
		Criteria: getCriteriaForListRequest(req)})

	if err != nil {
		return nil, err
	}

	return &model.ApiKeyList{Items: slice, Total: int32(n)}, nil
}
