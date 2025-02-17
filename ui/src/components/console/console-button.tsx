// Copyright (c) 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useRef } from 'react'
import { Link } from 'react-router-dom'
import { IconButton } from '@chakra-ui/react'
import { IconClose, IconAdminPanelSettings } from '@koupr/ui'

const ConsoleButton = () => {
  const buttonRef = useRef<HTMLButtonElement>(null)
  return (
    <>
      {location.pathname.startsWith('/console') ? (
        <Link to="/">
          <IconButton
            ref={buttonRef}
            icon={<IconClose />}
            title="Close console"
            aria-label="Close console"
          />
        </Link>
      ) : (
        <Link to="/console/dashboard">
          <IconButton
            ref={buttonRef}
            icon={<IconAdminPanelSettings />}
            title="Open console"
            aria-label="Open console"
          />
        </Link>
      )}
    </>
  )
}

export default ConsoleButton
