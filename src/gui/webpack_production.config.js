let path = require('path');
let webpack = require('webpack');

module.exports = {
    mode: 'production',
    entry: { 
        app: [            
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
        new webpack.optimize.OccurrenceOrderPlugin(), 
        new webpack.DefinePlugin({
            'process.env': {
                NODE_ENV: JSON.stringify('production')
            }
        })
    ]
}