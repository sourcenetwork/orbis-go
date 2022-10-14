package bulletin

import "github.com/samber/do"

type ProviderFn func(*do.Injector) Service

type Service interface{}
