"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const config_1 = require("../config/config");
/*
  This method deletes a resource identified by the URL.

  Example implementation:

  - Extract the file path from the URL.
  - Use fs.unlink() to delete the file.
  - Set the response status code to 204 if successful or an appropriate error code if the file is not found.
  - Return the response.
 */
async function handleDelete(req, res, token) {
    try {
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
        res.statusCode = 204;
        res.end();
    }
    catch (err) {
        console.error(err);
        res.statusCode = 500;
        res.end();
    }
}
exports.default = handleDelete;
//# sourceMappingURL=handle-delete.js.map