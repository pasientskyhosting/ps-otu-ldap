const proxy = require('express-http-proxy');
const express = require("express")

let app = express()

if (process.env.NODE_ENV !== 'production') {

    let webpackConfig = require('./webpack.config.js');    
    let compiler = require('webpack')(webpackConfig);

    app.use(require('webpack-dev-middleware')(compiler, {
        contentBase: 'public/',
        publicPath: '/js/',
        stats: { colors: true }
    }));

    app.use(require('webpack-hot-middleware')(compiler));
    
}

app.use('/', express.static('public'));

app.use('/', proxy('localhost:8081'));

app.listen(8080, () => {
    console.log("listening on port 8080")
})
