const path = require("path");
const TerserPlugin = require('terser-webpack-plugin');
const NodePolyfillPlugin = require('node-polyfill-webpack-plugin');

module.exports = {
    entry: './index.js',
    output: {
        path: path.join(__dirname, 'dist'),
        filename: 'main.js',
        library: 'amf-custom-validator',
        libraryTarget: 'umd',
        globalObject: 'this'
    },
    module: {
      rules: [
        {
         test: /main\.wasm\.gz/,
         type: 'asset/inline',
         generator: {
          dataUrl: content => {
            return Buffer.from(content).toString('base64');
          }
        }
        }
      ]
    },
    plugins: [
        new NodePolyfillPlugin({
            includeAliases: ['stream', 'Buffer']
        })
    ],
    optimization: {
      minimize: true,
      minimizer: [new TerserPlugin({
        extractComments: false,
      })],
    },
}
