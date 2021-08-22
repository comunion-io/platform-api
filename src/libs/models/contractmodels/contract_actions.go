package contractmodels

import (
	"context"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/models"
	"cos-backend-com/src/libs/sdk/cores"

	"github.com/jmoiron/sqlx"
)

var ContractActions = &contractActions{
	Connector: models.DefaultConnector,
}

type contractActions struct {
	dbconn.Connector
}

func (c *contractActions) Create(ctx context.Context, uid flake.ID, actionType cores.ContractActionType, input *cores.CreateContractActionInput) (err error) {
	stmt := `
		INSERT INTO contract_actions (id, uid, type) VALUES (${id}, ${uid}, ${type});
	`
	query, args := util.PgMapQueryV2(stmt, map[string]interface{}{
		"{id}":   input.Id,
		"{uid}":  uid,
		"{type}": actionType,
	})

	return c.Invoke(ctx, func(db *sqlx.DB) (er error) {
		_, err = db.ExecContext(ctx, query, args...)
		return err
	})
}

func (c *contractActions) List(ctx context.Context, uid flake.ID, actionType cores.ContractActionType, output interface{}) (err error) {
	stmt := `
		WITH res AS (
			SELECT *
			FROM contract_actions
			WHERE uid = ${uid} AND type = ${actionType}
			ORDER BY created_at DESC
		)
		SELECT COALESCE(json_agg(r.*), '[]'::json) FROM res r;
	`
	query, args := util.PgMapQueryV2(stmt, map[string]interface{}{
		"{uid}":        uid,
		"{actionType}": actionType,
	})

	return c.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, &util.PgJsonScanWrap{output}, query, args...)
	})
}

func (c *contractActions) Get(ctx context.Context, id string, uid flake.ID, actionType, output interface{}) (err error) {
	stmt := `
		WITH res AS (
			SELECT *
			FROM contract_actions
			WHERE uid = ${uid} AND type = ${actionType} AND id = ${id}
		)
		SELECT row_to_json(r) FROM res r;
	`
	query, args := util.PgMapQueryV2(stmt, map[string]interface{}{
		"{id}":         id,
		"{uid}":        uid,
		"{actionType}": actionType,
	})

	return c.Invoke(ctx, func(db dbconn.Q) (er error) {
		return db.GetContext(ctx, &util.PgJsonScanWrap{output}, query, args...)
	})
}
