require("/lib/wasm_exec_node");
let wasm_gz = require("../js/lib/main.wasm.gz")
const pako = require("pako");
const Buffer = require("buffer").Buffer;

let initialized = false
let go = undefined;
let wasm;

const validateKernel = function(profile, data, debug) {
    let before = new Date()
    const res = __AMF__validateCustomProfile(profile,data, debug);
    let after = new Date();
    if (debug) console.log("Elapsed : " + (after - before))
    return res;
}

const validateWithReportConfigurationKernel = function(profile, data, debug, reportConfig) {
    let before = new Date()
    const res = __AMF__validateCustomProfileWithConfiguration(profile,data, debug, undefined, reportConfig);
    let after = new Date();
    if (debug) console.log("Elapsed : " + (after - before))
    return res;
}

const validateCustomProfile = function(profile, data, debug, cb) {
    if (initialized) {
        let res = validateKernel(profile, data, debug);
        cb(res, undefined);
    } else {
        cb(undefined, new Error("WASM/GO not initialized"))
    }
}

const validateCustomProfileWithReportConfiguration = function(profile, data, debug, reportConfig, cb) {
    if (initialized) {
        let res = validateWithReportConfigurationKernel(profile, data, debug, reportConfig);
        cb(res,undefined);
    } else {
        cb(undefined,new Error("WASM/GO not initialized"))
    }
}

const initialize = function(cb) {
    if (initialized === true) {
        cb(undefined);
    }
    go = new Go();
    if(!wasm_gz || !wasm) {
        wasm = pako.ungzip(Buffer.from(wasm_gz, 'base64'))
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

const exit = function() {
    if(initialized) {
        __AMF__terminateValidator()
        go.exit(0)
        initialized = false;
    }
}

module.exports.initialize = initialize;
module.exports.validate = validateCustomProfile;
module.exports.validateWithReportConfiguration = validateCustomProfileWithReportConfiguration;
module.exports.exit = exit;
