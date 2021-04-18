import {Expression} from "./Expression";

export class Profile {
    public readonly validations: Expression[];
    public readonly name: string;
    public readonly violations: Expression[];
    public readonly warnings: Expression[];
    public readonly infos: Expression[];

    constructor(name: string, violations: Expression[], warnings: Expression[], infos: Expression[]) {
        this.name = name;
        this.violations = violations;
        this.warnings = warnings
        this.infos = infos;
        this.validations = this.violations.concat(this.warnings).concat(this.infos);
    }

    /**
     * From applies deMorgan
     */
    toRuleSet() {
        return this.validations.map((validation) => {
           return validation.toCanonicalForm();
        });
    }

    toString() {
        const expressionStrings = this.validations.map((validation) => {
            return validation.toString()
        });
        return expressionStrings.join("\n\n");
    }

    packageName() {
        return this.name.toLowerCase().replace(/ /g, "_");
    }

    entrypoint() {
        return this.packageName() + "/report";
    }

}