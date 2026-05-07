package handler

import (
	"net/http"

	"github.com/team/webchat-server/app/gateway/internal/middleware"
	"github.com/team/webchat-server/app/gateway/internal/svc"
	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	authMw := middleware.NewAuthMiddleware(serverCtx.TokenManager)
	auth := func(h http.HandlerFunc) http.HandlerFunc {
		return authMw.Handle(h)
	}

	// public routes (no auth)
	server.AddRoutes(
		[]rest.Route{
			{Method: http.MethodPost, Path: "/auth/login", Handler: loginHandler(serverCtx)},
			{Method: http.MethodPost, Path: "/auth/register", Handler: registerHandler(serverCtx)},
		},
		rest.WithPrefix("/api/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{Method: http.MethodPost, Path: "/auth/refresh", Handler: auth(refreshTokenHandler(serverCtx))},
			{Method: http.MethodGet, Path: "/users/me", Handler: auth(getUserHandler(serverCtx))},
			{Method: http.MethodPut, Path: "/users/me", Handler: auth(updateUserHandler(serverCtx))},
			{Method: http.MethodGet, Path: "/users/:id", Handler: auth(getUserByIdHandler(serverCtx))},
			{Method: http.MethodGet, Path: "/contacts", Handler: auth(listContactsHandler(serverCtx))},
			{Method: http.MethodPost, Path: "/contacts/request", Handler: auth(addContactHandler(serverCtx))},
			{Method: http.MethodPut, Path: "/contacts/request/:id", Handler: auth(handleFriendReqHandler(serverCtx))},
			{Method: http.MethodGet, Path: "/conversations", Handler: auth(getConversationsHandler(serverCtx))},
			{Method: http.MethodGet, Path: "/messages/:conv_id", Handler: auth(getMessagesHandler(serverCtx))},
			{Method: http.MethodPost, Path: "/messages/send", Handler: auth(sendMessageHandler(serverCtx))},
			{Method: http.MethodPost, Path: "/groups", Handler: auth(createGroupHandler(serverCtx))},
			{Method: http.MethodGet, Path: "/groups/:id", Handler: auth(getGroupHandler(serverCtx))},
			{Method: http.MethodPost, Path: "/groups/:id/members", Handler: auth(addGroupMemberHandler(serverCtx))},
			{Method: http.MethodPost, Path: "/files/upload", Handler: auth(getUploadURLHandler(serverCtx))},
			{Method: http.MethodGet, Path: "/files/:id/url", Handler: auth(getFileURLHandler(serverCtx))},
		},
		rest.WithPrefix("/api/v1"),
	)
}
