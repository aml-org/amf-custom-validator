import {Expression} from "./Expression";

export class Profile {
    public readonly validations: Expression[];
    public readonly name: string;

    constructor(name: string, validations: Expression[]) {
        this.name = name;
        this.validations = validations;
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
        return this.packageName() + "/violation";
    }
}