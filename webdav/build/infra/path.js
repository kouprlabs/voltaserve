"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.getTargetPath = void 0;
const path_1 = __importDefault(require("path"));
function getTargetPath(req) {
    const destinationHeader = req.headers.destination;
    if (!destinationHeader) {
        return null;
    }
    // Check if the destination header is a full URL
    if (destinationHeader.startsWith('http://') ||
        destinationHeader.startsWith('https://')) {
        const url = new URL(destinationHeader);
        return path_1.default.join(decodeURIComponent(url.pathname));
    }
    else {
        /* Extract the path from the destination header */
        const startIndex = destinationHeader.indexOf(req.headers.host) + req.headers.host.length;
        const value = destinationHeader.substring(startIndex);
        return path_1.default.join(decodeURIComponent(value));
    }
}
exports.getTargetPath = getTargetPath;
//# sourceMappingURL=path.js.map