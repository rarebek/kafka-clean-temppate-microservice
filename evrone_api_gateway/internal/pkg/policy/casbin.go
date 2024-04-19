package policy

//import (
//	"fmt"
//	"github.com/redis/go-redis/v9"
//
//	"github.com/casbin/casbin/v2"
//	"github.com/casbin/casbin/v2/model"
//	rediswatcher "github.com/casbin/redis-watcher/v2"
//	"go.uber.org/zap"
//
//	"evrone_api_gateway/internal/pkg/config"
//	"evrone_api_gateway/internal/pkg/postgres"
//)
//
//func NewCachedEnforcer(cfg *config.Config, logger *zap.Logger) (*casbin.CachedEnforcer, error) {
//	// initializing casbin model
//	m := model.NewModel()
//	m.AddDef("r", "r", "sub, obj, act")
//	m.AddDef("p", "p", "sub, obj, act")
//	m.AddDef("g", "g", "_, _")
//	m.AddDef("e", "e", "some(where (p.eft == allow))")
//	m.AddDef("m", "m", "g(r.sub, p.sub) && keyMatch4(r.obj, p.obj) && regexMatch(r.act, p.act)")
//	//initializing pgx adapter
//	adapter, err := postgres.GetAdapter(cfg)
//	if err != nil {
//		return nil, fmt.Errorf("NewCachedEnforcer GetAdapter: %w", err)
//	}
//	enforcer, err := casbin.NewCachedEnforcer(m, adapter)
//	if err != nil {
//		return nil, fmt.Errorf("NewCachedEnforcer: %w", err)
//	}
//	// initializing watcher
//	err = initializingWatcher(cfg, logger, enforcer)
//	if err != nil {
//		return nil, fmt.Errorf("NewCachedEnforcer: %w", err)
//	}
//	return enforcer, nil
//}
////
////func initializingWatcher(cfg *config.Config, logger *zap.Logger, enforcer *casbin.CachedEnforcer) error {
////	w, _ := rediswatcher.NewWatcher(fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port), rediswatcher.WatcherOptions{
////		Options: redis.Options{
////			Network:  "tcp",
////			Password: cfg.Redis.Password,
////		},
////		IgnoreSelf: true,
////		Channel:    "/casbin_watcher",
////	})
////	// set the watcher for the enforcer.
////	//err := enforcer.SetWatcher(w)
////	//if err != nil {
////	//	return fmt.Errorf("SetWatcher: %w", err)
////	//}
////	//set callback
////	//err = w.SetUpdateCallback(func(s string) {
////	//	if err := enforcer.LoadPolicy(); err != nil {
////	//		logger.Error("enforcer watcher LoadPolicy", zap.Error(err))
////	//	}
////	//	logger.Info("enforcer watcher", zap.String("callback", s))
////	//})
////	//if err != nil {
////	//	return fmt.Errorf("SetUpdateCallback: %w", err)
////	//}
////	//return nil
////}
