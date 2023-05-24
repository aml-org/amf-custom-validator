const fs = require("fs");
const { loadPolicy } = require("@open-policy-agent/opa-wasm");
const inputBuffer = fs.readFileSync("api.normalized.jsonld", {encoding: "utf-8"});

// Read the policy wasm file
const policyWasm = fs.readFileSync("policy.wasm");

(async () => {
    try {
        const policy = await loadPolicy(policyWasm);
        const input = JSON.parse(inputBuffer);
        policy.setData({})
        const result = policy.evaluate(input);
        console.log(JSON.stringify(result, null, 2));
    } catch(e) {
        console.error(e)
    }
})()