// Imports
const profile = require("../../../../test/data/integration/profile10/profile.yaml")
const data = require("../../../../test/data/integration/profile10/negative.data.jsonld")
const validator = require("../../dist/main")

// Test
function run() {
    validator.initialize(() => {
        const reportConfig = {
            "IncludeReportCreationTime": false
        }
        validator.validateWithReportConfiguration(profile, data.toString(), false, reportConfig, (r, err) => {
            if (err) {
                console.log(err);
            } else {
                let element = document.getElementById('report');
                let report = JSON.parse(r)
                let dateCreated = element.textContent = report[0]['doc:encodes'][0]['dateCreated']
                if (dateCreated) {
                    element.textContent = 'has date'
                } else {
                    element.textContent = 'does not have date'
                }
                element.report = report
            }
        });
    })
}

module.exports.run = run;

