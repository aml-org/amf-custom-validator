require(__dirname + "/lib/wasm_exec");
const fs = require("fs");
const pako = require("pako");
let wasm_gz
let wasm

let initialized = false
let go = undefined;

const run = function (profile, data, debug) {
    let before = new Date()
    const res = __AMF__validateCustomProfile(profile, data, debug);
    let after = new Date();
    if (debug) console.log("Elapsed : " + (after - before))
    return res;
}

const validateCustomProfile = function (profile, data, debug, cb) {
    if (initialized) {
        let res = run(profile, data, debug);
        cb(res, undefined);
    } else {
        cb(undefined, new Error("WASM/GO not initialized"))
    }
}

const generate = function (profile, cb) {
    if (initialized) {
        let res = __AMF__generate(profile);
        cb(res, undefined);
    } else {
        cb(undefined, new Error("WASM/GO not initialized"))
    }
}

const normalize = function (data, cb) {
    if (initialized) {
        let res = __AMF__normalize(data);
        cb(res, undefined);
    } else {
        cb(undefined, new Error("WASM/GO not initialized"))
    }
}
const initialize = function (cb) {
    if (initialized === true) {
        cb(undefined)
    }
    go = new Go();
    if (!wasm_gz || !wasm) {
        wasm_gz = fs.readFileSync(__dirname + "/lib/main.wasm.gz")
        wasm = pako.ungzip(wasm_gz)
    }
    if (WebAssembly) {
        WebAssembly.instantiate(wasm, go.importObject).then((result) => {
            go.run(result.instance);
            initialized = true;
            cb(undefined);
        });
    } else {
        cb(new Error("WebAssembly is not supported in your JS environment"));
    }
}

const exit = function () {
    if (initialized) {
        __AMF__terminateValidator()
        go.exit(0)
        initialized = false;
    }
}

module.exports.initialize = initialize;
module.exports.validate = validateCustomProfile;
module.exports.generate = generate;
module.exports.normalize = normalize;
module.exports.exit = exit;