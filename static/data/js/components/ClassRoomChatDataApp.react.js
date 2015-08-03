var React = require('react');
var PhoneDataStore = require('../stores/PhoneDataStore');
var PhoneDataList = require('./PhoneDataList.react');
var PhoneDataDetails = require('./PhoneDataDetails.react');

function getPhoneDataState() {
    return {
        allPhoneDataItems: PhoneDataStore.getAll(),
        selectedPhoneDataItem: PhoneDataStore.getSelected()
    };
}

var ClassRoomChatDataApp = React.createClass({
    getInitialState: function() {
        return getPhoneDataState();
    },
    componentDidMount: function() {
        PhoneDataStore.addChangeListener(this._onChange);
    },
    componentWillUnmount: function() {
        PhoneDataStore.removeChangeListener(this._onChange);
    },
    _onChange: function() {
        this.setState(getPhoneDataState());
    },
    render: function() {
        return (
            <div>
                <PhoneDataList phoneDataItems={this.state.allPhoneDataItems} />
                <PhoneDataDetails phoneDataItem={this.state.selectedPhoneDataItem} />
            </div>
        );
    }
});

module.exports = ClassRoomChatDataApp;
