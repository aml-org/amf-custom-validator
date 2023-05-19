// eval normal validation vs pre-compiled

// with "time" function from terminal

const fs = require("fs");
const {loadPolicySync} = require("@open-policy-agent/opa-wasm");
const validator = require('../../wrappers/js/index')
const assert = require("assert");


const ITERATIONS = 1;
const PROFILE = __dirname + '/../../test/data/production/best-practices/profile.yaml';
const DATA = __dirname + '/../../test/data/production/best-practices/negative1.raml.jsonld';

String.prototype.toUTF8ByteArray = function() {
    var bytes = [];

    var s = unescape(encodeURIComponent(this));

    for (var i = 0; i < s.length; i++) {
        var c = s.charCodeAt(i);
        bytes.push(c);/*  w  w  w.  j ava  2  s  . c o m*/
    }

    return bytes;
};

function main() {
    if (process.argv.length < 3) {
        noPrecompilation();
    } else if (process.argv[2] === '--pre-compiled') {
        precompiled();
    } else {
        throw new Error("Usage: 'node performance/main/performance.js' with optional --pre-compiled flag");
    }
}

function noPrecompilation() {
    const profile = read(PROFILE);
    const data = read(DATA);

    validator.initialize(() => {
        // Evaluate n times
        for (let i = 0; i < ITERATIONS; i++) {
            validator.validate(profile, data, false, (r, err) => {
                if (err) throw new Error(err);
            });
            console.log("validated " + i + 1 + " times")
        }
        validator.exit();
    });
}

function toBytesArray(str) {
    const encoder = new TextEncoder();
    return encoder.encode(str);
}

function precompiled() {
    const profile = read(PROFILE);
    // const data = JSON.parse(read(DATA));
    const inputJsonLD = read(DATA);

    validator.initialize(() => {
        validator.normalizeInput(inputJsonLD, (data, err) => {
            if (err) {
                console.log("FOUND ERROR")
                throw new Error(err);
            } else {
                console.log(profile)
                console.log(data)
                // generate rego policy as WASM
                validator.generateRegoWASM(profile, (wasm, err) => {
                    if (err) {
                        console.log("FOUND ERROR")
                        throw new Error(err);
                    } else {
                        const decoded = toBytesArray(wasm).buffer;
                        // let byteArray = utf8Encode.encode(wasm);
                        // console.log(wasm)
                        // Convert byte array to ArrayBuffer
                        // const arrayBuffer = byteArray


                        // Create WebAssembly.Module from the ArrayBuffer
                        // const wasmModule = new WebAssembly.Module(decoded);

                        const loadedPolicy = loadPolicySync(decoded);
                        // Evaluate n times with WASM
                        for (let i = 0; i < ITERATIONS; i++) {
                            loadedPolicy.evaluate(data);
                            console.log("validated " + i + 1 + " times")
                        }
                    }
                    validator.exit();
                });
            }
        });
    })


    // // Read the policy wasm file
    // const policyWasm = fs.readFileSync("policy.wasm");
}

function read(path) {
    return fs.readFileSync(path, {encoding: "utf-8"});
}

main();


