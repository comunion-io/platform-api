package usermodels

import (
	"context"
	accountEnv "cos-backend-com/src/account"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/models"
	"cos-backend-com/src/libs/sdk/account"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/wujiu2020/strip/utils"
)

const (
	accessKeyPrefix = "cn"
)

var (
	keyAlphabets   = []byte("0ictnbprkfs21zeovqj6wuldhymag748359x")
	defaultAvatars = []string{"default-avatar1.png", "default-avatar2.png", "default-avatar3.png", "default-avatar4.png", "default-avatar5.png"}
)

var Users = &users{
	Connector: models.DefaultConnector,
}

type users struct {
	dbconn.Connector
}

func (p *users) Get(ctx context.Context, id flake.ID, output interface{}) (err error) {
	stmt := `
		WITH res AS (
		    SELECT u.*,row_to_json(h.*) as hunter
		    FROM users u
		        LEFT JOIN hunters h ON u.id = h.user_id
		    WHERE u.id = ${id}
        )SELECT row_to_json(r.*) FROM res r;
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{id}": id,
	})

	return p.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, &util.PgJsonScanWrap{output}, query, args...)
	})
}

func (p *users) GetBypublicKey(ctx context.Context, publicKey string, output interface{}) (err error) {
	stmt := `
		SELECT *
		FROM users
		WHERE public_key ILIKE ${publicKey};
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{publicKey}": publicKey,
	})

	return p.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, output, query, args...)
	})
}

func (p *users) Create(ctx context.Context, publicKey string, output interface{}) (err error) {
	nonce := CreateNonce()
	uid, err := p.nextId(ctx)
	if err != nil {
		return err
	}
	stmt := `
		INSERT INTO users(avatar, public_key,private_secret,public_secret, nonce)
		VALUES (${avatar}, ${publicKey}, ${private_secret}, ${public_secret}, ${nonce})
		RETURNING *;
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{avatar}":         RangeAvatar(),
		"{publicKey}":      publicKey,
		"{nonce}":          nonce,
		"{private_secret}": CreateAccessKey(uid),
		"{public_secret}":  CreateSecretKey(),
	})

	return p.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, output, query, args...)
	})
}

func (p *users) UpdateNonce(ctx context.Context, id flake.ID, output interface{}) (err error) {
	nonce := CreateNonce()
	stmt := `
		UPDATE users SET nonce = ${nonce} WHERE id = ${id} RETURNING *;
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{id}":    id,
		"{nonce}": nonce,
	})

	return p.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, output, query, args...)
	})
}

func (p *users) FindOrCreate(ctx context.Context, publicKey string, output *account.UsersModel) (err error) {
	if err := p.GetBypublicKey(ctx, publicKey, output); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}
	if output.Id != flake.ID(0) {
		if err = p.UpdateNonce(ctx, output.Id, output); err != nil {
			return err
		}
		return
	}
	return p.Create(ctx, publicKey, output)
}

func (p *users) nextId(ctx context.Context) (netxtId flake.ID, err error) {
	err = p.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, &netxtId, `SELECT id_generator()`)
	})
	return
}

func CreateAccessKey(uid flake.ID) string {
	s, err := utils.RandomCreateString(6, keyAlphabets...)
	if err != nil {
		panic(err)
	}
	encoded := utils.NumberEncode(utils.ToStr(uid), keyAlphabets)
	return accessKeyPrefix + encoded + s
}

func CreateSecretKey() string {
	s, err := utils.RandomCreateString(32, keyAlphabets...)
	if err != nil {
		panic(err)
	}
	return s
}

func CreateNonce() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06v", rand.Intn(1000000000))
}

func RangeAvatar() string {
	rand.Seed(time.Now().UnixNano())
	return "https://" + accountEnv.Env.Minio.Endpoint + "/" + accountEnv.Env.Minio.StaticBucket + "/avatar/" + defaultAvatars[rand.Intn(4)]
}
