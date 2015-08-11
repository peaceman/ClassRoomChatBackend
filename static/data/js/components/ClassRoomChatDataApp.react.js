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
            <div className="row">
                <div className="col-md-6">
                    <div className="page-header">
                        <h4>PhoneDataList</h4>
                    </div>
                    <PhoneDataList phoneDataItems={this.state.allPhoneDataItems} />
                </div>

                <div className="col-md-6">
                    <div className="page-header">
                        <h4>PhoneDataDetails</h4>
                    </div>
                    <PhoneDataDetails phoneDataItem={this.state.selectedPhoneDataItem} />
                </div>
            </div>
        );
    }
});

module.exports = ClassRoomChatDataApp;
