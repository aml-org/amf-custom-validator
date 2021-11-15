const fs = require("fs");
var assert = require('assert');

describe('validator', () => {

    describe('validate', () => {

        function invalidReport(report) {
            return report[0]["doc:encodes"][0]["conforms"] === false
        }

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
    })
})