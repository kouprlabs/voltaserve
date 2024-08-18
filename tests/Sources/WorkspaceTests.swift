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
    var token = Token()

    func testFlow() async throws {
        let token = try await fetchTokenOrFail()

        /* Create organization */
        let organizationClient = VOOrganization(baseURL: config.apiURL, accessToken: token.accessToken)
        let organization = try await organizationClient.create(.init(name: "Test Organization", image: nil))

        /* Create workspaces */
        let workspaceClient = VOWorkspace(baseURL: config.apiURL, accessToken: token.accessToken)
        var workspaces: [VOWorkspace.Entity] = []

        var options: [VOWorkspace.CreateOptions] = []
        for index in 0 ..< 6 {
            options.append(.init(
                name: "Test Workspace \(index)",
                image: nil,
                organizationId: organization.id,
                storageCapacity: 100_000 + index
            ))
        }

        for index in 0 ..< options.count {
            let workspace = try await workspaceClient.create(options[index])
            workspaces.append(workspace)
        }

        /* Test creation */
        for index in 0 ..< options.count {
            XCTAssertEqual(workspaces[index].name, options[index].name)
            XCTAssertEqual(workspaces[index].organization.id, options[index].organizationId)
            XCTAssertEqual(workspaces[index].storageCapacity, options[index].storageCapacity)
        }

        /* Test list */

        /* Page 1 */
        let page1 = try await workspaceClient.fetchList(VOWorkspace.ListOptions(
            query: nil,
            size: 3,
            page: 1,
            sortBy: nil,
            sortOrder: nil
        ))
        XCTAssertGreaterThanOrEqual(page1.totalElements, options.count)
        XCTAssertEqual(page1.page, 1)
        XCTAssertEqual(page1.size, 3)
        XCTAssertEqual(page1.data.count, page1.size)

        /* Page 2 */
        let page2 = try await workspaceClient.fetchList(VOWorkspace.ListOptions(
            query: nil,
            size: 3,
            page: 2,
            sortBy: nil,
            sortOrder: nil
        ))
        XCTAssertGreaterThanOrEqual(page2.totalElements, options.count)
        XCTAssertEqual(page2.page, 2)
        XCTAssertEqual(page2.size, 3)
        XCTAssertEqual(page2.data.count, page2.size)

        /* Test fetch */
        let workspace = try await workspaceClient.fetch(workspaces[0].id)
        XCTAssertGreaterThan(workspace.id.count, 0)
        XCTAssertEqual(workspace.name, workspaces[0].name)
        XCTAssertEqual(workspace.organization.id, workspaces[0].organization.id)
        XCTAssertEqual(workspace.storageCapacity, workspaces[0].storageCapacity)

        /* Test delete */
        for workspace in workspaces {
            try await workspaceClient.delete(workspace.id)
        }
        for workspace in workspaces {
            do {
                _ = try await workspaceClient.fetch(workspace.id)
            } catch let error as VOErrorResponse {
                XCTAssertEqual(error.code, "workspace_not_found")
            }
        }

        /* Delete organization */
        try await organizationClient.delete(organization.id)
    }

    func fetchTokenOrFail() async throws -> VOToken.Value {
        if let token = try await token.fetch() {
            return token
        } else {
            throw FailedToFetchToken()
        }
    }
}
