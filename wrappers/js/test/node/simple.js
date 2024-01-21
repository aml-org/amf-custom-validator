const fs = require("fs");
var assert = require('assert');

function invalidReport(report) {
    return report[0]["doc:encodes"][0]["conforms"] === false
}

function requireValidator() {
    return require(__dirname + "/../../dist/bundle.js")
}

describe('validator', () => {

    describe('validate', () => {

        it("should load the WASM code, validate a profile, exit", async () => {
            const profile = fs.readFileSync(__dirname + "/../../../../test/data/integration/profile10/profile.yaml").toString()
            const data = fs.readFileSync(__dirname + "/../../../../test/data/integration/profile10/negative.data.jsonld").toString()
            const bundle = requireValidator()
            const validator = await bundle.CustomValidatorFactory.create()
            let result = await validator.validateCustomProfile(profile, data, false)
            let report = JSON.parse(result)
            assert.ok(invalidReport(report))
            result = await validator.validateCustomProfile(profile, data, false)
            report = JSON.parse(result)
            assert.ok(invalidReport(report))
            validator.exit();
        });

        it("should generate the Rego code for a profile and exit", async () => {
            const profile = fs.readFileSync(__dirname + "/../../../../test/data/integration/profile10/profile.yaml").toString()

            const bundle = requireValidator()
            const validator = await bundle.CustomValidatorFactory.create()
            const rego = validator.generateRego(profile)
            assert.ok(rego.indexOf("package profile_kiali") > -1)
            validator.exit();
        })

        it("should normalize input data", async () => {
            const data = fs.readFileSync(__dirname + "/../../../../test/data/integration/profile10/negative.data.jsonld").toString()

            const bundle = requireValidator()
            const validator = await bundle.CustomValidatorFactory.create()
            const normalized = await validator.normalizeInput(data)
            assert.ok(normalized.indexOf("@ids") > -1)
            validator.exit();
        })

        describe('validate with message expressions', () => {

            it("validate and compute message expression value", async () => {
                const profile = fs.readFileSync(__dirname + "/../../../../test/data/integration/profile26/profile.yaml").toString()
                const data = fs.readFileSync(__dirname + "/../../../../test/data/integration/profile26/negative.data.jsonld").toString()
                const bundle = requireValidator()
                const validator = await bundle.CustomValidatorFactory.create()
                const r = validator.validateCustomProfile(profile, data, false)
                let report = JSON.parse(r)
                assert.ok(report[0]["doc:encodes"][0]["result"][0]["resultMessage"] === "Movie 'Disaster Movie' has a rating of 1.9 but it does not have at least 10 reviews (actual reviews: 5) to support that rating")
                validator.exit();
            });
        })

        describe('validate with configuration', () => {

            it("must match expected report output", async () => {
                const profile = fs.readFileSync(__dirname + "/../../../../test/data/integration/profile26/profile.yaml").toString()
                const data = fs.readFileSync(__dirname + "/../../../../test/data/integration/profile26/negative.data.jsonld").toString()
                const bundle = requireValidator()
                const reportConfig = {
                    "IncludeReportCreationTime": false,
                    "ReportSchemaIri": "http://a.ml/report",
                    "LexicalSchemaIri": "http://a.ml/lexical"
                }
                const validator = await bundle.CustomValidatorFactory.create()
                const result = await validator.validateWithReportConfiguration(profile, data, false, reportConfig)
                let report = JSON.parse(result)
                assert.ok(report[0]["doc:encodes"][0]["dateCreated"] === undefined)
                assert.strictEqual(report[0]["@context"]["reportSchema"], "http://a.ml/report#/declarations/")
                assert.strictEqual(report[0]["@context"]["lexicalSchema"], "http://a.ml/lexical#/declarations/")
                validator.exit();
            });
        })
    })
})