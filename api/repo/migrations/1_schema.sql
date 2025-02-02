-- +goose Up

CREATE TABLE public.organization (
	id text NOT NULL,
	"name" text NOT NULL,
	create_time text NOT NULL,
	update_time text NULL,
	CONSTRAINT organization_pkey PRIMARY KEY (id)
);

CREATE TABLE public."snapshot" (
	id text NOT NULL,
	"version" int8 NULL,
	original jsonb NULL,
	preview jsonb NULL,
	"text" jsonb NULL,
	ocr jsonb NULL,
	entities jsonb NULL,
	mosaic jsonb NULL,
	thumbnail jsonb NULL,
	"language" text NULL,
	status text NULL,
	task_id text NULL,
	create_time text NOT NULL,
	update_time text NULL,
	CONSTRAINT snapshot_pkey PRIMARY KEY (id)
);

CREATE TABLE public."user" (
	id text NOT NULL,
	full_name text NOT NULL,
	username text NOT NULL,
	email text NULL,
	password_hash text NOT NULL,
	refresh_token_value text NULL,
	refresh_token_expiry text NULL,
	reset_password_token text NULL,
	email_confirmation_token text NULL,
	email_update_token text NULL,
	email_update_value text NULL,
	is_email_confirmed bool DEFAULT false NOT NULL,
	picture text NULL,
	create_time text NOT NULL,
	update_time text NULL,
	is_admin bool DEFAULT false NOT NULL,
	is_active bool DEFAULT true NOT NULL,
	failed_attempts int4 DEFAULT 0 NOT NULL,
	locked_until text NULL,
	CONSTRAINT user_email_confirmation_token_key UNIQUE (email_confirmation_token),
	CONSTRAINT user_email_key UNIQUE (email),
	CONSTRAINT user_email_update_token_key UNIQUE (email_update_token),
	CONSTRAINT user_email_update_value_key UNIQUE (email_update_value),
	CONSTRAINT user_pkey PRIMARY KEY (id),
	CONSTRAINT user_refresh_token_value_key UNIQUE (refresh_token_value),
	CONSTRAINT user_reset_password_token_key UNIQUE (reset_password_token),
	CONSTRAINT user_username_key UNIQUE (username)
);

CREATE TABLE public."group" (
	id text NOT NULL,
	"name" text NOT NULL,
	organization_id text NOT NULL,
	create_time text NOT NULL,
	update_time text NULL,
	CONSTRAINT group_pkey PRIMARY KEY (id),
	CONSTRAINT group_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE
);
CREATE INDEX group_organization_id_idx ON public."group" USING btree (organization_id);

CREATE TABLE public.grouppermission (
	id text NOT NULL,
	group_id text NULL,
	resource_id text NULL,
	"permission" text NULL,
	create_time text NOT NULL,
	CONSTRAINT grouppermission_pkey PRIMARY KEY (id),
	CONSTRAINT grouppermission_group_id_fkey FOREIGN KEY (group_id) REFERENCES public."group"(id) ON DELETE CASCADE
);
CREATE INDEX grouppermission_group_id_idx ON public.grouppermission USING btree (group_id);
CREATE UNIQUE INDEX grouppermission_group_id_resource_id_idx ON public.grouppermission USING btree (group_id, resource_id);
CREATE INDEX grouppermission_resource_id_idx ON public.grouppermission USING btree (resource_id);

CREATE TABLE public.invitation (
	id text NOT NULL,
	organization_id text NULL,
	owner_id text NULL,
	email text NULL,
	status text DEFAULT 'pending'::text NULL,
	create_time text NOT NULL,
	update_time text NULL,
	CONSTRAINT invitation_pkey PRIMARY KEY (id),
	CONSTRAINT invitation_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE,
	CONSTRAINT invitation_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public."user"(id) ON DELETE CASCADE
);
CREATE INDEX invitation_organization_id_idx ON public.invitation USING btree (organization_id);
CREATE INDEX invitation_user_id_idx ON public.invitation USING btree (owner_id);

