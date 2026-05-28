// Package zoobridgeshim tests the Zoo Bridge tenant configuration
// against the canonical upstream tenant.Config schema. The shim itself
// ships no Go binary — the production image overlays this tenant.yaml
// on top of the upstream OSS bridge container.
package zoobridgeshim

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/luxfi/bridge/pkg/tenant"
)

// TestTenantConfigValidates is the single guardrail this shim ships.
func TestTenantConfigValidates(t *testing.T) {
	_, thisFile, _, _ := runtime.Caller(0)
	cfgPath := filepath.Join(filepath.Dir(thisFile), "tenant.yaml")

	cfg, err := tenant.Load(cfgPath)
	if err != nil {
		t.Fatalf("tenant.Load(%s): %v", cfgPath, err)
	}

	if cfg.Brand.Slug != "zoo-bridge" {
		t.Errorf("Brand.Slug: got %q, want zoo-bridge", cfg.Brand.Slug)
	}
	if cfg.Network.ID != 200200 {
		t.Errorf("Network.ID: got %d, want 200200 (Zoo Mainnet)", cfg.Network.ID)
	}
	if cfg.IAM.Endpoint != "https://zoo.id" {
		t.Errorf("IAM.Endpoint: got %q, want https://zoo.id", cfg.IAM.Endpoint)
	}
	if cfg.PQProfile != "strict-pq" {
		t.Errorf("PQProfile: got %q, want strict-pq", cfg.PQProfile)
	}
	if cfg.Domain != "bridge.zoo.network" {
		t.Errorf("Domain: got %q, want bridge.zoo.network", cfg.Domain)
	}
}

// TestSupportedChainsScopedToZoo asserts the chain allowlist matches
// Zoo's deployment scope.
func TestSupportedChainsScopedToZoo(t *testing.T) {
	_, thisFile, _, _ := runtime.Caller(0)
	cfgPath := filepath.Join(filepath.Dir(thisFile), "tenant.yaml")
	cfg, err := tenant.Load(cfgPath)
	if err != nil {
		t.Fatalf("tenant.Load: %v", err)
	}
	for _, c := range []string{"eth", "btc", "sol", "ton", "xrp", "dot", "zoo"} {
		if !cfg.IsChainSupported(c) {
			t.Errorf("expected chain %q to be supported", c)
		}
	}
	for _, c := range []string{"liquid", "hanzo", "pars", "lux"} {
		if cfg.IsChainSupported(c) {
			t.Errorf("chain %q must NOT be supported by Zoo Bridge", c)
		}
	}
}
