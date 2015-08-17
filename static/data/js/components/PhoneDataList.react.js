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
            <table className="table table-condensed table-hover">
                <thead>
                    <tr>
                        <th>Email</th>
                        <th>Number</th>
                        <th>Model</th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    {phoneDataNodes}
                </tbody>
            </table>
        );
    }
});

module.exports = PhoneDataList;
