// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

@testable import Voltaserve
import XCTest

struct FailedToFetchToken: Error {}

final class WorkspaceTests: XCTestCase {
    var config = Config()

    func fetchTokenOrFail() async throws -> VOToken.Value {
        try await config.connect()
        if let token = config.token {
            return token
        } else {
            throw FailedToFetchToken()
        }
    }

    func testList() async throws {
        let token = try await fetchTokenOrFail()
        let client = VOWorkspace(baseURL: config.apiURL, accessToken: token.accessToken)
        let result = try await client.fetchList(VOWorkspace.ListOptions(query: nil, size: nil, page: nil, sortBy: nil, sortOrder: nil))

        // Ensure we receive at least one element
        XCTAssertGreaterThan(result.totalElements, 0)
    }

    func testFetch() async throws {
        let token = try await fetchTokenOrFail()
        let client = VOWorkspace(baseURL: config.apiURL, accessToken: token.accessToken)
        let list = try await client.fetchList(VOWorkspace.ListOptions(query: nil, size: nil, page: nil, sortBy: nil, sortOrder: nil))
        let result = try await client.fetch(list.data.first!.id)

        // Ensure we have a valid ID
        XCTAssertGreaterThan(result.id.count, 1)
    }
}
