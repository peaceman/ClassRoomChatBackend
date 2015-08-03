var React = require('react');
var PhoneDataActions = require('../actions/PhoneDataActions');

var PhoneDataListEntry = React.createClass({
    propTypes: {
        phoneData: React.PropTypes.object.isRequired
    },
    render: function() {
        var phone = this.props.phoneData;

        return (
            <li>
                <span onClick={this._onToggleSelection}>
                    {phone.PhoneNumber} >> {phone.Build.Manufacturer} {phone.Build.Model}
                </span>
                <button onClick={this._onDestroyClick}>Delete</button>
            </li>
        );
    },
    _onToggleSelection: function() {
        PhoneDataActions.toggleSelection(this.props.phoneData.id);
    },
    _onDestroyClick: function() {
        PhoneDataActions.destroy(this.props.phoneData.id);
    }
});

module.exports = PhoneDataListEntry;
