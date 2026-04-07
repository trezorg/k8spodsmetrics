package client

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"
)

func TestFindKubeConfig(t *testing.T) {
	t.Run("from env", func(t *testing.T) {
		expectedPath := "/custom/kube/config"
		t.Setenv("KUBECONFIG", expectedPath)

		path, err := FindKubeConfig()
		require.NoError(t, err)
		require.Equal(t, expectedPath, path)
	})

	t.Run("default path when file exists", func(t *testing.T) {
		t.Setenv("KUBECONFIG", "")
		homeDir := t.TempDir()
		t.Setenv("HOME", homeDir)
		configPath := filepath.Join(homeDir, ".kube", "config")
		require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0o755))
		require.NoError(t, os.WriteFile(configPath, []byte("apiVersion: v1\n"), 0o644))

		path, err := FindKubeConfig()
		require.NoError(t, err)
		require.Equal(t, configPath, path)
	})

	t.Run("empty when env is unset and default path is missing", func(t *testing.T) {
		t.Setenv("KUBECONFIG", "")
		homeDir := t.TempDir()
		t.Setenv("HOME", homeDir)

		path, err := FindKubeConfig()
		require.NoError(t, err)
		require.Empty(t, path)
	})
}

func TestRestConfig(t *testing.T) {
	t.Run("invalid kubeconfig", func(t *testing.T) {
		_, err := restConfig("/nonexistent/config", "")
		require.Error(t, err)
	})

	t.Run("in cluster config fallback", func(t *testing.T) {
		_, err := restConfig("", "")
		require.Error(t, err)
	})
}

func TestClients(t *testing.T) {
	t.Run("invalid config", func(t *testing.T) {
		mc, pc, err := Clients("/invalid/path", "")
		require.Error(t, err)
		require.Nil(t, mc)
		require.Nil(t, pc)
	})
}

func TestCoreV1Client(t *testing.T) {
	t.Run("invalid config", func(t *testing.T) {
		pc, err := CoreV1Client("/invalid/path", "")
		require.Error(t, err)
		require.Nil(t, pc)
	})
}

func TestForMetrics(t *testing.T) {
	t.Run("invalid config", func(t *testing.T) {
		mc, err := ForMetrics("/invalid/path", "")
		require.Error(t, err)
		require.Nil(t, mc)
	})
}

func TestApplyRateLimit(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		t.Setenv(clientQPSEnvVar, "")
		t.Setenv(clientBurstEnvVar, "")
		cfg := &rest.Config{}

		applyRateLimit(cfg)

		require.Equal(t, defaultClientQPS, cfg.QPS)
		require.Equal(t, defaultClientBurst, cfg.Burst)
		require.NotNil(t, cfg.RateLimiter)
	})

	t.Run("env overrides", func(t *testing.T) {
		t.Setenv(clientQPSEnvVar, "15.5")
		t.Setenv(clientBurstEnvVar, "30")
		cfg := &rest.Config{}

		applyRateLimit(cfg)

		require.Equal(t, float32(15.5), cfg.QPS)
		require.Equal(t, 30, cfg.Burst)
		require.NotNil(t, cfg.RateLimiter)
	})

	t.Run("invalid env values fallback to defaults", func(t *testing.T) {
		t.Setenv(clientQPSEnvVar, "not-a-number")
		t.Setenv(clientBurstEnvVar, "-1")
		cfg := &rest.Config{}

		applyRateLimit(cfg)

		require.Equal(t, defaultClientQPS, cfg.QPS)
		require.Equal(t, defaultClientBurst, cfg.Burst)
		require.NotNil(t, cfg.RateLimiter)
	})
}
