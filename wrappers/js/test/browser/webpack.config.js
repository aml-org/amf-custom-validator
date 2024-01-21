// noinspection JSPrimitiveTypeWrapperUsage
const path = require('path');
module.exports = {
    mode: 'none',
    output: {
        path: path.join(__dirname, 'out'),
        filename: '[name].js',
        library: 'test',
    },
    module: {
        rules: [
            {
                test: /.*\.yaml$/,
                type: 'asset/inline',
                generator: {
                    dataUrl: content => {
                        return content.toString()
                    }
                }

            },
            {
                test: /.*\.jsonld$/,
                type: 'asset/inline',
                generator: {
                    dataUrl: content => {
                        return content.toString()
                    }
                }
            }
        ]
    },
    entry: {
        simple: path.join(__dirname, 'scripts/simple.js'),
        double: path.join(__dirname, 'scripts/double.js'),
        withConfiguration: path.join(__dirname, 'scripts/withConfiguration.js'),
        messageExpression: path.join(__dirname, 'scripts/messageExpression.js')
    }
}