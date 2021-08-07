package app

import (
	"cos-backend-com/src/common/app"
	"cos-backend-com/src/common/providers/session"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/cores"
	"cos-backend-com/src/cores/routers/bounties"
	"cos-backend-com/src/cores/routers/categories"
	"cos-backend-com/src/cores/routers/discos"
	"cos-backend-com/src/cores/routers/exchanges"
	"cos-backend-com/src/cores/routers/files"
	"cos-backend-com/src/cores/routers/follows"
	"cos-backend-com/src/cores/routers/proposals"
	"cos-backend-com/src/cores/routers/startups"
	"cos-backend-com/src/cores/routers/tags"
	"cos-backend-com/src/libs/auth"
	"cos-backend-com/src/libs/filters"
	filesSdk "cos-backend-com/src/libs/sdk/files"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/mediocregopher/radix.v2/pool"
	s "github.com/wujiu2020/strip"
)

const (
	AppName = "cores"
)

func AppInit(tea *s.Strip, confPath string, files ...string) *appConfig {
	app := &appConfig{app.New(tea, AppName), cores.Env}
	app.ConfigLoad(app.Env, confPath, files...)
	app.ConfigCheck()
	app.ConfigDB()
	app.ConfigProviders()
	app.ConfigFilters()
	app.ConfigRoutes()
	app.ConfigDone()
	return app
}

type appConfig struct {
	app.AppConfig
	Env *cores.Config
}

// 配置完毕，做一些运行时检查与初始化
func (p *appConfig) Start() {
}

func (p *appConfig) ConfigCheck() {
}

func (p *appConfig) ConfigDB() *sqlx.DB {
	return p.ConnectDB()
}

func (p *appConfig) ConfigProviders() {
	redisPool, err := pool.NewCustom("tcp",
		p.Env.Redis.Addr,
		p.Env.Redis.PoolSize,
		util.RedisDialWithSecret(p.Env.Redis.Secret),
	)
	if err != nil {
		p.Logger().Error("redis pool.New:", err)
		os.Exit(1)
	}
	p.Provide(redisPool)

	if err := session.ConfigSessions(p.Strip, redisPool, p.Env.Session); err != nil {
		p.Logger().Error("config sessions:", err)
		os.Exit(1)
	}

	var rt http.RoundTripper
	if err := p.Injector().Find(&rt, ""); err != nil {
		panic(err)
	}
	oauth2TokenURL := p.Env.Service.Account + "/oauth2/token"
	p.Provide(auth.AuthTransportProvider(oauth2TokenURL))
	p.Provide(filesSdk.NewFileService(filesSdk.BaseConfig{
		cores.Env.Minio.Endpoint,
		cores.Env.Minio.Secure,
		cores.Env.Minio.AccessKey,
		cores.Env.Minio.SecretKey,
	}))
}

func (p *appConfig) ConfigFilters() {
	p.Filter(filters.Cors)
}

