## Membuat fitur login dengan jwt


### membuat method Login pada file handler ./handler/auth_handler.go

``` go
func (h *authHandler) Login(c *gin.Context) {

	var login *dto.LoginRequest

	err := c.ShouldBindJSON(&login)
	if err != nil {
		errorhandler.HandleError(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

}

```

### membuat definisi method GetUserByEmail() pada interface AuthRepository dan membuat method GetUserByEmail() pada file auth_repository.go ./repository/auth_repository.go

- `/repository/auth_repository.go`: Entry point for the auth repository

tambahakan definisi method GetUserByEmail() pada inteface AuthRepository

``` go
type AuthRepository interface {
	EmailExist(email string) bool
	Register(reqUser *entity.User) error
	GetUserByEmail(email string) (*entity.User, error) // tambah definisi method berikut
}

```

buat method untuk mencari data user berdasarkan email

``` go
// .....

func (r *authRepository) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User

	err := r.DB.First(&user, "email = ?", email).Error

	return &user, err
}

```



### Menambahkan struct LoginRequest dan LoginResponse pada file auth_dto.go ./dto/auth.go

tambahakan end point berikut untuk melakukan login

- `/dto/auth_dto.go`: Entry point for the auth dto
```go

    type LoginRequest struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    type LoginResponse struct {
        ID    int    `json:"id"`
        Name  string `json:"name"`
        Token string `json:"token"`
    }

```

### membuat function untuk verifikasi password pada file password.go ./helper/password.go



tambahakan definisi method GetUserByEmail() pada inteface AuthRepository

- `/helper/password.go`: helper for check password
``` go


package helper

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(passwordHash), err
}

// tambahkan function dibawah ini

func VerifyPassword(hashPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err
}

```

### membuat function untuk generate token jwt pada file token.go ./helper/token.go

import library berikut melalui terminal

```bash
go get github.com/golang-jwt/jwt/v4
```

tambahakan key pada struct berikut

- `/config/config.go`: helper for check generate token
```go

//...

type Config struct {
	PORT        string
	DB_USERNAME string
	DB_PASSWORD string
	DB_URL      string
	DB_DATABASE string
	JWT_KEY     string // tambahakan kode pada baris berikut 
}

// ...

```

kode berikut pada file token.go

- `/helper/token.go`: helper for check generate token
``` go

package helper

import (
	"go-socmed/config"
	"go-socmed/entity"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	ID int `json:"id"`
	jwt.RegisteredClaims
}

func GenerateToken(user *entity.User) (string, error) {

	claims := JWTClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(config.ENV.JWT_KEY))

	return ss, err

}

```


### membuat definisi method Login() pada interface AuthService dan membuat method Login() pada file auth_service ./service/auth_service.go

- `/service/auth_service.go`: Entry point for the auth service

tambahakan definisi method Login() pada inteface AuthService

``` go

type AuthService interface {
	Register(req *dto.RegisterRequest) error
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error) // tambah definisi method berikut
}

```

buat method untuk melakukan proses login

``` go
// .....

func (s *authService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {

	var data dto.LoginResponse

	// Pengecekan apakah email ada
	user, err := s.repository.GetUserByEmail(req.Email)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "wrong email or password!"}
	}

	// verifikasi password
	if err := helper.VerifyPassword(user.Password, req.Password); err != nil {
		return nil, &errorhandler.NotFoundError{Message: "wrong email or password!"}
	}

	token, err := helper.GenerateToken(user)
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	data = dto.LoginResponse{
		ID:    user.ID,
		Name:  user.Name,
		Token: token,
	}

	return &data, nil

}


```

### membuat method Login() pada file auth_handler ./handler/auth_handler.go

tambahakan handler untuk auth
- `/handler/auth_handler.go`: Entry point for the auth handler

```go

func (h *authHandler) Login(c *gin.Context) {

	var login dto.LoginRequest

	//  mengambil data request dan memasukan kedalam struct LoginRequest
	err := c.ShouldBindJSON(&login)
	if err != nil {
		errorhandler.HandleError(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	result, err := h.service.Login(&login)
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}

	res := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "Login successfuly",
		Data:       result,
	})

	c.JSON(http.StatusOK, res)

}

```

### menambahkan route /login pada file auth_router.go ./router/auth_router.go

tambahakan end point berikut untuk melakukan login

- `/router/auth_router.go`: Entry point for the auth route
```go

func AuthRouter(api *gin.RouterGroup) {
	authRepository := repository.NewAuthRepository(config.DB)
	authService := service.NewAuthService(authRepository)
	authHandler := handler.NewAuthHandler(authService)

	api.POST("/register", authHandler.Register)
	api.POST("/login", authHandler.Login) // tambahkan kode pada baris berikut
}

```

## Membuat fitur upload dan insert

