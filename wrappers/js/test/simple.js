const fs = require("fs");
var assert = require('assert');

function invalidReport(report) {
    return report[0]["doc:encodes"][0]["conforms"] === false
}

describe('validator', () => {

    describe('validate', () => {

        it("should load the WASM code, validate a profile, exit", (done) => {
            const profile = fs.readFileSync(__dirname + "/../../../test/data/integration/profile10/profile.yaml").toString()
            const data = fs.readFileSync(__dirname + "/../../../test/data/integration/profile10/negative.data.jsonld").toString()
            const validator = require(__dirname + "/../index")
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

        it ("should generate the Rego code for a profile and exit", (done) => {
            const profile = fs.readFileSync(__dirname + "/../../../test/data/integration/profile10/profile.yaml").toString()

            const validator = require(__dirname + "/../index")

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

        it ("should normalize input data", (done) => {
            const data = fs.readFileSync(__dirname + "/../../../test/data/integration/profile10/negative.data.jsonld").toString()

            const validator = require(__dirname + "/../index")

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
            const profile = fs.readFileSync(__dirname + "/../../../test/data/integration/profile26/profile.yaml").toString()
            const data = fs.readFileSync(__dirname + "/../../../test/data/integration/profile26/negative.data.jsonld").toString()
            const validator = require(__dirname + "/../index")
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
})