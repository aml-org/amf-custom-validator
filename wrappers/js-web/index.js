require("../js/lib/wasm_exec");
let wasm_gz = require("../js/lib/main.wasm.gz")
const pako = require("pako");

let initialized = false
let go = null;
let wasm;

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
        cb(res,null);
    } else {
        cb(null,new Error("WASM/GO not initialized"))
    }
}

const initialize = function(cb) {
    if (initialized === true) {
        cb(null);
    }
    go = new Go();
    if(!wasm_gz || !wasm) {
        wasm = pako.ungzip(Buffer.from(wasm_gz, 'base64'))
    }
    if (WebAssembly) {
        WebAssembly.instantiate(wasm, go.importObject).then((result) => {
            go.run(result.instance);
            initialized = true;
            cb(null);
        });
    } else {
        cb(new Error("WebAssembly is not supported in your JS environment"));
    }

}

const exit = function() {
    if(initialized) {
        __AMF__terminateValidator()
        go.exit(0)
        initialized = false;
    }
}

module.exports.initialize = initialize;
module.exports.validate = validateCustomProfile;
module.exports.exit = exit;