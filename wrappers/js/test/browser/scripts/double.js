// Imports
const profile = require("../../../../../test/data/integration/profile10/profile.yaml")
const data = require("../../../../../test/data/integration/profile10/negative.data.jsonld")
const bundle = require("../../../dist/bundle") // Use dist rather than src
const factory = bundle.CustomValidatorFactory

// Test
async function run() {
    const validator = await factory.create()
    let report1 = await validator.validateCustomProfile(profile, data.toString(), false)
    let report2 = await validator.validateCustomProfile(profile, data.toString(), false)
    let element = document.getElementById('report');
    element.textContent = report2
    element.report = JSON.parse(report2);
}

module.exports.run = run;

