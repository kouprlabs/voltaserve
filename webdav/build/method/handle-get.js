"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const fs_1 = __importDefault(require("fs"));
const http_1 = require("http");
const os_1 = __importDefault(require("os"));
const path_1 = __importDefault(require("path"));
const uuid_1 = require("uuid");
const config_1 = require("../config/config");
/*
  This method retrieves the content of a resource identified by the URL.

  Example implementation:

  - Extract the file path from the URL.
  - Create a read stream from the file and pipe it to the response stream.
  - Set the response status code to 200 if successful or an appropriate error code if the file is not found.
  - Return the response.
 */
async function handleGet(req, res, token) {
    try {
        const result = await fetch(`${config_1.API_URL}/v1/files/list?path=${req.url}`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token.access_token}`,
                'Content-Type': 'application/json',
            },
        });
        const file = (await result.json())[0];
        const filePath = path_1.default.join(os_1.default.tmpdir(), (0, uuid_1.v4)());
        (0, http_1.get)(`${config_1.API_URL}/v1/files/${file.id}/original${file.original.extension}?access_token=${token.access_token}`, (response) => {
            const ws = fs_1.default.createWriteStream(filePath);
            response.pipe(ws);
            ws.on('finish', () => {
                ws.close();
                const rs = fs_1.default.createReadStream(filePath);
                rs.on('error', (error) => {
                    console.error(error);
                    res.statusCode = 500;
                    res.end();
                });
                rs.on('end', () => {
                    fs_1.default.rmSync(filePath);
                });
                rs.pipe(res);
            });
        });
    }
    catch {
        res.statusCode = 500;
        res.end();
    }
}
exports.default = handleGet;
//# sourceMappingURL=handle-get.js.map