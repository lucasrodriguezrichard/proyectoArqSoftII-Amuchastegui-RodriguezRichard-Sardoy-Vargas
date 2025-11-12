// internal/domain/user.go
package domain

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"
)

/*
Dominio puro para users-api del sistema Restaurant:
- Entidad User
- DTOs de entrada/salida para Register/Login
- Errores de dominio
- Interfaces de puertos (Repository, PasswordHasher, TokenIssuer, Service)
No hay dependencias de infraestructura ni librerías criptográficas aquí.
*/

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID           uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Username     string    `json:"username" gorm:"size:32;uniqueIndex:ux_users_username;not null"`
	Email        string    `json:"email" gorm:"size:191;uniqueIndex:ux_users_email;not null"`
	FirstName    string    `json:"first_name" gorm:"size:100;not null"`
	LastName     string    `json:"last_name" gorm:"size:100;not null"`
	PasswordHash string    `json:"-" gorm:"size:255;not null"` // nunca exponer
	Role         Role      `json:"role" gorm:"type:varchar(20);default:user;not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	// Opcional: OwnerID si se requiere ownership cruzado en otras entidades
}

// Nombre de tabla explícito para evitar ambigüedad en GORM
func (User) TableName() string { return "users" }

// =========================
// DTOs
// =========================

type RegisterInput struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	// Rol fijo a user en registro público; admin sólo por backoffice/seeding.
}

type LoginInput struct {
	// Identificador puede ser email o username (lo resolvemos en servicio)
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type AuthTokens struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// =========================
// Validaciones de dominio
// =========================

var (
	ErrInvalidInput       = errors.New("invalid_input")
	ErrInvalidCredentials = errors.New("invalid_credentials")
	ErrUserExists         = errors.New("user_already_exists")
	ErrUserNotFound       = errors.New("user_not_found")
)

var (
	emailRx    = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	usernameRx = regexp.MustCompile(`^[a-zA-Z0-9._-]{3,32}$`)
)

func (in *RegisterInput) Normalize() {
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	in.Username = strings.TrimSpace(in.Username)
	in.FirstName = strings.TrimSpace(in.FirstName)
	in.LastName = strings.TrimSpace(in.LastName)
}

func (in RegisterInput) Validate() error {
	if !usernameRx.MatchString(in.Username) {
		return ErrInvalidInput
	}
	if !emailRx.MatchString(strings.ToLower(in.Email)) {
		return ErrInvalidInput
	}
	if len(in.Password) < 8 { // mínimo razonable
		return ErrInvalidInput
	}
	return nil
}

func (in *LoginInput) Normalize() {
	in.Identifier = strings.ToLower(strings.TrimSpace(in.Identifier))
}

func (in LoginInput) Validate() error {
	if in.Identifier == "" || in.Password == "" {
		return ErrInvalidInput
	}
	return nil
}

// =========================
// Puertos (interfaces)
// =========================

// Persistencia
type UserRepository interface {
	// Búsquedas
	GetByID(ctx context.Context, id uint64) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)

	// Altas/updates
	Create(ctx context.Context, u *User) error
	UpdatePasswordHash(ctx context.Context, id uint64, newHash string) error

	// Conflictos
	ExistsByEmailOrUsername(ctx context.Context, email, username string) (bool, error)
}

// Hash de contraseña (infra la implementa con bcrypt/argon2)
type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(hash string, plain string) bool
}

// Emisión de tokens (JWT)
type TokenIssuer interface {
	IssueAccessToken(u User) (token string, exp time.Time, err error)
	IssueRefreshToken(u User) (token string, err error) // opcional según enunciado
}

// Caso de uso (aplicación)
type UserService interface {
	// Registro público: crea usuario RoleUser y devuelve entidad (sin tokens).
	Register(ctx context.Context, in RegisterInput) (User, error)

	// Login: acepta email o username en Identifier, devuelve tokens y usuario.
	Login(ctx context.Context, in LoginInput) (AuthTokens, User, error)

	// GetByID obtiene un usuario por ID (usado por otros microservicios para validación).
	GetByID(ctx context.Context, id uint64) (User, error)

	// CreateAdmin crea un usuario admin (solo para admin).
	CreateAdmin(ctx context.Context, in RegisterInput) (User, error)
}

// =========================
// Reglas de dominio auxiliares
// =========================

func NewUserFromRegister(in RegisterInput, passwordHash string) User {
	now := time.Now().UTC()
	return User{
		Username:     in.Username,
		Email:        strings.ToLower(in.Email),
		FirstName:    in.FirstName,
		LastName:     in.LastName,
		PasswordHash: passwordHash,
		Role:         RoleUser,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// Política mínima para admin (no usada en registro público)
func (u *User) ElevateToAdmin() {
	u.Role = RoleAdmin
	u.UpdatedAt = time.Now().UTC()
}
