let path = require('path');
let webpack = require('webpack');

module.exports = {
    mode: 'development',
    entry: { 
        app: [
            'webpack-hot-middleware/client',
            './index.tsx'
        ] 
    },
    output: {
        path: path.resolve(__dirname, 'public', 'js'),
        filename: 'app.js',
        publicPath: '/js/'
    },
    module: {
        rules: [{
            test: /\.(js|ts|tsx)$/,            
            loaders: ['babel-loader']
        }]
    },
    resolve: {
        extensions: ['.js', '.jsx', '.ts', '.tsx']
    },
    plugins: [
        new webpack.HotModuleReplacementPlugin()
    ]
}