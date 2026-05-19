package bootstrap

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Tangyd893/TMS-Go/backend/internal/config"
	"github.com/Tangyd893/TMS-Go/backend/internal/middleware"
	authapi "github.com/Tangyd893/TMS-Go/backend/internal/modules/auth/api"
	authrepo "github.com/Tangyd893/TMS-Go/backend/internal/modules/auth/repository"
	authsvc "github.com/Tangyd893/TMS-Go/backend/internal/modules/auth/service"
	baseapi "github.com/Tangyd893/TMS-Go/backend/internal/modules/base/api"
	basemodel "github.com/Tangyd893/TMS-Go/backend/internal/modules/base/model"
	baserepo "github.com/Tangyd893/TMS-Go/backend/internal/modules/base/repository"
	basesvc "github.com/Tangyd893/TMS-Go/backend/internal/modules/base/service"
	dispatchapi "github.com/Tangyd893/TMS-Go/backend/internal/modules/dispatch/api"
	dispatchrepo "github.com/Tangyd893/TMS-Go/backend/internal/modules/dispatch/repository"
	dispatchsvc "github.com/Tangyd893/TMS-Go/backend/internal/modules/dispatch/service"
	healthapi "github.com/Tangyd893/TMS-Go/backend/internal/modules/health/api"
	healthsvc "github.com/Tangyd893/TMS-Go/backend/internal/modules/health/service"
	orderapi "github.com/Tangyd893/TMS-Go/backend/internal/modules/order/api"
	orderrepo "github.com/Tangyd893/TMS-Go/backend/internal/modules/order/repository"
	ordersvc "github.com/Tangyd893/TMS-Go/backend/internal/modules/order/service"
	sysapi "github.com/Tangyd893/TMS-Go/backend/internal/modules/system/api"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/system/dict"
	sysrepo "github.com/Tangyd893/TMS-Go/backend/internal/modules/system/repository"
	syssvc "github.com/Tangyd893/TMS-Go/backend/internal/modules/system/service"
	transportapi "github.com/Tangyd893/TMS-Go/backend/internal/modules/transport/api"
	transportrepo "github.com/Tangyd893/TMS-Go/backend/internal/modules/transport/repository"
	transportsvc "github.com/Tangyd893/TMS-Go/backend/internal/modules/transport/service"
	exceptionapi "github.com/Tangyd893/TMS-Go/backend/internal/modules/exception/api"
	exceptionrepo "github.com/Tangyd893/TMS-Go/backend/internal/modules/exception/repository"
	exceptionsvc "github.com/Tangyd893/TMS-Go/backend/internal/modules/exception/service"
	fileapi "github.com/Tangyd893/TMS-Go/backend/internal/modules/file/api"
	filerepo "github.com/Tangyd893/TMS-Go/backend/internal/modules/file/repository"
	filesvc "github.com/Tangyd893/TMS-Go/backend/internal/modules/file/service"
	financeapi "github.com/Tangyd893/TMS-Go/backend/internal/modules/finance/api"
	financerepo "github.com/Tangyd893/TMS-Go/backend/internal/modules/finance/repository"
	financesvc "github.com/Tangyd893/TMS-Go/backend/internal/modules/finance/service"
	reportapi "github.com/Tangyd893/TMS-Go/backend/internal/modules/report/api"
	reportsvc "github.com/Tangyd893/TMS-Go/backend/internal/modules/report/service"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/jwt"
	"gorm.io/gorm"
)

