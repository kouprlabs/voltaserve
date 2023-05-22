"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const http_1 = require("http");
const passport_1 = __importDefault(require("passport"));
const passport_http_1 = require("passport-http");
const config_1 = require("./config/config");
const handle_copy_1 = __importDefault(require("./method/handle-copy"));
const handle_delete_1 = __importDefault(require("./method/handle-delete"));
const handle_get_1 = __importDefault(require("./method/handle-get"));
const handle_head_1 = __importDefault(require("./method/handle-head"));
const handle_mkcol_1 = __importDefault(require("./method/handle-mkcol"));
const handle_move_1 = __importDefault(require("./method/handle-move"));
const handle_options_1 = __importDefault(require("./method/handle-options"));
const handle_propfind_1 = __importDefault(require("./method/handle-propfind"));
const handle_proppatch_1 = __importDefault(require("./method/handle-proppatch"));
const handle_put_1 = __importDefault(require("./method/handle-put"));
const tokens = new Map();
passport_1.default.use(new passport_http_1.BasicStrategy(async (username, password, done) => {
    const formBody = [];
    formBody.push('grant_type=password');
    formBody.push(`username=${encodeURIComponent(username)}`);
    formBody.push(`password=${encodeURIComponent(password)}`);
    try {
        const result = await fetch(`${config_1.IDP_URL}/v1/token`, {
            method: 'POST',
            body: formBody.join('&'),
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
        });
        const token = await result.json();
        tokens.set(username, token);
        return done(null, token);
    }
    catch (err) {
        return done(err, false);
    }
}));
const server = (0, http_1.createServer)((req, res) => {
    passport_1.default.authenticate('basic', { session: false }, async (err, token) => {
        if (err || !token) {
            res.statusCode = 401;
            res.setHeader('WWW-Authenticate', 'Basic realm="WebDAV Server"');
            res.end();
        }
        else {
            const method = req.method;
            switch (method) {
                case 'OPTIONS':
                    await (0, handle_options_1.default)(req, res);
                    break;
                case 'GET':
                    await (0, handle_get_1.default)(req, res, token);
                    break;
                case 'HEAD':
                    await (0, handle_head_1.default)(req, res, token);
                    break;
                case 'PUT':
                    await (0, handle_put_1.default)(req, res, token);
                    break;
                case 'DELETE':
                    await (0, handle_delete_1.default)(req, res, token);
                    break;
                case 'MKCOL':
                    await (0, handle_mkcol_1.default)(req, res, token);
                    break;
                case 'COPY':
                    await (0, handle_copy_1.default)(req, res, token);
                    break;
                case 'MOVE':
                    await (0, handle_move_1.default)(req, res, token);
                    break;
                case 'PROPFIND':
                    await (0, handle_propfind_1.default)(req, res, token);
                    break;
                case 'PROPPATCH':
                    await (0, handle_proppatch_1.default)(req, res);
                    break;
                default:
                    res.statusCode = 501;
                    res.end();
            }
        }
    })(req, res);
});
server.listen(config_1.PORT, () => {
    console.log(`WebDAV server is listening on port ${config_1.PORT}`);
});
//# sourceMappingURL=server.js.map