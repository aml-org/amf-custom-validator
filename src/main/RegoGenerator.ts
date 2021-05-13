import {Profile} from "./model/Profile";
import {Expression} from "./model/Expression";
import {ExpressionGenerator} from "./rego_generator/ExpressionGenerator";
import * as fs from "fs";

export class RegoGenerator {
    private ruleSet: Expression[];
    private profile: Profile;

    constructor(profile: Profile) {
        this.profile = profile;
        this.ruleSet = profile.toRuleSet();
    }

    generate() {
        const regoRules: string[] = [];
        this.ruleSet.forEach((expression) => {
            regoRules.push(new ExpressionGenerator(expression).generate());
        });

        return [this.package(), this.preamble()].concat(regoRules).join("\n")
    }

    private package() {
        return "package " + this.profile.packageName() + "\n".trim();
    }

    private preamble() {
        const preamble = [];
        preamble.push(fs.readFileSync("./src/main/resources/preamble.rego").toString().trim());
        if (this.profile.violations.length == 0) {
            preamble.push("default violation = []");
        }
        if (this.profile.warnings.length == 0) {
            preamble.push("default warning = []");
        }
        if (this.profile.infos.length == 0) {
            preamble.push("default info = []");
        }

        return preamble.join("\n\n");
    }
}