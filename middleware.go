package goblast

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Description: a helper method to register endpoints to controllers in the RegisterEndpoint method
func RegisterEndpointHelper(endpoint IEndpoint, group *echo.Group) {
	group.Match([]string{http.MethodPost}, endpoint.GetPath(), endpoint.GetHandler())
}

// Description: middleware for validating request body, initializing tracing ids, & creating contextful request data
// T - struct type with validation spec
type ContextfulReqEndpoint[T interface{}] struct {
	Endpoint IEndpoint
}

const (
	WELL_KNOWN_HEADER____AUTHORIZATION                      = "X-WellKnown-GoBlast-Authorization"
	WELL_KNOWN_HEADER____REFERENCE_ID                       = "X-WellKnown-GoBlast-Reference-Id"
	WELL_KNOWN_HEADER____AUTHORIZATION_CLAIMS____SUBJECT    = "X-WellKnown-GoBlast-Authorization-Claims-Subject"
	WELL_KNOWN_HEADER____AUTHORIZATION_CLAIMS____SESSION_ID = "X-WellKnown-GoBlast-Authorization-Claims-Session-Id"

	WELL_KNOWN_ERROR_MSG____AUTHORIZATION_FAILED       = "GoBlast:AuthorizationFailed"
	WELL_KNOWN_ERROR_MSG____UNAUTHORIZED               = "GoBlast:Unauthorized"
	WELL_KNOWN_ERROR_MSG____REQ_BODY_VALIDATION_FAILED = "GoBlast:ReqBodyValidationFailed"
)

// Metadata model for base model contextfulreq
type ContextfulReqMetadata struct {
	TracingId   string
	ReferenceId string
	Subject     string
	SessionId   string
	Ctx         context.Context
}

// Base Model to be received by the core endpoint
type ContextfulReq[T interface{}] struct {
	TracingId   string
	ReferenceId string
	Subject     string
	SessionId   string
	Ctx         context.Context
	ReqData     T
}

// For setting up metadata from passed-on contextfulreq data
func (c *ContextfulReq[T]) SetMetadata(metadata ContextfulReqMetadata) {
	c.TracingId = metadata.TracingId
	c.ReferenceId = metadata.ReferenceId
	c.Subject = metadata.Subject
	c.SessionId = metadata.SessionId
	c.Ctx = metadata.Ctx
}

// For exporting metadata from passed-on contextfulreq data
func (c *ContextfulReq[T]) GetMetadata() ContextfulReqMetadata {
	return ContextfulReqMetadata{
		TracingId:   c.TracingId,
		ReferenceId: c.ReferenceId,
		Subject:     c.Subject,
		SessionId:   c.SessionId,
		Ctx:         c.Ctx,
	}
}

func (cf *ContextfulReqEndpoint[T]) GetHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		tracingId := uuid.NewString()
		referenceId := c.Request().Header.Get(WELL_KNOWN_HEADER____REFERENCE_ID)
		if referenceId == "" {
			referenceId = uuid.NewString()
		}
		subject := c.Request().Header.Get(WELL_KNOWN_HEADER____AUTHORIZATION_CLAIMS____SUBJECT)
		sessionId := c.Request().Header.Get(WELL_KNOWN_HEADER____AUTHORIZATION_CLAIMS____SESSION_ID)

		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return NotOkResponse(c, WELL_KNOWN_ERROR_MSG____REQ_BODY_VALIDATION_FAILED)
		}

		var bodyData T
		err = json.Unmarshal(body, &bodyData)
		if err != nil {
			return NotOkResponse(c, WELL_KNOWN_ERROR_MSG____REQ_BODY_VALIDATION_FAILED)
		}

		validationChecker := validator.New()
		err = validationChecker.Struct(bodyData)
		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				LogError(referenceId, tracingId, "Request Body Validation Failed")
				return NotOkResponse(c, WELL_KNOWN_ERROR_MSG____REQ_BODY_VALIDATION_FAILED)
			}

			validationErrors := err.(validator.ValidationErrors)
			for _, validationError := range validationErrors {
				validationErrorMessage := "Field=" + validationError.Field() + ",Tag=" + validationError.Tag() + ",ActualTag=" + validationError.ActualTag() + ",Error=" + validationError.Error()
				LogError(referenceId, tracingId, validationErrorMessage)
			}

			return NotOkResponse(c, WELL_KNOWN_ERROR_MSG____REQ_BODY_VALIDATION_FAILED)
		}

		contextfulReq := ContextfulReq[T]{
			TracingId:   tracingId,
			ReferenceId: referenceId,
			Subject:     subject,
			SessionId:   sessionId,
			Ctx:         context.Background(),
			ReqData:     bodyData,
		}
		contextfulReqBytes := new(bytes.Buffer)
		json.NewEncoder(contextfulReqBytes).Encode(contextfulReq)

		newR := c.Request().Clone(c.Request().Context())
		c.Request().Body = io.NopCloser(bytes.NewReader(contextfulReqBytes.Bytes()))
		newR.Body = io.NopCloser(bytes.NewReader(contextfulReqBytes.Bytes()))
		err = c.Request().ParseForm()
		if err != nil {
			LogError(referenceId, tracingId, "Failed to clone request")
			return NotOkResponse(c, "Failed to clone request")
		}
		c.SetRequest(newR)
		return cf.Endpoint.GetHandler()(c)
	}
}

