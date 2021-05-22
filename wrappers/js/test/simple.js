const fs = require("fs");

const profile = fs.readFileSync(__dirname + "/../../../test/data/integration/profile10/profile.yaml").toString()
const data = fs.readFileSync(__dirname + "/../../../test/data/integration/profile10/negative.data.jsonld").toString()

const validator = require(__dirname + "/../index")

validator.validate(profile, data, false, (r, err) => {
    if (err) {
        console.log("Error running validation: " + err)
    } else {
        console.log("one")
        console.log(r)
        validator.validate(profile, data, false, (r,err) => {
            console.log("two")
            console.log(r);
            validator.exit();
        });
    }
});