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

        return [this.package(), this.prologue()].concat(regoRules).join("\n")
    }

    private package() {
        return "package " + this.profile.packageName() + "\n".trim();
    }

    private prologue() {
        const prologue = [];
        prologue.push(fs.readFileSync("./src/main/resources/prologue.rego").toString().trim());
        if (this.profile.violations.length == 0) {
            prologue.push("default violation = []");
        }
        if (this.profile.warnings.length == 0) {
            prologue.push("default warning = []");
        }
        if (this.profile.infos.length == 0) {
            prologue.push("default info = []");
        }

        return prologue.join("\n\n");
    }
}