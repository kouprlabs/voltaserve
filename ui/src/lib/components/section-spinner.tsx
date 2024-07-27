// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import cx from 'classnames'
import Spinner from './spinner'

type SectionSpinnerProps = {
  width?: string
  height?: string
}

const DEFAULT_WIDTH = '100%'
const DEFAULT_HEIGHT = '300px'

const SectionSpinner = ({ width, height }: SectionSpinnerProps) => (
  <div
    className={cx('flex', 'items-center', 'justify-center')}
    style={{ width: width || DEFAULT_WIDTH, height: height || DEFAULT_HEIGHT }}
  >
    <Spinner />
  </div>
)

export default SectionSpinner
