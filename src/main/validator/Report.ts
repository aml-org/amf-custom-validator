import {Profile} from "../model/Profile";

export class Report {
    private result: any;
    private violations: any[];
    private infos: any[];
    private warnings: any[];
    constructor(result: any) {
        this.result = result[0].result;
        if (result == undefined) {
            throw new Error("Invalid raw validation data: " + JSON.stringify(result, null, 2) )
        }
        this.violations = this.result.violation.map((data) => this.toValidationResult("violation", data))
        this.infos = this.result.info.map((data) => this.toValidationResult("info", data))
        this.warnings = this.result.warning.map((data) => this.toValidationResult("warning", data))

    }

    toJSON() {
        let json = {
            "@type": "http://www.w3.org/ns/shacl#ValidationReport",
            "http://www.w3.org/ns/shacl#conforms": this.conforms(),
        };

        const results = this.violations.concat(this.warnings).concat(this.infos);
        if (results.length > 0) {
            json["http://www.w3.org/ns/shacl#result"] = results;
        }

        return json;
    }

    private toValidationResult(level: string, result: any) {
        const trace = (<any[]>result.trace).map((traceElement) => {
            let value = traceElement.value;
            if (value["@id"]) {
                value = {"@id": value["@id"]};
            } else {
                value = {"@value": value}
            }

            return {
                "http://a.ml/vocabularies/validation#component": traceElement.component,
                "http://www.w3.org/ns/shacl#resultMessage": traceElement.message,
                "http://www.w3.org/ns/shacl#resultPath": {
                    "@id": traceElement.path
                },
                "http://www.w3.org/ns/shacl#focusNode": value
            }
        });
        return {
            "@type": "http://www.w3.org/ns/shacl#ValidationResult",
            "http://www.w3.org/ns/shacl#resultSeverity": {
                "@id": "http://www.w3.org/ns/shacl#" + this.capitalize(level)
            },
            "http://www.w3.org/ns/shacl#focusNode": {
                "@id": result.target
            },
            "http://a.ml/vocabularies/validation#trace": trace,
            "http://www.w3.org/ns/shacl#resultMessage": result.message,
            "http://www.w3.org/ns/shacl#sourceShape": {
                "@id": result.shapeId
            }
        }
    }

    private capitalize(s: string){
        return s.charAt(0).toUpperCase() + s.slice(1)
    }


    conforms() {
        return this.violations.length == 0;
    }
}