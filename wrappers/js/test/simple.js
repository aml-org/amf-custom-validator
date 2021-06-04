const fs = require("fs");
var assert = require('assert');

describe('validator', () => {

    describe('validate', () => {
        it("should load the WASM code, validate a profile, exit", (done) => {
            const profile = fs.readFileSync(__dirname + "/../../../test/data/integration/profile10/profile.yaml").toString()
            const data = fs.readFileSync(__dirname + "/../../../test/data/integration/profile10/negative.data.jsonld").toString()

            const validator = require(__dirname + "/../index")

            validator.validate(profile, data, false, (r, err) => {
                if (err) {
                    done(err);
                } else {
                    let report = JSON.parse(r)
                    assert.ok(report["http://www.w3.org/ns/shacl#conforms"] === false)
                    validator.validate(profile, data, false, (r,err) => {
                        if (err) {
                            done(err)
                        } else {
                            let report = JSON.parse(r)
                            assert.ok(report["http://www.w3.org/ns/shacl#conforms"] === false)
                            validator.exit();
                            done();
                        }
                    });
                }
            });
        });
    })
})