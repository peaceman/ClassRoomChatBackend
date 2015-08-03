var PhoneDataActions = require('../actions/PhoneDataActions');

var PhoneDataAPI = {
    init: function(apiEndpoint) {
        var ws = new WebSocket(apiEndpoint);

        ws.onmessage = function(event) {
            var phoneData = JSON.parse(event.data);
            PhoneDataActions.create(phoneData);
        }
    }
};

module.exports = PhoneDataAPI;
