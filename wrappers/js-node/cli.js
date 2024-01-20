const validator = require("./index.js");
const fs = require("fs");

const args = process.argv.slice(2);

const profile = fs.readFileSync(process.cwd() + "/" + args[0]).toString();
const data = fs.readFileSync(process.cwd() + "/" + args[1]).toString();
const debug = args[2] === 'true';

validator.initialize((any) => validator.validate(profile, data, debug, (result, other) => {
        console.log(result);
        validator.exit();
    })
)

