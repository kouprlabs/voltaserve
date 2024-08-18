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

struct Config {
    private(set) var token: VOToken.Value?
    let apiURL = "http://\(ProcessInfo.processInfo.environment["API_HOST"] ?? "localhost"):8080/v2"
    let idpURL = "http://\(ProcessInfo.processInfo.environment["IDP_HOST"] ?? "localhost"):8081/v2"

    mutating func connect() async throws {
        token = try await VOToken(baseURL: idpURL).exchange(VOToken.ExchangeOptions(
            grantType: .password,
            username: ProcessInfo.processInfo.environment["USERNAME"]!,
            password: ProcessInfo.processInfo.environment["PASSWORD"]!,
            refreshToken: nil,
            locale: nil
        ))
    }
}
