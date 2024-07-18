-- Copyright 2023 Anass Bouassaba.
--
-- Use of this software is governed by the Business Source License
-- included in the file licenses/BSL.txt.
--
-- As of the Change Date specified in that file, in accordance with
-- the Business Source License, use of this software will be governed
-- by the GNU Affero General Public License v3.0 only, included in the file
-- licenses/AGPL.txt.

SET DATABASE = voltaserve;

CREATE TABLE IF NOT EXISTS "file"
(
    id           text PRIMARY KEY,
    name         text NOT NULL,
    type         text NOT NULL,
    parent_id    text,
    workspace_id text REFERENCES workspace (id) ON DELETE CASCADE,
    snapshot_id  text,
    create_time  text NOT NULL DEFAULT (to_json(now())#>>'{}'),
    update_time  text ON UPDATE (to_json(now())#>>'{}')
);

CREATE INDEX IF NOT EXISTS file_parent_id_idx ON "file" (parent_id);
CREATE INDEX IF NOT EXISTS file_workspace_id_idx ON "file" (workspace_id);

CREATE TABLE IF NOT EXISTS "snapshot"
(
  id          text PRIMARY KEY,
  version     bigint,
  original    jsonb,
  preview     jsonb,
  text        jsonb,
  ocr         jsonb,
  entities    jsonb,
  mosaic      jsonb,
  thumbnail   jsonb,
  language    text,
  status      text,
  task_id     text,
  create_time text NOT NULL DEFAULT (to_json(now())#>>'{}'),
  update_time text ON UPDATE (to_json(now())#>>'{}')
);

CREATE TABLE IF NOT EXISTS snapshot_file
(
    snapshot_id text REFERENCES snapshot (id) ON DELETE CASCADE,
    file_id     text REFERENCES file (id) ON DELETE CASCADE,
    create_time text NOT NULL DEFAULT (to_json(now())#>>'{}'),
    PRIMARY KEY (snapshot_id, file_id)
);

CREATE INDEX IF NOT EXISTS snapshot_file_snapshot_id_idx ON snapshot_file (snapshot_id);
CREATE INDEX IF NOT EXISTS snapshot_file_file_id_idx ON snapshot_file (file_id);

CREATE TABLE IF NOT EXISTS "group"
(
    id              text PRIMARY KEY,
    name            text NOT NULL,
    organization_id text NOT NULL REFERENCES organization (id) ON DELETE CASCADE,
    create_time     text NOT NULL DEFAULT (to_json(now())#>>'{}'),
    update_time     text ON UPDATE (to_json(now())#>>'{}')
);

CREATE INDEX IF NOT EXISTS group_organization_id_idx ON "group" (organization_id);

CREATE TABLE IF NOT EXISTS group_user
(
    group_id    text REFERENCES "group" (id) ON DELETE CASCADE,
    user_id     text REFERENCES "user" (id) ON DELETE CASCADE,
    create_time text NOT NULL DEFAULT (to_json(now())#>>'{}'),
    PRIMARY KEY (group_id, user_id)
);

CREATE INDEX IF NOT EXISTS group_user_group_id_idx ON group_user (group_id);
CREATE INDEX IF NOT EXISTS group_user_user_id_idx ON group_user (user_id);

CREATE TABLE IF NOT EXISTS userpermission
(
    id          text PRIMARY KEY,
    user_id     text REFERENCES "user" (id) ON DELETE CASCADE,
    resource_id text,
    permission  text NOT NULL,
    create_time text NOT NULL DEFAULT (to_json(now())#>>'{}'),
    UNIQUE (user_id, resource_id)
);

CREATE INDEX IF NOT EXISTS userpermission_user_id_idx ON userpermission (user_id);
CREATE INDEX IF NOT EXISTS userpermission_resource_id_idx ON userpermission (resource_id);

CREATE TABLE IF NOT EXISTS grouppermission
(
    id          text PRIMARY KEY,
    group_id    text REFERENCES "group" (id) ON DELETE CASCADE,
    resource_id text,
    permission  text NOT NULL,
    create_time text NOT NULL DEFAULT (to_json(now())#>>'{}'),
    UNIQUE (group_id, resource_id)
);

CREATE INDEX IF NOT EXISTS grouppermission_group_id_idx ON grouppermission (group_id);
CREATE INDEX IF NOT EXISTS grouppermission_resource_id_idx ON grouppermission (resource_id);

CREATE TABLE IF NOT EXISTS invitation
(
  id              text PRIMARY KEY,
  organization_id text NOT NULL REFERENCES organization (id) ON DELETE CASCADE,
  owner_id        text NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
  email           text NOT NULL,
  status          text NOT NULL DEFAULT 'pending',
  create_time     text NOT NULL DEFAULT (to_json(now())#>>'{}'),
  update_time     text ON UPDATE (to_json(now())#>>'{}')
);

CREATE INDEX invitation_organization_id_idx ON invitation (organization_id);
CREATE INDEX invitation_user_id_idx ON invitation (owner_id);

CREATE TABLE IF NOT EXISTS task
(
  id                text PRIMARY KEY,
  name              text NOT NULL,
  error             text,
  percentage        smallint,
  is_complete       boolean NOT NULL DEFAULT FALSE,
  is_indeterminate  boolean NOT NULL DEFAULT FALSE,
  user_id           text NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
  status            text,
  payload           jsonb,
  task_id           text,
  create_time       text NOT NULL DEFAULT (to_json(now())#>>'{}'),
  update_time       text ON UPDATE (to_json(now())#>>'{}')
);