CREATE TABLE public.task (
	id text NOT NULL,
	"name" text NOT NULL,
	"error" text NULL,
	percentage int2 NULL,
	is_complete bool DEFAULT false NOT NULL,
	is_indeterminate bool DEFAULT false NOT NULL,
	user_id text NOT NULL,
	status text NULL,
	payload jsonb NULL,
	create_time text NOT NULL,
	update_time text NULL,
	CONSTRAINT task_pkey PRIMARY KEY (id),
	CONSTRAINT task_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE
);

CREATE TABLE public.userpermission (
	id text NOT NULL,
	user_id text NULL,
	resource_id text NULL,
	"permission" text NULL,
	create_time text NOT NULL,
	CONSTRAINT userpermission_pkey PRIMARY KEY (id),
	CONSTRAINT userpermission_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE
);
CREATE INDEX userpermission_resource_id_idx ON public.userpermission USING btree (resource_id);
CREATE INDEX userpermission_user_id_idx ON public.userpermission USING btree (user_id);
CREATE UNIQUE INDEX userpermission_user_id_resource_id_idx ON public.userpermission USING btree (user_id, resource_id);

CREATE TABLE public.file (
	id text NOT NULL,
	"name" text NOT NULL,
	"type" text NOT NULL,
	parent_id text NULL,
	workspace_id text NULL,
	snapshot_id text NULL,
	create_time text NOT NULL,
	update_time text NULL,
	CONSTRAINT file_pkey PRIMARY KEY (id)
);
CREATE INDEX file_parent_id_idx ON public.file USING btree (parent_id);
CREATE INDEX file_workspace_id_idx ON public.file USING btree (workspace_id);

CREATE TABLE public.snapshot_file (
	snapshot_id text NOT NULL,
	file_id text NOT NULL,
	create_time text NOT NULL,
	CONSTRAINT snapshot_file_pkey PRIMARY KEY (snapshot_id, file_id)
);
CREATE INDEX snapshot_file_file_id_idx ON public.snapshot_file USING btree (file_id);
CREATE INDEX snapshot_file_snapshot_id_idx ON public.snapshot_file USING btree (snapshot_id);

CREATE TABLE public.workspace (
	id text NOT NULL,
	"name" text NOT NULL,
	organization_id text NOT NULL,
	storage_capacity int8 NOT NULL,
	root_id text NULL,
	bucket text NOT NULL,
	create_time text NOT NULL,
	update_time text NULL,
	CONSTRAINT workspace_pkey PRIMARY KEY (id),
	CONSTRAINT workspace_root_id_key UNIQUE (root_id)
);
CREATE INDEX workspace_organization_id_idx ON public.workspace USING btree (organization_id);

ALTER TABLE public.file ADD CONSTRAINT file_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES public.file(id) ON DELETE SET NULL;
ALTER TABLE public.file ADD CONSTRAINT file_snapshot_id_fkey FOREIGN KEY (snapshot_id) REFERENCES public."snapshot"(id);
ALTER TABLE public.file ADD CONSTRAINT file_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES public.workspace(id) ON DELETE CASCADE;

ALTER TABLE public.snapshot_file ADD CONSTRAINT snapshot_file_file_id_fkey FOREIGN KEY (file_id) REFERENCES public.file(id) ON DELETE CASCADE;
ALTER TABLE public.snapshot_file ADD CONSTRAINT snapshot_file_snapshot_id_fkey FOREIGN KEY (snapshot_id) REFERENCES public."snapshot"(id);

ALTER TABLE public.workspace ADD CONSTRAINT workspace_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE;
ALTER TABLE public.workspace ADD CONSTRAINT workspace_root_id_fkey FOREIGN KEY (root_id) REFERENCES public.file(id);

-- +goose Down