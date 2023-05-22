"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const file_1 = require("../api/file");
const config_1 = require("../config/config");
/*
  This method is similar to GET but only retrieves the metadata of a resource, without returning the actual content.

  Example implementation:

  - Extract the file path from the URL.
  - Retrieve the file metadata using fs.stat().
  - Set the response status code to 200 if successful or an appropriate error code if the file is not found.
  - Set the Content-Length header with the file size.
  - Return the response.
*/
async function handleHead(req, res, token) {
    try {
        const result = await fetch(`${config_1.API_URL}/v1/files/list?path=${req.url}`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token.access_token}`,
                'Content-Type': 'application/json',
            },
        });
        const file = (await result.json())[0];
        if (file.type === file_1.FileType.File) {
            res.statusCode = 200;
            res.setHeader('Content-Length', file.original.size);
            res.end();
        }
        else {
            res.statusCode = 200;
            res.end();
        }
    }
    catch (err) {
        console.error(err);
        res.statusCode = 500;
        res.end();
    }
}
exports.default = handleHead;
//# sourceMappingURL=handle-head.js.map