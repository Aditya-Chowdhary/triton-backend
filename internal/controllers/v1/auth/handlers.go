package auth

import (
	"errors"
	"net/http"
	"triton-backend/internal/database"
	"triton-backend/internal/merrors"
	"triton-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthHandler struct {
	db *pgxpool.Pool
}

func Handler(db *pgxpool.Pool) *AuthHandler {
	return &AuthHandler{
		db: db,
	}
}

func (a *AuthHandler) RegisterOAuthUser(c *gin.Context) {
	var input struct {
		OAuthID string `json:"oauth_id" binding:"required"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		merrors.Validation(c, err.Error())
		return
	}

	tx, err := a.db.Begin(c)
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}
	defer tx.Rollback(c)

	qtx := database.New(a.db).WithTx(tx)

	// Create a new user UUID
	userUUID := uuid.New()

	// Try to create a new user in the database
	err = qtx.CreateUser(c, database.CreateUserParams{
		Uuid:     userUUID,
		AuthType: "oauth",
		OauthID:  input.OAuthID,
	})
	var e *pgconn.PgError
	if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
		merrors.Validation(c, "user already exists with this OAuth ID!")
		return
	} else if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	err = tx.Commit(c)
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BaseResponse{
		Success:    true,
		Message:    "OAuth user successfully registered",
		StatusCode: http.StatusOK,
	})
}

func (a *AuthHandler) GetUserByOAuthID(c *gin.Context) {
	var input struct {
		OAuthID string `json:"oauth_id" binding:"required"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		merrors.Validation(c, err.Error())
		return
	}

	q := database.New(a.db)
	userUUID, err := q.GetUserByOAuthID(c, database.GetUserByOAuthIDParams{
		AuthType: "oauth",
		OauthID:  input.OAuthID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		merrors.NotFound(c, "user not found!")
		return
	} else if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BaseResponse{
		Success:    true,
		Message:    "OAuth user successfully retrieved",
		Data:       userUUID,
		StatusCode: http.StatusOK,
	})
}

func (a *AuthHandler) RegisterAnonymousUser(c *gin.Context) {
	// For anonymous auth, we generate a UUID and register it as a user.
	userUUID := uuid.New()

	tx, err := a.db.Begin(c)
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}
	defer tx.Rollback(c)

	qtx := database.New(a.db).WithTx(tx)

	// Try to create a new user in the database
	err = qtx.CreateUser(c, database.CreateUserParams{
		Uuid:     userUUID,
		AuthType: "anonymous",
	})
	var e *pgconn.PgError
	if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
		merrors.Validation(c, "user already exists with this ID!")
		return
	} else if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	err = tx.Commit(c)
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BaseResponse{
		Success:    true,
		Message:    "Anonymous user successfully registered",
		Data:       userUUID,
		StatusCode: http.StatusOK,
	})
}

func (a *AuthHandler) GetUserByAnonymousID(c *gin.Context) {
	var input struct {
		UserID uuid.UUID `json:"user_id" binding:"required"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		merrors.Validation(c, err.Error())
		return
	}

	q := database.New(a.db)
	userUUID, err := q.GetUserByOAuthID(c, database.GetUserByOAuthIDParams{
		AuthType: "anonymous",
		OauthID:  input.UserID.String(), // You might need to adjust this part based on your actual query setup.
	})
	if errors.Is(err, pgx.ErrNoRows) {
		merrors.NotFound(c, "user not found!")
		return
	} else if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BaseResponse{
		Success:    true,
		Message:    "Anonymous user successfully retrieved",
		Data:       userUUID,
		StatusCode: http.StatusOK,
	})
}
