package bootstrap

import (
	"fmt"

	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/config"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/logging"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/tracing"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/metrics"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/cache"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/policy"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/i18n"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/flags"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/validate"
)

// EnhancedContainer holds all enterprise components
type EnhancedContainer struct {
	Config          *config.Config
	Logger          logging.Logger
	Tracer          tracing.Tracer
	Metrics         *metrics.Metrics
	Cache           cache.Cache
	PolicyEngine    *policy.PolicyEngine
	I18nManager     *i18n.I18nManager
	FeatureFlags    *flags.FeatureFlags
	SchemaValidator *validate.SchemaValidator
}

func NewEnhancedContainer(cfg *config.Config) (*EnhancedContainer, error) {
	// Initialize logger
	logger := logging.NewStructuredLogger("info")

	// Initialize tracer
	tracer, err := tracing.NewTracer("reference-master-geopolitical")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tracer: %w", err)
	}

	// Initialize metrics
	metrics, err := metrics.NewMetrics("reference-master-geopolitical")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize metrics: %w", err)
	}

	// Initialize cache
	cache := cache.NewCache(
		fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		cfg.Redis.Password,
		cfg.Redis.DB,
	)

	// Initialize policy engine
	policyEngine := policy.NewPolicyEngine()
	policyEngine.LoadDefaultRules()

	// Initialize i18n manager
	i18nManager := i18n.NewI18nManager("en-US")

	// Initialize feature flags
	featureFlags := flags.NewFeatureFlags()

	// Initialize schema validator
	schemaValidator := validate.NewSchemaValidator()

	return &EnhancedContainer{
		Config:          cfg,
		Logger:          logger,
		Tracer:          tracer,
		Metrics:         metrics,
		Cache:           cache,
		PolicyEngine:    policyEngine,
		I18nManager:     i18nManager,
		FeatureFlags:    featureFlags,
		SchemaValidator: schemaValidator,
	}, nil
}