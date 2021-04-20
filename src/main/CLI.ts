import {Validator} from './Validator';
import * as program from 'commander';

const usage = (error: string) =>{
    return "amf-op-validator: " + error + "\nOPTIONS: -in '' -mime-in 'application/yaml'|'application/ld+json'|'application/json' -cp CUSTOM_PROFILE_PATH [-cpw CUSTOM_PROFILE_WASM]";
}

const cmd = new program.Command("amf-opa-validator")
cmd.version("0.0.1");
cmd.usage("npm run [OPTIONS] <file>")
cmd.requiredOption("-in, --format-in <format>",  "Format of the input fileRAML 1.0'|'OAS 3.0'|'OAS 2.0'|'ASYNC 2.0'|'AML 1.0")
    .requiredOption("-mime-in, --media-type-in <mediaType>",  "Media type of the input file 'application/yaml', 'application/ld+json', 'application/json'")
    .option("-cp, --custom-profile [profile]", "Validation profile in AML language")
    .option("-cpw, --custom-profile-wasm [profile]", "Validation profile parsed as OPA web assembly bundle")
    .option("-d, --debug", "Traces transformations")
    .arguments("<file>");

cmd.parse(process.argv);

if (cmd.opts().customProfile == null && cmd.opts().customProfileWasm == null) {
    throw new Error("Missing profile or WASM profile\n" + cmd.usage());
}

const validator = new Validator(cmd.args[0], cmd.opts().formatIn, cmd.opts().mediaTypeIn, cmd.opts().customProfile, cmd.opts().customProfileWasm, cmd.opts().debug);
validator.validate().then((result) => {
    console.log(JSON.stringify(result.toJSON(), null,2));
    process.exit(0)
}).catch((e) => {
    console.error(e);
    process.exit(1)
})