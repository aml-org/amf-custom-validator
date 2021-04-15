import * as yaml from "yaml";
import * as fs from "fs";
import {Profile} from "./model/Profile";
import {ValidationParser} from "./profile_parser/ValidationParser";
import {Level} from "./model/Rule";

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
        let expressions = [];
        expressions = expressions.concat(this.parseViolations());
        return new Profile(this.data.profile, expressions);
    }

    private parseViolations() {
        const violations = this.data.violation || [];
        const validations = violations.map(violation => {
            const validation = this.findValidation(violation)
            return new ValidationParser(violation, validation, Level.Violation).parse();
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