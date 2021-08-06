package proposalmodels

import (
	"context"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/flake"
	json "cos-backend-com/src/common/pgencoding/json2"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/models"
	coresSdk "cos-backend-com/src/libs/sdk/cores"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

var Proposals = &proposals{
	Connector: models.DefaultConnector,
}

type proposals struct {
	dbconn.Connector
}

func (c *proposals) CreateProposal(ctx context.Context, input *coresSdk.CreateProposalInput, output *coresSdk.CreateProposalResult) (err error) {
	stmt := `
		INSERT INTO proposals(id, tx_id, startup_id, wallet_addr, contract_addr, status, title, type, user_id, contact, description,
							  voter_type, supporters, minimum_approval_percentage, duration, has_payment, payment_addr,
							  payment_type, payment_months, payment_date, payment_amount, total_payment_amount)
		VALUES (${id}, ${txId}, ${startupId}, ${walletAddr}, ${contractAddr}, ${status}, ${title}, ${type}, ${userId}, ${contact}, ${description},
				${voterType}, ${supporters}, ${minimumApprovalPercentage}, ${duration}, ${hasPayment}, ${paymentAddr},
				${paymentType}, ${paymentMonths}, ${paymentDate}, ${paymentAmount}, ${totalPaymentAmount})
		RETURNING id, status;
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{id}":                        input.Id,
		"{txId}":                      input.TxId,
		"{startupId}":                 input.StartupId,
		"{walletAddr}":                input.WalletAddr,
		"{contractAddr}":              input.ContractAddr,
		"{status}":                    input.Status,
		"{title}":                     input.Title,
		"{type}":                      input.Type,
		"{userId}":                    input.UserId,
		"{contact}":                   input.Contact,
		"{description}":               input.Description,
		"{voterType}":                 input.VoterType,
		"{supporters}":                input.Supporters,
		"{minimumApprovalPercentage}": input.MinimumApprovalPercentage,
		"{duration}":                  input.Duration,
		"{hasPayment}":                input.HasPayment,
		"{paymentAddr}":               input.PaymentAddr,
		"{paymentType}":               input.PaymentType,
		"{paymentMonths}":             input.PaymentMonths,
		"{paymentDate}":               input.PaymentDate,
		"{paymentAmount}":             input.PaymentAmount,
		"{totalPaymentAmount}":        input.TotalPaymentAmount,
	})

	return c.Invoke(ctx, func(db *sqlx.Tx) error {
		return db.GetContext(ctx, output, query, args...)
	})
}

func (c *proposals) UpdateProposalOver(ctx context.Context, output interface{}) (err error) {
	stmt := `
		UPDATE proposals set status = ${over}
		WHERE now() > created_at + cast(duration || 'days' as interval) and status = ${voting}
		RETURNING true AS done;
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{over}":   coresSdk.ProposalStatusOver,
		"{voting}": coresSdk.ProposalStatusVoting,
	})
	return c.Invoke(ctx, func(db *sqlx.Tx) error {
		return db.GetContext(ctx, output, query, args...)
	})
}

func (c *proposals) UpdateProposalStatus(ctx context.Context, input *coresSdk.UpdateProposalStatusInput, output *coresSdk.UpdateProposalStatusResult) (err error) {
	stmt := `
		UPDATE proposals set status=${status} where id=${id}
		RETURNING id;
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{status}": input.Status,
		"{id}":     input.Id,
	})
	return c.Invoke(ctx, func(db *sqlx.Tx) error {
		return db.GetContext(ctx, output, query, args...)
	})
}

func (c *proposals) VoteProposal(ctx context.Context, input *coresSdk.VoteProposalInput, output *coresSdk.VoteProposalResult) (err error) {
	stmt := `
		INSERT INTO proposal_votes(tx_id, proposal_id, amount, is_approved, wallet_addr)
		VALUES (${txId}, ${proposalId}, ${amount}, ${isApproved}, ${walletAddr})
		RETURNING id;
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{txId}":       input.TxId,
		"{proposalId}": input.Id,
		"{amount}":     input.Amount,
		"{isApproved}":   input.IsApproved,
		"{walletAddr}": input.WalletAddr,
	})
	return c.Invoke(ctx, func(db *sqlx.Tx) error {
		return db.GetContext(ctx, output, query, args...)
	})
}

