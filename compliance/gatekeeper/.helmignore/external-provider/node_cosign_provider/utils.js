const {v4 : uuidv4} = require('uuid')
function createRequest(body){
var msgid = uuidv4();
  return {
    uri: process.env.METRICS_URI,
    body : JSON.stringify(body),
    method: 'POST',
    headers : { 'Content-Type': 'application/json', 'Accept': 'application/json', 'msgid' : msgid}
  }
}


module.exports.createRequest = createRequest;