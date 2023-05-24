process.env['NODE_TLS_REJECT_UNAUTHORIZED'] = 0
const express = require('express');
const fs = require('fs');
const http = require('http');
const https = require("https");
const bodyParser = require('body-parser');
const methodOverride = require('method-override');
const helmet = require('helmet');
const async = require('async'); 
const app = express();
const createReq = require('./utils').createRequest;
const requestLibrary = require('request');
const moment = require('moment');

let started = false;

function start (done) {

  // app.use(helmet());

  // app.use(bodyParser.urlencoded({
  //   extended: true,
  // }));

  // app.use(methodOverride());
  
 app.use(bodyParser.json({ limit: '5mb' }));
 app.disable('x-powered-by');
 app.disable('ETag')
  
  // app.use((req, res, next) => {
  //   res.header('content-type', 'text/plain');
  //   next();
  // });
  

  app.post(`/validate`, (req, res)=> {
    const substr = "library";
    //Validar esquema
    
    console.log('----------Llega request----------');
   
    const response = {
      systemError : "Error de prueba"
        };
        res.setHeader('content-type','text/plain');
    res.send({apiVersion: 'externaldata.gatekeeper.sh/v1beta1', kind : 'ProviderResponse', response});
  });


  http.createServer(app).listen(3000).on('error', (err) => {
    if (err) {
      console.log('Cant start server'+err)
      process.exit(1);
    }
  });

  const options = {
    key: fs.readFileSync('certs/server.key'),
    cert: fs.readFileSync('certs/server.crt')
  }

  https.createServer(options, app).listen(3443).on('error', (err) => {
    if (err) {
      console.log("Cant start https server" + err)
      process.exit(1);
    }
  });

  const loaders = [];

  async.series(loaders,
    (err, result) => {
      if (err) {
        console.log('Cant load resources'+err);
        setTimeout(1000, process.exit(101));
      } else {
        console.log('Server started at port ' + 3000);
        started = true;
        done(null,started);
      }
    });
  
}

function active () {
  return started;
}

exports.start = start;
exports.active = active;