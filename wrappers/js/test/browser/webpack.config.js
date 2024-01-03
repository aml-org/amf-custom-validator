const path = require('path');

// Used for building tests only
module.exports = {
    entry: {
        simple: path.join(__dirname, 'src/simple.js'),
        double: path.join(__dirname, 'src/double.js'),
        withConfiguration: path.join(__dirname, 'src/withConfiguration.js'),
        messageExpression: path.join(__dirname, 'src/messageExpression.js')
    },
    output: {
        path: path.join(__dirname, 'out'),
        filename: '[name].js',
        library: 'test',
    },
    mode: "development",
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
    }
}
