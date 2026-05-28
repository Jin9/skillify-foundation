package config

import (
	"time"

	commonconfig "gitlab.com/b2c-e-commerce-platform/platform/backend/common/config"
	envlib "github.com/caarlos0/env/v11"
)

// Config holds all environment-derived configuration for the checkout service.
type Config struct {
	Server        ServerConfig
	Postgres      PostgresConfig
	JWT           JWTConfig
	Downstream    DownstreamConfig
	InternalAuth  InternalAuthConfig
	AccessControl AccessControlConfig
	Header        HeaderConfig
	HttpClient    HttpClientConfig
}

type ServerConfig struct {
	Port string `env:"SERVER_PORT" envDefault:"8080"`
}

type PostgresConfig struct {
	DSN string `env:"POSTGRES_DSN,required"`
}

type JWTConfig struct {
	PublicKey string `env:"JWT_PUBLIC_KEY,required"`
	Issuer    string `env:"JWT_ISSUER" envDefault:"shoppilot-identity"`
	Audience  string `env:"JWT_AUDIENCE" envDefault:"shoppilot-api"`
}

// DownstreamConfig holds base URLs and per-call timeouts for every upstream service.
// Timeouts are per the TD spec §timeouts_and_retries_summary.
type DownstreamConfig struct {
	// Cart service
	CartBaseURL           string        `env:"CART_BASE_URL,required"`
	CartReadTimeout       time.Duration `env:"CART_READ_TIMEOUT" envDefault:"800ms"`
	CartClearTimeout      time.Duration `env:"CART_CLEAR_TIMEOUT" envDefault:"800ms"`

	// Catalog service
	CatalogBaseURL          string        `env:"CATALOG_BASE_URL,required"`
	CatalogDetailTimeout    time.Duration `env:"CATALOG_DETAIL_TIMEOUT" envDefault:"500ms"`

	// Inventory service
	InventoryBaseURL             string        `env:"INVENTORY_BASE_URL,required"`
	InventoryBulkReadTimeout     time.Duration `env:"INVENTORY_BULK_READ_TIMEOUT" envDefault:"800ms"`
	InventoryReservationTimeout  time.Duration `env:"INVENTORY_RESERVATION_TIMEOUT" envDefault:"1500ms"`
	InventoryReleaseTimeout      time.Duration `env:"INVENTORY_RELEASE_TIMEOUT" envDefault:"1000ms"`

	// Payment service
	PaymentBaseURL        string        `env:"PAYMENT_BASE_URL,required"`
	PaymentIntentTimeout  time.Duration `env:"PAYMENT_INTENT_TIMEOUT" envDefault:"1500ms"`

	// Order service (CHK-AMBIG-001: implied internal contract; routes agreed as POST /api/v1/order/order/create-from-checkout)
	OrderBaseURL             string        `env:"ORDER_BASE_URL,required"`
	OrderCreateTimeout       time.Duration `env:"ORDER_CREATE_TIMEOUT" envDefault:"1000ms"`
	OrderCancelTimeout       time.Duration `env:"ORDER_CANCEL_TIMEOUT" envDefault:"1000ms"`

	// Identity service
	IdentityBaseURL         string        `env:"IDENTITY_BASE_URL,required"`
	IdentityAddressTimeout  time.Duration `env:"IDENTITY_ADDRESS_TIMEOUT" envDefault:"500ms"`
}

// InternalAuthConfig holds the shared secret used for service-to-service calls
// that require internal authentication (e.g. order.create-from-checkout).
type InternalAuthConfig struct {
	SharedSecret string `env:"INTERNAL_SHARED_SECRET,required"`
}

type AccessControlConfig struct {
	AllowOrigin string `env:"ACCESS_CONTROL_ALLOW_ORIGIN" envDefault:"*"`
}

type HeaderConfig struct {
	RefIDHeaderKey string `env:"HEADER_REF_ID_KEY" envDefault:"X-Ref-ID"`
}

type HttpClientConfig struct {
	EnableLogDebug bool `env:"HTTP_CLIENT_ENABLE_LOG_DEBUG" envDefault:"false"`
}

// C parses Config from environment variables.
// Uses caarlos0/env under the hood (same library as common/config.ParseEnv).
func C(_ string) Config {
	cfg, err := commonconfig.ParseEnv[Config](envlib.Options{})
	if err != nil {
		panic("checkout: config parse error: " + err.Error())
	}
	return cfg
}
