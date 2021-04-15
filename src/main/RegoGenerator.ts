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
        this.ruleSet.forEach((canonicalExpression) => {
            new ExpressionGenerator(canonicalExpression).generate().forEach((l) => regoRules.push(l));
        });

        return [this.package(), this.prologue()].concat(regoRules).join("\n")
    }

    private package() {
        return "package " + this.profile.packageName() + "\n".trim();
    }

    private prologue() {
        return fs.readFileSync("./src/main/resources/prologue.rego").toString().trim();
    }
}