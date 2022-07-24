// Imports
const data = require("../../../../test/data/integration/profile10/negative.data.jsonld")
const validator = require("../../dist/main") // Use dist rather than src

// Test
function run() {
    validator.initialize(() => {
        validator.normalize(data, (r, err) => {
            if (err) {
                console.log(err);
            } else {
                let element = document.getElementById('normalize');
                element.textContent = r
            }
        });
    })
}

module.exports.run = run;

