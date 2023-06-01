// eval normal validation vs pre-compiled

// with "time" function from terminal

const fs = require("fs");
const {loadPolicySync, loadPolicy} = require("@open-policy-agent/opa-wasm");
const validator = require('../../wrappers/js/index')


const ITERATIONS = 1;
const PROFILE = __dirname + '/../../test/data/production/best-practices/profile.yaml';
const DATA = __dirname + '/../../test/data/production/best-practices/negative1.raml.jsonld';
const WASM = __dirname + '../js-main/policy.wasm';

function hexStringToByteArray(hexString) {
    // if (hexString.length % 2 !== 0) {
    //     throw "Must have an even number of hex digits to convert to bytes";
    // }/* w w w.  jav  a2 s .  c o  m*/
    let numBytes = hexString.length / 2;
    let byteArray = new Uint8Array(numBytes);
    for (let i = 0; i < numBytes; i++) {
        byteArray[i] = parseInt(hexString.substr(i * 2, 2), 16);
    }
    return byteArray;
}

async function main() {
    if (process.argv.length < 3) {
        noPrecompilation();
    } else if (process.argv[2] === '--pre-compiled') {
        await precompiled();
    } else if (process.argv[2] === '--from-bundle') {
        await fromBundle()
    } else {
        throw new Error("Usage: 'node performance/main/performance.js' with optional --pre-compiled flag");
    }
}

function fromBundle() {
    const inputBuffer = read(DATA);
    const policyWasm = read(WASM);

    (async () => {
        try {
            const policy = await loadPolicy(policyWasm);
            const input = JSON.parse(inputBuffer);
            const result = policy.evaluate(input);
            console.log(JSON.stringify(result, null, 2));
        } catch (e) {
            console.error(e)
        }
    })()
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
                // generate rego policy as WASM
                validator.generateRegoWASM(profile, (wasm, err) => {
                    if (err) {
                        console.log("FOUND ERROR")
                        throw new Error(err);
                    } else {
                        const byteWasm = hexStringToByteArray(wasm)
                        fs.writeFileSync("policy.from-js.wasm", byteWasm)
                        const loadedPolicy = loadPolicySync(byteWasm);
                        // Evaluate n times with WASM
                        for (let i = 1; i < ITERATIONS + 1; i++) {
                            console.log(loadedPolicy.evaluate(data));
                            console.log("validated " + i + " times");
                        }
                    }
                    validator.exit();
                });
            }
        });
    })
}

function read(path) {
    return fs.readFileSync(path, {encoding: "utf-8"});
}

(async () => {
    await main()
})();


