
CREATE TABLE public.brc20_ticker_info (
    id bigserial NOT NULL,
    block_height int4 NOT NULL,
    tick varchar(4) NOT NULL,
	 max_supply numeric(40) NOT NULL,
	 decimals int4 NOT NULL,
	 limit_per_mint numeric(40) NOT NULL,
	 remaining_supply numeric(40) NOT NULL,
    pkscript_deployer text NOT NULL
);
CREATE INDEX brc20_ticker_info_block_height_idx ON public.brc20_ticker_info USING btree (block_height);
CREATE INDEX brc20_ticker_info_tick_idx ON public.brc20_ticker_info USING btree (tick);


CREATE TABLE public.brc20_user_balance (
    id bigserial NOT NULL,
    block_height int4 NOT NULL,
	 tick varchar(4) NOT NULL,
	 pkscript text NOT NULL,
	 available_balance numeric(40) NOT NULL,
	 transferable_balance numeric(40) NOT NULL
);
CREATE INDEX brc20_user_balance_block_height_idx ON public.brc20_user_balance USING btree (block_height);
CREATE INDEX brc20_user_balance_pkscript_tick_idx ON public.brc20_user_balance USING btree (pkscript, tick);


-- state of moved
CREATE TABLE public.brc20_transfer_state (
    id bigserial NOT NULL,
    block_height int4 NOT NULL,
    create_key text NOT NULL,
    moved boolean,
);

CREATE TABLE public.brc20_valid_transfer (
    id bigserial NOT NULL,
    block_height int4 NOT NULL,
	 tick varchar(4) NOT NULL,
	 pkscript text NOT NULL,
	 amount numeric(40) NOT NULL,
	 inscription_number int8 NOT NULL,
	 inscription_id text NOT NULL,
	 txid text NOT NULL,
    vout int4 NOT NULL,
	 output_value int8 NOT NULL,
	 output_offset int8 NOT NULL
);
CREATE INDEX brc20_valid_transfer_block_height_idx ON public.brc20_valid_transfer USING btree (block_height);
CREATE INDEX brc20_valid_transfer_pkscript_tick_idx ON public.brc20_valid_transfer USING btree (pkscript, tick);

CREATE TABLE public.brc20_history (
    id bigserial NOT NULL,
    block_height int4 NOT NULL,
    tick varchar(4) NOT NULL,
    history_type smallint NOT NULL,
    valid boolean,
    txid text NOT NULL,
    idx int4 NOT NULL,          -- inscription index in block
    vout int4 NOT NULL,
    output_value int8 NOT NULL,
    output_offset int8 NOT NULL,
    pkscript_from text NOT NULL,
    pkscript_to text NOT NULL,
    fee int8 NOT NULL,
    txidx int4 NOT NULL,        -- txidx in block
    block_time int4 NOT NULL,   -- equivalent to uint32 in Go
    inscription_number int8 NOT NULL,
    inscription_id text NOT NULL,
    inscription_content jsonb NOT NULL,
	 amount numeric(40) NOT NULL,
	 available_balance numeric(40) NOT NULL,
	 transferable_balance numeric(40) NOT NULL,

	 CONSTRAINT brc20_history_pk PRIMARY KEY (id)
);
CREATE INDEX brc20_history_block_height_idx ON public.brc20_history USING btree (block_height);
CREATE INDEX brc20_history_pkscript_from_tick_idx ON public.brc20_history USING btree (pkscript_from, tick);
CREATE INDEX brc20_history_pkscript_to_tick_idx ON public.brc20_history USING btree (pkscript_to, tick);

CREATE TABLE public.brc20_swap_info (
    id bigserial NOT NULL,
    block_height int4 NOT NULL,
    module_id text NOT NULL,
	 name text NOT NULL,
    pkscript_deployer text NOT NULL,
    pkscript_sequencer text NOT NULL,
    pkscript_gas_to text NOT NULL,
    pkscript_lp_fee text NOT NULL,
	 gas_tick varchar(4) NOT NULL,
    fee_rate_swap text NOT NULL
);
CREATE INDEX brc20_swap_info_block_height_idx ON public.brc20_swap_info USING btree (block_height);
CREATE INDEX brc20_swap_info_module_id_idx ON public.brc20_swap_info USING btree (module_id);


