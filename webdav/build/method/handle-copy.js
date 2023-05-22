"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const path_1 = __importDefault(require("path"));
const config_1 = require("../config/config");
const path_2 = require("../infra/path");
/*
  This method copies a resource from a source URL to a destination URL.

  Example implementation:

  - Extract the source and destination paths from the headers or request body.
  - Use fs.copyFile() to copy the file from the source to the destination.
  - Set the response status code to 204 if successful or an appropriate error code if the source file is not found or encountered an error.
  - Return the response.
 */
async function handleCopy(req, res, token) {
    try {
        const sourceResult = await fetch(`${config_1.API_URL}/v1/files/get?path=${req.url}`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token.access_token}`,
                'Content-Type': 'application/json',
            },
        });
        const sourceFile = await sourceResult.json();
        const targetResult = await fetch(`${config_1.API_URL}/v1/files/get?path=${path_1.default.dirname((0, path_2.getTargetPath)(req))}`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token.access_token}`,
                'Content-Type': 'application/json',
            },
        });
        const targetFile = await targetResult.json();
        if (sourceFile.workspaceId !== targetFile.workspaceId) {
            res.statusCode = 400;
            res.end();
            return;
        }
        const copyResponse = await fetch(`${config_1.API_URL}/v1/files/${targetFile.id}/copy`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token.access_token}`,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                ids: [sourceFile.id],
            }),
        });
        const clones = await copyResponse.json();
        await fetch(`${config_1.API_URL}/v1/files/${clones[0].id}/rename`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token.access_token}`,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                name: path_1.default.basename((0, path_2.getTargetPath)(req)),
            }),
        });
        res.statusCode = 204;
        res.end();
    }
    catch (err) {
        console.error(err);
        res.statusCode = 500;
        res.end();
    }
}
exports.default = handleCopy;
//# sourceMappingURL=handle-copy.js.map