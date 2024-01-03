
// Go glue code
require('./lib/wasm_exec');

// WASM
const wasm_gz = require('../js/lib/main.wasm.gz')

// Deps
const pako = require('pako');

// Vars
let initialized = false
let go = undefined;
let wasm;

const validateCustomProfile = function(profile, data, debug, cb) {
    if (initialized) {
        let before = new Date()
        const res = __AMF__validateCustomProfile(profile,data, debug);
        let after = new Date();
        if (debug) console.log('Elapsed : ' + (after - before))
        cb(res, undefined);
    } else {
        cb(undefined, new Error('WASM/GO not initialized'))
    }
}

const validateCustomProfileWithReportConfiguration = function(profile, data, debug, reportConfig, cb) {
    if (initialized) {
        let before = new Date()
        const res = __AMF__validateCustomProfileWithConfiguration(profile,data, debug, undefined, reportConfig);
        let after = new Date();
        if (debug) console.log('Elapsed : ' + (after - before))
        cb(res,undefined);
    } else {
        cb(undefined,new Error('WASM/GO not initialized'))
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
        cb(new Error('WebAssembly is not supported in your JS environment'));
    }
}

const generateRego = function(profile, cb) {
    if (initialized) {
        const res = __AMF__generateRego(profile);
        cb(res,undefined);
    } else {
        cb(undefined,new Error("WASM/GO not initialized"))
    }
}

const normalizeInput = function(data, cb) {
    if (initialized) {
        const res = __AMF__normalizeInput(data);
        cb(res,undefined);
    } else {
        cb(undefined,new Error("WASM/GO not initialized"))
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
module.exports.generateRego = generateRego;
module.exports.normalizeInput = normalizeInput;
module.exports.exit = exit;
