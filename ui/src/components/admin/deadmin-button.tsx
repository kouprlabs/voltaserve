// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useRef } from 'react'
import { Link } from 'react-router-dom'
import { IconButton } from '@chakra-ui/react'
import { IconDeAdmin } from '@/lib/components/icons'

const DeAdminButton = () => {
  const buttonRef = useRef<HTMLButtonElement>(null)
  return (
    <>
      <Link to="/" title="User dashbaord">
        <IconButton ref={buttonRef} icon={<IconDeAdmin />} aria-label="" />
      </Link>
    </>
  )
}

export default DeAdminButton
