// Imports
const profile = require("../../../../../test/data/integration/profile26/profile.yaml")
const data = require("../../../../../test/data/integration/profile26/negative.data.jsonld")
const bundle = require("../../../dist/bundle") // Use dist rather than src
const factory = bundle.CustomValidatorFactory

// Test
async function run() {
    const validator = await factory.create()
    let report = await validator.validateCustomProfile(profile, data.toString(), false)
    let element = document.getElementById('report');
    element.textContent = report
    element.report = JSON.parse(report);
}

module.exports.run = run;

