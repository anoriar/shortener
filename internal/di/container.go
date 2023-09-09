package di

import (
	"github.com/anoriar/shortener/internal/config"
	"github.com/anoriar/shortener/internal/handlers"
	"github.com/anoriar/shortener/internal/router"
	"github.com/anoriar/shortener/internal/storage"
	"github.com/anoriar/shortener/internal/util"
	"github.com/sarulabs/di"
)

type Container struct {
	ctn di.Container
}

func NewContainer(cnf *config.Config) (*Container, error) {
	builder, err := di.NewBuilder()
	if err != nil {
		return nil, err
	}
	err = builder.Add([]di.Def{
		{
			Name: "storage",
			Build: func(ctn di.Container) (interface{}, error) {
				return storage.GetInstance(), nil
			},
		},
		{
			Name: "keygen",
			Build: func(ctn di.Container) (interface{}, error) {
				return util.NewKeyGen(), nil
			},
		},
		{
			Name: "addHandler",
			Build: func(ctn di.Container) (interface{}, error) {
				storageVar := ctn.Get("storage").(storage.URLStorageInterface)
				keygen := ctn.Get("keygen").(util.KeyGenInterface)
				return handlers.NewAddHandler(storageVar, keygen, cnf.BaseURL), nil
			},
		},
		{
			Name: "getHandler",
			Build: func(ctn di.Container) (interface{}, error) {
				storageVar := ctn.Get("storage").(storage.URLStorageInterface)
				return handlers.NewGetHandler(storageVar), nil
			},
		},
		{
			Name: "router",
			Build: func(ctn di.Container) (interface{}, error) {
				addHandlerVar := ctn.Get("addHandler").(*handlers.AddHandler)
				getHandlerVar := ctn.Get("getHandler").(*handlers.GetHandler)
				return router.NewRouter(addHandlerVar, getHandlerVar), nil
			},
		},
	}...)

	if err != nil {
		return nil, err
	}

	return &Container{
		ctn: builder.Build(),
	}, nil
}

func (c *Container) Resolve(name string) interface{} {
	return c.ctn.Get(name)
}

func (c *Container) Clean() error {
	return c.ctn.Clean()
}
