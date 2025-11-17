package logic_test

import (
	"chainmaker_web/src/auth"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/logic"
	"chainmaker_web/src/utils"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestMain 设置全局配置func
func TestMain(m *testing.M) {
	_ = config.InitConfig("", "")
	// 初始化数据库配置
	db.InitRedisContainer()
	db.InitMySQLContainer()

	// 运行其他测试
	os.Exit(m.Run())
}

func TestCreateAndSaveToken_Success(t *testing.T) {
	// 创建token
	token, err := logic.CreateAndSaveToken("user123", "signature123")

	// 验证结果
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// 验证token是否有效
	_, _, valid := auth.GetUserAddrAndToken(&gin.Context{})
	assert.False(t, valid)
}

func TestUpdateTokenExpireTime(t *testing.T) {
	_ = logic.UpdateTokenExpireTime("user123", "token123")
}

func TestVerifyAccountLogin_Success(t *testing.T) {
	// 创建随机数和正确的MD5哈希
	randomNum := time.Now().Unix()
	correctMd5 := utils.GetMd5Hash(randomNum)

	// 验证登录
	isValid, addr, err := logic.VerifyAccountLogin(randomNum, correctMd5)

	// 验证结果
	assert.NoError(t, err)
	assert.True(t, isValid)

	// 管理员账户地址应该是密码的SHA256
	expectedAddr := utils.GetAccountHashStr()
	assert.Equal(t, expectedAddr, addr)
}

func TestCheckAdminLogin(t *testing.T) {
	ctx := &gin.Context{}
	ctx.Set("userAddr", utils.GetAccountHashStr())
	_ = logic.CheckAdminLogin(ctx)
}
