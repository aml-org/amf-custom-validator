// Imports
const profile = require("../../../../../test/data/integration/profile10/profile.yaml")
const data = require("../../../../../test/data/integration/profile10/negative.data.jsonld")
const validator = require("../../../dist/bundle") // Use dist rather than src

// Test
function run() {
    validator.initialize(() => {
        const reportConfig = {
            "IncludeReportCreationTime": false,
            "ReportSchemaIri": "http://a.ml/report",
            "LexicalSchemaIri": "http://a.ml/lexical"
        }
        validator.validateWithReportConfiguration(profile, data.toString(), false, reportConfig, (r, err) => {
            if (err) {
                console.log(err);
            } else {
                let element = document.getElementById('report');
                let report = JSON.parse(r)
                let dateCreated = element.textContent = report[0]['doc:encodes'][0]['dateCreated']
                let reportSchema = report[0]["@context"]["reportSchema"]
                let lexicalSchema = report[0]["@context"]["lexicalSchema"]

                let dateText
                if (dateCreated) {
                    dateText = 'has date'
                } else {
                    dateText = 'does not have date'
                }

                element.textContent = `${dateText}, ${reportSchema}, ${lexicalSchema}`
                element.report = report
            }
        });
    })
}

module.exports.run = run;

