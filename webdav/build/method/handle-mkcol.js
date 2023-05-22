"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const path_1 = __importDefault(require("path"));
const config_1 = require("../config/config");
/*
  This method creates a new collection (directory) at the specified URL.

  Example implementation:

  - Extract the directory path from the URL.
  - Use fs.mkdir() to create the directory.
  - Set the response status code to 201 if created or an appropriate error code if the directory already exists or encountered an error.
  - Return the response.
 */
async function handleMkcol(req, res, token) {
    let directory;
    try {
        const result = await fetch(`${config_1.API_URL}/v1/files/get?path=${path_1.default.dirname(req.url)}`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token.access_token}`,
                'Content-Type': 'application/json',
            },
        });
        directory = await result.json();
        await fetch(`${config_1.API_URL}/v1/files/create_folder`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token.access_token}`,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                workspaceId: directory.workspaceId,
                parentId: directory.id,
                name: decodeURIComponent(path_1.default.basename(req.url)),
            }),
        });
        res.statusCode = 201;
        res.end();
    }
    catch (err) {
        console.error(err);
        res.statusCode = 500;
        res.end();
    }
}
exports.default = handleMkcol;
//# sourceMappingURL=handle-mkcol.js.map