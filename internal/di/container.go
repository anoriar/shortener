package di

import (
	"github.com/anoriar/shortener/internal/config"
	"github.com/anoriar/shortener/internal/handlers/add_url_handler"
	"github.com/anoriar/shortener/internal/handlers/get_url_handler"
	"github.com/anoriar/shortener/internal/router"
	"github.com/anoriar/shortener/internal/storage"
	"github.com/anoriar/shortener/internal/util"
	"github.com/sarulabs/di"
)

const (
	StorageDef    = "storage"
	KeygenDef     = "keygen"
	AddHandlerDef = "add_url_handler"
	GetHandlerDef = "get_handler"
	RouterDef     = "router"
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
			Name: StorageDef,
			Build: func(ctn di.Container) (interface{}, error) {
				return storage.GetInstance(), nil
			},
		},
		{
			Name: KeygenDef,
			Build: func(ctn di.Container) (interface{}, error) {
				return util.NewKeyGen(), nil
			},
		},
		{
			Name: AddHandlerDef,
			Build: func(ctn di.Container) (interface{}, error) {
				storageVar := ctn.Get(StorageDef).(storage.URLStorageInterface)
				keygen := ctn.Get(KeygenDef).(util.KeyGenInterface)
				return add_url_handler.NewAddHandler(storageVar, keygen, cnf.BaseURL), nil
			},
		},
		{
			Name: GetHandlerDef,
			Build: func(ctn di.Container) (interface{}, error) {
				storageVar := ctn.Get(StorageDef).(storage.URLStorageInterface)
				return get_url_handler.NewGetHandler(storageVar), nil
			},
		},
		{
			Name: RouterDef,
			Build: func(ctn di.Container) (interface{}, error) {
				addHandlerVar := ctn.Get(AddHandlerDef).(*add_url_handler.AddHandler)
				getHandlerVar := ctn.Get(GetHandlerDef).(*get_url_handler.GetHandler)
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
