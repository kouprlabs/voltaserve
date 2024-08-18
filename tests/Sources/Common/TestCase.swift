// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

import Voltaserve
import XCTest

extension XCTestCase {
    func fetchTokenOrFail() async throws -> VOToken.Value {
        var token = Token()
        if let value = try await token.fetch() {
            return value
        } else {
            throw FailedToFetchToken()
        }
    }

    struct FailedToFetchToken: Error {}
}