CREATE TABLE public.brc20_swap_commit_state (
    id bigserial NOT NULL,
    block_height int4 NOT NULL,
    create_key text NOT NULL,
    moved boolean,
);
CREATE INDEX brc20_swap_commit_state_block_height_idx ON public.brc20_swap_commit_state USING btree (block_height);
CREATE INDEX brc20_swap_commit_state_create_key_idx ON public.brc20_swap_commit_state USING btree (create_key);

CREATE TABLE public.brc20_swap_valid_commit (
    id bigserial NOT NULL,
    block_height int4 NOT NULL,
    module_id text NOT NULL,
	 pkscript text NOT NULL,
	 inscription_number int8 NOT NULL,
	 inscription_id text NOT NULL,
	 txid text NOT NULL,
    vout int4 NOT NULL,
	 output_value int8 NOT NULL,
	 output_offset int8 NOT NULL
    inscription_content jsonb NOT NULL,
);
CREATE INDEX brc20_swap_valid_commit_block_height_idx ON public.brc20_swap_valid_commit USING btree (block_height);
CREATE INDEX brc20_swap_valid_commit_module_id_idx ON public.brc20_swap_valid_commit USING btree (module_id);


CREATE TABLE public.brc20_swap_commit_chain (
    id bigserial NOT NULL,
    block_height int4 NOT NULL,
    module_id text NOT NULL,
    commit_id text NOT NULL,
    valid boolean,
    connected boolean
);
CREATE INDEX brc20_swap_commit_chain_block_height_idx ON public.brc20_swap_commit_chain USING btree (block_height);
CREATE INDEX brc20_swap_commit_chain_module_id_idx ON public.brc20_swap_commit_chain USING btree (module_id);


CREATE TABLE public.brc20_swap_user_balance (
    id bigserial NOT NULL,
    block_height int4 NOT NULL,
    module_id text NOT NULL,
	 tick varchar(4) NOT NULL,
	 pkscript text NOT NULL,
	 swap_balance numeric(40) NOT NULL,
	 available_balance numeric(40) NOT NULL,
	 approveable_balance numeric(40) NOT NULL,
	 cond_approveable_balance numeric(40) NOT NULL,
	 withdrawable_balance numeric(40) NOT NULL
);
CREATE INDEX brc20_swap_user_balance_block_height_idx ON public.brc20_swap_user_balance USING btree (block_height);
CREATE INDEX brc20_swap_user_balance_module_id_idx ON public.brc20_swap_user_balance USING btree (module_id);
CREATE INDEX brc20_swap_user_balance_pkscript_tick_idx ON public.brc20_swap_user_balance USING btree (pkscript, tick);


-- state of moved
CREATE TABLE public.brc20_swap_approve_state (
    id bigserial NOT NULL,
    block_height int4 NOT NULL,
    create_key text NOT NULL,
    moved boolean,
);

CREATE TABLE public.brc20_swap_valid_approve (
    id bigserial NOT NULL,
    block_height int4 NOT NULL,
    module_id text NOT NULL,
	 tick varchar(4) NOT NULL,
	 pkscript text NOT NULL,
	 amount numeric(40) NOT NULL,
	 inscription_number int8 NOT NULL,
	 inscription_id text NOT NULL,
	 txid text NOT NULL,
    vout int4 NOT NULL,
	 output_value int8 NOT NULL,
	 output_offset int8 NOT NULL
);
CREATE INDEX brc20_swap_valid_approve_block_height_idx ON public.brc20_swap_valid_approve USING btree (block_height);
CREATE INDEX brc20_swap_valid_approve_module_id_idx ON public.brc20_swap_valid_approve USING btree (module_id);


-- state of moved
CREATE TABLE public.brc20_swap_cond_approve_state (
    id bigserial NOT NULL,
    block_height int4 NOT NULL,
    create_key text NOT NULL,
	 balance numeric(40) NOT NULL,
    moved boolean,
	 pkscript_owner text NOT NULL,
	 pkscript_delegator text NOT NULL,
);

CREATE TABLE public.brc20_swap_valid_cond_approve (
    id bigserial NOT NULL,
    block_height int4 NOT NULL,
    module_id text NOT NULL,
	 tick varchar(4) NOT NULL,
	 pkscript text NOT NULL,
	 amount numeric(40) NOT NULL,
	 inscription_number int8 NOT NULL,
	 inscription_id text NOT NULL,
	 txid text NOT NULL,
    vout int4 NOT NULL,
	 output_value int8 NOT NULL,
	 output_offset int8 NOT NULL
);
CREATE INDEX brc20_swap_valid_cond_approve_block_height_idx ON public.brc20_swap_valid_cond_approve USING btree (block_height);
CREATE INDEX brc20_swap_valid_cond_approve_module_id_idx ON public.brc20_swap_valid_cond_approve USING btree (module_id);


