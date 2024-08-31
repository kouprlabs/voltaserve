// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { CSSProperties, MouseEvent, ReactNode } from 'react'
import { Tr } from '@chakra-ui/react'
import { useColorModeValue } from '@chakra-ui/system'

export interface ConsoleHighlightableProps {
  onClick: (event: MouseEvent) => void
  style?: CSSProperties
  children: ReactNode
}

const ConsoleHighlightableTr = (props: ConsoleHighlightableProps) => {
  const hoverBg = useColorModeValue('gray.300', 'gray.700')

  return (
    <Tr
      _hover={{
        backgroundColor: hoverBg,
      }}
      style={{ ...props.style, cursor: 'pointer' }}
      onClick={props.onClick}
    >
      {props.children}
    </Tr>
  )
}

export default ConsoleHighlightableTr
