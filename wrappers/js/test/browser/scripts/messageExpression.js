// Imports
const profile = require("../../../../../test/data/integration/profile26/profile.yaml")
const data = require("../../../../../test/data/integration/profile26/negative.data.jsonld")
const validator = require("../../../dist/bundle") // Use dist rather than src

// Test
function run() {
    validator.initialize(() => {
        validator.validateCustomProfile(profile, data.toString(), false, (r, err) => {
            if (err) {
                console.log(err);
            } else {
                let element = document.getElementById('report');
                element.textContent = r
                element.report = JSON.parse(r);
            }
        });
    })
}

module.exports.run = run;

