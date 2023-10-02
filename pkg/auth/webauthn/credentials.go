package webauthn

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/go-webauthn/webauthn/webauthn"
	"golang.org/x/crypto/argon2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	cosmov1alpha1 "github.com/cosmo-workspace/cosmo/api/v1alpha1"
	"github.com/cosmo-workspace/cosmo/pkg/kosmo"
)

func GetUser(ctx context.Context, c kosmo.Client, userName string) (*User, error) {
	cosmoUser, err := c.GetUser(ctx, userName)
	if err != nil {
		return nil, err
	}
	u := User{User: *cosmoUser, client: c}
	l, err := NewCredentialList(ctx, c, userName)
	if err != nil {
		return nil, err
	}
	u.CredentialList = l
	return &u, nil
}

// User implements webauthn.User interface
// https://pkg.go.dev/github.com/go-webauthn/webauthn@v0.8.6/webauthn#User
type User struct {
	cosmov1alpha1.User
	CredentialList *CredentialList

	client kosmo.Client
}

func (u *User) WebAuthnID() []byte {
	id := make([]byte, 64)
	hashed := argon2.IDKey([]byte("tom"), nil, 1, 2048, 4, 32)
	n := hex.Encode(id, hashed)
	if n != 64 {
		panic(fmt.Errorf("invalid hash length: n=%d", n))
	}
	return id
}
func (u *User) WebAuthnName() string {
	return u.Spec.DisplayName
}
func (u *User) WebAuthnDisplayName() string {
	return u.Spec.DisplayName
}
func (u *User) WebAuthnCredentials() []webauthn.Credential {
	c := make([]webauthn.Credential, len(u.CredentialList.Creds))
	for i, cred := range u.CredentialList.Creds {
		c[i] = cred.Cred
	}
	return c
}

func (u *User) WebAuthnIcon() string {
	return ""
}

// RegisterCredential store credential to secret
func (u *User) RegisterCredential(ctx context.Context, cred *Credential) error {
	c, err := NewCredentialList(ctx, u.client, u.Name)
	if err != nil {
		return err
	}
	c.add(cred)
	return c.save(ctx)
}

// RemoveCredential removes credential in secret
func (u *User) RemoveCredential(ctx context.Context, credID string) error {
	c, err := NewCredentialList(ctx, u.client, u.Name)
	if err != nil {
		return err
	}
	ok := c.remove(credID)
	if !ok {
		return fmt.Errorf("credential not found")
	}
	return c.save(ctx)
}

// UpdateCredential updates credential in secret
func (u *User) UpdateCredential(ctx context.Context, credID string, displayName *string) error {
	c, err := NewCredentialList(ctx, u.client, u.Name)
	if err != nil {
		return err
	}

	// update display name if not nil
	if displayName != nil {
		for i, v := range c.Creds {
			if credID == base64.RawURLEncoding.EncodeToString(v.Cred.ID) {
				c.Creds[i].DisplayName = *displayName
				break
			}
		}
	}
	return c.save(ctx)
}

// ListCredentials returns list of registered credential
func (u *User) ListCredentials(ctx context.Context) ([]Credential, error) {
	l, err := NewCredentialList(ctx, u.client, u.Name)
	return l.Creds, err
}

type CredentialList struct {
	Creds []Credential

	client kosmo.Client
	sec    *corev1.Secret
}

type Credential struct {
	DisplayName string
	Timestamp   int64
	Cred        webauthn.Credential
}

const (
	CredentialSecretName string = "cosmo-user-creds"
	CredentialListKey    string = "credentials"
)

func NewCredentialList(ctx context.Context, c kosmo.Client, userName string) (*CredentialList, error) {
	cl := CredentialList{client: c}
	var sec corev1.Secret
	sec.SetName(CredentialSecretName)
	sec.SetNamespace(cosmov1alpha1.UserNamespace(userName))
	cosmov1alpha1.SetControllerManaged(&sec)

	if err := c.Get(ctx, types.NamespacedName{Name: sec.Name, Namespace: sec.Namespace}, &sec); err != nil {
		if !errors.IsNotFound(err) {
			return nil, fmt.Errorf("failed to get credential store: %w", err)
		}
	}
	if len(sec.Data) == 0 {
		sec.Data = map[string][]byte{CredentialListKey: []byte(`{"creds": []}`)}
	}
	if _, ok := sec.Data[CredentialListKey]; !ok {
		sec.Data[CredentialListKey] = []byte(`{"creds": []}`)
	}
	cl.sec = &sec

	if err := json.Unmarshal(sec.Data[CredentialListKey], &cl); err != nil {
		return nil, fmt.Errorf("failed to load credential list: %w", err)
	}
	return &cl, nil
}

func (c *CredentialList) add(cred *Credential) {
	notfound := true
	for i, v := range c.Creds {
		if base64.RawURLEncoding.EncodeToString(cred.Cred.ID) == base64.RawURLEncoding.EncodeToString(v.Cred.ID) {
			c.Creds[i] = *cred
			notfound = false
			break
		}
	}
	if notfound {
		c.Creds = append(c.Creds, *cred)
	}
}

func (c *CredentialList) remove(id string) bool {
	for i, v := range c.Creds {
		if id == base64.RawURLEncoding.EncodeToString(v.Cred.ID) {
			c.Creds = append(c.Creds[:i], c.Creds[i+1:]...)
			return true
		}
	}
	return false
}

func (c *CredentialList) save(ctx context.Context) error {
	raw, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to dump credential list: %w", err)
	}
	c.sec.Data[CredentialListKey] = raw
	return c.updateSecret(ctx)
}

func (c *CredentialList) updateSecret(ctx context.Context) error {
	if err := c.client.Update(ctx, c.sec); err != nil {
		if errors.IsNotFound(err) {
			if err := c.client.Create(ctx, c.sec); err != nil {
				return fmt.Errorf("failed to create credential secret: %w", err)
			}
		} else {
			return fmt.Errorf("failed to update credential store: %w", err)
		}
	}
	return nil
}