func (cf *ContextfulReqEndpoint[T]) GetPath() string {
	return cf.Endpoint.GetPath()
}

func (cf *ContextfulReqEndpoint[T]) Register(group *echo.Group) {
	RegisterEndpointHelper(cf, group)
}

// Description: middleware for validation contextful requests using auth
// AuthManager - interface for implementing custom authentication methods
type AuthEndpoint struct {
	Endpoint    IEndpoint
	AuthManager IAuthManager
}

// Description: interface for developers to implement their own authorization check methods
type IAuthManager interface {
	GetAuthorization(string) (AuthorizationClaims, error)
}

// Basic claims to be returned by IAuthManager implementation
type AuthorizationClaims struct {
	Subject   string
	SessionId string
}

func (a *AuthEndpoint) GetHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		authorizationHeader := c.Request().Header.Get(WELL_KNOWN_HEADER____AUTHORIZATION)
		if authorizationHeader == "" {
			return NotOkResponse(c, WELL_KNOWN_ERROR_MSG____UNAUTHORIZED)
		}
		authorizationClaims, err := a.AuthManager.GetAuthorization(authorizationHeader)
		if err != nil {
			return NotOkResponse(c, WELL_KNOWN_ERROR_MSG____AUTHORIZATION_FAILED)
		}
		if c.Request().Header.Get(WELL_KNOWN_HEADER____AUTHORIZATION_CLAIMS____SUBJECT) != "" {
			c.Request().Header.Del(WELL_KNOWN_HEADER____AUTHORIZATION_CLAIMS____SUBJECT)
		}
		c.Request().Header.Add(WELL_KNOWN_HEADER____AUTHORIZATION_CLAIMS____SUBJECT, authorizationClaims.Subject)
		if c.Request().Header.Get(WELL_KNOWN_HEADER____AUTHORIZATION_CLAIMS____SESSION_ID) != "" {
			c.Request().Header.Del(WELL_KNOWN_HEADER____AUTHORIZATION_CLAIMS____SESSION_ID)
		}
		c.Request().Header.Add(WELL_KNOWN_HEADER____AUTHORIZATION_CLAIMS____SESSION_ID, authorizationClaims.SessionId)
		return a.Endpoint.GetHandler()(c)
	}
}

func (a *AuthEndpoint) GetPath() string {
	return a.Endpoint.GetPath()
}

func (a *AuthEndpoint) Register(group *echo.Group) {
	RegisterEndpointHelper(a, group)
}

// Description: middleware for clearing up authentication headers sent from client without validation
type ClearAuthEndpoint struct {
	Endpoint IEndpoint
}

func (cl *ClearAuthEndpoint) GetHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get(WELL_KNOWN_HEADER____AUTHORIZATION) != "" {
			c.Request().Header.Del(WELL_KNOWN_HEADER____AUTHORIZATION)
		}
		if c.Request().Header.Get(WELL_KNOWN_HEADER____AUTHORIZATION_CLAIMS____SUBJECT) != "" {
			c.Request().Header.Del(WELL_KNOWN_HEADER____AUTHORIZATION_CLAIMS____SUBJECT)
		}
		if c.Request().Header.Get(WELL_KNOWN_HEADER____AUTHORIZATION_CLAIMS____SESSION_ID) != "" {
			c.Request().Header.Del(WELL_KNOWN_HEADER____AUTHORIZATION_CLAIMS____SESSION_ID)
		}
		return cl.Endpoint.GetHandler()(c)
	}
}

func (cl *ClearAuthEndpoint) GetPath() string {
	return cl.Endpoint.GetPath()
}

func (cl *ClearAuthEndpoint) Register(group *echo.Group) {
	RegisterEndpointHelper(cl, group)
}
