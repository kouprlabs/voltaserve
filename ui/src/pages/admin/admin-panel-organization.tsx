// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { Heading } from '@chakra-ui/react'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'

const AdminPanelOrganization = () => {
  // if (!list) {
  //   return <SectionSpinner />
  // }

  return (
    <>
      <Helmet>
        <title>Organization management</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>
          Organization management
        </Heading>
      </div>
    </>
  )
}

export default AdminPanelOrganization
