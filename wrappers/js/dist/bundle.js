(function (global, factory) {
    typeof exports === 'object' && typeof module !== 'undefined' ? factory(exports, require('fs'), require('pako')) :
    typeof define === 'function' && define.amd ? define(['exports', 'fs', 'pako'], factory) :
    (global = typeof globalThis !== 'undefined' ? globalThis : global || self, factory(global["amf-custom-validator"] = {}, global.fs, global.pako));
})(this, (function (exports, fs, pako) { 'use strict';

    require(__dirname + "/lib/wasm_exec_node");
    let wasm_gz;
    let wasm;

    let initialized = false;
    let go = undefined;

    const validateKernel = function(profile, data, debug) {
        let before = new Date();
        const res = __AMF__validateCustomProfile(profile,data, debug);
        let after = new Date();
        if (debug) console.log("Elapsed : " + (after - before));
        return res;
    };

    const validateWithReportConfigurationKernel = function(profile, data, debug, reportConfig) {
        let before = new Date();
        const res = __AMF__validateCustomProfileWithConfiguration(profile,data, debug, undefined, reportConfig);
        let after = new Date();
        if (debug) console.log("Elapsed : " + (after - before));
        return res;
    };

    const validateCustomProfile = function(profile, data, debug, cb) {
        if (initialized) {
            let res = validateKernel(profile, data, debug);
            cb(res,undefined);
        } else {
            cb(undefined,new Error("WASM/GO not initialized"));
        }
    };

    const validateCustomProfileWithReportConfiguration = function(profile, data, debug, reportConfig, cb) {
        if (initialized) {
            let res = validateWithReportConfigurationKernel(profile, data, debug, reportConfig);
            cb(res,undefined);
        } else {
            cb(undefined,new Error("WASM/GO not initialized"));
        }
    };

    const runGenerateRego = function(profile) {
        const res = __AMF__generateRego(profile);
        return res;
    };

    const generateRego = function(profile, cb) {
        if (initialized) {
            let res = runGenerateRego(profile);
            cb(res,undefined);
        } else {
            cb(undefined,new Error("WASM/GO not initialized"));
        }
    };

    const runNormalizeInput = function(data) {
        const res = __AMF__normalizeInput(data);
        return res;
    };

    const normalizeInput = function(data, cb) {
        if (initialized) {
            let res = runNormalizeInput(data);
            cb(res,undefined);
        } else {
            cb(undefined,new Error("WASM/GO not initialized"));
        }
    };

    const initialize = function(cb) {
        if (initialized === true) {
            cb(undefined);
        }
        go = new Go();
        if(!wasm_gz || !wasm) {
            wasm_gz = fs.readFileSync(__dirname + "/lib/main.wasm.gz");
            wasm = pako.ungzip(wasm_gz);
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
    };

    const exit = function() {
        if(initialized) {
            __AMF__terminateValidator();
            go.exit(0);
            initialized = false;
        }
    };

    exports.exit = exit;
    exports.generateRego = generateRego;
    exports.initialize = initialize;
    exports.normalizeInput = normalizeInput;
    exports.runGenerateRego = runGenerateRego;
    exports.runNormalizeInput = runNormalizeInput;
    exports.validateCustomProfile = validateCustomProfile;
    exports.validateCustomProfileWithReportConfiguration = validateCustomProfileWithReportConfiguration;
    exports.validateKernel = validateKernel;
    exports.validateWithReportConfigurationKernel = validateWithReportConfigurationKernel;

}));