func (p *appConfig) ConfigRoutes() {
	p.Routers(util.VersionRouter())
	p.Routers(
		s.Router("/startups",
			s.Get(startups.StartUpsHandler{}).Action("List"),
			s.Router("/:id",
				s.Get(startups.StartUpsHandler{}).Action("Get"),
			),
		),

		s.Router("/startups",
			s.Filter(filters.LoginRequiredInner),
			s.Post(startups.StartUpsHandler{}).Action("Create"),
			s.Router("/me",
				s.Get(startups.StartUpsHandler{}).Action("ListMe"),
				s.Router("/:id",
					s.Get(startups.StartUpsHandler{}).Action("GetMe"),
				),
				s.Router("/followed",
					s.Get(startups.StartUpsHandler{}).Action("ListMeFollowed"),
				),
			),
			s.Router("/prepareId",
				s.Get(startups.StartUpsHandler{}).Action("GetPrepareId"),
			),
			s.Router("/:id",
				s.Put(startups.StartUpsHandler{}).Action("Update"),
				s.Router("/hasFollowed",
					s.Get(startups.StartUpsHandler{}).Action("HasFollowed"),
				),
				s.Router("/settings",
					s.Put(startups.StartUpSettingsHandler{}).Action("Update"),
				),
				s.Router("/payTokens",
					s.Get(startups.StartUpsHandler{}).Action("GetPayTokens"),
				),
			),
			//restore startup
			s.Router("/:id:restore",
				s.Post(startups.StartUpsHandler{}).Action("Restore"),
			),
			s.Router("/tokenAddr",
				s.Router("/isBinding",
					s.Get(startups.StartUpsHandler{}).Action("IsTokenAddrBinding"),
				),

			),
			//restore startup settings
			s.Router("/:id/settings:restore",
				s.Post(startups.StartUpSettingsHandler{}).Action("Restore"),
			),

			s.Router("/:id/follows",
				s.Post(follows.FollowsHandler{}).Action("Create"),
				s.Delete(follows.FollowsHandler{}).Action("Delete"),
			),
			s.Router("/:id/exchange",
				s.Post(exchanges.ExchangesHandler{}).Action("CreateExchange"),
				s.Get(exchanges.ExchangesHandler{}).Action("GetExchangeByStartup"),
			),
		),

		s.Router("/categories",
			s.Get(categories.CategoriesHandler{}).Action("List"),
			s.Router("/:id",
				s.Get(categories.CategoriesHandler{}).Action("Get"),
			),
		),

		s.Router("/tags",
			s.Get(tags.TagsHandler{}).Action("List"),
		),

		s.Router("/files",
			s.Filter(filters.LoginRequiredInner),
			s.Post(files.FilesHandler{}).Action("SignUploadFile"),
		),

		s.Router("/prepareId",
			s.Get(startups.StartUpsHandler{}).Action("GetPrepareId"),
		),

		s.Router("/startups/:id/bounties",
			s.Get(bounties.BountiesHandler{}).Action("ListStartupBounties"),
			s.Router("/me",
				s.Get(bounties.BountiesHandler{}).Filter(filters.LoginRequiredInner).Action("ListStartupBountiesMe"),
			),
			s.Post(bounties.BountiesHandler{}).Filter(filters.LoginRequiredInner).Action("Create"),
		),

		s.Router("/bounties",
			s.Get(bounties.BountiesHandler{}).Action("ListBounties"),
			s.Router("/me",
				s.Filter(filters.LoginRequiredInner),
				s.Get(bounties.BountiesHandler{}).Action("ListBountiesMe"),
			),
			s.Router("/:id",
				s.Get(bounties.BountiesHandler{}).Action("GetBounty"),
				s.Router("/me",
					s.Get(bounties.BountiesHandler{}).Filter(filters.LoginRequiredInner).Action("GetBountyMe"),
				),
			),
			s.Router("/:id:closed",
				s.Put(bounties.BountiesHandler{}).Filter(filters.LoginRequiredInner).Action("Closed"),
			),
			s.Router("/:id:startWork",
				s.Filter(filters.LoginRequiredInner),
				s.Post(bounties.BountiesHandler{}).Action("StartWork"),
			),
			s.Router("/:id:submitted",
				s.Filter(filters.LoginRequiredInner),
				s.Put(bounties.BountiesHandler{}).Action("Submitted"),
			),
			s.Router("/:id:quited",
				s.Filter(filters.LoginRequiredInner),
				s.Put(bounties.BountiesHandler{}).Action("Quited"),
			),
			s.Router("/:id:rejected",
				s.Filter(filters.LoginRequiredInner),
				s.Put(bounties.BountiesHandler{}).Action("Rejected"),
			),
			s.Router("/:id:paid",
				s.Filter(filters.LoginRequiredInner),
				s.Put(bounties.BountiesHandler{}).Action("Paid"),
			),
			s.Router("/users/:userId",
				s.Get(bounties.BountiesHandler{}).Action("ListUserBounties"),
			),
		),
		s.Router("/exchanges",
			s.Get(exchanges.ExchangesHandler{}).Action("ListExchanges"),
			s.Router("/:id",
				s.Get(exchanges.ExchangesHandler{}).Action("GetExchange"),
				s.Router("/transactions",
					s.Post(exchanges.ExchangesHandler{}).Action("CreateExchangeTx"),
				),
				s.Router("/stats:total",
					s.Get(exchanges.ExchangesHandler{}).Action("GetExchangeOneStatsTotal"),
				),
				s.Router("/stats:priceChange",
					s.Get(exchanges.ExchangesHandler{}).Action("GetExchangeOneStatsPriceChange"),
				),
			),
		),
		s.Router("/exchange/transactions",
			s.Router("/:id",
				s.Get(exchanges.ExchangesHandler{}).Action("GetExchangeTx"),
			),
		),
		s.Router("/exchanges:stats",
			s.Get(exchanges.ExchangesHandler{}).Action("GetExchangeAllStatsTotal"),
		),

		s.Router("/swap",
			s.Router("/pairs",
				s.Post(exchanges.SwapEventsHandler{}).Action("CreatePair")),
			s.Router("/mints",
				s.Post(exchanges.SwapEventsHandler{}).Action("Mint")),
			s.Router("/burns",
				s.Post(exchanges.SwapEventsHandler{}).Action("Burn")),
			s.Router("/swaps",
				s.Post(exchanges.SwapEventsHandler{}).Action("Swap")),
			s.Router("/syncs",
				s.Post(exchanges.SwapEventsHandler{}).Action("Sync")),
		),

		s.Router("/proposals",
			s.Get(proposals.ProposalsHandler{}).Action("ListProposals"),
			s.Post(proposals.ProposalEventsHandler{}).Action("CreateProposal"),
			s.Router("/:id",
				s.Get(proposals.ProposalsHandler{}).Action("GetProposal")),
			s.Router("/over",
				s.Put(proposals.ProposalEventsHandler{}).Action("UpdateProposalOver")),
		),
		s.Router("/proposals",
			s.Filter(filters.LoginRequiredInner),
			s.Router("/me",
				s.Get(proposals.ProposalsHandler{}).Action("ListProposals")),
		),
		s.Router("/proposal/:id",
			// s.Filter(filters.LoginRequiredInner),
			s.Post(proposals.ProposalEventsHandler{}).Action("UpdateProposalStatus"),
			s.Router("/vote", s.Post(proposals.ProposalEventsHandler{}).Action("VoteProposal")),
		),

		s.Router("/startups",
			s.Router("/:id",
				s.Router("/discos",
					s.Post(discos.DiscosHandler{}).Filter(filters.LoginRequiredInner).Action("CreateStartupDisco"),
					s.Get(discos.DiscosHandler{}).Action("GetStartupDisco"),
					s.Router("/investors",
						s.Post(discos.DiscosInvestorsHandler{}).Filter(filters.LoginRequiredInner).Action("CreateStartupDiscoInvestor"),
						s.Get(discos.DiscosInvestorsHandler{}).Action("ListStartupDiscoInvestor"),
					),
				),
				s.Router("/discoSwapState",
					s.Get(discos.DiscosHandler{}).Action("GetDiscoSwapState"),
				),
			),
		),

		s.Router("/discos",
			s.Get(discos.DiscosHandler{}).Action("ListDisco"),
		),

		s.Router("/discos:statDiscoEthIncrease",
			s.Post(discos.DiscosHandler{}).Action("StatDiscoEthIncrease"),
		),
		s.Router("/discos:statDiscoEthTotal",
			s.Post(discos.DiscosHandler{}).Action("StatDiscoEthTotal"),
		),
		s.Router("/discos:statDiscoTotal",
			s.Post(discos.DiscosHandler{}).Action("StatDiscoTotal"),
		),
	)
}

func (p *appConfig) ConfigDone() {
}
