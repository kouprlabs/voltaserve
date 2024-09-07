// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'

const ConsolePanelDatabaseOverview = () => {
  return (
    <>
      <Helmet>
        <title>Database Management</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}></div>
    </>
  )
}

export default ConsolePanelDatabaseOverview
