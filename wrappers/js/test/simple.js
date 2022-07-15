const fs = require("fs");
var assert = require('assert');
const validator = require(__dirname + "/../index")
const profile = fs.readFileSync(__dirname + "/../../../test/data/integration/profile10/profile.yaml").toString()
const data = fs.readFileSync(__dirname + "/../../../test/data/integration/profile10/negative.data.jsonld").toString()


describe('validator', () => {
    it("should generate Rego from profile", (done) => {
        validator.initialize(() => {
            validator.generate(profile, (r, err) => {
                if (err) {
                    done(err);
                } else {
                    assert.ok(r.includes("package profile_kiali"))
                    done();
                }
            });
        })
    });

    it("should normalize input JSON-LD", (done) => {
        validator.initialize(() => {
            validator.normalize(data, (r, err) => {
                if (err) {
                    done(err);
                } else {
                    assert.ok(r.includes("@ids"))
                    done();
                }
            });
        })
    });


    describe('validate', () => {

        function invalidReport(report) {
            return report[0]["doc:encodes"][0]["conforms"] === false
        }

        it("should load the WASM code, validate a profile, exit", (done) => {
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
                                done();
                            }
                        });
                    }
                });
            })
        });
    })

    validator.exit();
})