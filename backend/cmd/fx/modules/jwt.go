package modules

import (
	"erp-system/pkg/jwt"

	"go.uber.org/fx"
)

// JWTModule provides JWT manager
var JWTModule = fx.Module("jwt",
	fx.Provide(jwt.NewManager),
)
