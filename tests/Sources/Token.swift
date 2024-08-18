// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

import Foundation
import Voltaserve

struct Token {
    private(set) var value: VOToken.Value?
    let config = Config()

    mutating func fetch() async throws -> VOToken.Value? {
        if let value {
            return value
        } else {
            value = try await VOToken(baseURL: config.idpURL).exchange(VOToken.ExchangeOptions(
                grantType: .password,
                username: config.username,
                password: config.password,
                refreshToken: nil,
                locale: nil
            ))
            return value
        }
    }
}
