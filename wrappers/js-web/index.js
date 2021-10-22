require("../js/lib/wasm_exec");
let wasm_gz = require("../js/lib/main.wasm.gz")
const pako = require("pako");
let wasm

let INIT = false
let go = new Go();

const run = function(profile, data, debug) {
    let before = new Date()
    const res = __AMF__validateCustomProfile(profile,data, debug);
    let after = new Date();
    if (debug) console.log("Elapsed : " + (after - before))
    return res;
}

const validateCustomProfile = function(profile, data, debug, cb) {
    if (INIT) {
        let res = run(profile, data, debug);
        cb(res,null);
    } else {
        if(!wasm_gz || !wasm) {
            wasm = pako.ungzip(Buffer.from(wasm_gz, 'base64'))
        }
        if (WebAssembly) {
            WebAssembly.instantiate(wasm, go.importObject).then((result) => {
                go.run(result.instance);
                INIT = true;
                let res = run(profile, data, debug);
                cb(res,null);
            });
        } else {
            cb(null,new Error("WebAssembly is not supported in your JS environment"));
        }
    }
}

const exit = function() {
    if(INIT) {
        __AMF__terminateValidator()
        go.exit(0)
        INIT = false;
        go = new Go();
    }
}

module.exports.validate = validateCustomProfile;
module.exports.exit = exit;