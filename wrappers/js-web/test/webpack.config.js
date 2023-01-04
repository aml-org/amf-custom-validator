// noinspection JSPrimitiveTypeWrapperUsage
const path = require('path');
// Use the top level config
const config = require('../webpack.config')
config.mode = 'none'
config.output = {
    path: path.join(__dirname, 'out'),
    filename: '[name].js',
    library: 'test',
}
config.module.rules.push(
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
)

config.entry = {
    simple: path.join(__dirname, 'src/simple.js'),
    double: path.join(__dirname, 'src/double.js'),
    messageExpression: path.join(__dirname, 'src/messageExpression.js')
}

module.exports = config