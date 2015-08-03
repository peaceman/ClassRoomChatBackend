var React = require('react');
var PhoneDataListEntry = require('./PhoneDataListEntry.react');

var PhoneDataList = React.createClass({
    propTypes: {
        phoneDataItems: React.PropTypes.object.isRequired
    },
    render: function() {
        var phoneDataNodes = [];

        for (var key in this.props.phoneDataItems) {
            phoneDataNodes.push(<PhoneDataListEntry key={key} phoneData={this.props.phoneDataItems[key]} />);
        }

        return (
            <ul id="phone-data-list">{phoneDataNodes}</ul>
        );
    }
});

module.exports = PhoneDataList;
