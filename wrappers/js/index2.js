require("./lib/wasm_exec");
const fs = require('fs')



function validate(profile, data) {
    const go = new Go(); // Defined in wasm_exec.js

    let wasm;
    const wasmBinary = fs.readFileSync("./lib/main.wasm")
    let report;

    function toJsString(ptr, len) {
        const buf = new Uint8Array(wasm.exports.memory.buffer, ptr, len)
        return new TextDecoder('utf8').decode(buf)
    }

    const env = { toJsString }
    go.importObject.env = Object.assign(go.importObject.env, env)

    WebAssembly.instantiate(wasmBinary, go.importObject).then(function (obj) {
        wasm = obj.instance;
        go.run(wasm);
        report = wasm.exports.validate(profile, data);
    })

    return report;
}

module.exports.validate = validate;
