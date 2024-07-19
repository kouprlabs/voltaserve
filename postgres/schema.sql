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
