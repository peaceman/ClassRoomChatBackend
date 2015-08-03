var AppDispatcher = require('../dispatcher/AppDispatcher');
var PhoneDataConstants = require('../constants/PhoneDataConstants');

var PhoneDataActions = {
    create: function(phoneData) {
        AppDispatcher.dispatch({
            actionType: PhoneDataConstants.PHONEDATA_CREATE,
            phoneData: phoneData
        });
    },

    destroy: function(id) {
        AppDispatcher.dispatch({
            actionType: PhoneDataConstants.PHONEDATA_DESTROY,
            id: id
        });
    },

    toggleSelection: function(id) {
        AppDispatcher.dispatch({
            actionType: PhoneDataConstants.PHONEDATA_SELECT,
            id: id
        });
    }
};

module.exports = PhoneDataActions;
