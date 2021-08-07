CREATE TABLE comunion.proposals (
    id bigint DEFAULT comunion.id_generator() NOT NULL
        CONSTRAINT proposals_id_pk
            PRIMARY KEY,
    tx_id text NOT NULL,
    startup_id bigint NOT NULL,
    wallet_addr text NOT NULL,
    contract_addr text NOT NULL,
    status int4 NOT NULL DEFAULT 0,
    title text NOT NULL,
    type int4 NOT NULL,
    user_id bigint NOT NULL,
    contact text NOT NULL,
    description text NOT NULL,
    voter_type int4 NOT NULL,
    supporters int4 NOT NULL,
    minimum_approval_percentage int4 NOT NULL,
    duration int4 NOT NULL,
    has_payment boolean NOT NULL,
    payment_addr text,
    payment_type int2,
    payment_months int4,
    payment_date text,
    payment_amount float8,
    total_payment_amount float8,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

COMMENT ON COLUMN comunion.proposals.status IS '0：待确认，1：未开始，2：进行中，3：已结束，4：未成案，5：提案被拒绝，6：提案被通过';
COMMENT ON COLUMN comunion.proposals.type IS '1：Finance，2：Governance，3：Strategy，4：Product，5：Media，6：Community，7：Node';
COMMENT ON COLUMN comunion.proposals.voter_type IS '1：FounderAssign，2：POS，3：All';
COMMENT ON COLUMN comunion.proposals.payment_type IS '1：一次性支付，2：按月支付';

CREATE TABLE comunion.proposal_terms (
    id bigint DEFAULT comunion.id_generator() NOT NULL
       CONSTRAINT proposal_terms_id_pk
           PRIMARY KEY,
    proposal_id bigint NOT NULL,
    amount float8 NOT NULL,
    content text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE comunion.proposal_votes (
    id bigint DEFAULT comunion.id_generator() NOT NULL
        CONSTRAINT proposal_votes_id_pk
            PRIMARY KEY,
    tx_id text NOT NULL,
    proposal_id bigint NOT NULL,
    amount float8 NOT NULL,
    is_approved boolean NOT NULL,
    wallet_addr text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);