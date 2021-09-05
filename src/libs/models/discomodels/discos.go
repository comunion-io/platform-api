package discomodels

import (
	"context"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/models"
	"cos-backend-com/src/libs/models/ethmodels"
	"cos-backend-com/src/libs/sdk/cores"
	ethSdk "cos-backend-com/src/libs/sdk/eth"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

var Discos = &discos{
	Connector: models.DefaultConnector,
}

type discos struct {
	dbconn.Connector
}

func (c *discos) CreateDisco(ctx context.Context, startupId flake.ID, input *cores.CreateDiscosInput) (err error) {
	stmt := `
		INSERT INTO discos(
			id, startup_id, wallet_addr, token_addr, description, fund_raising_started_at, fund_raising_ended_at,
			investment_reward, reward_decline_rate, share_token, min_fund_raising, add_liquidity_pool,
			total_deposit_token)
		VALUES (${id}, ${startupId}, ${walletAddr}, ${tokenAddr},${description}, ${fundRaisingStartedAt}, ${fundRaisingEndedAt},
			${investmentReward},${rewardDeclineRate}, ${shareToken}, ${minFundRaising}, ${addLiquidityPool}, ${totalDepositToken}
		);
	`
	query, args := util.PgMapQueryV2(stmt, map[string]interface{}{
		"{id}":                   input.Id,
		"{startupId}":            startupId,
		"{walletAddr}":           input.WalletAddr,
		"{tokenAddr}":            input.TokenAddr,
		"{description}":          input.Description,
		"{fundRaisingStartedAt}": input.FundRaisingStartedAt,
		"{fundRaisingEndedAt}":   input.FundRaisingEndedAt,
		"{investmentReward}":     input.InvestmentReward,
		"{rewardDeclineRate}":    input.RewardDeclineRate,
		"{shareToken}":           input.ShareToken,
		"{minFundRaising}":       input.MinFundRaising,
		"{addLiquidityPool}":     input.AddLiquidityPool,
		"{totalDepositToken}":    input.TotalDepositToken,
	})
	return c.Invoke(ctx, func(db *sqlx.Tx) error {
		newCtx := dbconn.WithDB(ctx, db)
		createTransactionsInput := ethSdk.CreateTransactionsInput{
			TxId:     input.TxId,
			Source:   ethSdk.TransactionSourceDisco,
			SourceId: input.Id,
		}

		if err := ethmodels.Transactions.Create(newCtx, &createTransactionsInput); err != nil {
			return err
		}
		_, err = db.ExecContext(newCtx, query, args...)
		return err
	})
}

func (c *discos) GetDisco(ctx context.Context, startupId flake.ID, output interface{}) (err error) {
	stmt := `
		WITH res AS (
		    SELECT *, t.state as transaction_state, d.state as state, d.id as id,
				(CASE
					WHEN d.state = 4 AND d.fund_raising_ended_at < NOW() THEN 8
					WHEN d.state = 4 AND d.fund_raising_started_at < NOW() THEN 7
					ELSE d.state
				END) as state
		    FROM discos d
		        INNER JOIN transactions t ON t.source = ${source} AND source_id = d.id
		    WHERE d.startup_id = ${startupId}
		)SELECT row_to_json(r) FROM res r;
	`
	query, args := util.PgMapQueryV2(stmt, map[string]interface{}{
		"{startupId}": startupId,
		"{source}":    ethSdk.TransactionSourceDisco,
	})
	return c.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, &util.PgJsonScanWrap{output}, query, args...)
	})
}