func (c *proposals) CreateProposalWithTerms(ctx context.Context, input *coresSdk.CreateProposalInput, output *coresSdk.CreateProposalResult) (err error) {
	return c.Invoke(ctx, func(db *sqlx.Tx) error {
		newCtx := dbconn.WithDB(ctx, db)
		if er := c.CreateProposal(newCtx, input, output); er != nil {
			return er
		}

		var outputTerm coresSdk.CreateTermResult
		for i := 0; i < len(input.Terms); i++ {
			if er := c.CreateTerm(newCtx, output.Id, &input.Terms[i], &outputTerm); er != nil {
				return er
			}
		}
		return nil
	})
}

func (c *proposals) GetProposal(ctx context.Context, input *coresSdk.GetProposalInput, output *coresSdk.ProposalResult) (err error) {
	where := ""
	if input.Id != 0 {
		where += `pr.id = ${id}`
	} else if input.TxId != "" {
		where += `pr.tx_id = ${txId}`
	} else {
		where += "1 = 2"
	}
	stmt := `
		WITH
			votes_cte AS (
				SELECT proposal_id, amount, is_approved, wallet_addr, created_at
			  	FROM proposal_votes
			  	WHERE proposal_id = ${id}
					ORDER BY created_at
			),
			votes_cte_group AS (
				SELECT vc.proposal_id, COALESCE(json_agg(vc), '[]'::json) AS votes
				FROM votes_cte vc
				GROUP BY vc.proposal_id
			),
			terms_cte AS (
				SELECT proposal_id, amount, content
			  	FROM proposal_terms
			  	WHERE proposal_id = ${id}
					ORDER BY created_at
			),
			terms_cte_group AS (
				SELECT tc.proposal_id, COALESCE(json_agg(tc), '[]'::json) AS terms
				FROM terms_cte tc
				GROUP BY tc.proposal_id
			),

			res AS (
	    		SELECT 
					pr.id,
					pr.tx_id,
					json_build_object('id',s.id,'name',s.name,'logo',sr.logo,'token_symbol',ssr.token_symbol) startup,
					json_build_object('id',COALESCE(ht.user_id, pr.user_id),'name',COALESCE(ht.name, pr.wallet_addr)) comer,
					pr.wallet_addr,
					pr.contract_addr,
					pr.created_at,
					pr.updated_at,
					pr.status,
					pr.title,
					pr.type,
					pr.contact,
					pr.description,
					pr.voter_type,
					pr.supporters,
					pr.minimum_approval_percentage,
					pr.duration,
					pr.has_payment,
					pr.payment_addr,
					pr.payment_type,
					pr.payment_months,
					pr.payment_date,
					pr.payment_amount,
					pr.total_payment_amount,
					COALESCE(vcg.votes, '[]'::json) AS votes,
					COALESCE(tcg.terms, '[]'::json) AS terms
	    		FROM proposals pr
					INNER JOIN startups s ON s.id = pr.startup_id
					INNER JOIN startup_revisions sr ON s.current_revision_id = sr.id
					INNER JOIN startup_settings ss ON s.id = ss.startup_id
					INNER JOIN startup_setting_revisions ssr ON ss.current_revision_id = ssr.id
					LEFT JOIN hunters ht ON pr.user_id = ht.user_id
					LEFT JOIN votes_cte_group vcg ON vcg.proposal_id = pr.id
					LEFT JOIN terms_cte_group tcg ON tcg.proposal_id = pr.id

	    		WHERE ` + where + `
				)
			SELECT row_to_json(res.*) FROM res
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{id}":   input.Id,
		"{txId}": input.TxId,
	})

	return c.Invoke(ctx, func(db dbconn.Q) (er error) {
		return db.GetContext(ctx, &util.PgJsonScanWrap{output}, query, args...)
	})
}

func (c *proposals) CreateTerm(ctx context.Context, proposalId flake.ID, input *coresSdk.CreateTermInput, output *coresSdk.CreateTermResult) (err error) {
	stmt := `
		INSERT INTO proposal_terms(proposal_id, amount, content)
		VALUES (${proposalId}, ${amount}, ${content})
		RETURNING id;
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{proposalId}": proposalId,
		"{amount}":     input.Amount,
		"{content}":    input.Content,
	})

	return c.Invoke(ctx, func(db *sqlx.Tx) error {
		return db.GetContext(ctx, output, query, args...)
	})
}

