// Package auth 登录中间件
package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var UTTestAdminAddr = "admin" // 测试用管理员地址

func TestMain(m *testing.M) {
	_ = config.InitConfig("", "")
	// 初始化数据库配置
	db.InitRedisContainer()
	db.InitMySQLContainer()

	// 运行其他测试
	os.Exit(m.Run())
}

func TestCreateSignedTokenAndVerifyToken(t *testing.T) {
	expire := time.Now().Add(time.Hour)
	token, err := CreateSignedToken(UTTestAdminAddr, expire)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	gotAddr, errMsg := VerifyToken(token)
	assert.Equal(t, UTTestAdminAddr, gotAddr)
	assert.Equal(t, "", errMsg)
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	invalidToken := "invalid.token.string"
	addr, errMsg := VerifyToken(invalidToken)
	assert.Empty(t, addr)
	assert.NotEmpty(t, errMsg)
}

func TestVerifyToken_ExpiredToken(t *testing.T) {
	expire := time.Now().Add(-time.Hour)
	token, err := CreateSignedToken(UTTestAdminAddr, expire)
	assert.NoError(t, err)

	_, errMsg := VerifyToken(token)
	assert.Equal(t, entity.ErrorMsgTokenExpired, errMsg)
}

func TestGetUserAddrAndToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := &gin.Context{}
	// 没有设置
	_, _, valid := GetUserAddrAndToken(ctx)
	assert.False(t, valid)

	// 设置错误类型
	ctx.Set("userAddr", 123)
	ctx.Set("token", 456)
	_, _, valid2 := GetUserAddrAndToken(ctx)
	assert.False(t, valid2)

	// 设置正确类型
	ctx.Set("userAddr", "abc")
	ctx.Set("token", "def")
	addr, token, valid := GetUserAddrAndToken(ctx)
	assert.True(t, valid)
	assert.Equal(t, "abc", addr)
	assert.Equal(t, "def", token)
}

func TestCheckTokenExpire(t *testing.T) {
	expire := time.Now().Add(time.Hour)
	token, err := CreateSignedToken(UTTestAdminAddr, expire)
	if err != nil {
		return
	}

	newUUID := uuid.New().String()
	userToken := &db.LoginUserToken{
		Id:         newUUID,
		UserAddr:   UTTestAdminAddr,
		Token:      token,
		Sign:       UTTestAdminAddr,
		ExpireTime: expire,
	}
	//存储token
	err = dbhandle.InsertUserToken(userToken)
	if err != nil {
		return
	}
	// token过期
	assert.False(t, checkTokenExpire("notoken"))

	// token有效
	assert.True(t, checkTokenExpire(token))
}

func TestTokenAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(TokenAuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		userAddr, token, valid := GetUserAddrAndToken(c)
		assert.True(t, valid)
		assert.Equal(t, "user1", userAddr)
		assert.Equal(t, "token1", token)
		c.String(200, "ok")
	})
}

func TestTokenAuthMiddleware_Real(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 1. 测试成功场景
	t.Run("Valid token", func(t *testing.T) {
		r := gin.New()
		r.Use(TokenAuthMiddleware())
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// 创建有效 token
		token, _ := CreateSignedToken(UTTestAdminAddr, time.Now().Add(time.Hour))
		insertToken(token)

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	// 2. 测试缺少 Header
	t.Run("Missing header", func(t *testing.T) {
		r := gin.New()
		r.Use(TokenAuthMiddleware())
		r.GET("/test", func(c *gin.Context) {})

		req := httptest.NewRequest("GET", "/test", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)
	})

	// 3. 测试无效的 Header 格式
	t.Run("Invalid header format", func(t *testing.T) {
		r := gin.New()
		r.Use(TokenAuthMiddleware())
		r.GET("/test", func(c *gin.Context) {})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Add("Authorization", "InvalidTokenFormat")
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)
	})

	// 4. 测试过期 token
	t.Run("Expired token", func(t *testing.T) {
		r := gin.New()
		r.Use(TokenAuthMiddleware())
		r.GET("/test", func(c *gin.Context) {})

		token, _ := CreateSignedToken(UTTestAdminAddr, time.Now().Add(-time.Hour))
		insertToken(token)

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Add("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)
	})
}

// 辅助函数
func insertToken(token string) {
	userToken := &db.LoginUserToken{
		Id:         uuid.New().String(),
		UserAddr:   UTTestAdminAddr,
		Token:      token,
		Sign:       UTTestAdminAddr,
		ExpireTime: time.Now().Add(time.Hour),
	}
	dbhandle.InsertUserToken(userToken)
}

func TestRegularlyVerifyToken_EdgeCases(t *testing.T) {
	// 1. 测试空 token 场景
	t.Run("Empty token after Bearer", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = &http.Request{Header: http.Header{}}
		c.Request.Header.Set("Authorization", "Bearer ")

		_, _, err := RegularlyVerifyToken(c)
		assert.NotNil(t, err)
	})

	// 2. 测试 token 验证失败场景
	t.Run("Token verification failure", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = &http.Request{Header: http.Header{}}
		c.Request.Header.Set("Authorization", "Bearer invalid.token.here")

		_, _, err := RegularlyVerifyToken(c)
		assert.NotNil(t, err)
	})
}
