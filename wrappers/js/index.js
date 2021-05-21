require(__dirname + "/lib/wasm_exec");
const fs = require("fs");
const wasm = fs.readFileSync(__dirname + "/lib/main.wasm")

let INIT = false
let go = new Go();

const run = function(profile, data, debug) {
    let before = new Date()
    const res = __AMF__validateCustomProfile(profile,data, debug);
    let after = new Date();
    if (debug) console.log("Ellapsed : " + (after - before))
    return res;
}

const validateCustomProfile = function(profile, data, debug, cb) {
    if (INIT) {
        let res = run(profile, data, debug);
        cb(res,null);
    } else {
        if (WebAssembly) {

            WebAssembly.instantiate(wasm, go.importObject).then((result) => {
                go.run(result.instance);
                INIT = true;
                let res = run(profile, data, debug);
                cb(res,null);
                //go.exit(0)
            });
        } else {
            cb(null,new Error("WebAssembly is not supported in your JS environment"));
        }
    }
}

const exit = function() {
    if(INIT) {
        go.exit(0)
        INIT = false;
        go = new Go();
    }
}

module.exports.validate = validateCustomProfile;
module.exports.exit = exit;