require(__dirname + "/lib/wasm_exec");
const fs = require("fs");
const pako = require("pako");
let wasm_gz
let wasm

let initialized = false
let go = undefined;

const validate = function (profile, data, debug, cb) {
    if (initialized) {
        let before = new Date()
        const res = __AMF__validate(profile, data, debug);
        let after = new Date();
        if (debug) console.log("Elapsed : " + (after - before))
        cb(res, undefined);
    } else {
        cb(undefined, new Error("WASM/GO not initialized"))
    }
}

const validateCompiled = function (compiledProfile, data, debug, cb) {
    if (initialized) {
        let before = new Date()
        const res = __AMF__validateCompiled(compiledProfile, data, debug);
        let after = new Date();
        if (debug) console.log("Elapsed : " + (after - before))
        cb(res, undefined);
    } else {
        cb(undefined, new Error("WASM/GO not initialized"))
    }
}

const compileProfile = function (profile, debug, cb) {
    if (initialized) {
        let before = new Date()
        const res = __AMF__compileProfile(profile, debug);
        let after = new Date();
        if (debug) console.log("Elapsed : " + (after - before))
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
        __AMF__exit()
        go.exit(0)
        initialized = false;
    }
}

module.exports.initialize = initialize;
module.exports.validate = validate;
module.exports.validateCompiled = validateCompiled;
module.exports.compileProfile = compileProfile;
module.exports.exit = exit;