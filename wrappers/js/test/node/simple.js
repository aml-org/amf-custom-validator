const fs = require('fs');
const assert = require('assert');
const path = require('path')
const validator = require('../../dist/main')

const root = path.resolve(__dirname, "..", "..", "..", "..")
function readIntegrationCaseFile(nr, file) {
    const filePath = path.join(root, 'test', 'data', 'integration', `profile${nr}`, file)
    return fs.readFileSync(filePath).toString()
}

function invalidReport(report) {
    return report[0]["doc:encodes"][0]["conforms"] === false
}

describe('validator', () => {

    describe('validate', () => {

        it("should load the WASM code, validate a profile, exit", (done) => {
            const profile = readIntegrationCaseFile('10', 'profile.yaml')
            const data = readIntegrationCaseFile('10', 'negative.data.jsonld')
            validator.initialize(() => {
                validator.validate(profile, data, false, (r, err) => {
                    if (err) {
                        done(err);
                    } else {
                        let report = JSON.parse(r)
                        assert.ok(invalidReport(report))
                        validator.validate(profile, data, false, (r, err) => {
                            if (err) {
                                done(err)
                            } else {
                                let report = JSON.parse(r)
                                assert.ok(invalidReport(report))
                                validator.exit();
                                done();
                            }
                        });
                    }
                });
            })
        });

        it("should generate the Rego code for a profile and exit", (done) => {
            const profile = readIntegrationCaseFile('10', 'profile.yaml')

            validator.initialize(() => {
                validator.generateRego(profile, (r, err) => {
                    if (err) {
                        done(err);
                    } else {
                        assert.ok(r.indexOf("package profile_kiali") > -1)
                        validator.exit();
                        done();
                    }
                });
            })
        })

        it("should normalize input data", (done) => {
            const data = readIntegrationCaseFile('10', 'negative.data.jsonld')

            validator.initialize(() => {
                validator.normalizeInput(data, (r, err) => {
                    if (err) {
                        done(err);
                    } else {
                        assert.ok(r.indexOf("@ids") > -1)
                        validator.exit();
                        done();
                    }
                });
            })
        })

    })

    describe('validate with message expressions', () => {

        it("validate and compute message expression value", (done) => {
            const profile = readIntegrationCaseFile('26', 'profile.yaml')
            const data = readIntegrationCaseFile('26', 'negative.data.jsonld')

            validator.initialize(() => {
                validator.validate(profile, data, false, (r, err) => {
                    if (err) {
                        done(err);
                    } else {
                        let report = JSON.parse(r)
                        assert.ok(report[0]["doc:encodes"][0]["result"][0]["resultMessage"] === "Movie 'Disaster Movie' has a rating of 1.9 but it does not have at least 10 reviews (actual reviews: 5) to support that rating")
                        validator.exit();
                        done();
                    }
                });
            })
        });
    })

    describe('validate with configuration', () => {

        it("must match expected report output", (done) => {
            const profile = readIntegrationCaseFile('26', 'profile.yaml')
            const data = readIntegrationCaseFile('26', 'negative.data.jsonld')

            validator.initialize(() => {
                const reportConfig = {
                    "IncludeReportCreationTime": false,
                    "ReportSchemaIri": "http://a.ml/report",
                    "LexicalSchemaIri": "http://a.ml/lexical"
                }
                validator.validateWithReportConfiguration(profile, data, false, reportConfig, (r, err) => {
                    if (err) {
                        done(err);
                    } else {
                        let report = JSON.parse(r)
                        assert.ok(report[0]["doc:encodes"][0]["dateCreated"] === undefined)
                        assert.strictEqual(report[0]["@context"]["reportSchema"], "http://a.ml/report#/declarations/")
                        assert.strictEqual(report[0]["@context"]["lexicalSchema"], "http://a.ml/lexical#/declarations/")
                        validator.exit();
                        done();
                    }
                });
            })
        });
    })
})
