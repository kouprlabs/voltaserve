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

final class VoltaserveTests: XCTestCase {
    var token = VOToken.Value(
        // swiftlint:disable:next line_length
        accessToken: "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJaeEtHcWJXTmIiLCJpYXQiOjE3MjM4NzM0ODYsImlzcyI6ImxvY2FsaG9zdCIsImF1ZCI6ImxvY2FsaG9zdCIsImV4cCI6MTcyNjQ2NTQ4Nn0.8v4tVsqBOduAzVmpTlFut-VG7XsfksXWCee8jl3eTOQ",
        expiresIn: 1_726_465_486,
        tokenType: "Bearer",
        refreshToken: "f7599a593043424eb74e1a3c3614146e"
    )
    var apiURL = "http://\(ProcessInfo.processInfo.environment["API_HOST"] ?? "localhost"):8080/v2"
    
    func testWorkspaceList() async throws {
        let client = VOWorkspace(baseURL: apiURL, accessToken: token.accessToken)
        let result = try await client.fetchList(VOWorkspace.ListOptions(query: nil, size: nil, page: nil, sortBy: nil, sortOrder: nil))
        XCTAssertGreaterThan(result.totalElements, 0)
    }
}
