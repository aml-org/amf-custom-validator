import * as yaml from "yaml";
import * as fs from "fs";
import {Profile} from "./model/Profile";
import {ValidationParser} from "./profile_parser/ValidationParser";
import {Level} from "./model/Rule";
import {ExpressionParser} from "./profile_parser/ExpressionParser";

export class ProfileParser {
    private path: string;
    private data: any;

    constructor(profilePath: string) {
        this.path = profilePath;
        if (this.path.startsWith("file://")) {
            this.path = this.path.replace("file://", "");
        }
    }

    async parse() {
        this.data = await this.parseYaml();
        const violations = this.parseValidations(this.data.violation || [], Level.Violation);
        const warnings = this.parseValidations(this.data.warning || [], Level.Warning);
        const infos = this.parseValidations(this.data.info || [], Level.Info);
        return new Profile(this.data.profile, violations, warnings,infos);
    }

    private parseValidations(validationNames: string[], level: string) {
        const validations = validationNames.map( validationName => {
            const validation = this.findValidation(validationName)
            return new ExpressionParser(validationName, validation, Level.Violation).parse();
        });

        return validations;
    }

    private findValidation(validatioName) {
        return this.data.validations[validatioName];
    }

    private async parseYaml(): Promise<any> {
        return new Promise<any>((resolve, reject) => {
            fs.readFile(this.path, (e, data) => {
                if (e) {
                    reject(e);
                } else {
                    const text = data.toString();
                    try {
                        const parsed = yaml.parse(text);
                        resolve(parsed);
                    } catch (e) {
                        reject(e);
                    }
                }
            });
        });
    }
}