-- state of moved
CREATE TABLE public.brc20_swap_withdraw_state (
    id bigserial NOT NULL,
    block_height int4 NOT NULL,
    create_key text NOT NULL,
    moved boolean,
);

CREATE TABLE public.brc20_swap_valid_withdraw (
    id bigserial NOT NULL,
    block_height int4 NOT NULL,
    module_id text NOT NULL,
	 tick varchar(4) NOT NULL,
	 pkscript text NOT NULL,
	 amount numeric(40) NOT NULL,
	 inscription_number int8 NOT NULL,
	 inscription_id text NOT NULL,
	 txid text NOT NULL,
    vout int4 NOT NULL,
	 output_value int8 NOT NULL,
	 output_offset int8 NOT NULL
);
CREATE INDEX brc20_swap_valid_withdraw_block_height_idx ON public.brc20_swap_valid_withdraw USING btree (block_height);
CREATE INDEX brc20_swap_valid_withdraw_module_id_idx ON public.brc20_swap_valid_withdraw USING btree (module_id);

CREATE TABLE public.brc20_swap_user_lp_balance (
    id bigserial NOT NULL,
    block_height int4 NOT NULL,
    module_id text NOT NULL,
	 pool varchar(9) NOT NULL,
	 pkscript text NOT NULL,
	 lp_balance numeric(40) NOT NULL
);
CREATE INDEX brc20_swap_user_lp_balance_block_height_idx ON public.brc20_swap_user_lp_balance USING btree (block_height);
CREATE INDEX brc20_swap_user_lp_balance_module_id_idx ON public.brc20_swap_user_lp_balance USING btree (module_id);

CREATE TABLE public.brc20_swap_pool_balance (
    id bigserial NOT NULL,
    block_height int4 NOT NULL,
    module_id text NOT NULL,
	 tick0 varchar(4) NOT NULL,
	 tick0_balance numeric(40) NOT NULL,
	 tick1 varchar(4) NOT NULL,
	 tick1_balance numeric(40) NOT NULL,
	 lp_balance numeric(40) NOT NULL
);
CREATE INDEX brc20_swap_pool_balance_block_height_idx ON public.brc20_swap_pool_balance USING btree (block_height);
CREATE INDEX brc20_swap_pool_balance_module_id_idx ON public.brc20_swap_pool_balance USING btree (module_id);

CREATE TABLE public.brc20_swap_history (
    id bigserial NOT NULL,
    block_height int4 NOT NULL,
    module_id text NOT NULL,
    history_type smallint NOT NULL,
    valid boolean,
    txid text NOT NULL,
    idx int4 NOT NULL,          -- inscription index in block
    vout int4 NOT NULL,
    output_value int8 NOT NULL,
    output_offset int8 NOT NULL,
    pkscript_from text NOT NULL,
    pkscript_to text NOT NULL,
    fee int8 NOT NULL,
    txidx int4 NOT NULL,        -- txidx in block
    block_time int4 NOT NULL,   -- equivalent to uint32 in Go
    inscription_number int8 NOT NULL,
    inscription_id text NOT NULL,
    inscription_content jsonb NOT NULL,
    extra_data jsonb,

	 CONSTRAINT brc20_swap_history_pk PRIMARY KEY (id)
);
CREATE INDEX brc20_swap_history_block_height_idx ON public.brc20_swap_history USING btree (block_height);
CREATE INDEX brc20_swap_history_module_id_idx ON public.brc20_swap_history USING btree (module_id);


CREATE TABLE public.brc20_swap_stats (
    id bigserial NOT NULL,
    block_height int4 NOT NULL,
    module_id text NOT NULL,
	 tick varchar(4) NOT NULL,
	 deposit_balance numeric(40) NOT NULL
);
CREATE INDEX brc20_swap_stats_block_height_idx ON public.brc20_swap_stats USING btree (block_height);
CREATE INDEX brc20_swap_stats_module_id_idx ON public.brc20_swap_stats USING btree (module_id);
