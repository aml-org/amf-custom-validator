import * as amf from '@api-modeling/amf-client-js';
import {model} from "@api-modeling/amf-client-js";
import * as jsonld from 'jsonld';

export class AmfParser {
    private context = {
        "shacl": "http://www.w3.org/ns/shacl#",
        "shapes": "http://a.ml/vocabularies/shapes#",
        "doc": "http://a.ml/vocabularies/document#",
        "apiContract": "http://a.ml/vocabularies/apiContract#",
        "core": "http://a.ml/vocabularies/core#"
    }
    private specUrl: string;


    public static  RAML1 = "RAML 1.0";
    public static OAS3 = "OAS 3.0";
    public static OAS2 = "OAS 2.0";
    public static AMF_GRAPH = "AMF Graph"
    public static ASYNC2 = "ASYNC 2.0"
    public static JSON_SCHEMA = "JSON Schema";

    public static YAML = "application/yaml";
    public static JSON = "application/json";
    public static JSONLD = "application/ld+json";
    private static formatMap : any = {

        "RAML 1.0": {'application/yaml' : amf.Raml10Parser},
        "OAS 2.0": {'application/json': amf.Oas20Parser, 'application/yaml': amf.Oas20YamlParser},
        "OAS 3.0": {'application/json': amf.Oas30Parser, 'application/yaml': amf.Oas30YamlParser},
        "OAS 3.0.0": {'application/json': amf.Oas30Parser, 'application/yaml': amf.Oas30YamlParser},
        // @ts-ignore
        "ASYNC 2.0": {'application/json': amf.Asyn20Parser, 'application/yaml': amf.Async20YamlParser},
        // @ts-ignore
        "Async 2.0": {'application/json': amf.Asyn20Parser, 'application/yaml': amf.Async20YamlParser},
        // @ts-ignore
        "AMF Graph": {'application/json' : amf.AmfGraphParser, 'application/ld+json': amf.AmfGraphParser},
        "AML 1.0": {'applciation/yaml': amf.Aml10Parser}
    }

    private initialized = false;
    private parsed = false;
    private syntax: string;
    private format: string;

    constructor(specUrl: string, format: string, syntax: string) {
        this.specUrl = specUrl;
        if (specUrl.indexOf("://") == -1) {
            this.specUrl = "file://" + this.specUrl;
        }
        this.format = format;
        this.syntax = syntax;
        if (syntax != AmfParser.YAML && syntax != AmfParser.JSON && syntax != AmfParser.JSONLD) {
            throw new Error(`Syntax must be either ${AmfParser.YAML}, ${AmfParser.JSON}, or ${AmfParser.JSONLD}`)
        }
        if (format != AmfParser.JSON_SCHEMA && format != AmfParser.RAML1 && format != AmfParser.OAS2 && format != AmfParser.OAS3 && format != AmfParser.AMF_GRAPH && format != AmfParser.ASYNC2 && format != "Async 2.0") {
            throw new Error(`Format must be either '${AmfParser.RAML1}', '${AmfParser.OAS2}', '${AmfParser.OAS3}', '${AmfParser.JSON_SCHEMA}' 'Async 2.0' or ${AmfParser.AMF_GRAPH}`);
        }

    }

    public async parse(): Promise<any> {
        const baseUnit = await this.parseInput();
        const output = await this.generateOutput(baseUnit);
        const flattened = await jsonld.flatten(JSON.parse(output));
        const normalized = await jsonld.compact(flattened, this.context);
        return normalized
    }

    protected async init() {
        if (!this.initialized) {
            amf.plugins.document.Vocabularies.register();
            amf.plugins.document.WebApi.register();
            await amf.Core.init();
            this.initialized = true;
        }
    }

    protected async parseInput(): Promise<amf.model.document.BaseUnit> {
        await this.init();
        const baseUnit = await amf.Core
            .parser(this.format, this.syntax)
            .parseFileAsync(this.specUrl)
        this.parsed = true;
        return baseUnit;
    }

    protected async generateOutput(baseUnit: amf.model.document.BaseUnit): Promise<string> {
        const generated = await amf.Core
            .generator(AmfParser.AMF_GRAPH, AmfParser.JSONLD)
            .generateString(baseUnit);

        return generated
    }

    private findParser(format: string): any {
        const parserChoices = AmfParser.formatMap[format];
        if (!parserChoices){
            throw new Error("Could not find parser choices for "+this.format);
        }
        const parser = parserChoices[this.syntax];
        if (!parser){
            throw new Error("Could not find parser choices for "+this.format+" with "+this.syntax);
        }

        return parser;
    }
}