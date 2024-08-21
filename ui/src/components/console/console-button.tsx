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
import { IconAdmin, IconDeAdmin } from '@/lib/components/icons'

const ConsoleButton = () => {
  const buttonRef = useRef<HTMLButtonElement>(null)
  return (
    <>
      {location.pathname.startsWith('/console') ? (
        <Link to="/" title="User dashbaord">
          <IconButton ref={buttonRef} icon={<IconDeAdmin />} aria-label="" />
        </Link>
      ) : (
        <Link to="/console/dashboard" title="Cloud console dashbaord">
          <IconButton ref={buttonRef} icon={<IconAdmin />} aria-label="" />
        </Link>
      )}
    </>
  )
}

export default ConsoleButton