func NewHTTPServer(cfg config.Config, log *slog.Logger, db *gorm.DB) *http.Server {
	accessTTL, _ := time.ParseDuration(cfg.JWT.AccessTokenTTL)
	if accessTTL == 0 { accessTTL = 2 * time.Hour }
	refreshTTL, _ := time.ParseDuration(cfg.JWT.RefreshTokenTTL)
	if refreshTTL == 0 { refreshTTL = 168 * time.Hour }

	jwtManager := jwt.NewManager(cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret, accessTTL, refreshTTL, cfg.JWT.Issuer)
	authMiddleware := middleware.Auth(jwtManager)

	userRepo := authrepo.NewUserRepository(db)
	authService := authsvc.NewAuthService(userRepo, jwtManager, log)
	authHandler := authapi.NewHandler(authService)

	sysUserRepo := sysrepo.NewUserRepository(db)
	sysUserService := syssvc.NewUserService(sysUserRepo)
	sysUserHandler := sysapi.NewHandler(sysUserService)

	customerRepo := baserepo.NewCrudRepository[basemodel.Customer](db)
	customerSvc := basesvc.NewEntityService(customerRepo, []string{"name", "code"})
	customerHandler := baseapi.NewCrudHandler(customerSvc)

	carrierRepo := baserepo.NewCrudRepository[basemodel.Carrier](db)
	carrierSvc := basesvc.NewEntityService(carrierRepo, []string{"name", "code"})
	carrierHandler := baseapi.NewCrudHandler(carrierSvc)

	vehicleRepo := baserepo.NewCrudRepository[basemodel.Vehicle](db)
	vehicleSvc := basesvc.NewEntityService(vehicleRepo, []string{"plate_no", "brand_model"})
	vehicleHandler := baseapi.NewCrudHandler(vehicleSvc)

	driverRepo := baserepo.NewCrudRepository[basemodel.Driver](db)
	driverSvc := basesvc.NewEntityService(driverRepo, []string{"name", "code", "phone"})
	driverHandler := baseapi.NewCrudHandler(driverSvc)

	routeRepo := baserepo.NewCrudRepository[basemodel.Route](db)
	routeSvc := basesvc.NewEntityService(routeRepo, []string{"name", "code"})
	routeHandler := baseapi.NewCrudHandler(routeSvc)

	stationRepo := baserepo.NewCrudRepository[basemodel.Station](db)
	stationSvc := basesvc.NewEntityService(stationRepo, []string{"name", "code", "city"})
	stationHandler := baseapi.NewCrudHandler(stationSvc)

	dictRepo := dict.NewRepository(db)
	dictService := dict.NewService(dictRepo)
	dictHandler := dict.NewHandler(dictService)

	orderRepo := orderrepo.NewOrderRepository(db)
	orderService := ordersvc.NewOrderService(orderRepo)
	orderHandler := orderapi.NewHandler(orderService)

	dispatchRepo := dispatchrepo.NewDispatchRepository(db)
	dispatchService := dispatchsvc.NewDispatchService(dispatchRepo)
	dispatchHandler := dispatchapi.NewHandler(dispatchService)

	transportRepo := transportrepo.NewTransportRepository(db)
	transportService := transportsvc.NewTransportService(transportRepo)
	transportHandler := transportapi.NewHandler(transportService)

	exceptionRepo := exceptionrepo.NewExceptionRepository(db)
	exceptionService := exceptionsvc.NewExceptionService(exceptionRepo)
	exceptionHandler := exceptionapi.NewHandler(exceptionService)

	fileRepo := filerepo.NewFileRepository(db)
	fileService := filesvc.NewFileService(fileRepo, cfg.MinIO)
	fileHandler := fileapi.NewHandler(fileService)

	financeRepo := financerepo.NewFinanceRepository(db)
	financeService := financesvc.NewFinanceService(financeRepo)
	financeHandler := financeapi.NewHandler(financeService)

	reportService := reportsvc.NewReportService(db, log)
	reportHandler := reportapi.NewHandler(reportService)

	mux := http.NewServeMux()

	healthService := healthsvc.NewService()
	healthHandler := healthapi.NewHandler(healthService)
	healthapi.RegisterRoutes(mux, healthHandler)

	authapi.RegisterRoutes(mux, authHandler)
	mux.Handle("GET /api/v1/auth/me", authMiddleware(http.HandlerFunc(authHandler.GetCurrentUser)))

	sysUserSubMux := http.NewServeMux()
	sysapi.RegisterRoutes(sysUserSubMux, sysUserHandler)
	mux.Handle("/api/v1/users", authMiddleware(sysUserSubMux))
	mux.Handle("/api/v1/users/", authMiddleware(sysUserSubMux))

	baseSubMux := http.NewServeMux()
	baseapi.RegisterCrudRoutes(baseSubMux, "/api/v1/customers", customerHandler)
	baseapi.RegisterCrudRoutes(baseSubMux, "/api/v1/carriers", carrierHandler)
	baseapi.RegisterCrudRoutes(baseSubMux, "/api/v1/vehicles", vehicleHandler)
	baseapi.RegisterCrudRoutes(baseSubMux, "/api/v1/drivers", driverHandler)
	baseapi.RegisterCrudRoutes(baseSubMux, "/api/v1/routes", routeHandler)
	baseapi.RegisterCrudRoutes(baseSubMux, "/api/v1/stations", stationHandler)
	for _, prefix := range []string{"/api/v1/customers", "/api/v1/carriers", "/api/v1/vehicles", "/api/v1/drivers", "/api/v1/routes", "/api/v1/stations"} {
		mux.Handle(prefix, authMiddleware(baseSubMux))
		mux.Handle(prefix+"/", authMiddleware(baseSubMux))
	}

	dictSubMux := http.NewServeMux()
	dict.RegisterRoutes(dictSubMux, dictHandler)
	mux.Handle("/api/v1/dict/", authMiddleware(dictSubMux))

	orderSubMux := http.NewServeMux()
	orderapi.RegisterRoutes(orderSubMux, orderHandler)
	mux.Handle("/api/v1/orders", authMiddleware(orderSubMux))
	mux.Handle("/api/v1/orders/", authMiddleware(orderSubMux))

	dispatchSubMux := http.NewServeMux()
	dispatchapi.RegisterRoutes(dispatchSubMux, dispatchHandler)
	mux.Handle("/api/v1/dispatch/", authMiddleware(dispatchSubMux))

	transportSubMux := http.NewServeMux()
	transportapi.RegisterRoutes(transportSubMux, transportHandler)
	exceptionapi.RegisterRoutes(transportSubMux, exceptionHandler)
	mux.Handle("/api/v1/transport/", authMiddleware(transportSubMux))

	fileSubMux := http.NewServeMux()
	fileapi.RegisterRoutes(fileSubMux, fileHandler)
	mux.Handle("/api/v1/files", authMiddleware(fileSubMux))
	mux.Handle("/api/v1/files/", authMiddleware(fileSubMux))

	financeSubMux := http.NewServeMux()
	financeapi.RegisterRoutes(financeSubMux, financeHandler)
	mux.Handle("/api/v1/finance/", authMiddleware(financeSubMux))

	reportSubMux := http.NewServeMux()
	reportapi.RegisterRoutes(reportSubMux, reportHandler)
	mux.Handle("/api/v1/reports/", authMiddleware(reportSubMux))

	handler := middleware.RequestID(middleware.AccessLog(log)(mux))

	return &http.Server{Addr: ":" + cfg.App.Port, Handler: handler}
}
