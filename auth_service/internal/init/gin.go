package init

import (
	hashImpl "auth_service/internal/hash/hasher"
	ginImpl "auth_service/internal/http/gin"
	"auth_service/internal/http/gin/middlewares/recovery"
	requestid "auth_service/internal/http/gin/middlewares/request-id"
	"auth_service/internal/http/gin/routes/auth/infoUser"
	"auth_service/internal/http/gin/routes/auth/login"
	"auth_service/internal/http/gin/routes/auth/refresh"
	"auth_service/internal/http/gin/routes/auth/register"
	"auth_service/internal/http/gin/routes/permissions/changeDescriptionPermissions"
	"auth_service/internal/http/gin/routes/role/createRole"
	"auth_service/internal/http/gin/routes/role/deleteRole"
	"auth_service/internal/http/gin/routes/role/getAllRoles"
	"auth_service/internal/http/gin/routes/role/infoRole"
	"auth_service/internal/http/gin/routes/user/deleteUser"
	"auth_service/internal/http/gin/routes/user/findUsersByRole"
	"auth_service/internal/jwt"
	"auth_service/internal/repo"
	"github.com/gin-gonic/gin"
)

const (
	authPath       = "/auth"
	rolePath       = "/role"
	permissionPath = "/permission"
	userPath       = "/user"
)

func Gin(db repo.DB, jwt jwt.Handler) *gin.Engine {
	ginRouter := ginImpl.NewGinServer()

	ginRouter.AddMiddleware(
		recovery.Middleware(),
		requestid.Middleware(),
	)

	ginRouter.AddRouters(
		getAuth(authPath, db, jwt),
		getRole(rolePath, db),
		getPermissions(permissionPath, db),
		getUser(userPath, db),
	)
	return ginRouter.Build()
}

func getAuth(path string, db repo.DB, jwtHandler jwt.Handler) *ginImpl.Group {
	authGroup := ginImpl.NewGroup(path)

	hasher := hashImpl.NewHasher(hashImpl.MinHashCost)

	authGroup.AddRouters(
		login.NewLoginHandler(db, hasher, jwtHandler),
		register.NewRegisterHandler(db, hasher),
		refresh.NewRefreshHandler(db, jwtHandler),
		infoUser.TakeInfoMe(db, jwtHandler),
	)

	return authGroup
}

func getUser(path string, db repo.DB) *ginImpl.Group {
	userGroup := ginImpl.NewGroup(path)

	userGroup.AddRouters(
		deleteUser.DeleteUser(db),
		findUsersByRole.FindUsersByRole(db),
	)

	return userGroup
}

func getRole(path string, db repo.DB) *ginImpl.Group {
	roleGroup := ginImpl.NewGroup(path)

	roleGroup.AddRouters(
		infoRole.TakeInfoRole(db),
		createRole.CreateRole(db),
		deleteRole.DeleteRole(db),
		getAllRoles.GetAllRoles(db),
	)

	return roleGroup
}

func getPermissions(path string, db repo.DB) *ginImpl.Group {
	permissionGroup := ginImpl.NewGroup(path)

	permissionGroup.AddRouters(
		changeDescriptionPermissions.NewDescriptionPermissions(db),
	)
	return permissionGroup
}
