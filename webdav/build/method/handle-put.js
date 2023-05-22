"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const fs_1 = __importDefault(require("fs"));
const promises_1 = require("fs/promises");
const os_1 = __importDefault(require("os"));
const path_1 = __importDefault(require("path"));
const uuid_1 = require("uuid");
const config_1 = require("../config/config");
/*
  This method creates or updates a resource with the provided content.

  Example implementation:

  - Extract the file path from the URL.
  - Create a write stream to the file.
  - Listen for the data event to write the incoming data to the file.
  - Listen for the end event to indicate the completion of the write stream.
  - Set the response status code to 201 if created or 204 if updated.
  - Return the response.
 */
async function handlePut(req, res, token) {
    try {
        /* Delete existing file (simulate an overwrite) */
        const result = await fetch(`${config_1.API_URL}/v1/files/get?path=${req.url}`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token.access_token}`,
                'Content-Type': 'application/json',
            },
        });
        const file = await result.json();
        await fetch(`${config_1.API_URL}/v1/files/${file.id}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token.access_token}`,
                'Content-Type': 'application/json',
            },
        });
    }
    catch (err) { }
    try {
        const result = await fetch(`${config_1.API_URL}/v1/files/get?path=${path_1.default.dirname(req.url)}`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token.access_token}`,
                'Content-Type': 'application/json',
            },
        });
        const directory = await result.json();
        const filePath = path_1.default.join(os_1.default.tmpdir(), (0, uuid_1.v4)());
        const ws = fs_1.default.createWriteStream(filePath);
        ws.on('error', (error) => {
            console.error(error);
            res.statusCode = 500;
            res.end();
        });
        req.on('data', (chunk) => {
            ws.write(chunk);
        });
        req.on('end', async () => {
            ws.end();
            const params = new URLSearchParams({
                workspace_id: directory.workspaceId,
            });
            params.append('parent_id', directory.id);
            const formData = new FormData();
            const blob = new Blob([await (0, promises_1.readFile)(filePath)]);
            formData.set('file', blob, decodeURIComponent(path_1.default.basename(req.url)));
            try {
                await fetch(`${config_1.API_URL}/v1/files?${params}`, {
                    method: 'POST',
                    headers: {
                        'Authorization': `Bearer ${token.access_token}`,
                    },
                    body: formData,
                });
                res.statusCode = 201;
            }
            catch (err) {
                console.error(err);
                res.statusCode = 500;
            }
            fs_1.default.rmSync(filePath);
            res.end();
        });
    }
    catch (err) {
        console.error(err);
        res.statusCode = 500;
        res.end();
        return;
    }
}
exports.default = handlePut;
//# sourceMappingURL=handle-put.js.map