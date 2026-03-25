package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/OpenJenie/goserve/business/web/v1/auth"
	"github.com/OpenJenie/goserve/foundation/keystore"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	actionGenerateKeys  = "generate-keys"
	actionGenerateToken = "generate-token"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	var cfg struct {
		action     string
		keysFolder string
		kid        string
		issuer     string
		subject    string
		roles      string
		ttl        time.Duration
	}

	flag.StringVar(&cfg.action, "action", actionGenerateToken, "action to perform: generate-keys or generate-token")
	flag.StringVar(&cfg.keysFolder, "keys-folder", ".local/keys", "folder containing PEM private keys")
	flag.StringVar(&cfg.kid, "kid", "", "key id to use; defaults to the first available key")
	flag.StringVar(&cfg.issuer, "issuer", "service project", "issuer claim")
	flag.StringVar(&cfg.subject, "subject", "12345678789", "subject claim")
	flag.StringVar(&cfg.roles, "roles", "ADMIN", "comma-separated roles")
	flag.DurationVar(&cfg.ttl, "ttl", 24*time.Hour, "token lifetime")
	flag.Parse()

	switch cfg.action {
	case actionGenerateKeys:
		return generateKeys(cfg.keysFolder, cfg.kid)
	case actionGenerateToken:
		return generateToken(cfg.keysFolder, cfg.kid, cfg.issuer, cfg.subject, cfg.roles, cfg.ttl)
	default:
		return fmt.Errorf("unknown action %q", cfg.action)
	}
}

func generateKeys(keysFolder string, kid string) error {
	if kid == "" {
		kid = uuid.NewString()
	}

	if err := os.MkdirAll(keysFolder, 0o755); err != nil {
		return fmt.Errorf("creating keys folder: %w", err)
	}

	privatePEM, publicPEM, err := keystore.GenerateRSAKeyPair(2048)
	if err != nil {
		return err
	}

	keyPath := filepath.Join(keysFolder, kid+".pem")
	if err := os.WriteFile(keyPath, privatePEM, 0o600); err != nil {
		return fmt.Errorf("writing private key: %w", err)
	}

	fmt.Printf("kid=%s\n", kid)
	fmt.Printf("key_path=%s\n", keyPath)
	fmt.Printf("%s", publicPEM)
	return nil
}

func generateToken(keysFolder string, kid string, issuer string, subject string, roles string, ttl time.Duration) error {
	if ttl <= 0 {
		return fmt.Errorf("ttl must be greater than zero")
	}

	if kid == "" {
		var err error
		kid, err = findFirstKID(keysFolder)
		if err != nil {
			return err
		}
	}

	ks, err := keystore.NewFS(os.DirFS(keysFolder))
	if err != nil {
		return fmt.Errorf("loading keys: %w", err)
	}

	authn, err := auth.New(auth.Config{
		KeyLookup: ks,
		Issuer:    issuer,
	})
	if err != nil {
		return fmt.Errorf("constructing auth: %w", err)
	}

	now := time.Now().UTC()
	claims := auth.Claims{
		RegisteredClaims: auth.Claims{}.RegisteredClaims,
		Roles:            parseRoles(roles),
	}
	claims.Subject = subject
	claims.Issuer = issuer
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.ExpiresAt = jwt.NewNumericDate(now.Add(ttl))

	token, err := authn.GenerateToken(kid, claims)
	if err != nil {
		return fmt.Errorf("generating token: %w", err)
	}

	fmt.Println(token)
	return nil
}

func findFirstKID(keysFolder string) (string, error) {
	entries, err := os.ReadDir(keysFolder)
	if err != nil {
		return "", fmt.Errorf("reading keys folder: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".pem" {
			continue
		}
		return strings.TrimSuffix(entry.Name(), ".pem"), nil
	}

	return "", fmt.Errorf("no PEM keys found in %s; run `make keys` first", keysFolder)
}

func parseRoles(raw string) []string {
	parts := strings.Split(raw, ",")
	roles := make([]string, 0, len(parts))
	for _, part := range parts {
		role := strings.TrimSpace(part)
		if role == "" {
			continue
		}
		roles = append(roles, role)
	}
	return roles
}
