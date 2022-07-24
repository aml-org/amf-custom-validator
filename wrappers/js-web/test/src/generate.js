// Imports
const profile = require("../../../../test/data/integration/profile10/profile.yaml")
const validator = require("../../dist/main") // Use dist rather than src

// Test
function run() {
    validator.initialize(() => {
        validator.generate(profile, (r, err) => {
            if (err) {
                console.log(err);
            } else {
                let element = document.getElementById('generate');
                element.textContent = r
            }
        });
    })
}

module.exports.run = run;