func (c *proposals) ListProposals(ctx context.Context, userId flake.ID, input *coresSdk.ListProposalsInput, outputs interface{}) (total int, err error) {
	filterStmt := ``
	if input.Type != "" {
		switch input.Type {
		case coresSdk.ListProposalsTypeAll:
			filterStmt = "1 = 1"
		case coresSdk.ListProposalsTypeCreated:
			filterStmt = "pr.user_id = ${userId}"
		case coresSdk.ListProposalsTypeVoted:
			filterStmt = "pr.id in (SELECT pv.proposal_id FROM proposal_votes pv INNER JOIN users u ON pv.wallet_addr = u.public_key WHERE u.id = ${userId})"
		}
	}

	if input.StartupId != 0 {
		filterStmt += ` AND pr.startup_id = ${startupId}`
	}

	var statuses []int
	err2 := json.Unmarshal([]byte(input.Statuses), &statuses)
	if err2 == nil {
		if len(statuses) > 0 {
			filterStmt += ` AND pr.status in ` + strings.Replace(
				strings.Replace(input.Statuses, "[", "(", -1),
				"]", ")", -1)
		}
	}

	keyword := ""
	if input.Keyword != "" {
		keyword = "%" + util.PgEscapeLike(input.Keyword) + "%"
		filterStmt += ` AND pr.title ILIKE ${keyword}`
	}

	orderStmt := "ORDER BY %v %v"
	if input.OrderBy == "" {
		orderStmt = fmt.Sprintf(orderStmt, "pr.created_at", "DESC")
	} else {
		orderField := "pr." + input.OrderBy
		orderSeq := ""
		if input.IsDesc {
			orderSeq = "DESC"
		} else {
			orderSeq = "ASC"
		}
		orderStmt = fmt.Sprintf(orderStmt, orderField, orderSeq)
	}

	stmt := ` 
		WITH
			res AS (
	    		SELECT 
					pr.id,
					json_build_object('id',s.id,'name',s.name,'logo',sr.logo,'token_symbol',ssr.token_symbol) startup,
					json_build_object('id',COALESCE(ht.user_id, pr.user_id),'name',COALESCE(ht.name, pr.wallet_addr)) comer,
					pr.status,
					pr.title,
					pr.duration,
					pr.has_payment,
					pr.total_payment_amount,
					pr.created_at,
					pr.updated_at
	    		FROM proposals pr
					INNER JOIN startups s ON s.id = pr.startup_id
					INNER JOIN startup_revisions sr ON s.current_revision_id = sr.id
					INNER JOIN startup_settings ss ON s.id = ss.startup_id
					INNER JOIN startup_setting_revisions ssr ON ss.current_revision_id = ssr.id
					LEFT JOIN hunters ht ON pr.user_id = ht.user_id
	    		WHERE ` + filterStmt + `
					` + orderStmt + `
				LIMIT ${limit} OFFSET ${offset}
			)

			SELECT COALESCE(json_agg(r.*), '[]'::json) FROM res r;
	`

	countStmt := `
		SELECT count(*)
		FROM proposals pr
					INNER JOIN startups s ON s.id = pr.startup_id
					INNER JOIN startup_revisions sr ON s.current_revision_id = sr.id
					INNER JOIN startup_settings ss ON s.id = ss.startup_id
					INNER JOIN startup_setting_revisions ssr ON ss.current_revision_id = ssr.id
					LEFT JOIN hunters ht ON pr.user_id = ht.user_id
	    		WHERE ` + filterStmt

	query, args := util.PgMapQuery(countStmt, map[string]interface{}{
		"{keyword}":   keyword,
		"{userId}":    userId,
		"{startupId}": input.StartupId,
	})

	if err = c.Invoke(ctx, func(db *sqlx.Tx) (er error) {
		return db.GetContext(ctx, &total, query, args...)
	}); err != nil {
		return
	}
	query, args = util.PgMapQuery(stmt, map[string]interface{}{
		"{keyword}":   keyword,
		"{offset}":    input.Offset,
		"{limit}":     input.GetLimit(),
		"{userId}":    userId,
		"{startupId}": input.StartupId,
	})

	return total, c.Invoke(ctx, func(db *sqlx.Tx) (er error) {
		return db.GetContext(ctx, &util.PgJsonScanWrap{outputs}, query, args...)
	})
}
