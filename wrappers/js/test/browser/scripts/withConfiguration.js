// Imports
const profile = require("../../../../../test/data/integration/profile10/profile.yaml")
const data = require("../../../../../test/data/integration/profile10/negative.data.jsonld")
const bundle = require("../../../dist/bundle") // Use dist rather than src
const factory = bundle.CustomValidatorFactory

const reportConfig = {
    "IncludeReportCreationTime": false,
    "ReportSchemaIri": "http://a.ml/report",
    "LexicalSchemaIri": "http://a.ml/lexical"
}
// Test
async function run() {
    const validator = await factory.create()
    const result = await validator.validateWithReportConfiguration(profile, data.toString(), false, reportConfig)
    let element = document.getElementById('report');
    let report = JSON.parse(result)
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

module.exports.run = run;