func (c *discos) ListDisco(ctx context.Context, input *cores.ListDiscosInput, outputs interface{}) (total int, err error) {
	filterStmt := ""
	orderStmt := "ORDER BY %v %v"
	keyword := ""
	if input.Keyword != "" {
		keyword = "%" + util.PgEscapeLike(input.Keyword) + "%"
		filterStmt += `AND s.name ILIKE ${keyword}`
	}
	orderByField := "d.created_at"
	switch *input.OrderBy {
	case cores.ListDiscosOrderByTime:
		orderByField = "d.created_at"
	case cores.ListDiscosOrderByName:
		orderByField = "s.name"
	case cores.ListDiscosOrderByInvestmentReward:
		orderByField = "d.investment_reward"
	case cores.ListDiscosOrderByLiquidityPool:
		orderByField = "d.add_liquidity_pool"
	}
	orderField := ""
	if input.IsAsc {
		orderField = "asc"
	} else {
		orderField = "desc"
	}
	orderStmt = fmt.Sprintf(orderStmt, orderByField, orderField)
	countStmt := `
		SELECT count(*)
		FROM discos d
			INNER JOIN startups s ON d.startup_id = s.id
			INNER JOIN startup_revisions sr ON s.current_revision_id = sr.id
			INNER JOIN startup_settings ss ON ss.startup_id = s.id
			INNER JOIN startup_setting_revisions ssr ON ss.current_revision_id = ssr.id
		WHERE d.state IN (${discoStateEnabled}, ${discoStateWaitingForStart}, ${discoStateFundraisingSuccess}, ${discoStateFundraisingFailed}, ${discoStateFundraising}, ${discoStateFundraisingComplete})` + filterStmt
	query, args := util.PgMapQueryV2(countStmt, map[string]interface{}{
		"{keyword}":                       keyword,
		"{discoStateEnabled}":             cores.DiscoStateEnabled,
		"{discoStateWaitingForStart}":     cores.DiscoStateWaitingForStart,
		"{discoStateFundraisingSuccess}":  cores.DiscoStateFundraisingSuccess,
		"{discoStateFundraisingFailed}":   cores.DiscoStateFundraisingFailed,
		"{discoStateFundraising}":         cores.DiscoStateFundraising,
		"{discoStateFundraisingComplete}": cores.DiscoStateFundraisingComplete,
	})

	if err = c.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, &total, query, args...)
	}); err != nil {
		return
	}

	if total > 0 {
		stmt := `
			WITH res AS (
				SELECT d.*,json_build_object('id', d.startup_id, 'name', s.name,'log', sr.logo, 'token_symbol', ssr.token_symbol) startup,
				(CASE
					WHEN d.state = 4 AND d.fund_raising_ended_at < NOW() THEN 8
					WHEN d.state = 4 AND d.fund_raising_started_at < NOW() THEN 7
					ELSE d.state
				END) as state
				FROM discos d
					INNER JOIN startups s ON d.startup_id = s.id
					INNER JOIN startup_revisions sr ON s.current_revision_id = sr.id
					INNER JOIN startup_settings ss ON ss.startup_id = s.id
					INNER JOIN startup_setting_revisions ssr ON ss.current_revision_id = ssr.id
				WHERE d.state IN (${discoStateEnabled}, ${discoStateWaitingForStart}, ${discoStateFundraisingSuccess}, ${discoStateFundraisingFailed}, ${discoStateFundraising}, ${discoStateFundraisingComplete})` + filterStmt + orderStmt + `
				LIMIT ${limit} OFFSET ${offset}
			)SELECT json_agg(r) FROM res r;
		`
		query, args = util.PgMapQueryV2(stmt, map[string]interface{}{
			"{keyword}":                       keyword,
			"{limit}":                         input.GetLimit(),
			"{offset}":                        input.Offset,
			"{discoStateEnabled}":             cores.DiscoStateEnabled,
			"{discoStateWaitingForStart}":     cores.DiscoStateWaitingForStart,
			"{discoStateFundraisingSuccess}":  cores.DiscoStateFundraisingSuccess,
			"{discoStateFundraisingFailed}":   cores.DiscoStateFundraisingFailed,
			"{discoStateFundraising}":         cores.DiscoStateFundraising,
			"{discoStateFundraisingComplete}": cores.DiscoStateFundraisingComplete,
		})
		err = c.Invoke(ctx, func(db dbconn.Q) error {
			return db.GetContext(ctx, &util.PgJsonScanWrap{outputs}, query, args...)
		})
		return
	}
	return
}

func (c *discos) StatDiscoEthTotal(ctx context.Context, endAt time.Time, output interface{}) (err error) {
	stmt := `
		SELECT coalesce(sum(di.eth_count),0)
		FROM disco_investors di
		INNER JOIN users u ON di.uid = u.id
		INNER JOIN discos d ON d.id=di.disco_id
		WHERE di.created_at <= ${endAt};
	`

	query, args := util.PgMapQueryV2(stmt, map[string]interface{}{
		"{endAt}": endAt,
	})

	return c.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, output, query, args...)
	})
}

func (c *discos) StatDiscoTotal(ctx context.Context, output interface{}) (err error) {
	stmt := `
		SELECT COUNT(*) count,
    	COUNT(*) FILTER (
    	    WHERE created_at <= CURRENT_TIMESTAMP
    	        AND created_at > DATE_TRUNC('day', CURRENT_TIMESTAMP)
    	) increase_count
		FROM discos;
	`

	return c.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, output, stmt)
	})
}

func (c *discos) StatDiscoEthIncrease(ctx context.Context, input *cores.StatDiscoEthIncreaseInput, outputs interface{}) (err error) {
	stmt := `
		WITH dates AS (
			SELECT * FROM generate_date_series(${timeFrom}::TIMESTAMPTZ, ${timeTO}::TIMESTAMPTZ, ${tz}::TEXT, 'day') date
		),disco_cte AS (
			SELECT sum(di.eth_count) count,date_trunc('day',di.created_at) date
			FROM disco_investors di
			GROUP BY date
		)SELECT d.date,coalesce(dc.count,0) count
		FROM dates d
		LEFT JOIN disco_cte dc on d.date = dc.date;
	`
	query, args := util.PgMapQueryV2(stmt, map[string]interface{}{
		"{timeFrom}": input.TimeFrom,
		"{timeTO}":   input.TimeTo,
		"{tz}":       input.TimeTo,
	})
	return c.Invoke(ctx, func(db dbconn.Q) error {
		return db.SelectContext(ctx, outputs, query, args...)
	})
}

func (c *discos) GetDiscoSwapState(ctx context.Context, startupId *flake.ID, output interface{}) (err error) {
	stmt := `
	SELECT 
		coalesce((
			SELECT
			(CASE
				WHEN state = 4 AND fund_raising_ended_at < NOW() THEN 8
				WHEN state = 4 AND fund_raising_started_at < NOW() THEN 7
				ELSE state
			END) as state
        	FROM discos
			WHERE startup_id = ${startupId}
		),-1) disco_state,
    	coalesce((
			SELECT status
        	FROM exchanges
			WHERE startup_id = ${startupId}
		),-1) swap_state;
	`

	query, args := util.PgMapQueryV2(stmt, map[string]interface{}{
		"{startupId}": startupId,
	})

	return c.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, output, query, args...)
	})
}

func (c *discos) DeleteDisco(ctx context.Context, startupId, uid flake.ID) (err error) {
	stmt := `
		DELETE FROM discos
		WHERE startup_id = ${startupId} AND state < 3;
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{startupId}": startupId,
	})

	println(startupId)
	return c.Invoke(ctx, func(db dbconn.Q) (er error) {
		_, er = db.ExecContext(ctx, query, args...)
		return er
	})
}
