require(__dirname + "/lib/polyfills");
require(__dirname + "/lib/wasm_exec");
const fs = require("fs");
const pako = require("pako");
let wasm_gz
let wasm

let initialized = false
let go = undefined;

const clearTimeouts = function() {
    go._scheduledTimeouts.forEach((value, key, map) => {
        clearTimeout(value);
        go._scheduledTimeouts.delete(key);
    })
}

const run = function(profile, data, debug) {
    let before = new Date()
    const res = __AMF__validateCustomProfile(profile,data, debug);
    let after = new Date();
    if (debug) console.log("Elapsed : " + (after - before))
    return res;
}

const validateCustomProfile = function(profile, data, debug, cb) {
    if (initialized) {
        let res = run(profile, data, debug);
        cb(res,undefined);
    } else {
        cb(undefined,new Error("WASM/GO not initialized"))
    }
}

const runGenerateRego = function(profile) {
    const res = __AMF__generateRego(profile);
    return res;
}

const generateRego = function(profile, cb) {
    if (initialized) {
        let res = runGenerateRego(profile);
        cb(res,undefined);
    } else {
        cb(undefined,new Error("WASM/GO not initialized"))
    }
}

const runNormalizeInput = function(data) {
    const res = __AMF__normalizeInput(data);
    return res;
}

const normalizeInput = function(data, cb) {
    if (initialized) {
        let res = runNormalizeInput(data);
        cb(res,undefined);
    } else {
        cb(undefined,new Error("WASM/GO not initialized"))
    }
}

const initialize = function(cb) {
    if (initialized === true) {
        cb(undefined)
    }
    go = new Go();
    if(!wasm_gz || !wasm) {
        wasm_gz = fs.readFileSync(__dirname + "/lib/main.wasm.gz")
        wasm = pako.ungzip(wasm_gz)
    }
    if (WebAssembly) {
        WebAssembly.instantiate(wasm, go.importObject).then((result) => {
            go.run(result.instance);
            initialized = true;
            cb(undefined);
        }).catch(rejection => console.log(rejection));
    } else {
        cb(new Error("WebAssembly is not supported in your JS environment"));
    }
}

const exit = function() {
    if(initialized) {
        clearTimeouts() // need to add this because some timeouts are triggered after Go program finished
        __AMF__terminateValidator()
        go.exit(0)
        initialized = false;
    }
}

module.exports.initialize = initialize;
module.exports.validate = validateCustomProfile;
module.exports.generateRego = generateRego;
module.exports.normalizeInput = normalizeInput;
module.exports.exit = exit;