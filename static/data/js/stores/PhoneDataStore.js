var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('events').EventEmitter;
var PhoneDataConstants = require('../constants/PhoneDataConstants');
var assign = require('object-assign');

var CHANGE_EVENT = 'phonedatastore-change';

var _phoneDataItems = {};
var _selectedPhoneDataId = null;

function create(phoneData) {
    // Using the current timestamp + random number in place of a real id.
    var id = (+new Date() + Math.floor(Math.random() * 999999)).toString(36);
    _phoneDataItems[id] = assign({}, phoneData, {id: id});
}

function destroy(id) {
    if (_selectedPhoneDataId === id) {
        _selectedPhoneDataId = null;
    }

    delete _phoneDataItems[id];
}

function select(id) {
    if (_selectedPhoneDataId === id) {
        _selectedPhoneDataId = null;
    } else {
        _selectedPhoneDataId = id;
    }
}

var PhoneDataStore = assign({}, EventEmitter.prototype, {
    getAll: function() {
        return _phoneDataItems;
    },
    getSelected: function() {
        return _selectedPhoneDataId === null
            ? null
            : _phoneDataItems[_selectedPhoneDataId];
    },
    emitChange: function() {
        this.emit(CHANGE_EVENT);
    },
    addChangeListener: function(callback) {
        this.on(CHANGE_EVENT, callback);
    },
    removeChangeListener: function(callback) {
        this.removeListener(CHANGE_EVENT, callback);
    }
});

AppDispatcher.register(function(action) {
    switch(action.actionType) {
        case PhoneDataConstants.PHONEDATA_CREATE:
            create(action.phoneData);
            PhoneDataStore.emitChange();
            break;
        case PhoneDataConstants.PHONEDATA_SELECT:
            select(action.id);
            PhoneDataStore.emitChange();
            break;
        case PhoneDataConstants.PHONEDATA_DESTROY:
            destroy(action.id);
            PhoneDataStore.emitChange();
            break;
        default:
            // no op
    }
});

module.exports = PhoneDataStore;
